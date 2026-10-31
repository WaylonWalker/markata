package cmd

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/searchapi"
	"github.com/WaylonWalker/markata/pkg/servecontrol"
	"github.com/WaylonWalker/markata/pkg/servefix"
	"github.com/WaylonWalker/markata/pkg/servetui"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// HTTP server timeout constants.
const (
	rebuildPageAction       = "rebuild-page"
	rebuildPageLegacyAction = "rebuild_page"
	serverReadHeaderTimeout = 10 * time.Second
	serveIndexFile          = "index.html"
)

// searchResult represents a search result for the no-JS fallback search.
type searchResult struct {
	Title       string
	Description string
	URL         string
	Score       int
}

var (
	// servePort is the port to serve on.
	servePort int

	// serveHost is the host to serve on.
	serveHost string
	serveBind string

	// serveSourceFile is the Markdown file passed to `serve <file>`. When set,
	// every build and rebuild uses single-page mode.
	serveSourceFile string

	// serveWatch enables file watching (default true).
	serveWatch bool

	// serveNoWatch disables file watching (legacy flag for backward compatibility).
	serveNoWatch bool

	// serveFast skips expensive non-essential plugins (minification, CSS purging)
	// for faster development iteration during serve.
	serveFast bool

	// serveIncremental reuses unchanged build results while retaining normal
	// production output processing.
	serveIncremental bool
	serveNoTUI       bool
	serveAdmin       bool
	serveNoAdmin     bool

	// serveOutputPath is the output directory path for filtering watch events.
	serveOutputPath string

	// serveInternalPaths are generated/cache directories that should never trigger rebuilds.
	serveInternalPaths []string

	// serveSearchEndpoint is the frontend search endpoint configured for this serve session.
	serveSearchEndpoint string

	// isRebuilding tracks whether a rebuild is in progress to avoid event loops.
	isRebuilding atomic.Bool

	// rebuildPending tracks whether changes happened during a rebuild.
	rebuildPending atomic.Bool

	// buildStatus holds the current build state for serve mode.
	buildStatus atomic.Value

	// serveChangedPaths tracks changed paths between rebuilds.
	serveChangedPathsMu   sync.Mutex
	serveChangedPaths     = make(map[string]fsnotify.Op)
	serveForceFullRebuild bool
	serveGlobDirty        bool
	serveCacheMu          sync.Mutex
	serveCache            *buildcache.Cache
	servePostsMu          sync.Mutex
	servePosts            map[string]*models.Post
	serveSearchPostsMu    sync.RWMutex
	serveSearchPosts      map[string]*models.Post
)

type expectedServeFixWrite struct {
	before  string
	after   string
	expires time.Time
}

var serveFixExpectedWrites = struct {
	sync.Mutex
	files map[string]expectedServeFixWrite
}{files: make(map[string]expectedServeFixWrite)}

const (
	serveBuildAction       = "build"
	markdownExtension      = ".markdown"
	markdownShortExtension = ".md"
	buildStatusBuilding    = "building"
	buildStatusSuccess     = "success"
	buildStatusError       = "error"
	buildStatusWarning     = "warning"

	buildStatusEventPrefix = "status:"
)

type BuildStatus struct {
	Status         string `json:"status"`
	Message        string `json:"message,omitempty"`
	LicenseWarning string `json:"license_warning,omitempty"`
}

func setBuildStatus(status, message, licenseWarning string) {
	buildStatus.Store(BuildStatus{Status: status, Message: message, LicenseWarning: licenseWarning})
}

func getBuildStatus() BuildStatus {
	if status, ok := buildStatus.Load().(BuildStatus); ok {
		return status
	}
	return BuildStatus{Status: buildStatusBuilding}
}

func buildStatusPayload(status BuildStatus) string {
	payload := struct {
		Status         string `json:"status"`
		Message        string `json:"message,omitempty"`
		LicenseWarning string `json:"license_warning,omitempty"`
	}{
		Status:         status.Status,
		Message:        status.Message,
		LicenseWarning: status.LicenseWarning,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return `{"status":"building"}`
	}

	return string(data)
}

func getModelsConfig(m *lifecycle.Manager) *models.Config {
	if m == nil || m.Config() == nil || m.Config().Extra == nil {
		return nil
	}
	raw, ok := m.Config().Extra["models_config"]
	if !ok {
		return nil
	}
	cfg, ok := raw.(*models.Config)
	if !ok {
		return nil
	}
	return cfg
}

// serveCmd represents the serve command.
var serveCmd = &cobra.Command{
	Use:        "serve [markdown-file]",
	Aliases:    []string{"s", "serv"},
	SuggestFor: []string{"dev", "preview"},
	Short:      "Build and serve locally with live reload",
	Long: `Serve builds the site and starts a local development server with file watching.

Features:
  - Initial build of the site
  - HTTP server serving the output directory
  - File watching for content, templates, and config
  - Automatic rebuild on file changes
  - Live reload support (injects reload script into HTML)
  - Theme picker "Bake" button: writes the chosen palette, style, font, and
    text size into the config file that holds [markata-go.theme]

Fast mode:
  --fast       Skip minification (JS/CSS) and CSS purging for faster builds.
               Applied to both the initial build and all subsequent rebuilds.

Incremental mode:
  --incremental Reuse unchanged posts while retaining normal output processing.

Example usage:
  markata-go serve              # Serve on localhost:8000 with file watching
  markata-go serve post.md      # Serve just post.md at the root index
  markata-go serve --fast       # Serve with fast mode (skip minification)
  markata-go serve -p 3000      # Serve on localhost:3000
  markata-go serve --bind 0.0.0.0 # Listen on all interfaces
  markata-go serve --no-admin  # Disable the local Control Center
  markata-go serve -m fast.toml # Serve with merged config overrides
  markata-go serve --watch      # Explicitly enable file watching (default)
  markata-go serve --watch=false # Disable file watching
  markata-go serve --no-watch   # Serve without file watching (legacy flag)
  markata-go serve -v           # Serve with verbose output`,
	Args: cobra.MaximumNArgs(1),
	RunE: runServeCommand,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8000, "port to serve on")
	serveCmd.Flags().StringVar(&serveHost, "host", "localhost", "host to serve on")
	serveCmd.Flags().StringVar(&serveBind, "bind", "localhost", "address to listen on (alias for --host)")
	serveCmd.Flags().BoolVar(&serveWatch, "watch", true, "enable file watching")
	serveCmd.Flags().BoolVar(&serveNoWatch, "no-watch", false, "disable file watching (legacy, overrides --watch)")
	serveCmd.Flags().BoolVar(&serveFast, "fast", false, "skip minification and CSS purging for faster builds")
	serveCmd.Flags().BoolVar(&serveIncremental, "incremental", false, "reuse unchanged posts without skipping production output processing")
	serveCmd.Flags().BoolVar(&serveNoTUI, "no-tui", false, "use line-oriented output instead of the interactive dashboard")
	serveCmd.Flags().BoolVar(&serveAdmin, "admin", false, "explicitly enable the local Control Center (loopback only; enabled by default on loopback)")
	serveCmd.Flags().BoolVar(&serveNoAdmin, "no-admin", false, "disable the local Control Center at /_markata/")
}

func runServeCommand(cmd *cobra.Command, args []string) error {
	clearExpectedServeFixWrites()
	defer clearExpectedServeFixWrites()
	currentCmd = cmd
	if err := validateServePort(servePort); err != nil {
		return err
	}
	if err := resolveServeHost(cmd); err != nil {
		return err
	}
	adminEnabled, err := resolveServeAdmin(serveHost, serveAdmin, serveNoAdmin)
	if err != nil {
		return err
	}
	serveSourceFile = ""
	if len(args) == 1 {
		serveSourceFile = args[0]
	}
	runtime := servecontrol.NewRuntime()
	tui := shouldRunServeTUI(serveNoTUI, noInput, inReader(), outWriter())
	paletteUpdates := make(chan *palettes.Palette, 1)
	setServeControl(runtime, tui, paletteUpdates)
	defer clearServeControl()
	configureServeLogger()
	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupServeSignals(cancel)

	// Create manager (config and plugin setup)
	m, err := createServeManager()
	if err != nil {
		runtime.SetServer(servecontrol.ServerState{Status: servecontrol.StateFailed, Message: err.Error()})
		return fmt.Errorf("initialization failed: %w", err)
	}
	// Do not let a previous serve session's post pointers become visible while
	// the new manager is still populating and transforming its posts.
	setServePosts(nil)
	setServeSearchPosts(nil)
	configureServeLoggerForManager(m)
	palette := publishServePalette(m)

	configureServeBuildModes(m)

	// Determine output directory
	outputPath, absOutputPath := resolveServeOutputPath(m)
	serveOutputPath = absOutputPath
	initServeInternalPaths()

	var wg sync.WaitGroup
	// Create HTTP server with search API
	addr := net.JoinHostPort(serveHost, strconv.Itoa(servePort))
	searchEndpoint, searchHandlerPath := configuredSearchEndpoints(getModelsConfig(m))
	serveSearchEndpoint = searchEndpoint
	handler := mountServeAdmin(createHandler(outputPath, m, searchHandlerPath), runtime, m.Config().ContentDir, adminEnabled)

	listener, err := listenWithPortRecovery(cmd, addr)
	if err != nil {
		return diagnoseServeListenError(err, addr)
	}
	defer listener.Close()
	addr = listener.Addr().String()

	rebuildCh, closeWatcher, err := setupWatcher(ctx, m, &wg)
	if err != nil {
		return err
	}
	defer closeWatcher()
	serveRequestFullRebuild = func() {
		serveChangedPathsMu.Lock()
		serveForceFullRebuild = true
		serveChangedPathsMu.Unlock()
		if isRebuilding.Load() {
			rebuildPending.Store(true)
			return
		}
		select {
		case rebuildCh <- struct{}{}:
		default:
		}
	}
	defer func() { serveRequestFullRebuild = nil }()
	configureServeActionHandler(ctx, runtime, rebuildCh, &wg)

	if rebuildCh != nil {
		infof("Watching for file changes...")
	}
	infof("\nServing at http://%s", addr)
	adminURL := ""
	if adminEnabled {
		adminURL = "http://" + addr + "/_markata/"
		infof("Admin at %s", adminURL)
	}
	infof("Press Ctrl+C to stop")
	server, serverErr := startHTTPServer(listener, handler)
	runtime.SetServer(servecontrol.ServerState{Status: servecontrol.StateSuccess, Address: addr, AdminURL: adminURL})

	startInitialBuild(m, rebuildCh, &wg)
	tuiErrors := make(chan error, 1)
	if tui {
		go func() {
			if tuiErr := servetui.RunWithSourceRoot(ctx, runtime, palette, m.Config().ContentDir, paletteUpdates); tuiErr != nil {
				runtime.AddDiagnostic(servecontrol.Diagnostic{Code: "serve.tui_error", Severity: "error", Message: tuiErr.Error(), SuggestedFix: "Run markata-go serve --no-tui to inspect the session as plain logs."})
				tuiErrors <- tuiErr
			}
			cancel()
		}()
	}

	if err := waitForShutdown(ctx, server, serverErr); err != nil {
		runtime.SetServer(servecontrol.ServerState{Status: servecontrol.StateFailed, Address: addr, Message: err.Error()})
		return err
	}
	runtime.SetServer(servecontrol.ServerState{Status: servecontrol.StateCancelled, Address: addr, AdminURL: adminURL})

	waitForGoroutines(&wg)
	select {
	case tuiErr := <-tuiErrors:
		return fmt.Errorf("serve terminal UI failed: %w; rerun with --no-tui", tuiErr)
	default:
	}

	errln("Server stopped")
	return nil
}

func configureServeBuildModes(m *lifecycle.Manager) {
	// Apply fast mode if requested
	if serveFast {
		applyFastMode(m)
		m.Config().Extra["cache_cleanup_async"] = true
	}
	if serveIncremental {
		lifecycle.SetServeIncremental(m, true)
		if m.Config().Extra == nil {
			m.Config().Extra = make(map[string]any)
		}
		m.Config().Extra["incremental_mode"] = true
	}

	if !serveFast && !serveIncremental {
		lifecycle.SetServeFullRebuild(m, true)
		lifecycle.SetServeChangedPaths(m, nil)
		lifecycle.SetServeRemovedPaths(m, nil)
		lifecycle.SetServeGlobDirty(m, true)
	}
}

func configureServeActionHandler(ctx context.Context, runtime *servecontrol.Runtime, rebuildCh chan struct{}, wg *sync.WaitGroup) {
	runtime.SetActionHandler(func(action servecontrol.ActionRequest) error {
		switch action.Kind {
		case serveBuildAction, "rerun", "rerun_job", rebuildPageLegacyAction, rebuildPageAction:
			if action.Kind == "rerun" || action.Kind == "rerun_job" {
				found := false
				snapshot := runtime.Snapshot()
				for i := range snapshot.Jobs {
					job := snapshot.Jobs[i]
					if job.ID == action.JobID {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("job %q was not found", action.JobID)
				}
			}
			if action.Kind == rebuildPageLegacyAction || action.Kind == rebuildPageAction {
				if !strings.HasSuffix(action.Page, ".md") && !strings.HasSuffix(action.Page, ".markdown") {
					return fmt.Errorf("page %q has no editable Markdown source", action.Page)
				}
				found := false
				for _, page := range runtime.Snapshot().Pages {
					if page.Path == action.Page {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("page %q was not found", action.Page)
				}
			}
			serveChangedPathsMu.Lock()
			serveForceFullRebuild = action.Kind != rebuildPageLegacyAction && action.Kind != rebuildPageAction
			if action.Kind == rebuildPageLegacyAction || action.Kind == rebuildPageAction {
				serveChangedPaths[action.Page] = fsnotify.Write
			}
			serveChangedPathsMu.Unlock()
			if isRebuilding.Load() {
				rebuildPending.Store(true)
				return nil
			}
			if rebuildCh == nil {
				wg.Add(1)
				go func() {
					defer wg.Done()
					doRebuild(ctx, make(chan struct{}, 1))
				}()
				return nil
			}
			if serveRequestFullRebuild != nil {
				serveRequestFullRebuild()
			}
			return nil
		default:
			return fmt.Errorf("unsupported serve action %q", action.Kind)
		}
	})
}

func validateServePort(port int) error {
	if port < 1 || port > 65535 {
		return newUsageError(&Diagnostic{Severity: "error", Code: "serve.invalid_port", Problem: fmt.Sprintf("port %d is outside the valid range", port), Details: []string{"Valid ports: 1 to 65535"}, StatusCode: exitCodeUsage})
	}
	return nil
}

func resolveServeHost(cmd *cobra.Command) error {
	if !cmd.Flags().Changed("bind") {
		return nil
	}
	if cmd.Flags().Changed("host") && serveHost != serveBind {
		return newUsageError(fmt.Errorf("--host and --bind specify different addresses; use one flag or give both the same value"))
	}
	serveHost = serveBind
	return nil
}

func setupServeSignals(cancel context.CancelFunc) {
	// Handle interrupt signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		errln("\nInterrupt received - shutting down...")
		if isRebuilding.Load() {
			errln("Rebuild in progress - canceling...")
		}
		cancel()
	}()
}

func resolveServeOutputPath(m *lifecycle.Manager) (outputPath, absOutputPath string) {
	outputPath = m.Config().OutputDir
	if outputPath == "" {
		outputPath = "output"
	}

	absOutputPath, err := filepath.Abs(outputPath)
	if err != nil {
		absOutputPath = outputPath
	}

	return outputPath, absOutputPath
}

func initServeInternalPaths() {
	paths := []string{".markata", ".markata-cache", "cache", "markout", "public", "output"}
	serveInternalPaths = serveInternalPaths[:0]
	for _, p := range paths {
		absPath, err := filepath.Abs(p)
		if err != nil {
			absPath = p
		}
		serveInternalPaths = append(serveInternalPaths, absPath)
	}
}

func setupWatcher(ctx context.Context, m *lifecycle.Manager, wg *sync.WaitGroup) (rebuildCh chan struct{}, closeWatcher func(), err error) {
	// Watch is enabled if: --watch is true (default) AND --no-watch is false
	// --no-watch takes precedence for backward compatibility
	shouldWatch := serveWatch && !serveNoWatch
	if !shouldWatch {
		return nil, func() {}, nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, func() {}, fmt.Errorf("failed to create file watcher: %w", err)
	}

	rebuildCh = make(chan struct{}, 1)

	// Start watcher goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		watchFiles(ctx, watcher, rebuildCh)
	}()

	// Start rebuild goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		handleRebuilds(ctx, rebuildCh)
	}()

	// Add paths to watch
	if err := addWatchPaths(watcher, m); err != nil {
		watcher.Close()
		return nil, func() {}, fmt.Errorf("failed to setup file watching: %w", err)
	}

	return rebuildCh, func() { _ = watcher.Close() }, nil
}

func startHTTPServer(listener net.Listener, handler http.Handler) (server *http.Server, serverErr <-chan error) {
	server = &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: serverReadHeaderTimeout,
	}

	// Start server in goroutine
	serverErrCh := make(chan error, 1)
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			serverErrCh <- err
		}
	}()

	return server, serverErrCh
}

func diagnoseServeListenError(err error, addr string) error {
	if isAddressInUseError(err) {
		port := servePort
		if _, portText, splitErr := net.SplitHostPort(addr); splitErr == nil {
			if parsedPort, parseErr := strconv.Atoi(portText); parseErr == nil {
				port = parsedPort
			}
		}
		available := 0
		for candidate := port + 1; candidate <= min(port+20, 65535); candidate++ {
			listener, listenErr := net.Listen("tcp", net.JoinHostPort(serveHost, strconv.Itoa(candidate)))
			if listenErr == nil {
				available = candidate
				_ = listener.Close()
				break
			}
		}
		diagnostic := &Diagnostic{Severity: "error", Code: "serve.address_in_use", Problem: fmt.Sprintf("port %d is already in use", port), Cause: err, Tips: []string{fmt.Sprintf("Another Markata server may already be serving at http://%s", addr)}, StatusCode: 1}
		if available > 0 {
			diagnostic.Command = fmt.Sprintf("markata-go serve --port %d", available)
		}
		return diagnostic
	}
	return &Diagnostic{Severity: "error", Code: "serve.listen_failed", Problem: fmt.Sprintf("could not listen on %s", addr), Cause: err, Tips: []string{"Check that the address is valid and available."}, StatusCode: 1}
}

func isAddressInUseError(err error) bool {
	if errors.Is(err, syscall.EADDRINUSE) {
		return true
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "address already in use") ||
		strings.Contains(message, "only one usage of each socket address")
}

func shouldRunServeTUI(disabled, noInput bool, input, output any) bool {
	if disabled || noInput || os.Getenv("TERM") == "dumb" {
		return false
	}
	in, inOK := input.(*os.File)
	out, outOK := output.(*os.File)
	return inOK && outOK && term.IsTerminal(int(in.Fd())) && term.IsTerminal(int(out.Fd()))
}

func serveLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func mountServeAdmin(site http.Handler, runtime *servecontrol.Runtime, sourceRoot string, enabled bool) http.Handler {
	if !enabled {
		return site
	}
	mux := http.NewServeMux()
	mux.HandleFunc(serveSubscriptionsPage, handleSubscriptionsPage)
	mux.HandleFunc(serveSubscriptionsAPI, handleSubscriptionsAPI)
	mux.HandleFunc(serveSubscriptionsFind, handleSubscriptionsDiscover)
	mux.Handle("/_markata/", servecontrol.NewWebHandlerWithFixHooks(runtime, sourceRoot, func(previews []servefix.FilePreview) {
		registerExpectedServeFixWrites(sourceRoot, previews)
	}))
	mux.Handle("/", site)
	return mux
}

func resolveServeAdmin(host string, force, disabled bool) (bool, error) {
	if force && disabled {
		return false, newUsageError(fmt.Errorf("--admin and --no-admin cannot be used together"))
	}
	if force && !serveLoopbackHost(host) {
		return false, newUsageError(fmt.Errorf("local admin requires a loopback address; use --bind localhost or --bind 127.0.0.1"))
	}
	return !disabled && serveLoopbackHost(host), nil
}

func startInitialBuild(m *lifecycle.Manager, rebuildCh chan struct{}, wg *sync.WaitGroup) {
	setBuildStatus(buildStatusBuilding, "", "")
	notifyBuildStatus()
	verbosef("Running initial build...")

	if serveFast || serveIncremental {
		lifecycle.SetServeFullRebuild(m, false)
	} else {
		lifecycle.SetServeFullRebuild(m, true)
	}
	lifecycle.SetServeGlobDirty(m, true)
	lifecycle.SetServeChangedPaths(m, nil)
	lifecycle.SetServeRemovedPaths(m, nil)

	wg.Add(1)
	go func() {
		defer wg.Done()
		isRebuilding.Store(true)
		defer func() {
			isRebuilding.Store(false)
			if rebuildPending.Swap(false) && rebuildCh != nil {
				select {
				case rebuildCh <- struct{}{}:
				default:
				}
			}
		}()

		jobID := startServeJob("Initial build", "startup", nil)
		result, buildErr := runServeBuild(m, jobID)
		finishServeJob(jobID, m, result, buildErr)
		if buildErr != nil {
			setBuildStatus(buildStatusError, buildErr.Error(), "")
			notifyBuildStatus()
			errlnf("Initial build failed: %v", buildErr)
			return
		}

		setBuildStatus(buildStatusSuccess, "", licenseWarningMessage(getModelsConfig(m)))
		notifyBuildStatus()
		infof("Built %d posts, %d feeds", result.PostsProcessed, result.FeedsGenerated)
		notifyLiveReload()
		posts := m.Posts()
		setServePosts(posts)
		setServeSearchPosts(posts)
		if serveFast || serveIncremental {
			if cached, ok := m.Cache().Get("build_cache"); ok {
				if bc, ok := cached.(*buildcache.Cache); ok {
					setServeCache(bc)
				}
			}
		}
	}()
}

func waitForShutdown(ctx context.Context, server *http.Server, serverErr <-chan error) error {
	select {
	case <-ctx.Done():
		// Graceful shutdown
		infof("Initiating graceful shutdown...")

		// Close live reload connections first
		closeAllLiveReloadConnections()

		// Shorter timeout for faster shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer shutdownCancel()

		verbosef("Shutting down HTTP server (timeout: 2s)...")
		if err := server.Shutdown(shutdownCtx); err != nil {
			errlnf("Server shutdown error: %v", err)
		} else {
			verbosef("HTTP server shutdown completed")
		}
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

func waitForGoroutines(wg *sync.WaitGroup) {
	// Wait for goroutines to finish with timeout
	activeConnections := liveReloadCount.Load()
	if verbose || activeConnections > 0 {
		verbosef("Waiting for goroutines to finish (active SSE connections: %d)...", activeConnections)
	} else {
		verbosef("Waiting for goroutines to finish...")
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		verbosef("All tracked goroutines finished")
	case <-time.After(1 * time.Second):
		errlnf("Shutdown timeout after 1s - forcing exit")
		if activeConnections > 0 {
			errlnf("Note: had %d active SSE connections that may not have closed cleanly", activeConnections)
		}
	}
}

// createHandler creates an HTTP handler that serves files with live reload injection
// and mounts the bleve search API at the configured endpoint.
//
//nolint:gocyclo // The handler dispatches independent dev endpoints and static-file fallbacks.
func createHandler(outputDir string, m *lifecycle.Manager, searchEndpoint string) http.Handler {
	fileServer := http.FileServer(http.Dir(outputDir))
	absOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		absOutputDir = outputDir
	}

	// Set up search API handler
	apiCfg := searchapi.DefaultConfig()
	apiCfg.IndexName = "serve"
	if mc := getModelsConfig(m); mc != nil {
		apiCfg.DefaultLimit = mc.Search.Bleve.DefaultLimit()
		apiCfg.MaxLimit = mc.Search.Bleve.GetMaxLimit()
		apiCfg.DefaultFuzzy = mc.Search.Bleve.IsFuzzy()
		if len(mc.Search.Bleve.CORSOrigins) > 0 {
			apiCfg.CORSOrigins = mc.Search.Bleve.CORSOrigins
		}
	}
	cacheDir := filepath.Join(filepath.Dir(absOutputDir), ".markata", "cache")
	searchHandler := searchapi.NewHandler(getServeSearchPostsSnapshot(), cacheDir, apiCfg)

	// Subscribe to post updates after rebuilds.
	// UpdatePosts is now change-aware: it only rebuilds the index
	// when the post set actually changes (cheap path-based fingerprint).
	go func() {
		for {
			time.Sleep(5 * time.Second)
			posts := getServeSearchPostsSnapshot()
			searchHandler.UpdatePosts(posts)
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status := getBuildStatus()
		// Log requests in verbose mode
		if verbose {
			verbosef("[%s] %s %s", time.Now().Format("15:04:05"), r.Method, r.URL.Path)
		}

		// Handle live reload endpoint
		if r.URL.Path == "/__livereload" {
			handleLiveReload(w, r)
			return
		}

		// Theme picker "Bake" writes picks into the site config (serve only).
		if r.URL.Path == serveThemeBakeEndpoint {
			handleThemeBake(w, r)
			return
		}

		// Settings sidebar: schema, bake endpoint, and its script (serve only).
		if r.URL.Path == serveSettingsEndpoint {
			handleSettings(w, r)
			return
		}
		if r.URL.Path == serveSettingsPreview {
			handleSettingsPreview(w, r)
			return
		}
		if r.URL.Path == serveSettingsScriptPath {
			handleSettingsScript(w, r)
			return
		}

		// Handle search API endpoint (bleve-powered, read-only)
		if r.URL.Path == searchEndpoint {
			searchHandler.ServeHTTP(w, r)
			return
		}

		// Handle search endpoint (no-JS fallback for 404 page search)
		if r.URL.Path == "/_search" && r.Method == http.MethodPost {
			handleSearchFallback(w, r, outputDir)
			return
		}

		// Determine the file path
		fullPath, requestPath, resolveErr := resolveRequestPath(absOutputDir, r.URL.Path)
		if resolveErr != nil {
			serve404Page(w, outputDir, status)
			return
		}

		// Check if file exists
		info, err := os.Stat(fullPath)
		if err == nil && info.IsDir() {
			// Try index.html in directory
			indexPath := filepath.Join(fullPath, serveIndexFile)
			if _, err := os.Stat(indexPath); err == nil {
				requestPath = path.Join(requestPath, serveIndexFile)
				fullPath = indexPath
			}
		}

		// Check if file exists - if not, serve 404 page
		if err != nil && os.IsNotExist(err) {
			serve404Page(w, outputDir, status)
			return
		}

		// Check if it's an HTML file and inject live reload script
		if strings.HasSuffix(requestPath, ".html") || (info != nil && !info.IsDir() && strings.HasSuffix(fullPath, ".html")) {
			serveHTMLWithLiveReload(w, fullPath, outputDir, status)
			return
		}

		// Serve with file server
		r.URL.Path = requestPath
		fileServer.ServeHTTP(w, r)
	})
}

// serveHTMLWithLiveReload reads an HTML file and injects the live reload script.
func serveHTMLWithLiveReload(w http.ResponseWriter, path, outputDir string, status BuildStatus) {
	content, err := os.ReadFile(path)
	if err != nil {
		// File not found - serve 404 page
		serve404Page(w, outputDir, status)
		return
	}

	html := injectDevScripts(string(content), status)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write([]byte(html)); err != nil && verbose {
		verbosef("Error writing response: %v", err)
	}
}

// serve404Page serves the static 404.html page with live reload injection.
// The 404 page uses client-side JavaScript for fuzzy search suggestions,
// so it works the same in dev server as in production.
func serve404Page(w http.ResponseWriter, outputDir string, status BuildStatus) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)

	// Try to serve the static 404.html
	notFoundPath := filepath.Join(outputDir, "404.html")
	content, err := os.ReadFile(notFoundPath)
	if err != nil {
		// Fallback to simple error message if 404.html doesn't exist
		fallback := fallback404HTML()
		fallback = injectDevScripts(fallback, status)
		//nolint:errcheck // Best effort write to HTTP response
		w.Write([]byte(fallback))
		return
	}

	html := injectDevScripts(string(content), status)

	//nolint:errcheck // Best effort write to HTTP response
	w.Write([]byte(html))
}

func fallback404HTML() string {
	return `<!DOCTYPE html>
<html>
<head><title>404 - Page Not Found</title></head>
<body>
<h1>404 - Page Not Found</h1>
<p>The requested page could not be found.</p>
<p><a href="/">Go to home page</a></p>
</body>
</html>`
}

func injectDevScripts(html string, status BuildStatus) string {
	devScript := buildDevScript(status)

	// Inject search endpoint override for dev mode.
	// This tells the frontend to use the local bleve API instead of pagefind.
	// It's only injected during `serve` — production builds never see this.
	// The theme bake endpoint is likewise serve-only; the picker shows its
	// Bake button only when this global exists. The settings sidebar script
	// is served by the dev server itself, so it never reaches output/.
	searchScript := `<script>window.__markataSearchEndpoint = "` + serveSearchEndpoint + `";` +
		`window.__markataThemeBakeEndpoint = "` + serveThemeBakeEndpoint + `";` +
		`window.__markataSettingsEndpoint = "` + serveSettingsEndpoint + `";</script>`
	bodyInjection := devScript + `<script src="` + serveSettingsScriptPath + `" defer></script>`
	if serveLoopbackHost(serveHost) && !serveNoAdmin {
		bodyInjection = `<a href="/_markata/subscriptions" aria-label="Manage reader subscriptions" style="position:fixed;bottom:14px;right:14px;z-index:2147482900;display:inline-block;padding:8px 12px;background:Canvas;color:CanvasText;border:1px solid currentColor;border-radius:20px;font:600 13px/1.3 system-ui,sans-serif;text-decoration:none;box-shadow:0 3px 12px #0003">RSS · Subscriptions</a>` + bodyInjection
	}

	headInjection := searchScript

	if strings.Contains(html, "</head>") {
		html = strings.Replace(html, "</head>", headInjection+"</head>", 1)
	} else {
		bodyInjection = headInjection + bodyInjection
	}

	switch {
	case strings.Contains(html, "</body>"):
		return strings.Replace(html, "</body>", bodyInjection+"</body>", 1)
	case strings.Contains(html, "</html>"):
		return strings.Replace(html, "</html>", bodyInjection+"</html>", 1)
	default:
		return html + bodyInjection
	}
}

func buildDevScript(status BuildStatus) string {
	payload := buildStatusPayload(status)

	return fmt.Sprintf(`<script>
(function() {
    var initialStatus = %s;
    var bannerId = 'markata-build-banner';
    var toastId = 'markata-license-toast';
    var styleId = 'markata-build-banner-style';

    function ensureStyle() {
        if (document.getElementById(styleId)) {
            return;
        }

        var style = document.createElement('style');
        style.id = styleId;
        style.textContent = '#'+bannerId+'{position:fixed;top:0;left:0;right:0;z-index:9999;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Ubuntu,"Helvetica Neue",Arial,sans-serif;font-size:14px;padding:8px 12px;text-align:center;box-shadow:0 4px 12px rgba(0,0,0,0.12);display:none;}'+
            '#'+bannerId+'[data-status="building"]{background:#f59e0b;color:#111827;}'+
            '#'+bannerId+'[data-status="error"]{background:#dc2626;color:#ffffff;}'+
            '#'+toastId+'{position:fixed;right:16px;bottom:16px;z-index:10000;max-width:420px;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Ubuntu,"Helvetica Neue",Arial,sans-serif;font-size:14px;line-height:1.35;background:#111827;color:#f9fafb;padding:10px 14px;border-radius:8px;box-shadow:0 10px 24px rgba(0,0,0,0.25);display:none;}';
        document.head.appendChild(style);
    }

    function ensureBanner() {
        var banner = document.getElementById(bannerId);
        if (banner) {
            return banner;
        }

        banner = document.createElement('div');
        banner.id = bannerId;
        document.body.appendChild(banner);
        return banner;
    }

    function ensureToast() {
        var toast = document.getElementById(toastId);
        if (toast) {
            return toast;
        }

        toast = document.createElement('div');
        toast.id = toastId;
        toast.setAttribute('role', 'status');
        toast.setAttribute('aria-live', 'polite');
        document.body.appendChild(toast);
        return toast;
    }

    function hideLicenseToast() {
        var toast = document.getElementById(toastId);
        if (!toast) {
            return;
        }
        toast.style.display = 'none';
    }

    function showLicenseToast(message) {
        if (!message) {
            hideLicenseToast();
            return;
        }
        var toast = ensureToast();
        toast.textContent = message;
        toast.style.display = 'block';
    }

    function applyStatus(state) {
        if (!state || !state.status) {
            return;
        }

        // Dev tools such as the settings sidebar follow rebuild progress.
        window.__markataBuildStatus = state;
        try {
            window.dispatchEvent(new CustomEvent('markata:build-status', { detail: state }));
        } catch (e) {}

        ensureStyle();
        var banner = ensureBanner();
        banner.setAttribute('data-status', state.status);

        if (state.status === '%s') {
            banner.style.display = 'none';
            showLicenseToast(state.license_warning || '');
            return;
        }

        hideLicenseToast();

        var text = 'Building...';
        if (state.status === '%s') {
            text = 'Build failed';
            if (state.message) {
                text += ': ' + state.message;
            }
        }

        banner.textContent = text;
        banner.style.display = 'block';
    }

    function parseStatusPayload(raw) {
        try {
            return JSON.parse(raw);
        } catch (e) {
            return null;
        }
    }

    applyStatus(initialStatus);

    var source = new EventSource('/__livereload');
    source.onmessage = function(e) {
        if (e.data === 'reload') {
            location.reload();
            return;
        }
        if (e.data.indexOf('%s') === 0) {
            var payload = e.data.slice(%d);
            var state = parseStatusPayload(payload);
            if (state) {
                applyStatus(state);
            }
        }
    };
    source.onerror = function() {
        source.close();
        setTimeout(function() {
            location.reload();
        }, 1000);
    };
})();
</script>`, payload, buildStatusSuccess, buildStatusError, buildStatusEventPrefix, len(buildStatusEventPrefix))
}

// handleSearchFallback handles POST requests to /_search for no-JS fallback search.
// It reads the posts index, performs fuzzy search, and renders results as HTML.
func handleSearchFallback(w http.ResponseWriter, r *http.Request, outputDir string) {
	// Parse the form to get the search query
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	query := strings.TrimSpace(r.FormValue("q"))
	if query == "" {
		// No query, redirect to home
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Load posts index
	indexPath := filepath.Join(outputDir, "_404-index.json")
	indexData, err := os.ReadFile(indexPath)
	if err != nil {
		// Index not available, redirect to home with query param
		http.Redirect(w, r, "/?q="+url.QueryEscape(query), http.StatusSeeOther)
		return
	}

	// Parse the index
	var posts []struct {
		Slug        string `json:"slug"`
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
	}
	if err := json.Unmarshal(indexData, &posts); err != nil {
		http.Redirect(w, r, "/?q="+url.QueryEscape(query), http.StatusSeeOther)
		return
	}

	// Perform simple search (case-insensitive substring matching + basic scoring)
	queryLower := strings.ToLower(query)
	queryWords := strings.Fields(queryLower)
	var results []searchResult

	for _, post := range posts {
		titleLower := strings.ToLower(post.Title)
		slugLower := strings.ToLower(post.Slug)
		descLower := strings.ToLower(post.Description)

		score := 0
		for _, word := range queryWords {
			if strings.Contains(titleLower, word) {
				score += 10
			}
			if strings.Contains(slugLower, word) {
				score += 8
			}
			if strings.Contains(descLower, word) {
				score += 5
			}
		}

		if score > 0 {
			results = append(results, searchResult{
				Title:       post.Title,
				Description: post.Description,
				URL:         post.URL,
				Score:       score,
			})
		}
	}

	// Sort by score (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Limit results
	if len(results) > 20 {
		results = results[:20]
	}

	// Render search results page
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	renderSearchResultsPage(w, query, results)
}

// renderSearchResultsPage renders HTML for search results (no-JS fallback).
func renderSearchResultsPage(w http.ResponseWriter, query string, results []searchResult) {
	var html strings.Builder
	html.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Search Results</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 2rem 1rem;
            background: #f9fafb;
            color: #1f2937;
        }
        @media (prefers-color-scheme: dark) {
            body { background: #111827; color: #f3f4f6; }
            .search-input { background: #1f2937; border-color: #374151; color: #f3f4f6; }
            .result-item { background: #1f2937; border-color: #374151; }
            .result-item:hover { background: #374151; }
        }
        h1 { font-size: 1.5rem; margin-bottom: 1.5rem; }
        .search-form { display: flex; gap: 0.5rem; margin-bottom: 2rem; }
        .search-input {
            flex: 1;
            padding: 0.75rem 1rem;
            font-size: 1rem;
            border: 2px solid #e5e7eb;
            border-radius: 0.5rem;
        }
        .search-button {
            padding: 0.75rem 1.5rem;
            font-size: 1rem;
            font-weight: 600;
            color: white;
            background: #3b82f6;
            border: none;
            border-radius: 0.5rem;
            cursor: pointer;
        }
        .search-button:hover { background: #2563eb; }
        .results-count { color: #6b7280; margin-bottom: 1rem; }
        .result-item {
            display: block;
            padding: 1rem;
            margin-bottom: 0.5rem;
            background: white;
            border: 1px solid #e5e7eb;
            border-radius: 0.5rem;
            text-decoration: none;
            color: inherit;
            transition: background 0.2s, transform 0.2s;
        }
        .result-item:hover { background: #f3f4f6; transform: translateX(4px); }
        .result-title { display: block; font-weight: 600; color: #3b82f6; margin-bottom: 0.25rem; }
        .result-desc { display: block; font-size: 0.875rem; color: #6b7280; }
        .no-results { color: #6b7280; font-style: italic; }
        .back-link { display: inline-block; margin-top: 1.5rem; color: #3b82f6; }
    </style>
</head>
<body>
    <h1>Search Results</h1>
    <form class="search-form" action="/_search" method="POST">
        <input type="text" name="q" class="search-input" value="`)
	html.WriteString(template.HTMLEscapeString(query))
	html.WriteString(`" placeholder="Search..." autocomplete="off">
        <button type="submit" class="search-button">Search</button>
    </form>
`)

	if len(results) == 0 {
		html.WriteString(`    <p class="no-results">No results found for "`)
		html.WriteString(template.HTMLEscapeString(query))
		html.WriteString(`"</p>
`)
	} else {
		html.WriteString(fmt.Sprintf(`    <p class="results-count">Found %d result(s) for "%s"</p>
`, len(results), template.HTMLEscapeString(query)))

		for _, result := range results {
			html.WriteString(`    <a href="`)
			html.WriteString(template.HTMLEscapeString(result.URL))
			html.WriteString(`" class="result-item">
        <span class="result-title">`)
			html.WriteString(template.HTMLEscapeString(result.Title))
			html.WriteString(`</span>
`)
			if result.Description != "" {
				desc := result.Description
				if len(desc) > 150 {
					desc = desc[:150] + "..."
				}
				html.WriteString(`        <span class="result-desc">`)
				html.WriteString(template.HTMLEscapeString(desc))
				html.WriteString(`</span>
`)
			}
			html.WriteString(`    </a>
`)
		}
	}

	html.WriteString(`    <a href="/" class="back-link">&larr; Back to home</a>
</body>
</html>`)

	//nolint:errcheck // Best effort write to HTTP response
	w.Write([]byte(html.String()))
}

// Live reload clients
var (
	liveReloadClients    = make(map[chan string]struct{})
	liveReloadClientsMu  sync.RWMutex
	liveReloadCount      atomic.Int32
	liveReloadShutdown   = make(chan struct{})
	liveReloadShutdownMu sync.Mutex
)

// handleLiveReload handles Server-Sent Events for live reload.
func handleLiveReload(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create channel for this client
	ch := make(chan string, 1)

	// Register client
	liveReloadClientsMu.Lock()
	liveReloadClients[ch] = struct{}{}
	liveReloadCount.Add(1)
	liveReloadClientsMu.Unlock()

	if verbose {
		verbosef("Live reload client connected (total: %d)", liveReloadCount.Load())
	}

	// Ensure client is removed on disconnect
	defer func() {
		liveReloadClientsMu.Lock()
		delete(liveReloadClients, ch)
		liveReloadCount.Add(-1)
		// Close channel if not already closed
		select {
		case <-ch:
			// Channel already closed
		default:
			close(ch)
		}
		liveReloadClientsMu.Unlock()

		if verbose {
			verbosef("Live reload client disconnected (total: %d)", liveReloadCount.Load())
		}
	}()

	// Get flusher for streaming
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Send initial connection message
	fmt.Fprintf(w, "data: connected\n\n")
	status := getBuildStatus()
	fmt.Fprintf(w, "data: %s%s\n\n", buildStatusEventPrefix, buildStatusPayload(status))
	flusher.Flush()

	// Wait for messages or disconnect
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				// Channel closed by global shutdown
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		case <-liveReloadShutdown:
			if verbose {
				verbosef("SSE handler received global shutdown signal")
			}
			return
		}
	}
}

// notifyLiveReload sends a reload message to all connected clients.
func notifyLiveReload() {
	notifyLiveReloadMessage("reload")
}

func notifyBuildStatus() {
	status := getBuildStatus()
	payload := buildStatusPayload(status)
	notifyLiveReloadMessage(buildStatusEventPrefix + payload)
}

func notifyLiveReloadMessage(message string) {
	liveReloadClientsMu.RLock()
	defer liveReloadClientsMu.RUnlock()

	for ch := range liveReloadClients {
		select {
		case ch <- message:
		case <-ch:
			// Channel closed, skip
		default:
			// Skip if channel is full
		}
	}
}

// closeAllLiveReloadConnections closes all SSE connections for shutdown.
func closeAllLiveReloadConnections() {
	liveReloadShutdownMu.Lock()
	defer liveReloadShutdownMu.Unlock()

	liveReloadClientsMu.Lock()
	defer liveReloadClientsMu.Unlock()

	count := len(liveReloadClients)
	if count > 0 {
		verbosef("Closing %d live reload connection(s)...", count)

		// Signal global shutdown - this will cause all SSE handlers to exit
		select {
		case <-liveReloadShutdown:
			// Already closed
		default:
			close(liveReloadShutdown)
		}

		// Close all client channels to force immediate exit
		for ch := range liveReloadClients {
			select {
			case <-ch:
				// Channel already closed
			default:
				close(ch)
			}
		}

		// Clear the map
		clear(liveReloadClients)
		liveReloadCount.Store(0)
		verbosef("Closed all live reload connections")
	}
}

// shouldIgnoreEvent returns true if the event should be ignored.
func shouldIgnoreEvent(event fsnotify.Event) bool {
	// Get absolute path for comparison
	absEventPath, err := filepath.Abs(event.Name)
	if err != nil {
		absEventPath = event.Name
	}

	// Ignore events for output directory
	if isPathWithinDir(absEventPath, serveOutputPath) {
		return true
	}

	for _, internalPath := range serveInternalPaths {
		if isPathWithinDir(absEventPath, internalPath) {
			return true
		}
	}

	// Ignore temporary/backup files and hidden files
	baseName := filepath.Base(event.Name)
	return strings.HasSuffix(event.Name, "~") ||
		strings.HasPrefix(baseName, ".") ||
		strings.HasSuffix(event.Name, ".swp") ||
		strings.HasSuffix(event.Name, ".swo") ||
		strings.HasSuffix(event.Name, ".tmp")
}

func registerExpectedServeFixWrites(root string, previews []servefix.FilePreview) {
	root, err := filepath.Abs(root)
	if err != nil {
		return
	}
	now := time.Now()
	serveFixExpectedWrites.Lock()
	for path, expected := range serveFixExpectedWrites.files {
		if now.After(expected.expires) {
			delete(serveFixExpectedWrites.files, path)
		}
	}
	for i := range previews {
		preview := previews[i]
		if len(serveFixExpectedWrites.files) >= 1024 {
			break
		}
		path, err := filepath.Abs(filepath.Join(root, filepath.FromSlash(preview.Path)))
		if err != nil {
			continue
		}
		serveFixExpectedWrites.files[filepath.Clean(path)] = expectedServeFixWrite{
			before:  contentDigest([]byte(preview.Before)),
			after:   contentDigest([]byte(preview.After)),
			expires: now.Add(5 * time.Second),
		}
	}
	serveFixExpectedWrites.Unlock()
}

func isExpectedServeFixEvent(eventPath string) bool {
	path, err := filepath.Abs(eventPath)
	if err != nil {
		return false
	}
	path = filepath.Clean(path)
	serveFixExpectedWrites.Lock()
	expected, ok := serveFixExpectedWrites.files[path]
	if !ok {
		serveFixExpectedWrites.Unlock()
		return false
	}
	if time.Now().After(expected.expires) {
		delete(serveFixExpectedWrites.files, path)
		serveFixExpectedWrites.Unlock()
		return false
	}
	serveFixExpectedWrites.Unlock()
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		// An absent path can be a user deletion. Let the watcher process it;
		// replacement remove/create bursts still coalesce in handleRebuilds.
		clearExpectedServeFixWrite(path, expected)
		return false
	}
	digest := contentDigest(content)
	if digest == expected.before || digest == expected.after {
		return true
	}
	serveFixExpectedWrites.Lock()
	if current, exists := serveFixExpectedWrites.files[path]; exists && current == expected {
		delete(serveFixExpectedWrites.files, path)
	}
	serveFixExpectedWrites.Unlock()
	return false
}

func clearExpectedServeFixWrite(path string, expected expectedServeFixWrite) {
	serveFixExpectedWrites.Lock()
	if current, exists := serveFixExpectedWrites.files[path]; exists && current == expected {
		delete(serveFixExpectedWrites.files, path)
	}
	serveFixExpectedWrites.Unlock()
}

func clearExpectedServeFixWrites() {
	serveFixExpectedWrites.Lock()
	clear(serveFixExpectedWrites.files)
	serveFixExpectedWrites.Unlock()
}

func contentDigest(content []byte) string {
	digest := sha256.Sum256(content)
	return hex.EncodeToString(digest[:])
}

var globMagicPattern = regexp.MustCompile(`[*?\[{]`)

func contentWatchRoots(config *lifecycle.Config) []string {
	contentDir := config.ContentDir
	if contentDir == "" {
		contentDir = "."
	}

	if contentDir != "." || len(config.GlobPatterns) == 0 {
		return []string{contentDir}
	}

	seen := map[string]struct{}{}
	roots := make([]string, 0, len(config.GlobPatterns))
	for _, pattern := range config.GlobPatterns {
		root := watchRootFromPattern(pattern)
		if root == "" {
			continue
		}
		if _, ok := seen[root]; ok {
			continue
		}
		seen[root] = struct{}{}
		roots = append(roots, root)
	}

	if len(roots) == 0 {
		return []string{contentDir}
	}

	sort.Strings(roots)
	return roots
}

func watchRootFromPattern(pattern string) string {
	cleaned := filepath.Clean(filepath.FromSlash(pattern))
	if cleaned == "." || cleaned == string(filepath.Separator) {
		return "."
	}

	parts := strings.Split(cleaned, string(filepath.Separator))
	rootParts := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" || part == "." {
			continue
		}
		if globMagicPattern.MatchString(part) {
			break
		}
		rootParts = append(rootParts, part)
	}

	if len(rootParts) == 0 {
		return "."
	}

	return filepath.Join(rootParts...)
}

// handleNewDirectory adds newly created directories to the watcher.
// This ensures new post folders, tag directories, etc. are watched.
func handleNewDirectory(watcher *fsnotify.Watcher, event fsnotify.Event) {
	if event.Op&fsnotify.Create == 0 {
		return
	}

	info, statErr := os.Stat(event.Name)
	if statErr != nil || !info.IsDir() {
		return
	}

	// Add the new directory and any subdirectories to the watcher
	if watchErr := addDirRecursive(watcher, event.Name); watchErr != nil {
		if verbose {
			verbosef("Failed to watch new directory %s: %v", event.Name, watchErr)
		}
	} else if verbose {
		verbosef("Now watching new directory: %s", event.Name)
	}
}

// watchFiles handles file system events.
func watchFiles(ctx context.Context, watcher *fsnotify.Watcher, rebuildCh chan<- struct{}) {
	watchFilesWithObserver(ctx, watcher, rebuildCh, nil)
}

func watchFilesWithObserver(ctx context.Context, watcher *fsnotify.Watcher, rebuildCh chan<- struct{}, observe func(fsnotify.Event, bool)) {
	for {
		select {
		case <-ctx.Done():
			if verbose {
				verbosef("File watcher canceled")
			}
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if shouldIgnoreEvent(event) {
				continue
			}

			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
				suppressed := isExpectedServeFixEvent(event.Name)
				if observe != nil {
					observe(event, suppressed)
				}
				if suppressed {
					continue
				}
			}

			// Handle newly created directories
			handleNewDirectory(watcher, event)

			// Trigger on write/create/remove/rename
			// Rename is important for editors that save via temp file + rename (vim, vscode, etc.)
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
				if verbose {
					verbosef("File changed: %s (%s)", event.Name, event.Op.String())
				}

				recordServeChangedPath(event)

				if isRebuilding.Load() {
					rebuildPending.Store(true)
					continue
				}

				// Debounce rebuilds
				select {
				case rebuildCh <- struct{}{}:
				default:
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			errlnf("Watcher error: %v", err)
		}
	}
}

func recordServeChangedPath(event fsnotify.Event) {
	if event.Name == "" {
		return
	}

	absPath, err := filepath.Abs(event.Name)
	if err != nil {
		absPath = event.Name
	}

	serveChangedPathsMu.Lock()
	serveChangedPaths[absPath] |= event.Op
	if event.Op&(fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
		serveGlobDirty = true
	}
	serveChangedPathsMu.Unlock()
}

func consumeServeChanges() (paths, removed []string, forceFull, globDirty bool) {
	serveChangedPathsMu.Lock()
	defer serveChangedPathsMu.Unlock()

	if len(serveChangedPaths) == 0 {
		forceFull = serveForceFullRebuild
		globDirty = serveGlobDirty
		serveForceFullRebuild = false
		serveGlobDirty = false
		return nil, nil, forceFull, globDirty
	}

	paths = make([]string, 0, len(serveChangedPaths))
	removed = make([]string, 0)
	for path := range serveChangedPaths {
		paths = append(paths, path)
	}

	forceFull = false
	for path, ops := range serveChangedPaths {
		if ops&fsnotify.Remove != 0 {
			removed = append(removed, path)
		}
		if ops&fsnotify.Rename != 0 {
			if _, err := os.Stat(path); err != nil {
				if os.IsNotExist(err) {
					removed = append(removed, path)
				}
			}
		}
	}
	if serveForceFullRebuild {
		forceFull = true
	}
	globDirty = serveGlobDirty

	serveChangedPaths = make(map[string]fsnotify.Op)
	serveForceFullRebuild = false
	serveGlobDirty = false
	return paths, removed, forceFull, globDirty
}

func getServeCache() *buildcache.Cache {
	serveCacheMu.Lock()
	defer serveCacheMu.Unlock()
	return serveCache
}

func setServeCache(cache *buildcache.Cache) {
	serveCacheMu.Lock()
	defer serveCacheMu.Unlock()
	serveCache = cache
}

func getServePosts() map[string]*models.Post {
	servePostsMu.Lock()
	defer servePostsMu.Unlock()
	return servePosts
}

func setServePosts(posts []*models.Post) {
	servePostsMu.Lock()
	defer servePostsMu.Unlock()
	if len(posts) == 0 {
		servePosts = nil
		return
	}
	newMap := make(map[string]*models.Post, len(posts))
	for _, post := range posts {
		if post != nil && post.Path != "" {
			newMap[post.Path] = post
		}
	}
	servePosts = newMap
}

// setServeSearchPosts publishes an immutable search snapshot. Build plugins
// mutate their manager-owned posts during a rebuild, so the search handler
// must never inspect those same pointers concurrently.
func setServeSearchPosts(posts []*models.Post) {
	newMap := make(map[string]*models.Post, len(posts))
	for _, post := range posts {
		if post != nil && post.Path != "" {
			newMap[post.Path] = cloneServeSearchPost(post)
		}
	}
	serveSearchPostsMu.Lock()
	defer serveSearchPostsMu.Unlock()
	if len(newMap) == 0 {
		serveSearchPosts = nil
		return
	}
	serveSearchPosts = newMap
}

func getServeSearchPostsSnapshot() []*models.Post {
	serveSearchPostsMu.RLock()
	defer serveSearchPostsMu.RUnlock()
	if len(serveSearchPosts) == 0 {
		return nil
	}
	posts := make([]*models.Post, 0, len(serveSearchPosts))
	for _, post := range serveSearchPosts {
		posts = append(posts, post)
	}
	return posts
}

func cloneServeSearchPost(post *models.Post) *models.Post {
	if post == nil {
		return nil
	}
	clone := *post
	if post.Title != nil {
		title := *post.Title
		clone.Title = &title
	}
	if post.Description != nil {
		description := *post.Description
		clone.Description = &description
	}
	if post.Date != nil {
		date := *post.Date
		clone.Date = &date
	}
	clone.Tags = append([]string(nil), post.Tags...)
	clone.Templates = make(map[string]string, len(post.Templates))
	for key, value := range post.Templates {
		clone.Templates[key] = value
	}
	clone.Extra = make(map[string]interface{}, len(post.Extra))
	for key, value := range post.Extra {
		clone.Extra[key] = value
	}
	return &clone
}

func normalizeServeChangedPaths(paths []string, contentDir string) (normalized []string, outside bool) {
	if len(paths) == 0 {
		return nil, false
	}

	absContentDir, err := filepath.Abs(contentDir)
	if err != nil {
		absContentDir = contentDir
	}

	seen := make(map[string]struct{}, len(paths))
	for _, p := range paths {
		if p == "" {
			continue
		}
		absPath, err := filepath.Abs(p)
		if err != nil {
			absPath = p
		}
		if !isPathWithinDir(absPath, absContentDir) {
			outside = true
			continue
		}
		rel, err := filepath.Rel(absContentDir, absPath)
		if err != nil {
			outside = true
			continue
		}
		rel = filepath.Clean(rel)
		if rel == "." {
			continue
		}
		if _, ok := seen[rel]; ok {
			continue
		}
		seen[rel] = struct{}{}
		normalized = append(normalized, rel)
	}

	return normalized, outside
}

// handleRebuilds processes rebuild requests with debouncing.
func handleRebuilds(ctx context.Context, rebuildCh chan struct{}) {
	handleRebuildsWith(ctx, rebuildCh, 300*time.Millisecond, func() { doRebuild(ctx, rebuildCh) })
}

func handleRebuildsWith(ctx context.Context, rebuildCh chan struct{}, debounceDelay time.Duration, rebuild func()) {
	var timer *time.Timer

	for {
		select {
		case <-ctx.Done():
			if timer != nil {
				if !timer.Stop() {
					drainTimer(timer)
				}
			}
			if verbose {
				verbosef("Rebuild handler canceled")
			}
			return
		case <-rebuildCh:
			// Reset debounce timer
			if timer != nil {
				if !timer.Stop() {
					drainTimer(timer)
				}
				timer.Reset(debounceDelay)
				continue
			}
			timer = time.NewTimer(debounceDelay)
		case <-timerChannel(timer):
			if timer != nil {
				timer.Stop()
				timer = nil
			}
			if ctx.Err() != nil {
				return
			}
			rebuild()
		}
	}
}

// doRebuild performs an incremental rebuild.
func doRebuild(ctx context.Context, rebuildCh chan<- struct{}) {
	// Set rebuilding flag to ignore events during build
	isRebuilding.Store(true)
	defer func() {
		isRebuilding.Store(false)
		if rebuildPending.Swap(false) {
			if ctx.Err() != nil {
				return
			}
			select {
			case rebuildCh <- struct{}{}:
			default:
			}
		}
	}()

	infof("\nRebuilding...")
	setBuildStatus(buildStatusBuilding, "", "")
	notifyBuildStatus()
	startTime := time.Now()
	jobID := startServeJob("Rebuild", "watch", nil)

	// Check if context is canceled before starting rebuild
	select {
	case <-ctx.Done():
		verbosef("Rebuild canceled")
		return
	default:
	}

	m, err := createServeManager()
	if err != nil {
		finishServeJob(jobID, nil, nil, err)
		setBuildStatus(buildStatusError, err.Error(), "")
		notifyBuildStatus()
		errlnf("Rebuild failed: %v", err)
		return
	}
	configureServeLoggerForManager(m)
	publishServePalette(m)
	changedPaths, removedPaths, forceFull, globDirty := consumeServeChanges()
	configureServeIncremental(m, changedPaths, removedPaths, forceFull, globDirty)
	configureRebuildManager(m)

	// Check for cancellation after creating manager
	select {
	case <-ctx.Done():
		verbosef("Rebuild canceled")
		return
	default:
	}

	result, err := runServeBuild(m, jobID)
	finishServeJob(jobID, m, result, err)
	if err != nil {
		setBuildStatus(buildStatusError, err.Error(), "")
		notifyBuildStatus()
		errlnf("Rebuild failed: %v", err)
		return
	}
	posts := m.Posts()
	setServePosts(posts)
	setServeSearchPosts(posts)
	if serveFast || serveIncremental {
		if cached, ok := m.Cache().Get("build_cache"); ok {
			if bc, ok := cached.(*buildcache.Cache); ok {
				setServeCache(bc)
			}
		}
	}

	// Check for cancellation after build
	select {
	case <-ctx.Done():
		verbosef("Rebuild canceled")
		return
	default:
	}

	duration := time.Since(startTime)
	setBuildStatus(buildStatusSuccess, "", licenseWarningMessage(getModelsConfig(m)))
	notifyBuildStatus()
	infof("Rebuilt in %.2fs (%d posts, %d feeds)", duration.Seconds(), result.PostsProcessed, result.FeedsGenerated)

	// Notify live reload clients
	notifyLiveReload()
}

func configureRebuildManager(m *lifecycle.Manager) {
	if serveFast {
		if m.Config().Extra == nil {
			m.Config().Extra = make(map[string]any)
		}
		m.Config().Extra["feeds_async"] = true
		applyFastMode(m)
		m.Config().Extra["cache_cleanup_async"] = true
	}
	if serveIncremental {
		lifecycle.SetServeIncremental(m, true)
		if m.Config().Extra == nil {
			m.Config().Extra = make(map[string]any)
		}
		m.Config().Extra["incremental_mode"] = true
	}
	if !serveFast && !serveIncremental {
		return
	}
	if cached := getServeCache(); cached != nil {
		m.Cache().Set("build_cache", cached)
	}
	if lifecycle.IsServeFullRebuild(m) {
		setServePosts(nil)
	} else if cachedPosts := getServePosts(); len(cachedPosts) > 0 {
		lifecycle.SetServeCachedPosts(m, cachedPosts)
	}
}

func configureServeIncremental(m *lifecycle.Manager, changedPaths, removedPaths []string, forceFull, globDirty bool) {
	if !serveFast && !serveIncremental {
		setFullServeRebuild(m)
		return
	}
	contentDir := m.Config().ContentDir
	if contentDir == "" {
		contentDir = "."
	}
	normalized, outside := normalizeServeChangedPaths(changedPaths, contentDir)
	normalizedRemoved, outsideRemoved := normalizeServeChangedPaths(removedPaths, contentDir)
	if (len(normalized) == 0 && len(normalizedRemoved) == 0) || forceFull || outside || outsideRemoved {
		setFullServeRebuild(m)
		if len(normalizedRemoved) > 0 {
			lifecycle.SetServeRemovedPaths(m, normalizedRemoved)
		}
		if verbose {
			verbosef("[serve] incremental disabled: normalized=%d removed=%d force_full=%t outside=%t content_dir=%s", len(normalized), len(normalizedRemoved), forceFull, outside || outsideRemoved, contentDir)
		}
		return
	}
	if incrementalPathsRequireFullRebuild(normalized, normalizedRemoved) {
		setFullServeRebuild(m)
		if len(normalizedRemoved) > 0 {
			lifecycle.SetServeRemovedPaths(m, normalizedRemoved)
		}
		if verbose {
			verbosef("[serve] incremental disabled: global input or removed content changed")
		}
		return
	}

	globDirty = serveGlobIsDirty(contentDir, normalized, globDirty)
	lifecycle.SetServeFullRebuild(m, false)
	lifecycle.SetServeChangedPaths(m, normalized)
	lifecycle.SetServeRemovedPaths(m, normalizedRemoved)
	lifecycle.SetServeGlobDirty(m, globDirty)
	if verbose {
		verbosef("[serve] incremental enabled: changed=%d removed=%d glob_dirty=%t content_dir=%s", len(normalized), len(normalizedRemoved), globDirty, contentDir)
	}
}

func serveGlobIsDirty(contentDir string, changed []string, dirty bool) bool {
	if dirty {
		return true
	}
	cached := getServeCache()
	if cached == nil {
		return false
	}
	cachedFiles, _ := cached.GetGlobCache()
	seen := make(map[string]struct{}, len(cachedFiles))
	for _, file := range cachedFiles {
		seen[file] = struct{}{}
	}
	for _, rel := range changed {
		if _, ok := seen[rel]; ok {
			continue
		}
		if _, err := os.Stat(filepath.Join(contentDir, rel)); err == nil {
			return true
		}
	}
	return false
}

func setFullServeRebuild(m *lifecycle.Manager) {
	lifecycle.SetServeFullRebuild(m, true)
	lifecycle.SetServeChangedPaths(m, nil)
	lifecycle.SetServeRemovedPaths(m, nil)
	lifecycle.SetServeGlobDirty(m, true)
}

func incrementalPathsRequireFullRebuild(changed, removed []string) bool {
	if len(removed) > 0 {
		return true
	}
	for _, path := range changed {
		ext := strings.ToLower(filepath.Ext(path))
		if ext != markdownShortExtension && ext != markdownExtension {
			return true
		}
	}
	return false
}

// addWatchPaths adds paths to the file watcher.
func addWatchPaths(watcher *fsnotify.Watcher, m *lifecycle.Manager) error {
	config := m.Config()

	// Watch current directory for config files
	if err := watcher.Add("."); err != nil {
		return err
	}

	// Watch content roots derived from configured glob patterns.
	// This avoids recursively watching generated trees when content_dir is ".".
	for _, contentRoot := range contentWatchRoots(config) {
		if _, err := os.Stat(contentRoot); err != nil {
			continue
		}
		if err := addDirRecursive(watcher, contentRoot); err != nil {
			return err
		}
	}

	// Watch templates directory
	templatesDir := "templates"
	if td, ok := config.Extra["templates_dir"].(string); ok && td != "" {
		templatesDir = td
	}
	if _, err := os.Stat(templatesDir); err == nil {
		if err := addDirRecursive(watcher, templatesDir); err != nil {
			return err
		}
	}

	// Watch static/assets directory
	assetsDir := "static"
	if ad, ok := config.Extra["assets_dir"].(string); ok && ad != "" {
		assetsDir = ad
	}
	if _, err := os.Stat(assetsDir); err == nil {
		if err := addDirRecursive(watcher, assetsDir); err != nil {
			return err
		}
	}

	return nil
}

// addDirRecursive recursively adds directories to the watcher.
func addDirRecursive(watcher *fsnotify.Watcher, root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Get absolute path for comparison
		absPath, pathErr := filepath.Abs(path)
		if pathErr != nil {
			absPath = path
		}

		// Skip output directory
		if isPathWithinDir(absPath, serveOutputPath) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		for _, internalPath := range serveInternalPaths {
			if isPathWithinDir(absPath, internalPath) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Skip hidden directories
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") && d.Name() != "." {
			return filepath.SkipDir
		}

		// Add directories to watcher
		if d.IsDir() {
			if verbose {
				verbosef("Watching: %s", path)
			}
			return watcher.Add(path)
		}

		return nil
	})
}

func resolveRequestPath(outputDir, requestPath string) (fullPath, cleanURLPath string, err error) {
	if requestPath == "" || requestPath == "/" {
		requestPath = "/" + serveIndexFile
	}

	cleanURLPath = path.Clean("/" + requestPath)
	relPath := strings.TrimPrefix(cleanURLPath, "/")
	if relPath == "" {
		relPath = serveIndexFile
		cleanURLPath = "/" + serveIndexFile
	}

	fullPath = filepath.Join(outputDir, filepath.FromSlash(relPath))
	if !isPathWithinDir(fullPath, outputDir) {
		return "", "", errors.New("resolved path escapes output directory")
	}

	return fullPath, cleanURLPath, nil
}

func isPathWithinDir(pathname, dir string) bool {
	if dir == "" {
		return false
	}
	if filepath.Clean(pathname) == filepath.Clean(dir) {
		return true
	}
	rel, err := filepath.Rel(dir, pathname)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func timerChannel(timer *time.Timer) <-chan time.Time {
	if timer == nil {
		return nil
	}
	return timer.C
}

func drainTimer(timer *time.Timer) {
	if timer == nil {
		return
	}
	select {
	case <-timer.C:
	default:
	}
}

// createServeManager builds the manager for serve and its rebuilds, keeping
// single-page mode when serve was given a Markdown file.
func createServeManager() (*lifecycle.Manager, error) {
	// Wait for any in-flight settings bake (and its rollback) to finish.
	serveConfigMu.Lock()
	defer serveConfigMu.Unlock()
	if serveSourceFile != "" {
		return createSinglePageManager(cfgFile, serveSourceFile)
	}
	return createManager(cfgFile)
}
