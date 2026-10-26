// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"encoding/json"
	"fmt"
	"html"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// ContributionGraphPlugin converts contribution-graph JSON code blocks into
// rendered Cal-Heatmap calendar heatmaps showing GitHub-style activity.
// It runs at the render stage (after markdown conversion).
type ContributionGraphPlugin struct {
	config    models.ContributionGraphConfig
	idCounter uint64
	assetURLs map[string]string
}

// NewContributionGraphPlugin creates a new ContributionGraphPlugin with default settings.
func NewContributionGraphPlugin() *ContributionGraphPlugin {
	return &ContributionGraphPlugin{
		config: models.NewContributionGraphConfig(),
	}
}

// Name returns the unique name of the plugin.
func (p *ContributionGraphPlugin) Name() string {
	return "contribution_graph"
}

// Priority returns the plugin's priority for a given stage.
// This plugin runs after render_markdown (which has default priority 0).
func (p *ContributionGraphPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageRender {
		return lifecycle.PriorityLate // Run after render_markdown
	}
	return lifecycle.PriorityDefault
}

// Configure reads configuration options for the plugin from config.Extra.
// Configuration is expected under the "contribution_graph" key.
func (p *ContributionGraphPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra == nil {
		return nil
	}

	p.assetURLs = assetURLsFromConfig(config)
	if url := resolveAssetURL(p.assetURLs, "cal-heatmap-css", ""); url != "" {
		p.config.CDNURL = strings.TrimSuffix(url, "/cal-heatmap.css")
	}

	// Check for contribution_graph config in Extra
	graphConfig, ok := config.Extra["contribution_graph"]
	if !ok {
		return nil
	}

	// Handle map configuration
	if cfgMap, ok := graphConfig.(map[string]interface{}); ok {
		if enabled, ok := cfgMap["enabled"].(bool); ok {
			p.config.Enabled = enabled
		}
		if cdnURL, ok := cfgMap["cdn_url"].(string); ok && cdnURL != "" {
			p.config.CDNURL = cdnURL
		}
		if containerClass, ok := cfgMap["container_class"].(string); ok && containerClass != "" {
			p.config.ContainerClass = containerClass
		}
		if theme, ok := cfgMap["theme"].(string); ok && theme != "" {
			p.config.Theme = theme
		}
		if percentile, ok := cfgMap["scale_max_percentile"].(float64); ok {
			p.config.ScaleMaxPercentile = percentile
		}
	}

	return nil
}

// Render processes contribution-graph code blocks in the rendered HTML for all posts.
func (p *ContributionGraphPlugin) Render(m *lifecycle.Manager) error {
	if !p.config.Enabled {
		return nil
	}

	posts := m.FilterPosts(func(post *models.Post) bool {
		if post.Skip || post.ArticleHTML == "" {
			return false
		}
		return strings.Contains(post.ArticleHTML, `class="language-contribution-graph"`)
	})

	return m.ProcessPostsSliceConcurrently(posts, func(post *models.Post) error {
		return p.processPost(m, post)
	})
}

// contributionGraphCodeBlockRegex matches <pre><code class="language-contribution-graph"> blocks.
// It captures the JSON content inside.
var contributionGraphCodeBlockRegex = regexp.MustCompile(
	`<pre><code class="language-contribution-graph"[^>]*>([\s\S]*?)</code></pre>`,
)

// processPost processes a single post's HTML for contribution-graph code blocks.
func (p *ContributionGraphPlugin) processPost(m *lifecycle.Manager, post *models.Post) error {
	// Skip posts marked as skip or with no HTML content
	if post.Skip || post.ArticleHTML == "" {
		return nil
	}

	// Check if there are any contribution-graph code blocks
	if !strings.Contains(post.ArticleHTML, `class="language-contribution-graph"`) {
		return nil
	}

	// Track if we found any contribution-graph blocks and collect initialization scripts
	var initScripts []string

	// Replace contribution-graph code blocks with div elements
	result := contributionGraphCodeBlockRegex.ReplaceAllStringFunc(post.ArticleHTML, func(match string) string {
		// Extract the JSON content
		submatches := contributionGraphCodeBlockRegex.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}

		// Decode HTML entities in the JSON (goldmark encodes them)
		jsonContent := html.UnescapeString(submatches[1])
		jsonContent = strings.TrimSpace(jsonContent)

		// Validate and parse the JSON content
		var graphConfig map[string]interface{}
		if err := json.Unmarshal([]byte(jsonContent), &graphConfig); err != nil {
			// Invalid JSON, return with an error comment
			return fmt.Sprintf(`<div class="%s contribution-graph-error">
  <p>Contribution Graph Error: Invalid JSON configuration</p>
  <pre>%s</pre>
</div>`, p.config.ContainerClass, html.EscapeString(err.Error()))
		}

		// Generate a unique ID for this graph
		graphID := fmt.Sprintf("contribution-graph-%d", atomic.AddUint64(&p.idCounter, 1))

		// Extract data and options from the config
		options := graphConfig["options"]
		if options == nil {
			options = map[string]interface{}{}
		}
		optionsMap, ok := options.(map[string]interface{})
		if !ok || optionsMap == nil {
			optionsMap = map[string]interface{}{}
		}

		data := graphConfig["data"]
		if data == nil {
			data = p.buildContributionData(m, optionsMap)
		}

		// Auto-detect year from first data point if not specified
		if _, hasYear := optionsMap["year"]; !hasYear {
			if dataSlice, ok := data.([]interface{}); ok && len(dataSlice) > 0 {
				if firstItem, ok := dataSlice[0].(map[string]interface{}); ok {
					if dateStr, ok := firstItem["date"].(string); ok && len(dateStr) >= 4 {
						// Extract year from date string (YYYY-MM-DD format)
						if year, err := strconv.Atoi(dateStr[:4]); err == nil {
							optionsMap["year"] = float64(year)
						}
					}
				}
			}
		}

		maxValue := p.computeColorScaleMax(data, optionsMap)

		// Marshal data and options back to JSON for the script
		dataJSON, err := json.Marshal(data)
		if err != nil {
			return fmt.Sprintf(`<div class="%s contribution-graph-error">
  <p>Contribution Graph Error: Failed to serialize data</p>
  <pre>%s</pre>
</div>`, p.config.ContainerClass, html.EscapeString(err.Error()))
		}
		optionsJSON, err := json.Marshal(optionsMap)
		if err != nil {
			return fmt.Sprintf(`<div class="%s contribution-graph-error">
  <p>Contribution Graph Error: Failed to serialize options</p>
  <pre>%s</pre>
</div>`, p.config.ContainerClass, html.EscapeString(err.Error()))
		}

		// Create the initialization script for this graph
		// Uses Cal-Heatmap Tooltip plugin for hover information
		// Stores paint function for theme change re-rendering
		initScript := fmt.Sprintf(`
  (function() {
    const graphId = '%s';
    const data = %s;
    const options = {%s};
    const maxValue = %d;
    const displayData = data.map(function(point) {
      const value = point.value || 0;
      if (options.maxValue && value > options.maxValue) {
        return Object.assign({}, point, { value: options.maxValue });
      }
      return point;
    });

    function fitGraph() {
      const inner = document.getElementById(graphId);
      if (!inner) return;

      const outer = inner.parentElement;
      if (!outer) return;

      if (!inner.dataset.baseWidth) {
        inner.dataset.baseWidth = String(inner.scrollWidth || inner.getBoundingClientRect().width || 0);
      }

      const baseWidth = Number(inner.dataset.baseWidth) || inner.scrollWidth || inner.getBoundingClientRect().width || 0;
      const scale = baseWidth > 0 ? Math.min(1, outer.clientWidth / baseWidth) : 1;
      inner.style.zoom = String(scale);
    }

    function paintGraph() {
      // Clear existing graph
      const container = document.getElementById(graphId);
      if (!container) return;
      container.innerHTML = '';
      delete container.dataset.baseWidth;

      // Get theme colors from CSS variables
      const styles = getComputedStyle(document.documentElement);
      const bgColor = styles.getPropertyValue('--color-background').trim();
      const surfaceColor = styles.getPropertyValue('--color-surface').trim();
      const primaryColor = styles.getPropertyValue('--color-primary').trim();

      // Use surface color as base, primary as accent
      const baseColor = surfaceColor || bgColor || '#ebedf0';
      const accentColor = primaryColor || '#216e39';

      const cal = new CalHeatmap();
      cal.paint(
        {
          itemSelector: '#' + graphId,
          data: {
            source: displayData,
            x: 'date',
            y: 'value'
          },
          date: options.date,
          domain: options.domain || { type: 'year' },
          subDomain: options.subDomain || { type: 'day' },
          range: options.range,
          scale: {
            color: {
              type: 'linear',
              range: [baseColor, accentColor],
              domain: [0, maxValue]
            }
          }
        },
        [
          [
            Tooltip,
            {
              text: function (date, value, dayjsDate) {
                const original = data.find(function(point) {
                  return point.date === dayjsDate.format('YYYY-MM-DD');
                });
                const originalValue = original ? (original.value || 0) : (value || 0);
                return (originalValue ? originalValue : 'No') + ' posts on ' + dayjsDate.format('MMM D, YYYY');
              },
            },
          ],
        ]
      );

      fitGraph();
    }

    // Initial paint
    paintGraph();

    // Register for theme changes
    if (!window._contributionGraphPainters) {
      window._contributionGraphPainters = [];
    }
    window._contributionGraphPainters.push(paintGraph);

    if (!window._contributionGraphFitters) {
      window._contributionGraphFitters = [];
    }
    window._contributionGraphFitters.push(fitGraph);
  })();`, graphID, string(dataJSON), p.buildOptionsObject(optionsJSON), maxValue)
		initScripts = append(initScripts, initScript)

		// Return the container div
		return fmt.Sprintf(`<div class="%s">
  <div id="%s"></div>
</div>`, p.config.ContainerClass, graphID)
	})

	// If we found contribution-graph blocks, inject the scripts
	if len(initScripts) > 0 {
		result = p.injectCalHeatmapScripts(result, initScripts)
	}

	post.ArticleHTML = result
	return nil
}

type contributionGraphDataPoint struct {
	Date  string `json:"date"`
	Value int    `json:"value"`
}

func (p *ContributionGraphPlugin) buildContributionData(m *lifecycle.Manager, options map[string]interface{}) interface{} {
	year := time.Now().Year()
	if value, ok := options["year"].(float64); ok {
		year = int(value)
	} else {
		options["year"] = float64(year)
	}

	counts := map[string]int{}
	posts := m.FilterPosts(func(post *models.Post) bool {
		return post != nil && !post.Skip && post.Published && post.Date != nil && post.Date.Year() == year
	})

	for _, post := range posts {
		date := post.Date.Format("2006-01-02")
		counts[date]++
	}

	data := make([]contributionGraphDataPoint, 0, len(counts))
	for date, count := range counts {
		data = append(data, contributionGraphDataPoint{Date: date, Value: count})
	}

	for i := 0; i < len(data)-1; i++ {
		for j := i + 1; j < len(data); j++ {
			if data[i].Date > data[j].Date {
				data[i], data[j] = data[j], data[i]
			}
		}
	}

	return data
}

func (p *ContributionGraphPlugin) computeColorScaleMax(data interface{}, options map[string]interface{}) int {
	values := p.extractContributionValues(data)
	if len(values) == 0 {
		return 1
	}

	if value, ok := options["maxValue"].(float64); ok && value >= 1 {
		return int(value)
	}

	if percentile, ok := p.resolveScaleMaxPercentile(options); ok {
		return percentileValue(values, percentile)
	}

	return values[len(values)-1]
}

func (p *ContributionGraphPlugin) resolveScaleMaxPercentile(options map[string]interface{}) (float64, bool) {
	if value, ok := options["maxPercentile"].(float64); ok && value > 0 {
		return minFloat(value, 100), true
	}
	if p.config.ScaleMaxPercentile > 0 {
		return minFloat(p.config.ScaleMaxPercentile, 100), true
	}
	return 0, false
}

func (p *ContributionGraphPlugin) extractContributionValues(data interface{}) []int {
	encoded, err := json.Marshal(data)
	if err != nil {
		return []int{1}
	}

	var points []contributionGraphDataPoint
	if err := json.Unmarshal(encoded, &points); err != nil {
		return []int{1}
	}

	values := make([]int, 0, len(points))
	for _, point := range points {
		if point.Value > 0 {
			values = append(values, point.Value)
		}
	}
	if len(values) == 0 {
		return []int{1}
	}

	sort.Ints(values)
	return values
}

func percentileValue(values []int, percentile float64) int {
	if len(values) == 0 {
		return 1
	}

	index := int(math.Ceil((percentile/100)*float64(len(values)))) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(values) {
		index = len(values) - 1
	}
	if values[index] < 1 {
		return 1
	}
	return values[index]
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// buildOptionsObject converts options JSON into a JavaScript object literal for Cal-Heatmap.
// This is used by the reactive theme system to store options that can be re-applied.
func (p *ContributionGraphPlugin) buildOptionsObject(optionsJSON []byte) string {
	var options map[string]interface{}
	if err := json.Unmarshal(optionsJSON, &options); err != nil {
		return ""
	}

	var parts []string

	// Handle date.start configuration from "year" option
	if year, ok := options["year"].(float64); ok {
		parts = append(parts, fmt.Sprintf(`date: { start: new Date('%d-01-01') }`, int(year)))
	}

	// Handle domain configuration
	if domain, ok := options["domain"].(string); ok {
		parts = append(parts, fmt.Sprintf(`domain: { type: '%s' }`, domain))
	}

	// Handle subDomain configuration
	if subDomain, ok := options["subDomain"].(string); ok {
		parts = append(parts, fmt.Sprintf(`subDomain: { type: '%s' }`, subDomain))
	}

	// Handle range
	if rangeVal, ok := options["range"].(float64); ok {
		parts = append(parts, fmt.Sprintf(`range: %d`, int(rangeVal)))
	}

	return strings.Join(parts, ", ")
}

// injectCalHeatmapScripts adds the Cal-Heatmap library and initialization scripts to the HTML.
func (p *ContributionGraphPlugin) injectCalHeatmapScripts(htmlContent string, initScripts []string) string {
	containerSelector := "." + p.config.ContainerClass

	// Build the combined script
	// Cal-Heatmap v4 requires d3 as a dependency
	// Tooltip plugin requires popper.js
	script := fmt.Sprintf(`
<style>
%s {
  width: 100%%;
  overflow: hidden;
  margin: 1rem 0;
  display: flex;
  justify-content: center;
}
%s > div {
  flex-shrink: 0;
  transform-origin: top center;
}
#ch-tooltip {
  background: var(--color-surface, #333);
  color: var(--color-text, #fff);
  padding: 0.5rem 0.75rem;
  border-radius: 4px;
  font-size: 0.875rem;
  box-shadow: 0 2px 8px rgba(0,0,0,0.2);
  z-index: 10000 !important;
}
</style>
<link rel="stylesheet" href="%s">
<script src="%s"></script>
<script src="%s"></script>
<script src="%s"></script>
<script src="%s"></script>
<script>
document.addEventListener('DOMContentLoaded', function() {
  // Initialize graphs
  %s

  // Watch for theme/palette changes and re-paint graphs
  const observer = new MutationObserver(function(mutations) {
    mutations.forEach(function(mutation) {
      if (mutation.attributeName === 'data-palette' || mutation.attributeName === 'class') {
        // Small delay to let CSS variables update
        setTimeout(function() {
          if (window._contributionGraphPainters) {
            window._contributionGraphPainters.forEach(function(paint) {
              paint();
            });
          }
        }, 50);
      }
    });
  });

  observer.observe(document.documentElement, { attributes: true });
  observer.observe(document.body, { attributes: true });

  window.addEventListener('resize', function() {
    if (window._contributionGraphFitters) {
      window._contributionGraphFitters.forEach(function(fit) {
        fit();
      });
    }
  });
});
</script>`, containerSelector, containerSelector, p.resolveAssetURL("cal-heatmap-css", p.config.CDNURL+"/cal-heatmap.css"), p.resolveAssetURL("d3", "https://d3js.org/d3.v7.min.js"), p.resolveAssetURL("popper", "https://unpkg.com/@popperjs/core@2"), p.resolveAssetURL("cal-heatmap-js", p.config.CDNURL+"/cal-heatmap.min.js"), p.resolveAssetURL("cal-heatmap-tooltip", p.config.CDNURL+"/plugins/Tooltip.min.js"), strings.Join(initScripts, "\n"))

	// Append the script to the end of the content
	return htmlContent + script
}

func (p *ContributionGraphPlugin) resolveAssetURL(name, fallback string) string {
	return resolveAssetURL(p.assetURLs, name, fallback)
}

// SetConfig sets the contribution graph configuration directly.
// This is useful for testing or programmatic configuration.
func (p *ContributionGraphPlugin) SetConfig(config models.ContributionGraphConfig) {
	p.config = config
}

// Config returns the current contribution graph configuration.
func (p *ContributionGraphPlugin) Config() models.ContributionGraphConfig {
	return p.config
}

// Ensure ContributionGraphPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*ContributionGraphPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*ContributionGraphPlugin)(nil)
	_ lifecycle.RenderPlugin    = (*ContributionGraphPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*ContributionGraphPlugin)(nil)
)
