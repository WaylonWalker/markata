package plugins

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestHorizontalBarDefaults(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		wantDefault bool
	}{
		{
			name:        "horizontal bar",
			source:      `{"type":"bar","options":{"indexAxis":"y"}}`,
			wantDefault: true,
		},
		{
			name:        "vertical bar",
			source:      `{"type":"bar","options":{"indexAxis":"x"}}`,
			wantDefault: false,
		},
		{
			name:        "line chart",
			source:      `{"type":"line","options":{"indexAxis":"y"}}`,
			wantDefault: false,
		},
		{
			name:        "bar with no options",
			source:      `{"type":"bar"}`,
			wantDefault: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config map[string]interface{}
			if err := json.Unmarshal([]byte(tt.source), &config); err != nil {
				t.Fatal(err)
			}
			applyHorizontalBarDefaults(config)
			options, _ := config["options"].(map[string]interface{})
			datasets, _ := options["datasets"].(map[string]interface{})
			bar, present := datasets["bar"].(map[string]interface{})
			if present != tt.wantDefault {
				t.Fatalf("bar defaults present=%v, want %v", present, tt.wantDefault)
			}
			if !tt.wantDefault {
				return
			}
			if bar["borderRadius"] != 6 || bar["borderSkipped"] != false || bar["maxBarThickness"] != 24 {
				t.Fatalf("unexpected horizontal defaults: %#v", bar)
			}
		})
	}
}

func TestHorizontalBarDefaultsRespectAuthoredOptions(t *testing.T) {
	var config map[string]interface{}
	source := `{
		"type":"bar",
		"data":{"datasets":[{"borderRadius":0,"maxBarThickness":50}]},
		"options":{
			"indexAxis":"y",
			"elements":{"bar":{"borderSkipped":true}},
			"datasets":{"bar":{"borderRadius":0,"maxBarThickness":40}}
		}
	}`
	if err := json.Unmarshal([]byte(source), &config); err != nil {
		t.Fatal(err)
	}
	applyHorizontalBarDefaults(config)
	options := config["options"].(map[string]interface{})
	bar := options["datasets"].(map[string]interface{})["bar"].(map[string]interface{})
	if bar["borderRadius"] != float64(0) || bar["maxBarThickness"] != float64(40) {
		t.Fatalf("authored bar styling overwritten: %#v", bar)
	}
	if _, exists := bar["borderSkipped"]; exists {
		t.Fatalf("authored elements.bar.borderSkipped must take precedence: %#v", bar)
	}
	dataset := config["data"].(map[string]interface{})["datasets"].([]interface{})[0].(map[string]interface{})
	if dataset["borderRadius"] != float64(0) || dataset["maxBarThickness"] != float64(50) {
		t.Fatalf("authored per-dataset styling overwritten: %#v", dataset)
	}
}

func TestHorizontalBarDefaultsPreserveMalformedAuthoredOptions(t *testing.T) {
	config := map[string]interface{}{
		"type": "bar",
		"options": map[string]interface{}{
			"indexAxis": "y",
			"datasets": "custom",
		},
	}
	applyHorizontalBarDefaults(config)
	options := config["options"].(map[string]interface{})
	if options["datasets"] != "custom" {
		t.Fatalf("authored options replaced: %#v", options)
	}
}

func TestChartJSPluginHorizontalDefaultsInRenderedConfig(t *testing.T) {
	post := &models.Post{ArticleHTML: `<pre><code class="language-chartjs">{"type":"bar","data":{"labels":["A","B"],"datasets":[{"data":[3,5]}]},"options":{"indexAxis":"y"}}</code></pre>`}
	if err := NewChartJSPlugin().processPost(post); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"borderRadius":6`, `"borderSkipped":false`, `"maxBarThickness":24`, `<details class="chartjs-data" open>`} {
		if !strings.Contains(post.ArticleHTML, want) {
			t.Errorf("rendered chart missing %q", want)
		}
	}
}
