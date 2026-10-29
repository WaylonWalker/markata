package plugins

// StandaloneSchemaVersion is the version of the existing standalone JSON result.
// New transports must not change this value without a compatibility plan.
const StandaloneSchemaVersion = 1

// StandaloneCapabilities describes the features of the shared standalone
// renderer, independent of any CLI, Python, HTTP, or WASM transport.
//
// Capabilities describe the implementation, not the options selected for a
// particular render. In particular, local bundled assets do not imply that
// authored images or links are free of external dependencies.
type StandaloneCapabilities struct {
	SchemaVersion    int                 `json:"schema_version"`
	Renderer         string              `json:"renderer"`
	Formats          []string            `json:"formats"`
	Options          []string            `json:"options"`
	Features         []string            `json:"features"`
	DiagramModes     map[string][]string `json:"diagram_modes"`
	AssetModes       []string            `json:"asset_modes"`
	LoadsSiteConfig  bool                `json:"loads_site_config"`
	ExternalURLAudit bool                `json:"external_url_audit"`
}

// GetStandaloneCapabilities returns a fresh, deterministic description of the
// currently implemented standalone renderer. Keep this conservative: only
// advertise behavior supported by RenderStandalone, not the full site builder.
func GetStandaloneCapabilities() StandaloneCapabilities {
	return StandaloneCapabilities{
		SchemaVersion: StandaloneSchemaVersion,
		Renderer:      "markata-standalone",
		Formats:       []string{"html_fragment"},
		Options:       []string{"asset_base", "inline_assets", "unsafe"},
		Features: []string{
			"admonitions", "anchors", "chartjs", "checklists", "code",
			"figures", "mermaid", "roads", "tables",
		},
		DiagramModes: map[string][]string{
			"chartjs": {"client"},
			"mermaid": {"client"},
		},
		AssetModes:       []string{"inline", "linked"},
		LoadsSiteConfig:  false,
		ExternalURLAudit: false,
	}
}
