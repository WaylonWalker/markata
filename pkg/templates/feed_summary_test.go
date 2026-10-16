package templates

import (
	"strings"
	"testing"
	"time"
)

func TestFeedSummary_BoundedTextWithoutMedia(t *testing.T) {
	for _, tree := range []string{"../../templates", "../themes/default/templates"} {
		t.Run(tree, func(t *testing.T) {
			engine, err := NewEngine(tree)
			if err != nil {
				t.Fatal(err)
			}
			ctx := map[string]interface{}{
				"variant": "summary",
				"posts": []map[string]interface{}{{
					"title": "A useful tip", "href": "/useful-tip/",
					"date":         time.Date(2026, 10, 15, 0, 0, 0, 0, time.UTC),
					"article_html": "<p>" + strings.Repeat("Useful context. ", 100) + "</p><img src=\"/large.jpg\"><video src=\"/large.mp4\"></video><pre><code>command</code></pre>",
				}},
			}
			output, err := engine.RenderToString("partials/feed_preview.html", ctx)
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range []string{`href="/useful-tip/"`, "A useful tip", "Useful context.", "Oct 15, 2026"} {
				if !strings.Contains(output, want) {
					t.Errorf("summary missing %q: %s", want, output)
				}
			}
			for _, unwanted := range []string{"<img", "<video", "<pre", "<code", strings.Repeat("Useful context. ", 30)} {
				if strings.Contains(output, unwanted) {
					t.Errorf("summary contains full article content: %s", unwanted)
				}
			}
		})
	}
}
