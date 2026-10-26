// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"encoding/csv"
	"html"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// CSVFencePlugin converts CSV code blocks into HTML tables.
// It processes article_html during the render stage, after render_markdown.
type CSVFencePlugin struct {
	// enabled controls whether the plugin processes CSV blocks
	enabled bool

	// tableClass is the CSS class for the generated table (default: "csv-table")
	tableClass string

	// hasHeader indicates whether the first row is a header (default: true)
	hasHeader bool

	// delimiter is the CSV field delimiter (default: ",")
	delimiter rune
}

// NewCSVFencePlugin creates a new CSVFencePlugin with default settings.
func NewCSVFencePlugin() *CSVFencePlugin {
	return &CSVFencePlugin{
		enabled:    true,
		tableClass: "csv-table",
		hasHeader:  true,
		delimiter:  ',',
	}
}

// Name returns the unique name of the plugin.
func (p *CSVFencePlugin) Name() string {
	return "csv_fence"
}

// Priority returns the plugin's priority for the render stage.
// Returns a positive priority to run after render_markdown.
func (p *CSVFencePlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageRender {
		return lifecycle.PriorityLate // Run after render_markdown
	}
	return lifecycle.PriorityDefault
}

// Configure reads configuration options for the plugin from config.Extra.
// Configuration is expected in config.Extra["csv_fence"] as a map.
func (p *CSVFencePlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra == nil {
		return nil
	}

	// Look for csv_fence configuration
	cfgRaw, ok := config.Extra["csv_fence"]
	if !ok {
		return nil
	}

	// Handle map[string]interface{} configuration
	if cfgMap, ok := cfgRaw.(map[string]interface{}); ok {
		if enabled, ok := cfgMap["enabled"].(bool); ok {
			p.enabled = enabled
		}
		if tableClass, ok := cfgMap["table_class"].(string); ok {
			p.tableClass = tableClass
		}
		if hasHeader, ok := cfgMap["has_header"].(bool); ok {
			p.hasHeader = hasHeader
		}
		if delimiter, ok := cfgMap["delimiter"].(string); ok && delimiter != "" {
			p.delimiter = rune(delimiter[0])
		}
	}

	return nil
}

// csvBlockRegex matches <pre><code class="language-csv">...</code></pre> blocks.
// It captures:
// - Group 1: optional attributes after language-csv (e.g., delimiter=";" has_header="false")
// - Group 2: the CSV content inside the code block
var csvBlockRegex = regexp.MustCompile(`(?s)<pre><code class="language-csv([^"]*)">(.*?)</code></pre>`)

// blockOptionRegex extracts key="value" options from block attributes.
var blockOptionRegex = regexp.MustCompile(`(\w+)="([^"]*)"`)

// Render processes article_html and converts CSV blocks to HTML tables.
// Posts with Skip=true or empty ArticleHTML are skipped.
func (p *CSVFencePlugin) Render(m *lifecycle.Manager) error {
	if !p.enabled {
		return nil
	}

	posts := m.FilterPosts(func(post *models.Post) bool {
		if post.Skip || post.ArticleHTML == "" {
			return false
		}
		return strings.Contains(post.ArticleHTML, "language-csv")
	})

	return m.ProcessPostsSliceConcurrently(posts, p.processPost)
}

// processPost converts CSV blocks to HTML tables in a single post's ArticleHTML.
func (p *CSVFencePlugin) processPost(post *models.Post) error {
	if post.Skip || post.ArticleHTML == "" {
		return nil
	}

	post.ArticleHTML = csvBlockRegex.ReplaceAllStringFunc(post.ArticleHTML, func(match string) string {
		return p.processCSVBlock(match)
	})

	return nil
}

// processCSVBlock converts a single CSV code block to an HTML table.
func (p *CSVFencePlugin) processCSVBlock(match string) string {
	submatches := csvBlockRegex.FindStringSubmatch(match)
	if len(submatches) < 3 {
		return match
	}

	attrs := submatches[1]
	csvContent := submatches[2]

	// Parse per-block options
	delimiter := p.delimiter
	hasHeader := p.hasHeader
	tableClass := p.tableClass

	// Unescape HTML entities in attributes (goldmark may have encoded them)
	attrs = html.UnescapeString(attrs)

	// Extract options from attributes
	options := blockOptionRegex.FindAllStringSubmatch(attrs, -1)
	for _, opt := range options {
		if len(opt) >= 3 {
			key := opt[1]
			value := opt[2]
			switch key {
			case "delimiter":
				if value != "" {
					delimiter = rune(value[0])
				}
			case "has_header":
				hasHeader = value == BoolTrue
			case "table_class":
				tableClass = value
			}
		}
	}

	// Decode HTML entities in the CSV content (goldmark may have encoded them)
	csvContent = html.UnescapeString(csvContent)

	// Parse CSV content
	records, err := p.parseCSV(csvContent, delimiter)
	if err != nil || len(records) == 0 {
		// Return original if parsing fails or empty
		return match
	}

	// Generate HTML table
	return p.generateTable(records, hasHeader, tableClass)
}

// parseCSV parses CSV content into a slice of records.
func (p *CSVFencePlugin) parseCSV(content string, delimiter rune) ([][]string, error) {
	// Trim leading/trailing whitespace
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, nil
	}

	reader := csv.NewReader(strings.NewReader(content))
	reader.Comma = delimiter
	reader.FieldsPerRecord = -1 // Allow variable number of fields
	reader.TrimLeadingSpace = true

	return reader.ReadAll()
}

// generateTable creates an HTML table from CSV records.
func (p *CSVFencePlugin) generateTable(records [][]string, hasHeader bool, tableClass string) string {
	if len(records) == 0 {
		return ""
	}

	var sb strings.Builder

	// Open table tag with class
	sb.WriteString(`<table class="`)
	sb.WriteString(html.EscapeString(tableClass))
	sb.WriteString("\">\n")

	startRow := 0

	// Generate header if configured
	if hasHeader && len(records) > 0 {
		sb.WriteString("  <thead>\n    <tr>\n")
		for _, cell := range records[0] {
			sb.WriteString("      <th>")
			sb.WriteString(html.EscapeString(cell))
			sb.WriteString("</th>\n")
		}
		sb.WriteString("    </tr>\n  </thead>\n")
		startRow = 1
	}

	// Generate body
	if startRow < len(records) {
		sb.WriteString("  <tbody>\n")
		for i := startRow; i < len(records); i++ {
			sb.WriteString("    <tr>\n")
			for _, cell := range records[i] {
				sb.WriteString("      <td>")
				sb.WriteString(html.EscapeString(cell))
				sb.WriteString("</td>\n")
			}
			sb.WriteString("    </tr>\n")
		}
		sb.WriteString("  </tbody>\n")
	}

	sb.WriteString("</table>")

	return sb.String()
}

// SetEnabled enables or disables the plugin.
func (p *CSVFencePlugin) SetEnabled(enabled bool) {
	p.enabled = enabled
}

// SetTableClass sets the CSS class for generated tables.
func (p *CSVFencePlugin) SetTableClass(class string) {
	p.tableClass = class
}

// SetHasHeader sets whether the first row is treated as a header.
func (p *CSVFencePlugin) SetHasHeader(hasHeader bool) {
	p.hasHeader = hasHeader
}

// SetDelimiter sets the CSV field delimiter.
func (p *CSVFencePlugin) SetDelimiter(delimiter rune) {
	p.delimiter = delimiter
}

// Ensure CSVFencePlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*CSVFencePlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*CSVFencePlugin)(nil)
	_ lifecycle.RenderPlugin    = (*CSVFencePlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*CSVFencePlugin)(nil)
)
