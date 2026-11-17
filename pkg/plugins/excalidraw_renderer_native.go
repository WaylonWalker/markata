//go:build !js

package plugins

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/visualizations"

	"github.com/chromedp/chromedp"
)

const maxExcalidrawSceneBytes = 4 << 20

type chromiumExcalidrawRenderer struct {
	config   models.ExcalidrawConfig
	darkMode bool
	version  string
	server   *httptest.Server
	url      string

	mu             sync.Mutex
	allocatorClose context.CancelFunc
	browserCtx     context.Context //nolint:containedctx // chromedp requires retaining a browser context
	browserClose   context.CancelFunc
	closed         bool
	renderGate     chan struct{}
}

func newExcalidrawRenderer(config models.ExcalidrawConfig, darkMode bool) (excalidrawSVGRenderer, error) {
	switch config.Mode {
	case excalidrawModeAuto, excalidrawModeCached:
		return newCachedExcalidrawRenderer(config.Mode, darkMode)
	case excalidrawModeChromium:
		return newChromiumExcalidrawRenderer(config, darkMode)
	default:
		return nil, fmt.Errorf("unsupported Excalidraw mode %q (supported: auto, cached, chromium)", config.Mode)
	}
}

func newChromiumExcalidrawRenderer(config models.ExcalidrawConfig, darkMode bool) (excalidrawSVGRenderer, error) {
	info := checkChromiumDependency(config.Chromium.BrowserPath)
	if !info.IsInstalled {
		message := "static Excalidraw rendering requires Chrome or Chromium; install a browser or set [markata-go.excalidraw.chromium].browser_path"
		if config.Chromium.BrowserPath != "" {
			return nil, fmt.Errorf("%s (configured path %q is not usable)", message, config.Chromium.BrowserPath)
		}
		return nil, errors.New(message)
	}

	assets, err := visualizations.ExcalidrawAssetFS()
	if err != nil {
		return nil, err
	}
	version, err := visualizations.ExcalidrawRendererVersion()
	if err != nil {
		return nil, err
	}
	server := httptest.NewServer(http.FileServer(http.FS(assets)))
	renderer := &chromiumExcalidrawRenderer{
		config:     config,
		darkMode:   darkMode,
		version:    excalidrawVersionForTheme(version, darkMode),
		server:     server,
		url:        server.URL + "/renderer.html",
		renderGate: make(chan struct{}, 1),
	}
	return renderer, nil
}

func (r *chromiumExcalidrawRenderer) ID() string      { return excalidrawName }
func (r *chromiumExcalidrawRenderer) Version() string { return r.version }

func (r *chromiumExcalidrawRenderer) RenderSVG(ctx context.Context, req visualizations.Request) ([]byte, error) {
	if len(req.Source) > maxExcalidrawSceneBytes {
		return nil, fmt.Errorf("scene JSON exceeds the %d-byte Excalidraw input limit", maxExcalidrawSceneBytes)
	}
	if _, err := validateExcalidrawScene(req.Source); err != nil {
		return nil, err
	}
	if err := r.ensureBrowser(ctx); err != nil {
		return nil, err
	}
	select {
	case r.renderGate <- struct{}{}:
		defer func() { <-r.renderGate }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	timeCtx, cancel := context.WithTimeout(ctx, time.Duration(r.config.Chromium.Timeout)*time.Second)
	defer cancel()
	mergedCtx, mergeCancel := mergeContextLifetime(r.browserCtx, timeCtx)
	defer mergeCancel()

	sceneLiteral, err := json.Marshal(req.Source)
	if err != nil {
		return nil, fmt.Errorf("encode Excalidraw scene: %w", err)
	}
	script := "window.markataRenderExcalidrawSvg(JSON.parse(" + string(sceneLiteral) + ")," + fmt.Sprint(r.darkMode) + ")"
	svg, err := chromedp.Run(mergedCtx, chromedp.Evaluate[string](script, chromedp.EvalAwaitPromise))
	if err != nil {
		return nil, fmt.Errorf("export Excalidraw scene to SVG: %w", err)
	}
	return []byte(strings.TrimSpace(svg)), nil
}

func (r *chromiumExcalidrawRenderer) ensureBrowser(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return errors.New("excalidraw renderer is closed")
	}
	if r.browserCtx != nil {
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	info := checkChromiumDependency(r.config.Chromium.BrowserPath)
	if !info.IsInstalled {
		return errors.New("browser Chrome or Chromium became unavailable before SVG rendering")
	}

	allocatorOptions := append([]chromedp.ExecAllocatorOption(nil), chromedp.DefaultExecAllocatorOptions[:]...)
	allocatorOptions = append(allocatorOptions,
		chromedp.ExecPath(filepath.Clean(info.BinaryPath)),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("host-resolver-rules", "MAP * ~NOTFOUND, EXCLUDE localhost, EXCLUDE 127.0.0.1"),
	)
	if r.config.Chromium.NoSandbox {
		allocatorOptions = append(allocatorOptions, chromedp.NoSandbox)
	}
	allocatorCtx, allocatorClose := chromedp.NewExecAllocator(context.Background(), allocatorOptions...)
	browserCtx, browserClose := chromedp.NewContext(allocatorCtx)
	timeCtx, timeoutCancel := context.WithTimeout(ctx, time.Duration(r.config.Chromium.Timeout)*time.Second)
	stopTimeout := context.AfterFunc(timeCtx, browserClose)
	_, err := chromedp.Run(browserCtx, chromedp.Navigate(r.url))
	if err == nil {
		var ready bool
		ready, err = chromedp.Run(browserCtx, chromedp.Evaluate[bool]("window.markataExcalidrawReady", chromedp.EvalAwaitPromise))
		if err == nil && !ready {
			reason, reasonErr := chromedp.Run(browserCtx, chromedp.Evaluate[string]("window.markataExcalidrawError || 'unknown module loading error'"))
			if reasonErr != nil {
				err = fmt.Errorf("read local Excalidraw module error: %w", reasonErr)
			} else {
				err = fmt.Errorf("local Excalidraw module failed to load: %s", reason)
			}
		}
	}
	stopTimeout()
	if err == nil && timeCtx.Err() != nil {
		err = context.DeadlineExceeded
	}
	timeoutCancel()
	if err != nil {
		browserClose()
		allocatorClose()
		return fmt.Errorf("initialize local Excalidraw browser page: %w", err)
	}
	r.allocatorClose = allocatorClose
	r.browserCtx = browserCtx
	r.browserClose = browserClose
	return nil
}

func (r *chromiumExcalidrawRenderer) Close() error {
	r.mu.Lock()
	if r.closed {
		r.mu.Unlock()
		return nil
	}
	r.closed = true
	browserClose, allocatorClose := r.browserClose, r.allocatorClose
	r.mu.Unlock()
	if browserClose != nil {
		browserClose()
	}
	if allocatorClose != nil {
		allocatorClose()
	}
	r.server.Close()
	return nil
}

func availableExcalidrawModes() []string {
	return []string{excalidrawModeAuto, excalidrawModeCached, excalidrawModeChromium}
}
