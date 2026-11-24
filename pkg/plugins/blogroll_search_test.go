package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

// Regression for #1588: the blogroll reader supplies its own template context,
// so it must preserve backend selection when rendering the shared theme base.
func TestBlogrollReaderSearchBackendAssets(t *testing.T) {
	disabled := false
	tests := []struct {
		name            string
		search          models.SearchConfig
		wantEndpoint    string
		wantPagefindDir string
	}{
		{
			name: "bleve with configured endpoint",
			search: models.SearchConfig{
				Backend: "bleve",
				Bleve: models.BleveSearchConfig{
					Endpoint: "https://search.example.com/api/search",
				},
			},
			wantEndpoint: "https://search.example.com/api/search",
		},
		{
			name:         "bleve uses fallback endpoint",
			search:       models.SearchConfig{Backend: "bleve"},
			wantEndpoint: "/api/search",
		},
		{
			name: "pagefind uses selected bundle",
			search: models.SearchConfig{
				Backend: "pagefind",
				Pagefind: models.PagefindConfig{
					BundleDir: "_search",
				},
			},
			wantPagefindDir: "_search",
		},
		{
			name:   "disabled search includes no loader",
			search: models.SearchConfig{Enabled: &disabled, Backend: "bleve"},
		},
	}

	engine, err := templates.NewEngine("")
	if err != nil {
		t.Fatalf("create theme template engine: %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := lifecycle.NewManager()
			manager.Config().Extra["search"] = tt.search
			manager.Cache().Set("templates.engine", engine)

			blogroll := NewBlogrollPlugin()
			configMap := blogroll.configToMap(manager.Config())
			searchMap, ok := configMap["search"].(map[string]interface{})
			if !ok {
				t.Fatalf("blogroll config.search missing or wrong type: %#v", configMap["search"])
			}
			if got := searchMap["backend"]; got != tt.search.SearchBackend() {
				t.Fatalf("config.search.backend = %v, want %s", got, tt.search.SearchBackend())
			}

			html, err := blogroll.renderTemplate(manager, "reader.html", map[string]interface{}{
				"title":  "Reader",
				"config": configMap,
				"page":   map[string]interface{}{"number": 1},
			})
			if err != nil {
				t.Fatalf("render actual reader theme: %v", err)
			}

			bleveAssignment := "window.__markataSearchEndpoint = window.__markataSearchEndpoint || '"
			switch {
			case tt.wantEndpoint != "":
				if !strings.Contains(html, bleveAssignment+tt.wantEndpoint+"';") {
					t.Errorf("reader did not initialize configured Bleve endpoint %q", tt.wantEndpoint)
				}
				if !strings.Contains(html, "/js/bleve-search") {
					t.Error("reader omitted Bleve client script")
				}
				if strings.Contains(html, "pagefind-ui.css") || strings.Contains(html, "pagefind-ui.js") {
					t.Error("Bleve reader still renders Pagefind asset requests")
				}
			case tt.wantPagefindDir != "":
				if strings.Contains(html, bleveAssignment) {
					t.Error("Pagefind reader unexpectedly sets the Bleve endpoint")
				}
				if !strings.Contains(html, "var bundleDir = '/"+tt.wantPagefindDir+"';") {
					t.Errorf("reader omitted configured Pagefind bundle directory %q", tt.wantPagefindDir)
				}
				if !strings.Contains(html, "pagefind-ui.css") || !strings.Contains(html, "pagefind-ui.js") {
					t.Error("Pagefind reader omitted lazy CSS/JS asset loader")
				}
			default:
				if strings.Contains(html, bleveAssignment) ||
					strings.Contains(html, "pagefind-ui.css") ||
					strings.Contains(html, "pagefind-ui.js") {
					t.Error("disabled search still renders a search loader")
				}
			}
		})
	}
}
