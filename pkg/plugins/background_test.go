package plugins

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestBackgroundPlugin_Name(t *testing.T) {
	p := NewBackgroundPlugin()
	if p.Name() != "background" {
		t.Errorf("expected name 'background', got %q", p.Name())
	}
}

func TestBackgroundPlugin_Interfaces(t *testing.T) {
	_ = t // suppress unused parameter warning
	p := NewBackgroundPlugin()

	// Verify plugin implements required interfaces
	var _ lifecycle.Plugin = p
	var _ lifecycle.ConfigurePlugin = p
}

func TestBackgroundPlugin_DefaultDisabled(t *testing.T) {
	p := NewBackgroundPlugin()
	if p.IsEnabled() {
		t.Error("expected background to be disabled by default")
	}
}

func TestBackgroundPlugin_Configure(t *testing.T) {
	tests := []struct {
		name        string
		extra       map[string]interface{}
		wantEnabled bool
		wantCSS     string
		wantScripts []string
		wantBgs     int
		wantErr     bool
	}{
		{
			name:        "empty config",
			extra:       nil,
			wantEnabled: false,
		},
		{
			name: "enabled with backgrounds",
			extra: map[string]interface{}{
				"theme": map[string]interface{}{
					"background": map[string]interface{}{
						"enabled": true,
						"backgrounds": []interface{}{
							map[string]interface{}{
								"html":    `<snow-fall count="200"></snow-fall>`,
								"z_index": int64(-10),
							},
						},
						"scripts": []interface{}{"/static/js/snow-fall.js"},
						"css":     ".snow { color: white; }",
					},
				},
			},
			wantEnabled: true,
			wantCSS:     ".snow { color: white; }",
			wantScripts: []string{"/static/js/snow-fall.js"},
			wantBgs:     1,
		},
		{
			name: "multiple backgrounds",
			extra: map[string]interface{}{
				"theme": map[string]interface{}{
					"background": map[string]interface{}{
						"enabled": true,
						"backgrounds": []interface{}{
							map[string]interface{}{
								"html": `<div class="stars"></div>`,
							},
							map[string]interface{}{
								"html":    `<div class="clouds"></div>`,
								"z_index": int64(-5),
							},
						},
					},
				},
			},
			wantEnabled: true,
			wantBgs:     2,
		},
		{
			name: "empty html in background",
			extra: map[string]interface{}{
				"theme": map[string]interface{}{
					"background": map[string]interface{}{
						"enabled": true,
						"backgrounds": []interface{}{
							map[string]interface{}{
								"html": "",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "script tag in html",
			extra: map[string]interface{}{
				"theme": map[string]interface{}{
					"background": map[string]interface{}{
						"enabled": true,
						"backgrounds": []interface{}{
							map[string]interface{}{
								"html": `<script>alert("bad")</script>`,
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewBackgroundPlugin()
			m := lifecycle.NewManager()

			config := lifecycle.NewConfig()
			config.Extra = tt.extra
			m.SetConfig(config)

			err := p.Configure(m)
			if (err != nil) != tt.wantErr {
				t.Errorf("Configure() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if p.IsEnabled() != tt.wantEnabled {
				t.Errorf("IsEnabled() = %v, want %v", p.IsEnabled(), tt.wantEnabled)
			}

			cfg := p.Config()
			if cfg.CSS != tt.wantCSS {
				t.Errorf("CSS = %q, want %q", cfg.CSS, tt.wantCSS)
			}

			if len(cfg.Scripts) != len(tt.wantScripts) {
				t.Errorf("Scripts length = %d, want %d", len(cfg.Scripts), len(tt.wantScripts))
			}

			if len(cfg.Backgrounds) != tt.wantBgs {
				t.Errorf("Backgrounds length = %d, want %d", len(cfg.Backgrounds), tt.wantBgs)
			}
		})
	}
}

func TestBackgroundPlugin_GenerateBackgroundHTML(t *testing.T) {
	tests := []struct {
		name     string
		config   models.BackgroundConfig
		contains []string
		empty    bool
	}{
		{
			name: "disabled",
			config: models.BackgroundConfig{
				Enabled:     bgBoolPtr(false),
				Backgrounds: []models.BackgroundElement{{HTML: "<div></div>"}},
			},
			empty: true,
		},
		{
			name: "no backgrounds",
			config: models.BackgroundConfig{
				Enabled:     bgBoolPtr(true),
				Backgrounds: []models.BackgroundElement{},
			},
			empty: true,
		},
		{
			name: "single background",
			config: models.BackgroundConfig{
				Enabled: bgBoolPtr(true),
				Backgrounds: []models.BackgroundElement{
					{HTML: `<snow-fall count="200"></snow-fall>`},
				},
			},
			contains: []string{
				"background-layer",
				`<snow-fall count="200"></snow-fall>`,
				"z-index: -1",
			},
		},
		{
			name: "background with z-index",
			config: models.BackgroundConfig{
				Enabled: bgBoolPtr(true),
				Backgrounds: []models.BackgroundElement{
					{HTML: `<div class="stars"></div>`, ZIndex: -10},
				},
			},
			contains: []string{
				"z-index: -10",
				`<div class="stars"></div>`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewBackgroundPlugin()
			p.SetConfig(tt.config)

			result := p.GenerateBackgroundHTML()

			if tt.empty {
				if result != "" {
					t.Errorf("expected empty result, got %q", result)
				}
				return
			}

			for _, want := range tt.contains {
				if !bgContainsString(result, want) {
					t.Errorf("result missing %q:\n%s", want, result)
				}
			}
		})
	}
}

func TestBackgroundPlugin_GenerateBackgroundCSS(t *testing.T) {
	tests := []struct {
		name     string
		config   models.BackgroundConfig
		contains string
		empty    bool
	}{
		{
			name: "disabled",
			config: models.BackgroundConfig{
				Enabled: bgBoolPtr(false),
				CSS:     ".test { color: red; }",
			},
			empty: true,
		},
		{
			name: "no css",
			config: models.BackgroundConfig{
				Enabled: bgBoolPtr(true),
				CSS:     "",
			},
			empty: true,
		},
		{
			name: "with css",
			config: models.BackgroundConfig{
				Enabled: bgBoolPtr(true),
				CSS:     ".snow { opacity: 0.8; }",
			},
			contains: ".snow { opacity: 0.8; }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewBackgroundPlugin()
			p.SetConfig(tt.config)

			result := p.GenerateBackgroundCSS()

			if tt.empty {
				if result != "" {
					t.Errorf("expected empty result, got %q", result)
				}
				return
			}

			if !bgContainsString(result, tt.contains) {
				t.Errorf("result missing %q:\n%s", tt.contains, result)
			}

			if !bgContainsString(result, "<style>") {
				t.Error("result missing <style> tag")
			}
		})
	}
}

func TestBackgroundPlugin_GenerateBackgroundScripts(t *testing.T) {
	tests := []struct {
		name     string
		config   models.BackgroundConfig
		contains []string
		empty    bool
	}{
		{
			name: "disabled",
			config: models.BackgroundConfig{
				Enabled: bgBoolPtr(false),
				Scripts: []string{"/js/test.js"},
			},
			empty: true,
		},
		{
			name: "no scripts",
			config: models.BackgroundConfig{
				Enabled: bgBoolPtr(true),
				Scripts: []string{},
			},
			empty: true,
		},
		{
			name: "single script",
			config: models.BackgroundConfig{
				Enabled: bgBoolPtr(true),
				Scripts: []string{"/static/js/snow-fall.js"},
			},
			contains: []string{
				`<script src="/static/js/snow-fall.js"></script>`,
			},
		},
		{
			name: "multiple scripts",
			config: models.BackgroundConfig{
				Enabled: bgBoolPtr(true),
				Scripts: []string{"/js/particles.js", "/js/snow.js"},
			},
			contains: []string{
				`<script src="/js/particles.js"></script>`,
				`<script src="/js/snow.js"></script>`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewBackgroundPlugin()
			p.SetConfig(tt.config)

			result := p.GenerateBackgroundScripts()

			if tt.empty {
				if result != "" {
					t.Errorf("expected empty result, got %q", result)
				}
				return
			}

			for _, want := range tt.contains {
				if !bgContainsString(result, want) {
					t.Errorf("result missing %q:\n%s", want, result)
				}
			}
		})
	}
}

func TestBackgroundConfig_IsEnabled(t *testing.T) {
	tests := []struct {
		name    string
		enabled *bool
		want    bool
	}{
		{"nil defaults to false", nil, false},
		{"explicit true", bgBoolPtr(true), true},
		{"explicit false", bgBoolPtr(false), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := models.BackgroundConfig{Enabled: tt.enabled}
			if got := cfg.IsEnabled(); got != tt.want {
				t.Errorf("IsEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function
func bgContainsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || s != "" && bgContainsSubstring(s, substr))
}

func bgContainsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// bgBoolPtr returns a pointer to the given bool value
func bgBoolPtr(b bool) *bool {
	return &b
}
