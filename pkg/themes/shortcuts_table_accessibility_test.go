package themes

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestShortcutsTablesHaveAccessibleNamesAndHeaders(t *testing.T) {
	root, err := os.ReadFile("../../templates/partials/shortcuts-modal.html")
	if err != nil {
		t.Fatal(err)
	}
	bundled, err := ReadTemplate("partials/shortcuts-modal.html")
	if err != nil {
		t.Fatal(err)
	}

	tablePattern := regexp.MustCompile(`(?s)<table class="shortcuts-table">(.*?)</table>`)
	sections := []string{"Scrolling", "Navigation", "Theme", "Utility", "Custom Navigation"}
	for name, source := range map[string]string{
		"root":    string(root),
		"bundled": string(bundled),
	} {
		t.Run(name, func(t *testing.T) {
			tables := tablePattern.FindAllStringSubmatch(source, -1)
			if len(tables) != len(sections) {
				t.Fatalf("expected %d shortcut tables, got %d", len(sections), len(tables))
			}
			for i, match := range tables {
				table := match[1]
				caption := "<caption class=\"visually-hidden\">" + sections[i] + " keyboard shortcuts</caption>"
				if !strings.Contains(table, caption) {
					t.Errorf("%s table is missing its accessible caption", sections[i])
				}
				if !strings.Contains(table, `<th scope="col">Keys</th>`) ||
					!strings.Contains(table, `<th scope="col">Action</th>`) {
					t.Errorf("%s table is missing column headers", sections[i])
				}
			}
		})
	}
}
