package lifecycle

import (
	"errors"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

// TestPlugin is a test plugin that implements all stage interfaces.
type TestPlugin struct {
	name        string
	priority    int
	stagesRun   []Stage
	shouldError Stage
	errorMsg    string
	configureFn func(*Manager) error
	validateFn  func(*Manager) error
	globFn      func(*Manager) error
	loadFn      func(*Manager) error
	transformFn func(*Manager) error
	renderFn    func(*Manager) error
	collectFn   func(*Manager) error
	writeFn     func(*Manager) error
	cleanupFn   func(*Manager) error
}

func NewTestPlugin(name string) *TestPlugin {
	return &TestPlugin{
		name:      name,
		stagesRun: make([]Stage, 0),
	}
}

func (p *TestPlugin) Name() string { return p.name }

func (p *TestPlugin) Priority(_ Stage) int { return p.priority }

func (p *TestPlugin) Configure(m *Manager) error {
	p.stagesRun = append(p.stagesRun, StageConfigure)
	if p.shouldError == StageConfigure {
		return errors.New(p.errorMsg)
	}
	if p.configureFn != nil {
		return p.configureFn(m)
	}
	return nil
}

func (p *TestPlugin) Validate(m *Manager) error {
	p.stagesRun = append(p.stagesRun, StageValidate)
	if p.shouldError == StageValidate {
		return errors.New(p.errorMsg)
	}
	if p.validateFn != nil {
		return p.validateFn(m)
	}
	return nil
}

func (p *TestPlugin) Glob(m *Manager) error {
	p.stagesRun = append(p.stagesRun, StageGlob)
	if p.shouldError == StageGlob {
		return errors.New(p.errorMsg)
	}
	if p.globFn != nil {
		return p.globFn(m)
	}
	return nil
}

func (p *TestPlugin) Load(m *Manager) error {
	p.stagesRun = append(p.stagesRun, StageLoad)
	if p.shouldError == StageLoad {
		return errors.New(p.errorMsg)
	}
	if p.loadFn != nil {
		return p.loadFn(m)
	}
	return nil
}

func (p *TestPlugin) Transform(m *Manager) error {
	p.stagesRun = append(p.stagesRun, StageTransform)
	if p.shouldError == StageTransform {
		return errors.New(p.errorMsg)
	}
	if p.transformFn != nil {
		return p.transformFn(m)
	}
	return nil
}

func (p *TestPlugin) Render(m *Manager) error {
	p.stagesRun = append(p.stagesRun, StageRender)
	if p.shouldError == StageRender {
		return errors.New(p.errorMsg)
	}
	if p.renderFn != nil {
		return p.renderFn(m)
	}
	return nil
}

func (p *TestPlugin) Collect(m *Manager) error {
	p.stagesRun = append(p.stagesRun, StageCollect)
	if p.shouldError == StageCollect {
		return errors.New(p.errorMsg)
	}
	if p.collectFn != nil {
		return p.collectFn(m)
	}
	return nil
}

func (p *TestPlugin) Write(m *Manager) error {
	p.stagesRun = append(p.stagesRun, StageWrite)
	if p.shouldError == StageWrite {
		return errors.New(p.errorMsg)
	}
	if p.writeFn != nil {
		return p.writeFn(m)
	}
	return nil
}

func (p *TestPlugin) Cleanup(m *Manager) error {
	p.stagesRun = append(p.stagesRun, StageCleanup)
	if p.shouldError == StageCleanup {
		return errors.New(p.errorMsg)
	}
	if p.cleanupFn != nil {
		return p.cleanupFn(m)
	}
	return nil
}

func TestStageOrder(t *testing.T) {
	expected := []Stage{
		StageConfigure,
		StageValidate,
		StageGlob,
		StageLoad,
		StageTransform,
		StageRender,
		StageCollect,
		StageWrite,
		StageCleanup,
	}

	if len(StageOrder) != len(expected) {
		t.Errorf("StageOrder has %d stages, expected %d", len(StageOrder), len(expected))
	}

	for i, s := range expected {
		if StageOrder[i] != s {
			t.Errorf("StageOrder[%d] = %s, expected %s", i, StageOrder[i], s)
		}
	}
}

func TestStageIndex(t *testing.T) {
	tests := []struct {
		stage Stage
		want  int
	}{
		{StageConfigure, 0},
		{StageValidate, 1},
		{StageGlob, 2},
		{StageLoad, 3},
		{StageTransform, 4},
		{StageRender, 5},
		{StageCollect, 6},
		{StageWrite, 7},
		{StageCleanup, 8},
		{Stage("invalid"), -1},
	}

	for _, tt := range tests {
		got := StageIndex(tt.stage)
		if got != tt.want {
			t.Errorf("StageIndex(%s) = %d, want %d", tt.stage, got, tt.want)
		}
	}
}

func TestIsValidStage(t *testing.T) {
	for _, s := range StageOrder {
		if !IsValidStage(s) {
			t.Errorf("IsValidStage(%s) = false, want true", s)
		}
	}

	if IsValidStage(Stage("invalid")) {
		t.Error("IsValidStage(invalid) = true, want false")
	}
}

func TestStagesUpTo(t *testing.T) {
	stages := StagesUpTo(StageLoad)
	expected := []Stage{StageConfigure, StageValidate, StageGlob, StageLoad}

	if len(stages) != len(expected) {
		t.Errorf("StagesUpTo(StageLoad) returned %d stages, expected %d", len(stages), len(expected))
	}

	for i, s := range expected {
		if stages[i] != s {
			t.Errorf("StagesUpTo(StageLoad)[%d] = %s, expected %s", i, stages[i], s)
		}
	}
}

func TestNewManager(t *testing.T) {
	m := NewManager()

	if m == nil {
		t.Fatal("NewManager() returned nil")
	}

	if m.config == nil {
		t.Error("Manager.config is nil")
	}

	if m.posts == nil {
		t.Error("Manager.posts is nil")
	}

	if m.stagesRun == nil {
		t.Error("Manager.stagesRun is nil")
	}

	if m.cache == nil {
		t.Error("Manager.cache is nil")
	}
}

func TestManagerRegisterPlugin(t *testing.T) {
	m := NewManager()
	p := NewTestPlugin("test")

	m.RegisterPlugin(p)

	plugins := m.Plugins()
	if len(plugins) != 1 {
		t.Errorf("Expected 1 plugin, got %d", len(plugins))
	}

	if plugins[0].Name() != "test" {
		t.Errorf("Expected plugin name 'test', got %s", plugins[0].Name())
	}
}

func TestManagerRunAllStages(t *testing.T) {
	m := NewManager()
	p := NewTestPlugin("test")

	m.RegisterPlugin(p)

	err := m.Run()
	if err != nil {
		t.Errorf("Run() returned error: %v", err)
	}

	// Check all stages were run
	if len(p.stagesRun) != len(StageOrder) {
		t.Errorf("Expected %d stages to run, got %d", len(StageOrder), len(p.stagesRun))
	}

	for i, s := range StageOrder {
		if p.stagesRun[i] != s {
			t.Errorf("Stage %d: expected %s, got %s", i, s, p.stagesRun[i])
		}
	}

	// Check all stages are marked as run
	for _, s := range StageOrder {
		if !m.HasRun(s) {
			t.Errorf("Stage %s not marked as run", s)
		}
	}
}

func TestManagerRunTo(t *testing.T) {
	m := NewManager()
	p := NewTestPlugin("test")

	m.RegisterPlugin(p)

	err := m.RunTo(StageLoad)
	if err != nil {
		t.Errorf("RunTo(StageLoad) returned error: %v", err)
	}

	expected := []Stage{StageConfigure, StageValidate, StageGlob, StageLoad}
	if len(p.stagesRun) != len(expected) {
		t.Errorf("Expected %d stages to run, got %d", len(expected), len(p.stagesRun))
	}

	// Check stages after Load are not marked as run
	if m.HasRun(StageTransform) {
		t.Error("StageTransform should not be marked as run")
	}
}

func TestManagerRunToTwice(t *testing.T) {
	m := NewManager()
	p := NewTestPlugin("test")

	m.RegisterPlugin(p)

	// Run to Load
	err := m.RunTo(StageLoad)
	if err != nil {
		t.Errorf("First RunTo(StageLoad) returned error: %v", err)
	}

	stagesAfterFirst := len(p.stagesRun)

	// Run to Load again - should not re-run stages
	err = m.RunTo(StageLoad)
	if err != nil {
		t.Errorf("Second RunTo(StageLoad) returned error: %v", err)
	}

	if len(p.stagesRun) != stagesAfterFirst {
		t.Errorf("Second RunTo re-ran stages: had %d, now %d", stagesAfterFirst, len(p.stagesRun))
	}
}

func TestManagerRunToContinue(t *testing.T) {
	m := NewManager()
	p := NewTestPlugin("test")

	m.RegisterPlugin(p)

	// Run to Load
	if err := m.RunTo(StageLoad); err != nil {
		t.Fatalf("RunTo(StageLoad) failed: %v", err)
	}
	stagesAfterLoad := len(p.stagesRun)

	// Continue to Render
	if err := m.RunTo(StageRender); err != nil {
		t.Fatalf("RunTo(StageRender) failed: %v", err)
	}

	// Should have run Transform and Render (2 more stages)
	expectedAdditional := 2
	if len(p.stagesRun) != stagesAfterLoad+expectedAdditional {
		t.Errorf("Expected %d additional stages, got %d",
			expectedAdditional, len(p.stagesRun)-stagesAfterLoad)
	}
}

func TestManagerErrorStopsCriticalStage(t *testing.T) {
	m := NewManager()
	p := NewTestPlugin("test")
	p.shouldError = StageLoad
	p.errorMsg = "load error"

	m.RegisterPlugin(p)

	err := m.Run()
	if err == nil {
		t.Error("Expected error from Run(), got nil")
	}

	// Should not have run stages after Load
	for _, s := range p.stagesRun {
		if StageIndex(s) > StageIndex(StageLoad) {
			t.Errorf("Stage %s should not have run after Load error", s)
		}
	}
}

func TestManagerPriorityOrdering(t *testing.T) {
	m := NewManager()

	order := make([]string, 0)

	p1 := NewTestPlugin("last")
	p1.priority = PriorityLast
	p1.configureFn = func(_ *Manager) error {
		order = append(order, "last")
		return nil
	}

	p2 := NewTestPlugin("first")
	p2.priority = PriorityFirst
	p2.configureFn = func(_ *Manager) error {
		order = append(order, "first")
		return nil
	}

	p3 := NewTestPlugin("default")
	p3.priority = PriorityDefault
	p3.configureFn = func(_ *Manager) error {
		order = append(order, "default")
		return nil
	}

	// Register in non-priority order
	m.RegisterPlugin(p1)
	m.RegisterPlugin(p2)
	m.RegisterPlugin(p3)

	if err := m.RunTo(StageConfigure); err != nil {
		t.Fatalf("RunTo(StageConfigure) failed: %v", err)
	}

	expected := []string{"first", "default", "last"}
	if len(order) != len(expected) {
		t.Fatalf("Expected %d plugins to run, got %d", len(expected), len(order))
	}

	for i, name := range expected {
		if order[i] != name {
			t.Errorf("Plugin %d: expected %s, got %s", i, name, order[i])
		}
	}
}

func TestManagerFilter(t *testing.T) {
	m := NewManager()

	title1 := "Post 1"
	title2 := "Post 2"
	title3 := "Post 3"

	m.SetPosts([]*models.Post{
		{Path: "a.md", Title: &title1, Published: true, Tags: []string{"go", "web"}},
		{Path: "b.md", Title: &title2, Published: false, Tags: []string{"go"}},
		{Path: "c.md", Title: &title3, Published: true, Tags: []string{"rust"}},
	})

	tests := []struct {
		expr     string
		expected int
	}{
		{"published==true", 2},
		{"published==false", 1},
		{"tags contains go", 2},
		{"tags contains rust", 1},
		{"tags contains python", 0},
		{"published==true and tags contains go", 1},
		{"published==true or tags contains go", 3},
		{"", 3}, // Empty filter returns all
	}

	for _, tt := range tests {
		posts, err := m.Filter(tt.expr)
		if err != nil {
			t.Errorf("Filter(%q) returned error: %v", tt.expr, err)
			continue
		}
		if len(posts) != tt.expected {
			t.Errorf("Filter(%q) returned %d posts, expected %d", tt.expr, len(posts), tt.expected)
		}
	}
}

func TestManagerMap(t *testing.T) {
	m := NewManager()

	title1 := "Alpha"
	title2 := "Beta"
	title3 := "Gamma"

	m.SetPosts([]*models.Post{
		{Path: "a.md", Title: &title1, Published: true},
		{Path: "b.md", Title: &title2, Published: false},
		{Path: "c.md", Title: &title3, Published: true},
	})

	// Map titles of published posts
	titles, err := m.Map("Title", "published==true", "", false)
	if err != nil {
		t.Fatalf("Map() returned error: %v", err)
	}

	if len(titles) != 2 {
		t.Errorf("Expected 2 titles, got %d", len(titles))
	}
}

func TestManagerMap_TitleUsesSemanticRepresentation(t *testing.T) {
	source := "Good **places** to ==think=="
	post := &models.Post{Title: &source, TitleText: "Good places to think", TitleTextDerived: true}
	m := NewManager()
	m.SetPosts([]*models.Post{post})
	values, err := m.Map("title", "", "", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(values) != 1 || values[0] != "Good places to think" {
		t.Fatalf("Map(title) = %#v", values)
	}
}

func TestManagerCache(t *testing.T) {
	m := NewManager()

	// Set a value
	m.Cache().Set("key", "value")

	// Get it back
	v, ok := m.Cache().Get("key")
	if !ok {
		t.Error("Cache.Get returned false for existing key")
	}
	if v != "value" {
		t.Errorf("Cache.Get returned %v, expected 'value'", v)
	}

	// Delete it
	m.Cache().Delete("key")
	_, ok = m.Cache().Get("key")
	if ok {
		t.Error("Cache.Get returned true after Delete")
	}

	// Clear
	m.Cache().Set("a", 1)
	m.Cache().Set("b", 2)
	m.Cache().Clear()
	_, ok = m.Cache().Get("a")
	if ok {
		t.Error("Cache.Get returned true after Clear")
	}
}

func TestManagerReset(t *testing.T) {
	m := NewManager()
	p := NewTestPlugin("test")
	m.RegisterPlugin(p)

	m.AddPost(&models.Post{Path: "test.md"})
	m.AddFile("test.md")
	m.Cache().Set("key", "value")

	if err := m.RunTo(StageLoad); err != nil {
		t.Fatalf("RunTo(StageLoad) failed: %v", err)
	}

	// Reset
	m.Reset()

	if len(m.Posts()) != 0 {
		t.Error("Posts not cleared after Reset")
	}
	if len(m.Files()) != 0 {
		t.Error("Files not cleared after Reset")
	}
	if m.HasRun(StageLoad) {
		t.Error("StagesRun not cleared after Reset")
	}
	if _, ok := m.Cache().Get("key"); ok {
		t.Error("Cache not cleared after Reset")
	}
}

func TestManagerSetFilesNormalizesAbsoluteContentPaths(t *testing.T) {
	contentDir := t.TempDir()
	absolutePath := filepath.Join(contentDir, "posts", "entry.md")

	m := NewManager()
	m.Config().ContentDir = contentDir
	m.SetFiles([]string{absolutePath})

	if got := m.Files(); len(got) != 1 || got[0] != filepath.ToSlash(filepath.Join("posts", "entry.md")) {
		t.Fatalf("Files() = %v, want [posts/entry.md]", got)
	}
	snapshot := m.ContentDiagnostics()
	if len(snapshot.Entries) != 1 || snapshot.Entries[0].Path != "posts/entry.md" {
		t.Fatalf("content entries = %+v, want relative source path", snapshot.Entries)
	}
}

func TestManagerSetFilesRedactsAbsolutePathsOutsideContentRoot(t *testing.T) {
	contentDir := t.TempDir()
	outsideDir := t.TempDir()
	absolutePath := filepath.Join(outsideDir, "entry.md")

	m := NewManager()
	m.Config().ContentDir = contentDir
	m.SetFiles([]string{absolutePath})

	snapshot := m.ContentDiagnostics()
	if len(snapshot.Entries) != 1 {
		t.Fatalf("content entries = %+v, want one entry", snapshot.Entries)
	}
	path := snapshot.Entries[0].Path
	if filepath.IsAbs(path) || strings.HasPrefix(path, "../") || strings.Contains(path, outsideDir) {
		t.Fatalf("content path exposes outside root: %q", path)
	}
}

func TestManagerResetInvalidatesPostIndex(t *testing.T) {
	m := NewManager()
	post := &models.Post{Slug: "post-a", Href: "/post-a/", Path: "post-a.md"}
	m.SetPosts([]*models.Post{post})

	idx1 := m.PostIndex()
	if got := idx1.LookupBySlug("post-a"); got != post {
		t.Fatalf("PostIndex().LookupBySlug(%q) = %v, want post A", "post-a", got)
	}

	m.Reset()

	if got := m.Posts(); len(got) != 0 {
		t.Fatalf("Posts() after Reset() = %d posts, want empty", len(got))
	}

	idx2 := m.PostIndex()
	if idx2 == idx1 {
		t.Fatal("PostIndex() after Reset() returned the stale index")
	}
	if got := idx2.LookupBySlug("post-a"); got != nil {
		t.Fatalf("PostIndex().LookupBySlug(%q) after Reset() = %v, want nil", "post-a", got)
	}
}

func TestManagerConcurrentProcessing(t *testing.T) {
	m := NewManager()
	m.SetConcurrency(2)

	posts := make([]*models.Post, 10)
	for i := 0; i < 10; i++ {
		posts[i] = &models.Post{Path: "test.md"}
	}
	m.SetPosts(posts)

	processed := make(chan struct{}, 10)

	err := m.ProcessPostsConcurrently(func(_ *models.Post) error {
		processed <- struct{}{}
		return nil
	})

	if err != nil {
		t.Errorf("ProcessPostsConcurrently returned error: %v", err)
	}

	close(processed)
	count := 0
	for range processed {
		count++
	}

	if count != 10 {
		t.Errorf("Expected 10 posts processed, got %d", count)
	}
}

func TestManagerConcurrentProcessingErrorHandling(t *testing.T) {
	m := NewManager()
	m.SetConcurrency(4)

	posts := make([]*models.Post, 10)
	for i := 0; i < 10; i++ {
		posts[i] = &models.Post{Path: "test.md"}
	}
	m.SetPosts(posts)

	// Track processed posts with atomic counter for thread safety
	var callCount int64
	// Track error count atomically
	var errCount int64
	err := m.ProcessPostsConcurrently(func(_ *models.Post) error {
		count := atomic.AddInt64(&callCount, 1)
		// Every 3rd call fails
		if count%3 == 0 {
			atomic.AddInt64(&errCount, 1)
			return errors.New("simulated error")
		}
		return nil
	})

	if err == nil {
		t.Error("Expected error from ProcessPostsConcurrently, got nil")
	}

	// All 10 posts should be processed (error doesn't stop processing)
	finalCount := atomic.LoadInt64(&callCount)
	if finalCount != 10 {
		t.Errorf("Expected 10 posts processed, got %d", finalCount)
	}
}

func TestManagerConcurrentProcessingErrorOrderIsDeterministic(t *testing.T) {
	m := NewManager()
	m.SetConcurrency(2)
	m.SetPosts([]*models.Post{
		{Path: "first.md"},
		{Path: "second.md"},
	})

	firstStarted := make(chan struct{})
	secondFinished := make(chan struct{})
	releaseFirst := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- m.ProcessPostsConcurrently(func(post *models.Post) error {
			if post.Path == "first.md" {
				close(firstStarted)
				<-releaseFirst
				return errors.New("first failure")
			}
			close(secondFinished)
			return errors.New("second failure")
		})
	}()

	<-firstStarted
	<-secondFinished
	close(releaseFirst)
	err := <-done
	if err == nil || !strings.Contains(err.Error(), "first.md") {
		t.Fatalf("ProcessPostsConcurrently() error = %v, want first input error", err)
	}
}

// TestProcessPostsConcurrentlyGoroutineBound verifies that ProcessPostsConcurrently
// uses a bounded worker pool and does not spawn a goroutine per post.
// This is critical for large builds (5k+ posts) to avoid scheduler overhead.
func TestProcessPostsConcurrentlyGoroutineBound(t *testing.T) {
	const postCount = 1000
	const concurrency = 4
	// Allow some overhead for test infrastructure, GC goroutines, etc.
	const maxOverhead = 50

	m := NewManager()
	m.SetConcurrency(concurrency)

	posts := make([]*models.Post, postCount)
	for i := 0; i < postCount; i++ {
		posts[i] = &models.Post{Path: "test.md"}
	}
	m.SetPosts(posts)

	// Capture baseline goroutine count
	baselineGoroutines := runtime.NumGoroutine()

	// Channel to block workers so we can measure goroutine count during processing
	block := make(chan struct{})
	started := make(chan struct{}, postCount)

	// Start processing in a separate goroutine
	done := make(chan error)
	go func() {
		done <- m.ProcessPostsConcurrently(func(_ *models.Post) error {
			started <- struct{}{}
			<-block // Wait for signal to continue
			return nil
		})
	}()

	// Wait for workers to start processing
	// We expect exactly `concurrency` workers to be active
	for i := 0; i < concurrency; i++ {
		<-started
	}

	// Measure goroutine count while workers are blocked
	peakGoroutines := runtime.NumGoroutine()

	// Unblock all workers
	close(block)

	// Wait for completion
	if err := <-done; err != nil {
		t.Fatalf("ProcessPostsConcurrently failed: %v", err)
	}

	// Calculate goroutine increase
	goroutineIncrease := peakGoroutines - baselineGoroutines

	// The increase should be bounded by concurrency + overhead (main goroutine, test goroutine, etc.)
	// It should NOT be close to postCount (which would indicate goroutine-per-post)
	maxExpected := concurrency + maxOverhead
	if goroutineIncrease > maxExpected {
		t.Errorf("Goroutine count increased by %d (peak=%d, baseline=%d), expected <= %d (concurrency=%d + overhead=%d)",
			goroutineIncrease, peakGoroutines, baselineGoroutines, maxExpected, concurrency, maxOverhead)
	}

	// Sanity check: should be significantly less than post count
	if goroutineIncrease > postCount/10 {
		t.Errorf("Goroutine increase %d is too close to post count %d - possible goroutine-per-post leak",
			goroutineIncrease, postCount)
	}

	t.Logf("Goroutine count: baseline=%d, peak=%d, increase=%d (expected <= %d)",
		baselineGoroutines, peakGoroutines, goroutineIncrease, maxExpected)
}

func TestInvalidStage(t *testing.T) {
	m := NewManager()

	err := m.RunTo(Stage("invalid"))
	if err == nil {
		t.Error("Expected error for invalid stage, got nil")
	}
}
