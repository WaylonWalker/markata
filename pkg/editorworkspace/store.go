package editorworkspace

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/lsp"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
	"github.com/fsnotify/fsnotify"
)

type sourceStamp struct {
	Modified int64
	Size     int64
}
type metadataEntry struct {
	file File
	info *lsp.PostInfo
}

// StoreStats reports work performed by a long-lived workspace store.
type StoreStats struct{ GlobScans, MetadataLoads, SourceReads uint64 }

// Store retains discovery, metadata and immutable LSP snapshots for one server.
// Close it, or cancel its parent context, when the server shuts down.
type Store struct {
	workspace        *Workspace
	indexRoot        string
	ctx              context.Context
	cancel           context.CancelFunc
	watcher          *fsnotify.Watcher
	done             chan struct{}
	refreshMu        sync.Mutex
	loadMu           sync.Mutex
	mu               sync.RWMutex
	cfg              *models.Config
	lifecycleConfig  *lifecycle.Config
	configRevision   uint64
	revision         uint64
	metadataRevision uint64
	metadataErr      error
	metadataRunning  bool
	jobs             sync.WaitGroup
	paths            map[string]sourceStamp
	entries          map[string]metadataEntry
	index            *lsp.Index
	indexRevision    uint64
	stats            StoreStats
}

// Open discovers content once and starts filesystem invalidation for a server.
func Open(ctx context.Context, root string) (*Store, error) {
	workspace, err := New(root)
	if err != nil {
		return nil, err
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	child, cancel := context.WithCancel(ctx)
	indexRoot, err := filepath.Abs(root)
	if err != nil {
		cancel()
		_ = watcher.Close()
		return nil, err
	}
	store := &Store{indexRoot: indexRoot, workspace: workspace, ctx: child, cancel: cancel, watcher: watcher, done: make(chan struct{}), paths: make(map[string]sourceStamp), entries: make(map[string]metadataEntry), index: lsp.NewIndex(log.Default())}
	if err = store.refresh(nil, false); err != nil {
		cancel()
		_ = watcher.Close()
		return nil, err
	}
	go store.watch()
	return store, nil
}

// Close releases the watcher and reconciliation loop. It is safe to call repeatedly.
func (s *Store) Close() error {
	s.mu.Lock()
	s.cancel()
	s.mu.Unlock()
	<-s.done
	s.jobs.Wait()
	s.loadMu.Lock()
	defer s.loadMu.Unlock()
	return nil
}

// Config returns the immutable current site configuration.
func (s *Store) Config() (*models.Config, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg, nil
}

// Configuration returns one immutable configuration version atomically.
func (s *Store) Configuration() (configuration *models.Config, revision uint64) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg, s.configRevision
}

// ConfigRevision changes only when configuration or appearance dependencies change.
func (s *Store) ConfigRevision() uint64 { s.mu.RLock(); defer s.mu.RUnlock(); return s.configRevision }

// Stats returns operation counts without changing cached state.
func (s *Store) Stats() StoreStats { s.mu.RLock(); defer s.mu.RUnlock(); return s.stats }

// Paths returns cached path membership without scanning or reading document bodies.
func (s *Store) Paths() ([]File, error) {
	if err := s.ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	files := make([]File, 0, len(s.paths))
	for path := range s.paths {
		files = append(files, File{Path: path})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

// Read checks cached membership and current filesystem containment before reading.
func (s *Store) Read(path string) (string, error) {
	if err := s.ctx.Err(); err != nil {
		return "", err
	}
	s.mu.Lock()
	_, ok := s.paths[path]
	if ok {
		s.stats.SourceReads++
	}
	s.mu.Unlock()
	if !ok {
		return "", fmt.Errorf("source file is not in the workspace content globs")
	}
	return s.workspace.readSource(path)
}

// Write atomically replaces an already-discovered source file.
func (s *Store) Write(path, source string) error {
	if err := s.ctx.Err(); err != nil {
		return err
	}
	if len(source) > 200000 {
		return fmt.Errorf("source exceeds 200,000 bytes")
	}
	s.mu.RLock()
	_, ok := s.paths[path]
	s.mu.RUnlock()
	if !ok {
		return fmt.Errorf("source file is not in the workspace content globs")
	}
	return s.workspace.writeSource(path, source)
}

// Create adds a Markdown file in a directory that already contains discovered content.
func (s *Store) Create(path, source string) error {
	if err := s.ctx.Err(); err != nil {
		return err
	}
	if len(source) > 200000 {
		return fmt.Errorf("source exceeds 200,000 bytes")
	}
	local := filepath.Clean(filepath.FromSlash(path))
	if !filepath.IsLocal(local) || !isMarkdownSource(local) {
		return fmt.Errorf("new workspace path must be a local Markdown file")
	}
	path = filepath.ToSlash(local)
	directory := filepath.ToSlash(filepath.Dir(local))
	s.mu.RLock()
	_, exists := s.paths[path]
	allowedDirectory := false
	for existing := range s.paths {
		if filepath.ToSlash(filepath.Dir(filepath.FromSlash(existing))) == directory {
			allowedDirectory = true
			break
		}
	}
	s.mu.RUnlock()
	if exists {
		return fmt.Errorf("workspace file already exists")
	}
	if !allowedDirectory {
		return fmt.Errorf("new files must use a directory with discovered workspace content")
	}
	if err := s.workspace.Create(path, source); err != nil {
		return err
	}
	return s.refresh(nil, false)
}

// Snapshot shares parsed workspace information, without sharing mutable drafts.
func (s *Store) Snapshot() (index *lsp.Index, revision uint64) {
	s.mu.Lock()
	if s.metadataRevision != s.revision && !s.metadataRunning && s.ctx.Err() == nil {
		s.metadataRunning = true
		s.jobs.Add(1)
		go func() {
			defer s.jobs.Done()
			if _, err := s.Files(); err != nil && s.ctx.Err() == nil {
				log.Printf("editor workspace metadata: %v", err)
			}
			s.mu.Lock()
			s.metadataRunning = false
			s.mu.Unlock()
		}()
	}
	index, revision = s.index, s.indexRevision
	s.mu.Unlock()
	return index, revision
}

// Files loads only invalidated metadata and reuses all unchanged parsed entries.
func (s *Store) Files() ([]File, error) {
	s.loadMu.Lock()
	defer s.loadMu.Unlock()
	for {
		if err := s.ctx.Err(); err != nil {
			return nil, err
		}
		s.mu.RLock()
		revision := s.revision
		if s.metadataRevision == revision {
			files, err := s.filesLocked(), s.metadataErr
			s.mu.RUnlock()
			return files, err
		}
		cfg := s.lifecycleConfig
		changed := make([]string, 0)
		for path := range s.paths {
			if _, ok := s.entries[path]; !ok {
				changed = append(changed, filepath.FromSlash(path))
			}
		}
		s.mu.RUnlock()
		sort.Strings(changed)
		manager := lifecycle.NewManager()
		manager.SetConfig(cfg)
		manager.SetFiles(changed)
		loader := plugins.NewLoadPlugin()
		err := loader.Configure(manager)
		if err == nil {
			err = loader.Load(manager)
		}
		if cancelErr := s.ctx.Err(); cancelErr != nil {
			return nil, cancelErr
		}
		parsed := make(map[string]metadataEntry)
		for _, post := range manager.Posts() {
			path := filepath.ToSlash(post.Path)
			parsed[path] = metadataEntry{file: File{Path: path, Slug: post.Slug, Title: post.Title, Tags: append([]string(nil), post.Tags...)}, info: lsp.LoadedPostInfo(s.indexRoot, post)}
		}
		s.mu.Lock()
		s.stats.MetadataLoads += uint64(len(changed))
		if s.revision != revision {
			s.mu.Unlock()
			continue
		}
		for path, entry := range parsed {
			s.entries[path] = entry
		}
		infos := make([]*lsp.PostInfo, 0, len(s.entries))
		for _, entry := range s.entries {
			infos = append(infos, entry.info)
		}
		s.mu.Unlock()
		sort.Slice(infos, func(i, j int) bool { return infos[i].Path < infos[j].Path })
		index := lsp.NewWorkspaceIndex(log.Default(), s.indexRoot, infos)
		s.mu.Lock()
		if s.revision != revision {
			s.mu.Unlock()
			continue
		}
		if cancelErr := s.ctx.Err(); cancelErr != nil {
			s.mu.Unlock()
			return nil, cancelErr
		}
		s.index = index
		s.indexRevision++
		s.metadataRevision = revision
		s.metadataErr = err
		files := s.filesLocked()
		s.mu.Unlock()
		return files, err
	}
}

func (s *Store) filesLocked() []File {
	files := make([]File, 0, len(s.entries))
	for _, entry := range s.entries {
		file := entry.file
		file.Tags = append([]string(nil), file.Tags...)
		if file.Title != nil {
			title := *file.Title
			file.Title = &title
		}
		files = append(files, file)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files
}

func (s *Store) refresh(events map[string]bool, forceConfig bool) error {
	s.refreshMu.Lock()
	defer s.refreshMu.Unlock()
	if err := s.ctx.Err(); err != nil {
		return err
	}
	manager, err := s.workspace.manager()
	if err != nil {
		return err
	}
	_, ok := manager.Config().Extra["models_config"].(*models.Config)
	if !ok {
		return fmt.Errorf("workspace configuration missing")
	}
	// Preserve standalone defaults when the workspace has no config file.
	authored, err := s.workspace.Config()
	if err != nil {
		return err
	}
	paths := make(map[string]sourceStamp)
	for _, path := range manager.Files() {
		info, statErr := os.Stat(filepath.Join(s.workspace.root, path))
		if statErr == nil {
			paths[filepath.ToSlash(path)] = sourceStamp{Modified: info.ModTime().UnixNano(), Size: info.Size()}
		}
	}
	s.mu.Lock()
	s.stats.GlobScans++
	configChanged := s.configRevision == 0 || forceConfig || !reflect.DeepEqual(s.cfg, authored)
	changed := configChanged || !reflect.DeepEqual(s.paths, paths)
	for path := range events {
		if _, ok := paths[path]; ok {
			delete(s.entries, path)
			changed = true
		}
	}
	for path, previous := range s.paths {
		if next, ok := paths[path]; !ok || next != previous {
			delete(s.entries, path)
		}
	}
	if configChanged {
		s.entries = make(map[string]metadataEntry)
		s.cfg = authored
		s.configRevision++
	}
	s.paths = paths
	s.lifecycleConfig = manager.Config()
	if changed {
		s.revision++
		s.metadataErr = nil
	}
	s.mu.Unlock()
	return s.addWatches()
}

func (s *Store) addWatches() error {
	return filepath.WalkDir(s.workspace.root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if !entry.IsDir() {
			return nil
		}
		if path != s.workspace.root {
			switch entry.Name() {
			case ".git", "node_modules", "output", "cache", ".markata", ".markata-cache", ".venv", "venv":
				return filepath.SkipDir
			}
		}
		if err := s.watcher.Add(path); err != nil {
			return fmt.Errorf("watch workspace directory: %w", err)
		}
		return nil
	})
}

func (s *Store) watch() {
	defer close(s.done)
	defer func() { _ = s.watcher.Close() }()
	reconcile := time.NewTicker(30 * time.Second)
	defer reconcile.Stop()
	timer := time.NewTimer(time.Hour)
	if !timer.Stop() {
		<-timer.C
	}
	defer timer.Stop()
	var pending <-chan time.Time
	events := make(map[string]bool)
	forceConfig := false
	for {
		select {
		case <-s.ctx.Done():
			return
		case event, ok := <-s.watcher.Events:
			if !ok {
				return
			}
			relative, err := filepath.Rel(s.workspace.root, event.Name)
			if err == nil {
				events[filepath.ToSlash(relative)] = true
			}
			extension := strings.ToLower(filepath.Ext(event.Name))
			switch extension {
			case ".toml", ".yaml", ".yml", ".json", ".css", ".woff", ".woff2", ".ttf", ".otf":
				forceConfig = true
			}
			timer.Reset(100 * time.Millisecond)
			pending = timer.C
		case _, ok := <-s.watcher.Errors:
			if !ok {
				return
			}
			forceConfig = true
			timer.Reset(100 * time.Millisecond)
			pending = timer.C
		case <-pending:
			if err := s.refresh(events, forceConfig); err != nil && s.ctx.Err() == nil {
				log.Printf("editor workspace refresh: %v", err)
			}
			events = make(map[string]bool)
			forceConfig = false
			pending = nil
			s.Snapshot()
		case <-reconcile.C:
			if err := s.refresh(nil, false); err != nil && s.ctx.Err() == nil {
				log.Printf("editor workspace reconcile: %v", err)
			}
			s.Snapshot()
		}
	}
}

// SourceVersion identifies the exact bytes read, independent of filesystem clocks.
func SourceVersion(source string) string {
	sum := sha256.Sum256([]byte(source))
	return fmt.Sprintf("%x", sum)
}
