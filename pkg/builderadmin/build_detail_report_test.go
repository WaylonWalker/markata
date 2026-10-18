package builderadmin

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/servecontrol"
)

func TestBuildDetailIncludesPortableDiagnosticHandoff(t *testing.T) {
	t.Parallel()

	build := BuildRecord{
		ID:            "build-123",
		Status:        "failed",
		TriggerType:   "file-watch",
		TriggerDetail: "Debounced file-watch build (2 paths)",
		ChangedPaths:  []string{"posts/one.md", "templates/feed.html"},
		QueueWaitMS:   1250,
		PrepareMS:     2500,
		BuildMS:       810400,
		PromoteMS:     3100,
		PruneMS:       900,
		TotalMS:       818150,
		PerfSummary: []string{
			"Duration: 13m30.4s",
			"Hotspots: write/publish_feeds 184.2s",
		},
		LogDiagnostics: []servecontrol.Diagnostic{{
			Severity:     "warning",
			Code:         "MARKATA-W123",
			Message:      "example diagnostic",
			File:         "posts/one.md",
			Line:         7,
			SuggestedFix: "fix the example",
		}},
		Error:      "command failed with exit code 1",
		FinishedAt: time.Now().UTC(),
	}

	tmpl := template.Must(template.New("build-detail").Funcs(template.FuncMap{
		"msToSeconds": func(ms int64) string { return fmt.Sprintf("%.2fs", float64(ms)/1000) },
		"since":       func(time.Time) string { return "now" },
		"statusClass": func(string) string { return "" },
		"browserTokens": func(uiTheme) template.CSS {
			return template.CSS("")
		},
	}).Parse(buildDetailHTML))

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, struct {
		BuildRecord
		PreviewURL string
		Theme      uiTheme
	}{BuildRecord: build}); err != nil {
		t.Fatalf("render build detail: %v", err)
	}
	body := rendered.String()

	for _, want := range []string{
		"Copy diagnostic report",
		"Download report",
		"Markata-Go Builder Diagnostic Report",
		"build-123",
		"810.40s",
		"posts/one.md",
		"Hotspots: write/publish_feeds 184.2s",
		"MARKATA-W123",
		"fix the example",
		"command failed with exit code 1",
		"Raw logs, environment values, credentials, and absolute release paths are not included",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("build detail missing %q", want)
		}
	}
}
