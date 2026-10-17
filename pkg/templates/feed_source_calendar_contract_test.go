package templates

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Source templates and packaged default-theme templates must use the same
// on-demand calendar metadata. Embedding all posts in a hidden <ul> turns a
// small paginated feed into a multi-thousand-node document.
func TestFeedTemplatesDeferCalendarHistory(t *testing.T) {
	for _, root := range []string{"../../templates", "../themes/default/templates"} {
		for _, file := range []string{"feed.html", "feed-photo-grid.html"} {
			t.Run(filepath.Join(root, file), func(t *testing.T) {
				data, err := os.ReadFile(filepath.Join(root, file))
				if err != nil {
					t.Fatal(err)
				}
				html := string(data)
				for _, required := range []string{
					`data-calendar-src="`,
					`calendar.json`,
					`data-calendar-years`,
					`data-calendar-error`,
				} {
					if !strings.Contains(html, required) {
						t.Errorf("%s must include %s", file, required)
					}
				}
				for _, forbidden := range []string{
					`data-calendar-list`,
					`data-calendar-post`,
				} {
					if strings.Contains(html, forbidden) {
						t.Errorf("%s embeds full-history calendar rows (%s)", file, forbidden)
					}
				}
			})
		}
	}
}
