package themes_test

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/templates"
	"github.com/WaylonWalker/markata/pkg/themes"
	"github.com/chromedp/chromedp"
)

// TestFeedPhotoResponsiveLayout_Browser checks the displayed geometry *and*
// downloaded srcset resolution. Density-based candidate selection must never
// make the archive's images render as tiny, density-adjusted thumbnails.
// This also covers optional card_classes wrappers in mixed-type feeds.
func TestFeedPhotoResponsiveLayout_Browser(t *testing.T) {
	browserPath := ""
	for _, candidate := range []string{"chromium", "chromium-browser", "google-chrome", "google-chrome-stable"} {
		if path, err := exec.LookPath(candidate); err == nil {
			browserPath = path
			break
		}
	}
	if browserPath == "" {
		t.Skip("Chromium is unavailable")
	}

	engine, err := templates.NewEngine("")
	if err != nil {
		t.Fatal(err)
	}
	render := func(name string) string {
		t.Helper()
		got, renderErr := engine.RenderToString("partials/cards/photo-card.html", map[string]interface{}{
			"post": map[string]interface{}{
				"image": "https://dropper.waylonwalker.com/" + name + ".png?token=abc",
				"href":  "/shots/" + name + "/",
				"title": name,
			},
		})
		if renderErr != nil {
			t.Fatal(renderErr)
		}
		return got
	}

	// Use decoded raster fixtures with known natural aspect ratios. An SVG
	// response behind a .webp URL can expose its fallback 300x150 viewport
	// instead, testing the SVG decoder rather than responsive CSS geometry.
	var imageBytes [2][]byte
	for i, dimensions := range [][2]int{{60, 40}, {40, 60}} {
		var encoded bytes.Buffer
		if encodeErr := png.Encode(&encoded, image.NewRGBA(image.Rect(0, 0, dimensions[0], dimensions[1]))); encodeErr != nil {
			t.Fatal(encodeErr)
		}
		imageBytes[i] = encoded.Bytes()
	}

	var fixture string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/css/") {
			css, readErr := themes.ReadStatic(strings.TrimPrefix(r.URL.Path, "/"))
			if readErr != nil {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/css")
			if _, writeErr := w.Write(css); writeErr != nil {
				t.Errorf("write CSS fixture: %v", writeErr)
			}
			return
		}
		if strings.HasSuffix(r.URL.Path, ".png") {
			imageIndex := 0
			if strings.Contains(r.URL.Path, "portrait") {
				imageIndex = 1
			}
			w.Header().Set("Content-Type", "image/png")
			if _, writeErr := w.Write(imageBytes[imageIndex]); writeErr != nil {
				t.Errorf("write PNG fixture: %v", writeErr)
			}
			return
		}
		if r.URL.Path == "/favicon.ico" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		if _, writeErr := fmt.Fprint(w, fixture); writeErr != nil {
			t.Errorf("write HTML fixture: %v", writeErr)
		}
	}))
	defer server.Close()

	// Mirror the actual stylesheet order of base.html, keeping the CSS from
	// the embedded theme rather than copying a special test-only CSS rule.
	direct := render("landscape")
	wrapped := render("portrait")
	fixture = `<!doctype html><html><head><meta name="viewport" content="width=device-width,initial-scale=1">
<link rel="stylesheet" href="/css/main.css">
<link rel="stylesheet" href="/css/cards.css">
<link rel="stylesheet" href="/css/feeds.css">
</head><body><div class="feed h-feed"><div class="posts posts-list">` +
		direct + `<div class="col-span-2">` + wrapped + `</div></div></div></body></html>`
	fixture = strings.ReplaceAll(fixture, "https://dropper.waylonwalker.com", server.URL)

	allocator, cancelAllocator := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(browserPath),
		chromedp.NoFirstRun,
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)...)
	defer cancelAllocator()
	browser, cancelBrowser := chromedp.NewContext(allocator)
	defer cancelBrowser()
	browser, cancelTimeout := context.WithTimeout(browser, 90*time.Second)
	defer cancelTimeout()

	type geometry struct {
		Width         int    `json:"width"`
		Height        int    `json:"height"`
		Figure        int    `json:"figure"`
		NaturalWidth  int    `json:"naturalWidth"`
		NaturalHeight int    `json:"naturalHeight"`
		CSSHeight     string `json:"cssHeight"`
		Contain       string `json:"contain"`
		Chosen        string `json:"chosen"`
	}
	const inspect = `[...document.querySelectorAll(".photo-figure img")].map(img => ({
		width: Math.round(img.getBoundingClientRect().width),
		height: Math.round(img.getBoundingClientRect().height),
		figure: Math.round(img.closest(".photo-figure").getBoundingClientRect().width),
		naturalWidth: img.naturalWidth,
		naturalHeight: img.naturalHeight,
		cssHeight: getComputedStyle(img).height,
		contain: getComputedStyle(img).contain,
		chosen: img.currentSrc
	}))`

	for _, tc := range []struct {
		name     string
		viewport int64
		minWidth int
	}{
		{name: "desktop-2x", viewport: 1280, minWidth: 590},
		{name: "mobile-2x", viewport: 390, minWidth: 300},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var measured []geometry
			err := chromedp.Do(browser,
				setDeviceMetricsOverride(tc.viewport, 900, 2, false),
				chromedp.Navigate(server.URL),
				pollVoid(`[...document.querySelectorAll(".photo-figure img")].length === 2 && [...document.querySelectorAll(".photo-figure img")].every(img => img.complete && img.naturalWidth > 0)`),
				// decoding="async" deliberately allows the decoder to complete
				// after fetch; inspect only after both raster decodes finish.
				evaluateVoid(`Promise.all([...document.querySelectorAll(".photo-figure img")].map(img => img.decode())).then(() => true)`),
				evaluateInto(inspect, &measured),
			)
			if err != nil {
				t.Fatal(err)
			}
			if len(measured) != 2 {
				t.Fatalf("got %d rendered photos, want 2", len(measured))
			}
			for i, img := range measured {
				if img.Width < tc.minWidth || img.Width > 600 || img.Figure != img.Width {
					t.Errorf("photo %d has collapsed or overflowing display slot: %+v", i, img)
				}
				if img.Contain == "size" {
					t.Errorf("photo %d is size-contained by sizes=auto, forcing a 300x150 fallback: %+v", i, img)
				}
				if img.Height <= 0 {
					t.Errorf("photo %d has zero rendered height", i)
				} else {
					ratio := float64(img.Height) / float64(img.Width)
					want := 2.0 / 3
					if i == 1 {
						want = 1.5
					}
					if ratio < want-0.03 || ratio > want+0.03 {
						t.Errorf("photo %d aspect ratio %0.2f, want %0.2f: %+v", i, ratio, want, img)
					}
				}
				// The CDN variant should approximately track the actual display
				// slot at deviceScaleFactor=2 rather than downloading originals.
				// All widths in the card srcset top out at 1200 pixels.
				idx := strings.LastIndex(img.Chosen, "w=")
				if idx < 0 {
					t.Errorf("photo %d did not use an optimized width variant: %s", i, img.Chosen)
					continue
				}
				selected, parseErr := strconv.Atoi(strings.SplitN(img.Chosen[idx+2:], "&", 2)[0])
				minPixels := img.Width * 3 / 2
				maxPixels := img.Width * 5 / 2
				if maxPixels > 1200 {
					maxPixels = 1200
				}
				if parseErr != nil || selected < minPixels || selected > maxPixels {
					t.Errorf("photo %d chose inappropriate source %q for %d CSS px at 2x density (want %d..%d source px)", i, img.Chosen, img.Width, minPixels, maxPixels)
				}
			}
		})
	}
}
