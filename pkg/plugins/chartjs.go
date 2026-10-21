// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"encoding/json"
	"fmt"
	"html"
	"regexp"
	"strings"
	"sync/atomic"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// ChartJSPlugin converts Chart.js JSON code blocks into rendered charts.
// It runs at the render stage (after markdown conversion).
type ChartJSPlugin struct {
	config    models.ChartJSConfig
	idCounter uint64
}

// NewChartJSPlugin creates a new ChartJSPlugin with default settings.
func NewChartJSPlugin() *ChartJSPlugin {
	return &ChartJSPlugin{
		config: models.NewChartJSConfig(),
	}
}

// Name returns the unique name of the plugin.
func (p *ChartJSPlugin) Name() string {
	return "chartjs"
}

// Priority returns the plugin's priority for a given stage.
// This plugin runs after render_markdown (which has default priority 0).
func (p *ChartJSPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageRender {
		return lifecycle.PriorityLate // Run after render_markdown
	}
	return lifecycle.PriorityDefault
}

// Configure reads configuration options for the plugin from config.Extra.
// Configuration is expected under the "chartjs" key.
func (p *ChartJSPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra == nil {
		return nil
	}

	if url := resolveAssetURL(assetURLsFromConfig(config), "chartjs", ""); url != "" {
		p.config.CDNURL = url
	}

	// Check for chartjs config in Extra
	chartjsConfig, ok := config.Extra["chartjs"]
	if !ok {
		return nil
	}

	// Handle map configuration
	if cfgMap, ok := chartjsConfig.(map[string]interface{}); ok {
		if enabled, ok := cfgMap["enabled"].(bool); ok {
			p.config.Enabled = enabled
		}
		if cdnURL, ok := cfgMap["cdn_url"].(string); ok && cdnURL != "" {
			p.config.CDNURL = cdnURL
		}
		if containerClass, ok := cfgMap["container_class"].(string); ok && containerClass != "" {
			p.config.ContainerClass = containerClass
		}
	}

	return nil
}

// Render processes chartjs code blocks in the rendered HTML for all posts.
func (p *ChartJSPlugin) Render(m *lifecycle.Manager) error {
	if !p.config.Enabled {
		return nil
	}

	posts := m.FilterPosts(func(post *models.Post) bool {
		if post.Skip || post.ArticleHTML == "" {
			return false
		}
		return strings.Contains(post.ArticleHTML, `class="language-chartjs"`)
	})

	return m.ProcessPostsSliceConcurrently(posts, p.processPost)
}

// chartjsCodeBlockRegex matches <pre><code class="language-chartjs"> blocks.
// It captures the JSON content inside.
var chartjsCodeBlockRegex = regexp.MustCompile(
	`<pre><code class="language-chartjs"[^>]*>([\s\S]*?)</code></pre>`,
)

// processPost processes a single post's HTML for chartjs code blocks.
func (p *ChartJSPlugin) processPost(post *models.Post) error {
	// Skip posts marked as skip or with no HTML content
	if post.Skip || post.ArticleHTML == "" {
		return nil
	}

	// Check if there are any chartjs code blocks
	if !strings.Contains(post.ArticleHTML, `class="language-chartjs"`) {
		return nil
	}

	// Track if we found any chartjs blocks and collect initialization scripts
	var initScripts []string

	// Replace chartjs code blocks with canvas elements
	result := chartjsCodeBlockRegex.ReplaceAllStringFunc(post.ArticleHTML, func(match string) string {
		// Extract the JSON content
		submatches := chartjsCodeBlockRegex.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}

		// Decode HTML entities in the JSON (goldmark encodes them)
		jsonContent := html.UnescapeString(submatches[1])
		jsonContent = strings.TrimSpace(jsonContent)

		// Validate that it's valid JSON
		var chartConfig map[string]interface{}
		if err := json.Unmarshal([]byte(jsonContent), &chartConfig); err != nil {
			// Invalid JSON, return with an error comment
			return fmt.Sprintf(`<div class="%s chartjs-error">
  <p>Chart.js Error: Invalid JSON configuration</p>
  <pre>%s</pre>
</div>`, p.config.ContainerClass, html.EscapeString(err.Error()))
		}

		// Re-encode decoded JSON using encoding/json's HTML-safe escaping.
		// Inline scripts must never contain an authored </script> sequence.
		safeJSON, err := json.Marshal(chartConfig)
		if err != nil {
			return match
		}

		// Generate a unique ID for this chart.
		chartID := fmt.Sprintf("chartjs-%d", atomic.AddUint64(&p.idCounter, 1))

		// Render a real data table (or the configuration for unusual charts)
		// so the data remains available with JavaScript disabled or offline.
		fallback := chartStaticFallback(chartConfig, string(safeJSON))

		// Keep the fallback present and inspectable. Collapse it only after
		// Chart.js successfully initializes; never hide it when scripts fail.
		initScript := fmt.Sprintf(`
  (function() {
    const ctx = document.getElementById('%s');
    if (!ctx || typeof Chart !== 'function') return;
    ctx.hidden = false;
    try {
      new Chart(ctx, %s);
      const data = document.getElementById('%s-data');
      if (data) data.open = false;
    } catch (error) {
      ctx.hidden = true;
    }
  })();`, chartID, string(safeJSON), chartID)
		initScripts = append(initScripts, initScript)

		return fmt.Sprintf(`<div class="%s">
  <canvas id="%s" hidden role="img" aria-label="Chart visualization; expand Chart data for exact values"></canvas>
  <div id="%s-data">%s</div>
</div>`, html.EscapeString(p.config.ContainerClass), chartID, chartID, fallback)
	})

	// If we found chartjs blocks, inject the scripts
	if len(initScripts) > 0 {
		result = p.injectChartJSScripts(result, initScripts)
	}

	post.ArticleHTML = result
	return nil
}

// chartStaticFallback provides semantic source data for non-JavaScript readers.
// Common categorical chart types become an accessible table. Less-structured
// Chart.js configurations preserve their JSON for inspection instead.
func chartStaticFallback(config map[string]interface{}, source string) string {
	data, _ := config["data"].(map[string]interface{})
	labels, _ := data["labels"].([]interface{})
	series, _ := data["datasets"].([]interface{})
	if len(labels) > 0 && len(labels) <= 200 && len(series) > 0 && len(series) <= 12 {
		var table strings.Builder
		table.WriteString(`<details class="chartjs-data" open><summary>Chart data</summary><table><caption>Chart data</caption><thead><tr><th scope="col">Label</th>`)
		for i, raw := range series {
			dataset, ok := raw.(map[string]interface{})
			if !ok {
				return chartConfigFallback(source)
			}
			values, ok := dataset["data"].([]interface{})
			if !ok || len(values) != len(labels) {
				return chartConfigFallback(source)
			}
			label, _ := dataset["label"].(string)
			if label == "" {
				label = fmt.Sprintf("Series %d", i+1)
			}
			table.WriteString(`<th scope="col">`)
			table.WriteString(html.EscapeString(label))
			table.WriteString("</th>")
		}
		table.WriteString("</tr></thead><tbody>")
		for i, label := range labels {
			table.WriteString(`<tr><th scope="row">`)
			table.WriteString(html.EscapeString(fmt.Sprint(label)))
			table.WriteString("</th>")
			for _, raw := range series {
				dataset := raw.(map[string]interface{})
				values := dataset["data"].([]interface{})
				table.WriteString("<td>")
				if values[i] != nil {
					table.WriteString(html.EscapeString(fmt.Sprint(values[i])))
				}
				table.WriteString("</td>")
			}
			table.WriteString("</tr>")
		}
		table.WriteString("</tbody></table></details>")
		return table.String()
	}
	return chartConfigFallback(source)
}

func chartConfigFallback(source string) string {
	return `<details class="chartjs-data" open><summary>Chart configuration</summary><pre><code>` +
		html.EscapeString(source) + "</code></pre></details>"
}

// injectChartJSScripts adds the Chart.js library and initialization scripts to the HTML.
func (p *ChartJSPlugin) injectChartJSScripts(htmlContent string, initScripts []string) string {
	// Build the combined script
	chartURL := p.config.CDNURL
	if strings.HasSuffix(chartURL, "/") {
		chartURL = strings.TrimRight(chartURL, "/") + "/chart.min.js"
	}

	script := fmt.Sprintf(`
<script src="%s"></script>
<script>
document.addEventListener('DOMContentLoaded', function() {%s
});
</script>`, chartURL, strings.Join(initScripts, ""))

	// Append the script to the end of the content
	return htmlContent + script
}

// SetConfig sets the chartjs configuration directly.
// This is useful for testing or programmatic configuration.
func (p *ChartJSPlugin) SetConfig(config models.ChartJSConfig) {
	p.config = config
}

// Config returns the current chartjs configuration.
func (p *ChartJSPlugin) Config() models.ChartJSConfig {
	return p.config
}

// Ensure ChartJSPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*ChartJSPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*ChartJSPlugin)(nil)
	_ lifecycle.RenderPlugin    = (*ChartJSPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*ChartJSPlugin)(nil)
)
