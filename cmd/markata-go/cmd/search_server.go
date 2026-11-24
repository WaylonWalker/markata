package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/searchapi"
	"github.com/spf13/cobra"
)

var (
	searchServerPort          int
	searchServerHost          string
	searchServerMode          string
	searchServerIndexDir      string
	searchServerHashPath      string
	searchServerIndexName     string
	searchServerRebuild       bool
	searchServerWatchDebounce time.Duration
)

const (
	searchServerModeRuntime  = "runtime-index"
	searchServerModeReadOnly = "read-only-index"
	searchServerModeWatch    = "watch-content"
)

var searchServerCmd = &cobra.Command{
	Use:   "search-server",
	Short: "Start a standalone search API server",
	Long: `Start a standalone bleve-backed search API server.

The server provides a single GET endpoint that returns JSON search results.
Private post content is never indexed — only explicit frontmatter title and description, plus safe navigation fields.

Modes:
  runtime-index    Load content and build or refresh a local index
	read-only-index  Open an existing prebuilt index without rebuilding
	watch-content    Load content and refresh the local index when content changes

Example:
  markata-go search-server
  markata-go search-server --port 8081
  markata-go search-server --mode read-only-index --index-dir /data/search.bleve
  curl "http://localhost:3001/api/search?q=golang&fuzzy=true&limit=10"`,
	RunE: runSearchServer,
}

func init() {
	rootCmd.AddCommand(searchServerCmd)
	searchServerCmd.Flags().IntVar(&searchServerPort, "port", 3001, "Port to listen on")
	searchServerCmd.Flags().StringVar(&searchServerHost, "host", "localhost", "Host to bind to")
	searchServerCmd.Flags().StringVar(&searchServerMode, "mode", searchServerModeRuntime, "server mode: runtime-index, watch-content, or read-only-index")
	searchServerCmd.Flags().StringVar(&searchServerIndexDir, "index-dir", "", "directory of the bleve index")
	searchServerCmd.Flags().StringVar(&searchServerHashPath, "hash-path", "", "path for the content hash file in runtime-index mode")
	searchServerCmd.Flags().StringVar(&searchServerIndexName, "index-name", "server", "named index suffix inside the default cache directory")
	searchServerCmd.Flags().BoolVar(&searchServerRebuild, "rebuild-index", false, "force a rebuild in runtime-index mode")
	searchServerCmd.Flags().DurationVar(&searchServerWatchDebounce, "watch-debounce", 750*time.Millisecond, "debounce duration for watch-content rebuilds")
}

func runSearchServer(cmd *cobra.Command, _ []string) error {
	if searchServerMode != searchServerModeRuntime && searchServerMode != searchServerModeReadOnly && searchServerMode != searchServerModeWatch {
		return newUsageError(fmt.Errorf("invalid --mode %q (expected runtime-index, watch-content, or read-only-index)", searchServerMode))
	}

	apiCfg := searchapi.DefaultConfig()
	cacheDir := filepath.Join(".markata", "cache")
	clientSearchEndpoint := defaultSearchEndpoint
	handlerSearchPath := defaultSearchEndpoint
	apiCfg.IndexName = searchServerIndexName
	apiCfg.HashPath = searchServerHashPath
	apiCfg.Rebuild = searchServerRebuild

	var (
		handler *searchapi.Handler
		posts   []*models.Post
	)

	if searchServerMode == searchServerModeRuntime || searchServerMode == searchServerModeWatch {
		app, err := loadListApp(cmd.Context())
		if err != nil {
			return fmt.Errorf("initialization failed: %w", err)
		}

		posts = app.Manager.Posts()
		clientSearchEndpoint, handlerSearchPath = applySearchServerModelsConfig(getModelsConfig(app.Manager), &apiCfg, cacheDir)
		handler = searchapi.NewHandler(posts, cacheDir, apiCfg)

		if searchServerMode == searchServerModeWatch {
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				return fmt.Errorf("create watcher: %w", err)
			}
			defer watcher.Close()

			for _, root := range searchContentWatchRoots(app.Manager.Config()) {
				if _, statErr := os.Stat(root); statErr != nil {
					continue
				}
				if addErr := searchAddDirRecursive(watcher, root); addErr != nil {
					return fmt.Errorf("watch %s: %w", root, addErr)
				}
			}

			ctx, cancelWatch := context.WithCancel(cmd.Context())
			defer cancelWatch()
			go watchSearchContent(ctx, watcher, handler)
		}
	} else {
		if searchServerIndexDir == "" {
			return newUsageError(fmt.Errorf("--index-dir is required in read-only-index mode"))
		}
		apiCfg.IndexDir = searchServerIndexDir
		handler = searchapi.NewReadOnlyHandler(searchServerIndexDir, apiCfg)
	}
	defer handler.Close()
	if err := handler.Warm(); err != nil {
		return fmt.Errorf("warm search index: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle(handlerSearchPath, handler)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","mode":%q,"posts":%d}`, searchServerMode, len(posts))
	})

	addr := fmt.Sprintf("%s:%d", searchServerHost, searchServerPort)
	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Fprintf(os.Stderr, "Search API listening on http://%s%s\n", addr, handlerSearchPath)
		if clientSearchEndpoint != handlerSearchPath {
			fmt.Fprintf(os.Stderr, "Configured client endpoint: %s\n", clientSearchEndpoint)
		}
		switch searchServerMode {
		case searchServerModeRuntime:
			fmt.Fprintf(os.Stderr, "Indexed %d posts (%d published)\n", len(posts), countPublished(posts))
		case searchServerModeWatch:
			fmt.Fprintf(os.Stderr, "Watching content for index updates (%d posts, %d published)\n", len(posts), countPublished(posts))
		default:
			fmt.Fprintf(os.Stderr, "Serving read-only index from %s\n", searchServerIndexDir)
		}
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			os.Exit(1)
		}
	}()

	<-stop
	fmt.Fprintln(os.Stderr, "\nShutting down...")
	return server.Close()
}

func applySearchServerModelsConfig(mc *models.Config, apiCfg *searchapi.Config, cacheDir string) (clientEndpoint, handlerPath string) {
	clientEndpoint = defaultSearchEndpoint
	handlerPath = defaultSearchEndpoint
	if mc == nil {
		return clientEndpoint, handlerPath
	}

	clientEndpoint, handlerPath = configuredSearchEndpoints(mc)
	apiCfg.DefaultLimit = mc.Search.Bleve.DefaultLimit()
	apiCfg.MaxLimit = mc.Search.Bleve.GetMaxLimit()
	apiCfg.DefaultFuzzy = mc.Search.Bleve.IsFuzzy()
	apiCfg.IndexDir = searchServerIndexDir
	if searchServerIndexDir == "" && searchServerIndexName != "" {
		apiCfg.IndexDir = filepath.Join(cacheDir, "search-"+searchServerIndexName+".bleve")
	}
	if len(mc.Search.Bleve.CORSOrigins) > 0 {
		apiCfg.CORSOrigins = mc.Search.Bleve.CORSOrigins
	}

	return clientEndpoint, handlerPath
}

func watchSearchContent(ctx context.Context, watcher *fsnotify.Watcher, handler *searchapi.Handler) {
	var timer *time.Timer
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if searchShouldIgnorePath(event.Name) {
				continue
			}
			searchHandleNewDirectory(watcher, event)
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) == 0 {
				continue
			}
			if timer != nil {
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
			}
			timer = time.AfterFunc(searchServerWatchDebounce, func() {
				app, err := loadListApp(ctx)
				if err != nil {
					errlnf("Search watcher reload failed: %v", err)
					return
				}
				handler.UpdatePosts(app.Manager.Posts())
				if err := handler.Warm(); err != nil {
					errlnf("Search watcher warm failed: %v", err)
				}
			})
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			errlnf("Search watcher error: %v", err)
		}
	}
}

func countPublished(posts []*models.Post) int {
	n := 0
	for _, p := range posts {
		if p.Published && !p.Draft && !p.Skip {
			n++
		}
	}
	return n
}
