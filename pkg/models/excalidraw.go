package models

// ExcalidrawConfig configures static rendering of Excalidraw code fences.
type ExcalidrawConfig struct {
	Enabled  bool                     `json:"enabled" yaml:"enabled" toml:"enabled"`
	Chromium ExcalidrawChromiumConfig `json:"chromium" yaml:"chromium" toml:"chromium"`
}

// ExcalidrawChromiumConfig configures the local browser used for SVG export.
type ExcalidrawChromiumConfig struct {
	BrowserPath string `json:"browser_path" yaml:"browser_path" toml:"browser_path"`
	Timeout     int    `json:"timeout" yaml:"timeout" toml:"timeout"`
	NoSandbox   bool   `json:"no_sandbox" yaml:"no_sandbox" toml:"no_sandbox"`
}

// NewExcalidrawConfig returns defaults for the static Excalidraw renderer.
func NewExcalidrawConfig() ExcalidrawConfig {
	return ExcalidrawConfig{
		Enabled: true,
		Chromium: ExcalidrawChromiumConfig{
			Timeout: 30,
		},
	}
}
