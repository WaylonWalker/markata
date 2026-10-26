package themes

import (
	"io/fs"
	"regexp"
	"strings"
	"testing"
)

var (
	cssCommentRe = regexp.MustCompile(`(?s)/\*.*?\*/`)
	cssStringRe  = regexp.MustCompile(`"[^"\n]*"|'[^'\n]*'`)
)

// stripCSSNoise removes comments and string literals so brace counting only
// sees structural braces.
func stripCSSNoise(css string) string {
	css = cssCommentRe.ReplaceAllString(css, "")
	return cssStringRe.ReplaceAllString(css, `""`)
}

// TestCSSLayerBlocksStayBalanced guards against a stray brace closing a
// top-level @layer block early. When that happens every rule after the stray
// brace becomes unlayered and outranks all layered rules across the theme,
// which silently breaks the cascade for unrelated components.
func TestCSSLayerBlocksStayBalanced(t *testing.T) {
	static := DefaultStatic()
	err := fs.WalkDir(static, "css", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(p, ".css") {
			return nil
		}
		raw, readErr := fs.ReadFile(static, p)
		if readErr != nil {
			return readErr
		}
		css := stripCSSNoise(string(raw))

		t.Run(p, func(t *testing.T) {
			depth := 0
			line := 1
			topLevelLayerOpenedAt := 0
			for _, r := range css {
				switch r {
				case '\n':
					line++
				case '{':
					depth++
				case '}':
					depth--
					if depth < 0 {
						t.Fatalf("%s:%d: unmatched closing brace", p, line)
					}
					if depth == 0 && topLevelLayerOpenedAt > 0 {
						topLevelLayerOpenedAt = 0
					}
				}
			}
			if depth != 0 {
				t.Fatalf("%s: unbalanced braces (depth %d at EOF)", p, depth)
			}

			// Any file that opens a top-level `@layer name {` block must keep
			// all of its rules inside layer blocks.
			if !regexp.MustCompile(`(?m)^@layer\s+[a-z-]+\s*\{`).MatchString(css) {
				return
			}
			depth = 0
			line = 1
			inLayerBlock := false
			for i := 0; i < len(css); i++ {
				c := css[i]
				switch c {
				case '\n':
					line++
				case '{':
					if depth == 0 {
						stmtStart := strings.LastIndexAny(css[:i], "};")
						stmt := strings.TrimSpace(css[stmtStart+1 : i])
						inLayerBlock = strings.HasPrefix(stmt, "@layer")
						if !inLayerBlock {
							t.Errorf("%s:%d: rule %q is outside any @layer block", p, line, firstLine(stmt))
						}
					}
					depth++
				case '}':
					depth--
				}
			}
		})
		return nil
	})
	if err != nil {
		t.Fatalf("walk css: %v", err)
	}
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return strings.TrimSpace(s[:i])
	}
	return s
}

// TestPhotoGridUsesAvailableLayoutWidth guards the full-bleed Shots layout.
// Sidebars may exist in the DOM while collapsed, so sizing must not depend on
// :not(:has(.feed-sidebar)). Pinned drawers still reserve usable space.
func TestPhotoGridUsesAvailableLayoutWidth(t *testing.T) {
	raw, err := fs.ReadFile(DefaultStatic(), "css/feeds.css")
	if err != nil {
		t.Fatal(err)
	}
	css := strings.ReplaceAll(string(raw), "\r\n", "\n")

	for _, want := range []string{
		".page-wrapper.page-wrapper--feed-only:has(> .main-content > .feed.h-feed.feed--photo-grid) > .main-content {",
		"width: max(0px, calc(100% - var(--push-left, 0px) - var(--push-right, 0px)));",
		"margin-left: var(--push-left, 0px);",
		".main-content > .feed.h-feed.feed--photo-grid {\n  width: 100%;\n  max-width: 100%;",
		"container-name: shots-feed;",
		"container-type: inline-size;",
		"@container shots-feed (min-width: 700px)",
		"@container shots-feed (min-width: 1200px)",
		"@container shots-feed (min-width: 1500px)",
		"@container shots-feed (min-width: 1800px)",
	} {
		if !strings.Contains(css, want) {
			t.Errorf("photo grid sizing missing %q", want)
		}
	}

	// The older selector failed as soon as a collapsed/hidden sidebar existed
	// in the page. Keep the new width rule independent of the DOM sidebar count.
	legacy := ".page-wrapper:has(> .main-content > .feed.h-feed.feed--photo-grid):not(:has("
	if strings.Contains(css, legacy) {
		t.Error("photo-grid main-content sizing still depends on the absence of sidebar elements")
	}

	// Both classes appear together on photo-grid feeds that enable the Calendar
	// peer view. The eager calendar sizing rule must exclude photo grids, and
	// the deferred stylesheet must not reintroduce a competing wrapper width.
	calendarCSS, err := fs.ReadFile(DefaultStatic(), "css/calendar-feed.css")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(css, "calendar-feed:not(.feed--photo-grid)") {
		t.Error("calendar feed width rule must exclude photo grids")
	}
	if strings.Contains(string(calendarCSS), ".page-wrapper.page-wrapper--feed-only:has(") {
		t.Error("deferred calendar CSS must not change the page wrapper width")
	}
}
