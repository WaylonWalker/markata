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

const (
	chartDataDisplayKey = "data_display"
	chartDataOpen       = "open"
	chartDataCollapsed  = "collapsed"
)

// ChartJSPlugin converts Chart.js JSON code blocks into rendered charts.
// It runs at the render stage (after markdown conversion).
type ChartJSPlugin struct {
	scriptCollector func(string) // standalone rendering collects generated scripts as assets
	config          models.ChartJSConfig
	idCounter       uint64
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
		if value, exists := cfgMap[chartDataDisplayKey]; exists {
			display, ok := value.(string)
			if !ok || (display != chartDataOpen && display != chartDataCollapsed) {
				return fmt.Errorf("chartjs.data_display must be \"open\" or \"collapsed\"")
			}
			p.config.DataDisplay = display
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

		options, err := parseChartFigureOptions(chartConfig)
		if err != nil {
			return fmt.Sprintf(`<div class="%s chartjs-error"><p>Chart.js Error: %s</p></div>`,
				html.EscapeString(p.config.ContainerClass), html.EscapeString(err.Error()))
		}

		// Preserve authored configuration for the fallback before applying presentation defaults.
		authoredJSON, err := json.Marshal(chartConfig)
		if err != nil {
			return match
		}
		applyHorizontalBarDefaults(chartConfig)

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
		fallback := chartStaticFallback(chartConfig, string(authoredJSON))
		display := options.DataDisplay
		if display == "" {
			display = p.config.DataDisplay
		}
		if display == chartDataCollapsed {
			fallback = strings.Replace(fallback, `class="chartjs-data" open`, `class="chartjs-data"`, 1)
		}
		// Preserve legacy output byte-for-byte unless the chart opts into the
		// new presentation/disclosure options. Legacy initialization targets
		// the wrapper div; opted-in charts target the actual disclosure.
		dataLookup := fmt.Sprintf("document.getElementById('%s-data')", chartID)
		if options.Title != "" || options.Caption != "" || options.DataDisplay != "" || display == chartDataCollapsed {
			dataLookup = fmt.Sprintf("document.querySelector('#%s-data > details')", chartID)
		}

		// Keep the fallback present and inspectable. Collapse it only after
		// Chart.js successfully initializes; never hide it when scripts fail.
		initScript := fmt.Sprintf(`
  (function() {
    const ctx = document.getElementById('%s');
    if (!ctx || typeof Chart !== 'function') return;
    ctx.hidden = false;
    try {
      new Chart(ctx, %s);
      const data = %s;
      if (data) data.open = false;
    } catch (error) {
      ctx.hidden = true;
    }
  })();`, chartID, string(safeJSON), dataLookup)
		initScripts = append(initScripts, initScript)

		label := "Chart visualization; expand Chart data for exact values"
		description := ""
		if options.Title != "" || options.Caption != "" {
			if options.Title != "" {
				label = options.Title
			}
			description = fmt.Sprintf(` aria-describedby="%s-caption"`, chartID)
		}
		container := fmt.Sprintf(`<div class="%s">
  <canvas id="%s" hidden role="img" aria-label="%s"%s></canvas>
  <div id="%s-data">%s</div>
</div>`, html.EscapeString(p.config.ContainerClass), chartID, html.EscapeString(label), description, chartID, fallback)
		return chartFigureHTML(options, chartConfig, chartID, container)
	})

	// If we found chartjs blocks, inject the scripts
	if len(initScripts) > 0 {
		result = p.injectChartJSScripts(result, initScripts)
	}

	post.ArticleHTML = result
	return nil
}

// chartFigureOptions are authoring metadata, never Chart.js configuration.
type chartFigureOptions struct {
	Title       string
	Caption     string
	DataDisplay string
}

func parseChartFigureOptions(config map[string]interface{}) (chartFigureOptions, error) {
	var options chartFigureOptions
	for _, field := range []struct {
		key   string
		value *string
	}{
		{"title", &options.Title},
		{"caption", &options.Caption},
		{chartDataDisplayKey, &options.DataDisplay},
	} {
		raw, exists := config[field.key]
		if !exists {
			continue
		}
		value, ok := raw.(string)
		if !ok {
			return options, fmt.Errorf("%s must be a string", field.key)
		}
		if field.key == chartDataDisplayKey && value != chartDataOpen && value != chartDataCollapsed {
			return options, fmt.Errorf("data_display must be \"open\" or \"collapsed\"")
		}
		*field.value = value
		delete(config, field.key)
	}
	return options, nil
}

// applyHorizontalBarDefaults adds compact geometry only to simple horizontal bars.
func applyHorizontalBarDefaults(config map[string]interface{}) {
	const barType = "bar"
	if config["type"] != barType {
		return
	}
	options, ok := config["options"].(map[string]interface{})
	if _, exists := config["options"]; exists && !ok {
		return
	}
	data, ok := config["data"].(map[string]interface{})
	if !ok {
		return
	}
	labels, ok := data["labels"].([]interface{})
	if !ok || len(labels) == 0 || len(labels) > 5 {
		return
	}
	datasets, ok := data["datasets"].([]interface{})
	if !ok || len(datasets) != 1 {
		return
	}
	dataset, ok := datasets[0].(map[string]interface{})
	if !ok || (dataset["type"] != nil && dataset["type"] != barType) {
		return
	}
	if horizontalBarIndexAxis(dataset, options) != "y" {
		return
	}
	if options == nil {
		options = make(map[string]interface{})
		config["options"] = options
	}
	if _, a := options["aspectRatio"]; !a {
		if _, b := options["maintainAspectRatio"]; !b {
			options["aspectRatio"] = 3
		}
	}
	if hasBarGeometry(dataset, options) {
		return
	}
	dataset["maxBarThickness"] = 14
	dataset["borderRadius"] = 7
	dataset["borderSkipped"] = false
}

// horizontalBarIndexAxis follows Chart.js dataset, dataset-options, then chart precedence.
func horizontalBarIndexAxis(dataset, options map[string]interface{}) string {
	if axis, ok := dataset["indexAxis"].(string); ok && axis != "" {
		return axis
	}
	if groups, ok := options["datasets"].(map[string]interface{}); ok {
		if bar, ok := groups["bar"].(map[string]interface{}); ok {
			if axis, ok := bar["indexAxis"].(string); ok && axis != "" {
				return axis
			}
		}
	}
	if axis, ok := options["indexAxis"].(string); ok && axis != "" {
		return axis
	}
	return "x"
}

func hasBarGeometry(dataset, options map[string]interface{}) bool {
	for _, key := range []string{"barThickness", "maxBarThickness", "barPercentage", "categoryPercentage", "borderRadius", "borderSkipped"} {
		if _, exists := dataset[key]; exists {
			return true
		}
		if _, exists := options[key]; exists {
			return true
		}
		for _, section := range []string{"datasets", "elements"} {
			group, ok := options[section].(map[string]interface{})
			if !ok {
				continue
			}
			bar, ok := group["bar"].(map[string]interface{})
			if !ok {
				continue
			}
			if _, exists := bar[key]; exists {
				return true
			}
		}
	}
	return false
}

func isCompactHorizontalBar(config map[string]interface{}) bool {
	data, ok := config["data"].(map[string]interface{})
	if !ok {
		return false
	}
	labels, ok := data["labels"].([]interface{})
	if !ok {
		return false
	}
	chartOptions, ok := config["options"].(map[string]interface{})
	if !ok {
		return false
	}
	return config["type"] == "bar" && chartOptions["indexAxis"] == "y" && len(labels) > 0 && len(labels) <= 4
}

func chartFigureHTML(options chartFigureOptions, config map[string]interface{}, chartID, container string) string {
	if options.Title == "" && options.Caption == "" {
		return container
	}
	class := "chartjs-figure"
	if isCompactHorizontalBar(config) {
		class += " chartjs-figure-compact"
	}
	var caption strings.Builder
	if options.Title != "" {
		fmt.Fprintf(&caption, `<strong class="chartjs-figure-title">%s</strong>`, html.EscapeString(options.Title))
	}
	if options.Caption != "" {
		fmt.Fprintf(&caption, `<p class="chartjs-figure-caption">%s</p>`, html.EscapeString(options.Caption))
	}
	return fmt.Sprintf(`<figure class="%s"><figcaption id="%s-caption">%s</figcaption>
%s</figure>`, class, chartID, caption.String(), container)
}

// chartStaticFallback provides semantic source data for non-JavaScript readers.
// Common categorical chart types become an accessible table. Less-structured
// Chart.js configurations preserve their JSON for inspection instead.
func chartStaticFallback(config map[string]interface{}, source string) string {
	data, ok := config["data"].(map[string]interface{})
	if !ok {
		return chartConfigFallback(source)
	}
	labels, ok := data["labels"].([]interface{})
	if !ok {
		return chartConfigFallback(source)
	}
	series, ok := data["datasets"].([]interface{})
	if !ok {
		return chartConfigFallback(source)
	}
	if len(labels) > 0 && len(labels) <= 200 && len(series) > 0 && len(series) <= 12 {
		var table strings.Builder
		_, _ = table.WriteString(`<details class="chartjs-data" open><summary>Chart data</summary><table><caption>Chart data</caption><thead><tr><th scope="col">Label</th>`)
		for i, raw := range series {
			dataset, ok := raw.(map[string]interface{})
			if !ok {
				return chartConfigFallback(source)
			}
			values, ok := dataset["data"].([]interface{})
			if !ok || len(values) != len(labels) {
				return chartConfigFallback(source)
			}
			label, hasLabel := dataset["label"].(string)
			if !hasLabel || label == "" {
				label = fmt.Sprintf("Series %d", i+1)
			}
			_, _ = table.WriteString(`<th scope="col">`)
			_, _ = table.WriteString(html.EscapeString(label))
			_, _ = table.WriteString("</th>")
		}
		_, _ = table.WriteString("</tr></thead><tbody>")
		for i, label := range labels {
			_, _ = table.WriteString(`<tr><th scope="row">`)
			_, _ = table.WriteString(html.EscapeString(fmt.Sprint(label)))
			_, _ = table.WriteString("</th>")
			for _, raw := range series {
				dataset, ok := raw.(map[string]interface{})
				if !ok {
					return chartConfigFallback(source)
				}
				values, ok := dataset["data"].([]interface{})
				if !ok || len(values) != len(labels) {
					return chartConfigFallback(source)
				}
				_, _ = table.WriteString("<td>")
				if values[i] != nil {
					_, _ = table.WriteString(html.EscapeString(fmt.Sprint(values[i])))
				}
				_, _ = table.WriteString("</td>")
			}
			_, _ = table.WriteString("</tr>")
		}
		_, _ = table.WriteString("</tbody></table></details>")
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
(function() {
const initCharts = function() {%s
};
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', initCharts, {once: true});
} else {
  initCharts();
}
})();
</script>`, chartURL, strings.Join(initScripts, ""))

	// Append the script to the end of the content
	if p.scriptCollector != nil {
		// Standalone hosts can insert a fragment after DOMContentLoaded.
		// Site output retains the existing initialization contract.
		script = strings.Replace(script, "document.addEventListener('DOMContentLoaded', function() {", "(() => { const init = () => {", 1)
		script = strings.Replace(script, "\n});\n</script>", "\n}; if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', init, {once:true}); else init(); })();\n</script>", 1)
		p.scriptCollector(script)
		return htmlContent
	}
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
