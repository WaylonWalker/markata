package agentskill

import (
	"encoding/json"
	"io/fs"
	"strings"
	"testing"
)

// The skill's example must contain renderable Chart.js JSON, not chat-UI-only
// chart syntax or an illustrative Markdown table standing in for a chart.
func TestSiteSkill_ChartCompositionExamples(t *testing.T) {
	skillFS, err := SiteSkill()
	if err != nil {
		t.Fatal(err)
	}
	data, err := fs.ReadFile(skillFS, "examples/chart-composition.md")
	if err != nil {
		t.Fatal(err)
	}
	src := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(src, "\n")
	inChart := false
	start := 0
	var source strings.Builder
	types := make(map[string]int)

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "```chartjs" {
			if inChart {
				t.Fatalf("nested chartjs fence on line %d", i+1)
			}
			inChart = true
			start = i + 1
			source.Reset()
			continue
		}
		if !inChart {
			continue
		}
		if line != "```" {
			source.WriteString(line)
			source.WriteByte('\n')
			continue
		}

		var config struct {
			Type string `json:"type"`
			Data struct {
				Labels   []string `json:"labels"`
				Datasets []struct {
					Label string    `json:"label"`
					Data  []float64 `json:"data"`
				} `json:"datasets"`
			} `json:"data"`
		}
		if err := json.Unmarshal([]byte(source.String()), &config); err != nil {
			t.Fatalf("chart beginning on line %d has invalid JSON: %v", start, err)
		}
		if config.Type == "" || len(config.Data.Labels) == 0 || len(config.Data.Datasets) == 0 {
			t.Fatalf("chart beginning on line %d is missing type, labels, or datasets", start)
		}
		for _, dataset := range config.Data.Datasets {
			if dataset.Label == "" || len(dataset.Data) != len(config.Data.Labels) {
				t.Fatalf("chart on line %d cannot generate the accessible data table", start)
			}
		}
		types[config.Type]++
		inChart = false
	}
	if inChart {
		t.Fatal("unclosed chartjs fence in example")
	}
	for _, kind := range []string{"bar", "line", "doughnut"} {
		if types[kind] != 1 {
			t.Errorf("expected one %s chart, got %d", kind, types[kind])
		}
	}
	if !strings.Contains(src, "Illustrative") {
		t.Error("example must mark fabricated values as illustrative")
	}
}
