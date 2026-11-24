// Package sourcegit provides the canonical Git state used for source builds.
package sourcegit

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/WaylonWalker/markata/internal/fileidentity"
)

// State describes the checked-out source tree. Dirty is nil when Git state is
// unavailable, rather than pretending that an unavailable tree is clean.
type State struct {
	Commit string
	Dirty  *bool
	// fingerprint is internal build-snapshot data. It is not part of the wire
	// format because a dirty source is never revision-equivalent.
	fingerprint string
}

// ReadOptions controls which ignored files contribute to the source
// fingerprint. Ignored files are build inputs when Git-ignore filtering is
// disabled, but generated ignored trees must not be walked when filtering is
// enabled.
type ReadOptions struct {
	IncludeIgnoredContent bool
	// IgnoredPathFilter reports whether an ignored path is excluded by the
	// caller's build input rules. It is used only by the filtered metadata path
	// to reject Git-only excludes that could hide build inputs.
	IgnoredPathFilter func(path string) (bool, error)
}

// Command creates a Git command rooted at sourceDir using the same safe
// directory handling as the builder-admin source checkout.
func Command(ctx context.Context, sourceDir string, args ...string) *exec.Cmd {
	gitArgs := append([]string{"-c", "safe.directory=" + sourceDir, "-C", sourceDir}, args...)
	return exec.CommandContext(ctx, "git", gitArgs...)
}

// Read returns HEAD and whether Git reports any tracked, staged, deleted,
// untracked, or ignored files in the source tree.
func Read(ctx context.Context, sourceDir string) (State, error) {
	return ReadWithOptions(ctx, sourceDir, ReadOptions{IncludeIgnoredContent: true})
}

// ReadWithOptions returns the source state using the requested ignored-file
// policy.
func ReadWithOptions(ctx context.Context, sourceDir string, options ReadOptions) (State, error) {
	return read(ctx, sourceDir, options, false)
}

// ReadMetadataWithOptions returns a file-identity fingerprint for changed
// files instead of reading their contents. Callers must use this only when
// fileidentity.Available is true and the build environment keeps source files
// on a filesystem whose inode and change-time metadata updates on writes.
//
// This is intended for fast, local equivalence checks. It is not a replacement
// for ReadWithOptions when a content-addressed fingerprint is required.
func ReadMetadataWithOptions(ctx context.Context, sourceDir string, options ReadOptions) (State, error) {
	if !fileidentity.Available {
		return State{}, fmt.Errorf("file identity metadata is unavailable")
	}
	if !options.IncludeIgnoredContent {
		return readFilteredMetadata(ctx, sourceDir, options)
	}
	return read(ctx, sourceDir, options, true)
}

func readFilteredMetadata(ctx context.Context, sourceDir string, options ReadOptions) (State, error) {
	commit, err := Head(ctx, sourceDir)
	if err != nil {
		return State{}, err
	}
	if err := rejectNoOpIndexFlags(ctx, sourceDir); err != nil {
		return State{}, err
	}
	if err := rejectUnsafeGitStatConfig(ctx, sourceDir); err != nil {
		return State{}, err
	}
	status, err := Command(ctx, sourceDir, "status", "--porcelain=v1", "-z", "--ignored=matching", "--untracked-files=all").Output()
	if err != nil {
		return State{}, fmt.Errorf("read Git source status: %w", err)
	}
	if err := rejectUnfilteredIgnoredPaths(status, options.IgnoredPathFilter); err != nil {
		return State{}, err
	}
	statusOutput, sourcePaths := filteredMetadataStatus(status)
	ignoredSources := ignoredMarkdownFiles(status)
	fingerprint, err := metadataSnapshotFingerprint(ctx, sourceDir, statusOutput, sourcePaths, ignoredSources)
	if err != nil {
		return State{}, err
	}
	currentCommit, err := Head(ctx, sourceDir)
	if err != nil {
		return State{}, err
	}
	if currentCommit != commit {
		return State{}, fmt.Errorf("git HEAD changed during source snapshot: %s -> %s", commit, currentCommit)
	}
	dirty := len(statusOutput) > 0 || len(ignoredSources) > 0
	return State{Commit: commit, Dirty: &dirty, fingerprint: fingerprint}, nil
}

func rejectUnsafeGitStatConfig(ctx context.Context, sourceDir string) error {
	trustCtime, present, err := readGitConfig(ctx, sourceDir, "core.trustctime")
	if err != nil {
		return err
	}
	if present {
		value, err := strconv.ParseBool(trustCtime)
		if err != nil || !value {
			return fmt.Errorf("unsupported core.trustctime setting %q", trustCtime)
		}
	}

	checkStat, present, err := readGitConfig(ctx, sourceDir, "core.checkStat")
	if err != nil {
		return err
	}
	if present && checkStat != "default" && checkStat != "full" {
		return fmt.Errorf("unsupported core.checkStat setting %q", checkStat)
	}
	fsmonitor, present, err := readGitConfig(ctx, sourceDir, "core.fsmonitor")
	if err != nil {
		return err
	}
	if present && fsmonitor != "false" && fsmonitor != "0" && !strings.EqualFold(fsmonitor, "off") {
		return fmt.Errorf("unsupported core.fsmonitor setting %q", fsmonitor)
	}
	return nil
}

func readGitConfig(ctx context.Context, sourceDir, key string) (value string, present bool, err error) {
	output, err := Command(ctx, sourceDir, "config", "--get", key).Output()
	if err == nil {
		return strings.TrimSpace(string(output)), true, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return "", false, nil
	}
	return "", false, fmt.Errorf("read Git config %s: %w", key, err)
}

func rejectUnfilteredIgnoredPaths(status []byte, filter func(string) (bool, error)) error {
	if filter == nil {
		return nil
	}
	for _, record := range splitNUL(status) {
		if !strings.HasPrefix(record, "!! ") {
			continue
		}
		path := strings.TrimPrefix(record, "!! ")
		ignoredByBuild, err := filter(path)
		if err != nil {
			return fmt.Errorf("check ignored source path %q: %w", path, err)
		}
		if !ignoredByBuild {
			return fmt.Errorf("git excludes %q but the build input rules do not", path)
		}
	}
	return nil
}

func rejectNoOpIndexFlags(ctx context.Context, sourceDir string) error {
	output, err := Command(ctx, sourceDir, "ls-files", "-v", "-z").Output()
	if err != nil {
		return fmt.Errorf("inspect Git index flags: %w", err)
	}
	for _, record := range splitNUL(output) {
		if len(record) < 3 || record[0] != 'H' || record[1] != ' ' {
			return fmt.Errorf("unsupported Git index flag in source snapshot")
		}
	}
	return nil
}

type metadataSourcePath struct {
	name string
	kind string
}

func filteredMetadataStatus(status []byte) ([]byte, []metadataSourcePath) {
	records := splitNUL(status)
	filtered := make([]byte, 0, len(status))
	paths := make([]metadataSourcePath, 0, len(records))
	for index := 0; index < len(records); index++ {
		record := records[index]
		if len(record) < 3 {
			continue
		}
		code := record[:2]
		if code == "!!" {
			continue
		}
		name := record[3:]
		filtered = append(filtered, record...)
		filtered = append(filtered, 0)
		kind := "tracked"
		if code == "??" {
			kind = "untracked"
		}
		paths = append(paths, metadataSourcePath{name: name, kind: kind})
		if strings.ContainsAny(code, "RC") && index+1 < len(records) {
			index++
			oldName := records[index]
			filtered = append(filtered, records[index]...)
			filtered = append(filtered, 0)
			paths = append(paths, metadataSourcePath{name: oldName, kind: "tracked"})
		}
	}
	return filtered, paths
}

func metadataSnapshotFingerprint(ctx context.Context, sourceDir string, statusOutput []byte, sourcePaths []metadataSourcePath, ignoredSources []string) (string, error) {
	hash := sha256.New()
	_, _ = hash.Write(statusOutput)
	for _, sourcePath := range sourcePaths {
		if err := hashSourceMetadata(ctx, hash, sourceDir, sourcePath.name, sourcePath.kind, ReadOptions{}); err != nil {
			return "", err
		}
	}
	for _, sourcePath := range ignoredSources {
		if err := hashSourceMetadata(ctx, hash, sourceDir, sourcePath, "ignored", ReadOptions{}); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func read(ctx context.Context, sourceDir string, options ReadOptions, metadataOnly bool) (State, error) {
	commit, err := Head(ctx, sourceDir)
	if err != nil {
		return State{}, err
	}
	statusOutput, err := Command(ctx, sourceDir, "status", "--porcelain", "--untracked-files=all").Output()
	if err != nil {
		return State{}, fmt.Errorf("read git worktree status: %w", err)
	}
	dirty := strings.TrimSpace(string(statusOutput)) != ""
	// A build may disable Markata's .gitignore filtering. Treat ignored files
	// as dirty as well, because Git cannot prove that they are not inputs.
	ignoredMode := "--ignored=matching"
	if options.IncludeIgnoredContent {
		ignoredMode = "--ignored"
	}
	ignoredOutput, err := Command(ctx, sourceDir, "status", "--porcelain=v1", "-z", ignoredMode, "--untracked-files=all").Output()
	if err != nil {
		return State{}, fmt.Errorf("read ignored git worktree status: %w", err)
	}
	ignoredSources := ignoredMarkdownFiles(ignoredOutput)
	if options.IncludeIgnoredContent && strings.Contains(string(ignoredOutput), "!! ") {
		dirty = true
	} else if !options.IncludeIgnoredContent && len(ignoredSources) > 0 {
		// An ignored directory such as generated output is not a build input
		// when Git-ignore filtering is enabled. Direct ignored Markdown files
		// remain visible so the source contract can still report them.
		dirty = true
	}
	fingerprint, err := snapshotFingerprint(ctx, sourceDir, statusOutput, ignoredSources, options, metadataOnly)
	if err != nil {
		return State{}, err
	}
	finalCommit, err := Head(ctx, sourceDir)
	if err != nil {
		return State{}, err
	}
	if finalCommit != commit {
		return State{}, fmt.Errorf("source Git HEAD changed while reading status")
	}
	return State{Commit: commit, Dirty: &dirty, fingerprint: fingerprint}, nil
}

// Equal reports whether two complete source snapshots are identical.
func (s State) Equal(other State) bool {
	if s.Commit != other.Commit {
		return false
	}
	if s.fingerprint != other.fingerprint {
		return false
	}
	if s.Dirty == nil || other.Dirty == nil {
		return s.Dirty == nil && other.Dirty == nil
	}
	return *s.Dirty == *other.Dirty
}

// Fingerprint returns the source fingerprint. ReadWithOptions fingerprints
// changed file contents; ReadMetadataWithOptions fingerprints file identities
// and requires filesystem change-time guarantees.
func (s State) Fingerprint() string {
	return s.fingerprint
}

func ignoredMarkdownFiles(status []byte) []string {
	var result []string
	for _, record := range splitNUL(status) {
		if !strings.HasPrefix(record, "!! ") {
			continue
		}
		name := strings.TrimPrefix(record, "!! ")
		// Ignore whole directories without walking them. They are already
		// represented as dirty, and walking generated output trees can dominate
		// build time. Direct ignored Markdown files retain byte-level tracking.
		if !strings.HasSuffix(name, "/") && isMarkdownSource(name) {
			result = append(result, name)
		}
	}
	return result
}

func snapshotFingerprint(ctx context.Context, sourceDir string, statusOutput []byte, ignoredSources []string, options ReadOptions, metadataOnly bool) (string, error) {
	hash := sha256.New()
	_, _ = hash.Write(statusOutput)
	// Hash the changed paths and their current bytes instead of asking Git to
	// render a complete binary patch. The latter is needlessly expensive for a
	// large dirty checkout, while path + bytes preserves the same snapshot
	// distinction for tracked content changes.
	rawDiff, err := Command(ctx, sourceDir, "diff", "HEAD", "--raw", "-z", "--").Output()
	if err != nil {
		return "", fmt.Errorf("read changed source metadata: %w", err)
	}
	_, _ = hash.Write(rawDiff)
	changed, err := Command(ctx, sourceDir, "diff", "HEAD", "--name-only", "-z", "--").Output()
	if err != nil {
		return "", fmt.Errorf("read changed source paths: %w", err)
	}
	for _, name := range splitNUL(changed) {
		if name == "" {
			continue
		}
		if err := hashSourceInput(ctx, hash, sourceDir, name, "tracked", options, metadataOnly); err != nil {
			return "", err
		}
	}
	untracked, err := Command(ctx, sourceDir, "ls-files", "--others", "--exclude-standard", "-z").Output()
	if err != nil {
		return "", fmt.Errorf("read untracked source files: %w", err)
	}
	for _, name := range splitNUL(untracked) {
		if name == "" {
			continue
		}
		if err := hashSourceInput(ctx, hash, sourceDir, name, "untracked", options, metadataOnly); err != nil {
			return "", err
		}
	}
	for _, path := range ignoredSources {
		if err := hashSourceInput(ctx, hash, sourceDir, path, "ignored", options, metadataOnly); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func hashSourceInput(ctx context.Context, hash io.Writer, sourceDir, name, kind string, options ReadOptions, metadataOnly bool) error {
	if metadataOnly {
		return hashSourceMetadata(ctx, hash, sourceDir, name, kind, options)
	}
	return hashSourceFile(ctx, hash, sourceDir, name, kind, options)
}

func hashSourceMetadata(ctx context.Context, hash io.Writer, sourceDir, name, kind string, options ReadOptions) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if _, err := io.WriteString(hash, kind+"\x00"+name+"\x00"); err != nil {
		return fmt.Errorf("hash %s source path %q: %w", kind, name, err)
	}
	path := filepath.Join(sourceDir, filepath.FromSlash(name))
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		if _, err := io.WriteString(hash, "missing\x00"); err != nil {
			return fmt.Errorf("hash missing %s source path %q: %w", kind, name, err)
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("stat %s source file %q: %w", kind, name, err)
	}
	if info.IsDir() {
		// Submodules have independent worktrees. Keep their existing content-aware
		// handling until a metadata-only submodule snapshot is available.
		return hashSubmodule(ctx, hash, path, name, options)
	}
	if err := hashFileIdentity(hash, info, path); err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(path)
		if err != nil {
			return fmt.Errorf("readlink %s source file %q: %w", kind, name, err)
		}
		if _, err := io.WriteString(hash, "symlink\x00"+target+"\x00"); err != nil {
			return fmt.Errorf("hash %s source symlink %q: %w", kind, name, err)
		}
		targetInfo, err := os.Stat(path)
		if err == nil && targetInfo.Mode().IsRegular() {
			return hashFileIdentity(hash, targetInfo, path+" -> "+target)
		}
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("stat %s source symlink target %q: %w", kind, name, err)
		}
		if err == nil && targetInfo.IsDir() {
			return fmt.Errorf("metadata-only source snapshot does not support directory symlink %q", name)
		}
	}
	return nil
}

func hashFileIdentity(hash io.Writer, info os.FileInfo, name string) error {
	identity := fileidentity.Read(info)
	if identity.Inode == 0 {
		return fmt.Errorf("file identity unavailable for source input %q", name)
	}
	_, err := fmt.Fprintf(hash, "identity\x00%d\x00%d\x00%d\x00%d\x00%d\x00%d\x00", identity.Device, identity.Inode, identity.Size, identity.Mode, identity.Modified, identity.Changed)
	if err != nil {
		return fmt.Errorf("hash file identity for source input %q: %w", name, err)
	}
	return nil
}

func hashSourceFile(ctx context.Context, hash io.Writer, sourceDir, name, kind string, options ReadOptions) error {
	if _, err := io.WriteString(hash, kind+"\x00"+name+"\x00"); err != nil {
		return fmt.Errorf("hash %s source file %q: %w", kind, name, err)
	}
	path := filepath.Join(sourceDir, filepath.FromSlash(name))
	info, lstatErr := os.Lstat(path)
	if lstatErr != nil && !os.IsNotExist(lstatErr) {
		return fmt.Errorf("stat %s source file %q: %w", kind, name, lstatErr)
	}
	if lstatErr == nil && info.Mode()&os.ModeSymlink != 0 {
		if _, err := fmt.Fprintf(hash, "mode\x00%#o\x00", info.Mode()); err != nil {
			return fmt.Errorf("hash mode for source file %q: %w", name, err)
		}
		target, err := os.Readlink(path)
		if err != nil {
			return fmt.Errorf("readlink %s source file %q: %w", kind, name, err)
		}
		if _, err := io.WriteString(hash, "symlink\x00"+target+"\x00"); err != nil {
			return fmt.Errorf("hash symlink source file %q: %w", name, err)
		}
		if targetInfo, err := os.Stat(path); err == nil && targetInfo.Mode().IsRegular() {
			return hashFileContents(hash, path, name, kind)
		}
		return nil
	}
	if lstatErr == nil && info.IsDir() {
		return hashSubmodule(ctx, hash, path, name, options)
	}
	if lstatErr == nil {
		if _, err := fmt.Fprintf(hash, "mode\x00%#o\x00", info.Mode()); err != nil {
			return fmt.Errorf("hash mode for source file %q: %w", name, err)
		}
	}
	return hashFileContents(hash, path, name, kind)
}

func hashSubmodule(ctx context.Context, hash io.Writer, path, name string, options ReadOptions) error {
	head, err := Command(ctx, path, "rev-parse", "HEAD").Output()
	if err != nil {
		return fmt.Errorf("read submodule HEAD %q: %w", name, err)
	}
	status, err := Command(ctx, path, "status", "--porcelain", "--untracked-files=all").Output()
	if err != nil {
		return fmt.Errorf("read submodule status %q: %w", name, err)
	}
	diff, err := Command(ctx, path, "diff", "HEAD", "--binary").Output()
	if err != nil {
		return fmt.Errorf("read submodule diff %q: %w", name, err)
	}
	for _, part := range [][]byte{head, status, diff} {
		if _, err := hash.Write(part); err != nil {
			return fmt.Errorf("hash submodule %q: %w", name, err)
		}
	}
	ignoredMode := "--ignored=matching"
	if options.IncludeIgnoredContent {
		ignoredMode = "--ignored"
	}
	ignoredStatus, err := Command(ctx, path, "status", "--porcelain=v1", "-z", ignoredMode, "--untracked-files=all").Output()
	if err != nil {
		return fmt.Errorf("read submodule ignored files %q: %w", name, err)
	}
	ignoredSources := ignoredMarkdownFiles(ignoredStatus)
	if _, err := hash.Write(ignoredStatus); err != nil {
		return fmt.Errorf("hash submodule ignored files %q: %w", name, err)
	}
	for _, child := range ignoredSources {
		if err := hashSourceFile(ctx, hash, path, child, "submodule-ignored", options); err != nil {
			return err
		}
	}
	untracked, err := Command(ctx, path, "ls-files", "--others", "--exclude-standard", "-z").Output()
	if err != nil {
		return fmt.Errorf("read submodule untracked files %q: %w", name, err)
	}
	for _, child := range splitNUL(untracked) {
		if child == "" {
			continue
		}
		if err := hashSourceFile(ctx, hash, path, child, "submodule-untracked", options); err != nil {
			return err
		}
	}
	return nil
}

func hashFileContents(hash io.Writer, path, name, kind string) error {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read %s source file %q: %w", kind, name, err)
	}
	if _, err := io.Copy(hash, file); err != nil {
		_ = file.Close()
		return fmt.Errorf("hash %s source file %q: %w", kind, name, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close %s source file %q: %w", kind, name, err)
	}
	return nil
}

func isMarkdownSource(name string) bool {
	return strings.HasSuffix(name, ".md") || strings.HasSuffix(name, ".markdown") || strings.HasSuffix(name, ".mdx")
}

func splitNUL(value []byte) []string {
	return strings.Split(strings.TrimSuffix(string(value), "\x00"), "\x00")
}

// Head returns only the checked-out commit without scanning worktree status.
func Head(ctx context.Context, sourceDir string) (string, error) {
	output, err := Command(ctx, sourceDir, "rev-parse", "HEAD").Output()
	if err != nil {
		return "", fmt.Errorf("read git HEAD: %w", err)
	}
	commit := strings.TrimSpace(string(output))
	if commit == "" {
		return "", fmt.Errorf("git returned an empty HEAD")
	}
	return commit, nil
}
