package csspurge

import (
	"strings"
	"testing"
)

func TestPurgeCSS(t *testing.T) {
	tests := []struct {
		name          string
		css           string
		usedClasses   []string
		usedIDs       []string
		usedElements  []string
		preserve      []string
		wantKeptRules int
		wantRemoved   bool
		checkOutput   func(t *testing.T, output string)
	}{
		{
			name:          "keep used class",
			css:           `.foo { color: red; } .bar { color: blue; }`,
			usedClasses:   []string{"foo"},
			wantKeptRules: 1,
			wantRemoved:   true,
			checkOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, ".foo") {
					t.Error("output should contain .foo")
				}
				if strings.Contains(output, ".bar") {
					t.Error("output should not contain .bar")
				}
			},
		},
		{
			name:          "keep used ID",
			css:           `#header { height: 60px; } #footer { height: 40px; }`,
			usedIDs:       []string{"header"},
			wantKeptRules: 1,
			wantRemoved:   true,
		},
		{
			name:          "keep used element",
			css:           `div { margin: 0; } span { padding: 0; }`,
			usedElements:  []string{"div"},
			wantKeptRules: 1,
			wantRemoved:   true,
		},
		{
			name:          "keep preserved pattern",
			css:           `.js-toggle { display: none; } .unused { color: red; }`,
			usedClasses:   []string{},
			preserve:      []string{"js-*"},
			wantKeptRules: 1,
			wantRemoved:   true,
			checkOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, ".js-toggle") {
					t.Error("output should contain .js-toggle")
				}
			},
		},
		{
			name:          "keep keyframes always",
			css:           `@keyframes spin { from { transform: rotate(0); } to { transform: rotate(360deg); } }`,
			usedClasses:   []string{},
			wantKeptRules: 1,
			wantRemoved:   false,
		},
		{
			name:          "keep font-face always",
			css:           `@font-face { font-family: 'Custom'; src: url('font.woff2'); }`,
			usedClasses:   []string{},
			wantKeptRules: 1,
			wantRemoved:   false,
		},
		{
			name:          "keep import always",
			css:           `@import url('other.css');`,
			usedClasses:   []string{},
			wantKeptRules: 1,
			wantRemoved:   false,
		},
		{
			name:          "purge inside media query",
			css:           `@media (max-width: 600px) { .used { display: block; } .unused { display: none; } }`,
			usedClasses:   []string{"used"},
			wantKeptRules: 1, // The @media rule itself
			wantRemoved:   true,
			checkOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, ".used") {
					t.Error("output should contain .used")
				}
				if strings.Contains(output, ".unused") {
					t.Error("output should not contain .unused")
				}
			},
		},
		{
			name:          "remove empty media query",
			css:           `@media (max-width: 600px) { .unused { display: none; } } .used { color: red; }`,
			usedClasses:   []string{"used"},
			wantKeptRules: 1,
			wantRemoved:   true,
			checkOutput: func(t *testing.T, output string) {
				if strings.Contains(output, "@media") {
					t.Error("output should not contain empty @media")
				}
			},
		},
		{
			name:          "keep universal selector",
			css:           `* { box-sizing: border-box; } .unused { color: red; }`,
			usedClasses:   []string{},
			wantKeptRules: 1,
			wantRemoved:   true,
			checkOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "*") {
					t.Error("output should contain universal selector")
				}
			},
		},
		{
			name:          "comma-separated selectors - one used",
			css:           `h1, h2, .title { font-weight: bold; }`,
			usedElements:  []string{"h1"},
			wantKeptRules: 1,
			wantRemoved:   false,
		},
		{
			name:          "comma-separated selectors - none used",
			css:           `.a, .b, .c { color: red; }`,
			usedClasses:   []string{},
			wantKeptRules: 0,
			wantRemoved:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			used := NewUsedSelectors()
			for _, c := range tt.usedClasses {
				used.Classes[c] = true
			}
			for _, id := range tt.usedIDs {
				used.IDs[id] = true
			}
			for _, elem := range tt.usedElements {
				used.Elements[elem] = true
			}

			opts := PurgeOptions{
				Preserve: tt.preserve,
			}

			output, stats := PurgeCSS(tt.css, used, opts)

			if stats.KeptRules != tt.wantKeptRules {
				t.Errorf("kept %d rules, want %d", stats.KeptRules, tt.wantKeptRules)
			}

			hasRemoved := stats.RemovedRules > 0
			if hasRemoved != tt.wantRemoved {
				t.Errorf("removed = %v, want %v (removed %d rules)", hasRemoved, tt.wantRemoved, stats.RemovedRules)
			}

			if tt.checkOutput != nil {
				tt.checkOutput(t, output)
			}
		})
	}
}

func TestPurgeStats(t *testing.T) {
	used := NewUsedSelectors()
	used.Classes["used"] = true

	css := `.used { color: red; } .unused1 { color: blue; } .unused2 { color: green; }`
	_, stats := PurgeCSS(css, used, PurgeOptions{})

	if stats.TotalRules != 3 {
		t.Errorf("TotalRules = %d, want 3", stats.TotalRules)
	}
	if stats.KeptRules != 1 {
		t.Errorf("KeptRules = %d, want 1", stats.KeptRules)
	}
	if stats.RemovedRules != 2 {
		t.Errorf("RemovedRules = %d, want 2", stats.RemovedRules)
	}
	if stats.OriginalSize == 0 {
		t.Error("OriginalSize should not be 0")
	}
	if stats.PurgedSize >= stats.OriginalSize {
		t.Error("PurgedSize should be less than OriginalSize")
	}

	savings := stats.SavingsPercent()
	if savings <= 0 || savings >= 100 {
		t.Errorf("SavingsPercent() = %f, want between 0 and 100", savings)
	}
}

func TestDefaultPurgePreservesInteractiveSidebarState(t *testing.T) {
	// Generated static HTML ships the drawer closed. The open and expanded
	// classes are added by JS, so they never occur in build-time HTML.
	used := NewUsedSelectors()
	for _, class := range []string{"feed-sidebar", "doc-sidebar", "sidebar-toggle", "page-wrapper"} {
		used.Classes[class] = true
	}

	css := `
.feed-sidebar, .doc-sidebar { visibility: hidden; transform: translateX(-100%); }
.feed-sidebar.sidebar--pinned, .doc-sidebar.sidebar--pinned { visibility: visible; transform: none; }
.feed-sidebar.feed-sidebar--expanded .feed-nav-drawer { max-height: 12rem; }
.sidebar--no-motion > .sidebar-toggle { transition: none; }
.unused-selector { display: none; }
`
	// feed-nav-drawer is statically present as a closed element.
	used.Classes["feed-nav-drawer"] = true

	out, _ := PurgeCSS(css, used, PurgeOptions{Preserve: DefaultPreservePatterns()})
	for _, selector := range []string{
		".feed-sidebar.sidebar--pinned",
		".doc-sidebar.sidebar--pinned",
		".feed-sidebar.feed-sidebar--expanded",
		".sidebar--no-motion",
	} {
		if !strings.Contains(out, selector) {
			t.Errorf("purge removed JS-only sidebar selector %q: %s", selector, out)
		}
	}
	if strings.Contains(out, ".unused-selector") {
		t.Errorf("purge kept unrelated CSS: %s", out)
	}
}

// The article Copy control is created after HTML has been published.
// Purging based exclusively on build-time HTML must retain the CSS rules
// that position it, including the language-label and focus styles.
func TestDefaultPurgePreservesRuntimeCodeChrome(t *testing.T) {
	used := NewUsedSelectors()
	used.Elements["pre"] = true
	used.Elements["code"] = true
	used.Classes["post-content"] = true

	css := `@layer components {
.post-content pre.md-code-enhanced { position: relative; }
.post-content pre.md-code-enhanced .md-code-copy { position: absolute; inset-inline-end: .45rem; }
.post-content pre.md-code-enhanced .md-code-copy:focus-visible { outline: 2px solid blue; }
.post-content pre.md-code-enhanced[data-language]::before { inset-inline-end: 5.4rem; }
.unused-css-test-only { color: red; }
}`
	out, _ := PurgeCSS(css, used, PurgeOptions{Preserve: DefaultPreservePatterns()})
	for _, want := range []string{
		"pre.md-code-enhanced { position: relative; }",
		".md-code-copy { position: absolute;",
		".md-code-copy:focus-visible",
		"[data-language]::before",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("purge removed runtime code chrome rule %q: %s", want, out)
		}
	}
	if strings.Contains(out, ".unused-css-test-only") {
		t.Errorf("purge preserved unrelated unused selector: %s", out)
	}
}

func TestDefaultPurgePreservesEditorialRoadEnhancements(t *testing.T) {
	used := NewUsedSelectors()
	used.Classes["road"] = true
	used.Classes["story"] = true
	used.Classes["reveal"] = true
	css := `.road[data-road="editorial"] .road-entry-icon { color: inherit; }
.road[data-road-ready].reveal .road-entry-revealed { animation: road-entry-reveal .6s; }
.unused-road-test { color: red; }`
	out, _ := PurgeCSS(css, used, PurgeOptions{
		Preserve:           DefaultPreservePatterns(),
		PreserveAttributes: DefaultPreserveAttributes(),
	})
	for _, want := range []string{".road-entry-icon", "[data-road-ready]", ".road-entry-revealed"} {
		if !strings.Contains(out, want) {
			t.Errorf("purge removed editorial road selector %q: %s", want, out)
		}
	}
	if strings.Contains(out, ".unused-road-test") {
		t.Errorf("purge preserved unrelated selector: %s", out)
	}
}

func TestDefaultPreservePatterns(t *testing.T) {
	patterns := DefaultPreservePatterns()
	if len(patterns) == 0 {
		t.Error("DefaultPreservePatterns() should return patterns")
	}

	// Check some expected patterns including theme-related ones
	expected := []string{
		"js-*", "htmx-*", "active", "hidden",
		"dark", "light", "dark-mode", "light-mode",
		"theme-*", "palette-*", // Theme class patterns
		"sidebar--pinned", "sidebar--no-motion", "feed-sidebar--expanded",
		"md-code-*", // Runtime article code-copy enhancement
	}
	patternSet := make(map[string]bool)
	for _, p := range patterns {
		patternSet[p] = true
	}

	for _, e := range expected {
		if !patternSet[e] {
			t.Errorf("expected pattern %q not found", e)
		}
	}
}

func TestMatchesPreservePatterns(t *testing.T) {
	patterns := []string{"js-*", "htmx-*", "active"}

	tests := []struct {
		name string
		want bool
	}{
		{"js-toggle", true},
		{"js-", true},
		{"htmx-request", true},
		{"active", true},
		{"inactive", false},
		{"other", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesPreservePatterns(tt.name, patterns)
			if got != tt.want {
				t.Errorf("matchesPreservePatterns(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestIsSelectorUsed(t *testing.T) {
	used := NewUsedSelectors()
	used.Classes["foo"] = true
	used.IDs["bar"] = true
	used.Elements["div"] = true

	tests := []struct {
		selector string
		want     bool
	}{
		{".foo", true},
		{".bar", false},
		{"#bar", true},
		{"#baz", false},
		{"div", true},
		{"span", false},
		{"div.foo", true},
		{"div.unused", false},
		{":root", true},
		{"::selection", true},
		{":where(*)", true},
		{".foo, .unused", true},       // One of multiple matches
		{".unused1, .unused2", false}, // None match
	}

	for _, tt := range tests {
		t.Run(tt.selector, func(t *testing.T) {
			got := isSelectorUsed(tt.selector, used, PurgeOptions{})
			if got != tt.want {
				t.Errorf("isSelectorUsed(%q) = %v, want %v", tt.selector, got, tt.want)
			}
		})
	}
}

func TestDefaultPreserveAttributes(t *testing.T) {
	attrs := DefaultPreserveAttributes()
	if len(attrs) == 0 {
		t.Error("DefaultPreserveAttributes() should return patterns")
	}

	// Check expected theme attributes
	expected := []string{"data-theme", "data-palette", "data-text-size", "data-mode", "data-color-scheme"}
	attrSet := make(map[string]bool)
	for _, a := range attrs {
		attrSet[a] = true
	}

	for _, e := range expected {
		if !attrSet[e] {
			t.Errorf("expected attribute %q not found", e)
		}
	}
}

// TestThemeSelectorsPreserved verifies that theme-related selectors are preserved
// when css_purge runs, even if they're not present in the static HTML.
// This addresses GitHub issue #710.
func TestThemeSelectorsPreserved(t *testing.T) {
	tests := []struct {
		name            string
		css             string
		usedClasses     []string
		usedElements    []string
		usedAttributes  []string
		preserve        []string
		preserveAttrs   []string
		wantKept        bool
		wantContains    []string
		wantNotContains []string
	}{
		{
			name:          "preserve data-theme attribute selector",
			css:           `[data-theme="dark"] { background: #000; } .unused { color: red; }`,
			usedElements:  []string{"html"},
			preserveAttrs: []string{"data-theme"},
			wantKept:      true,
			wantContains:  []string{"data-theme"},
		},
		{
			name:          "preserve data-palette attribute selector",
			css:           `[data-palette="blue"] .card { border-color: blue; }`,
			usedClasses:   []string{"card"},
			preserveAttrs: []string{"data-palette"},
			wantKept:      true,
			wantContains:  []string{"data-palette"},
		},
		{
			name:         "preserve theme-* class pattern",
			css:          `.theme-dark { background: #1a1a1a; } .theme-light { background: #fff; }`,
			preserve:     []string{"theme-*"},
			wantKept:     true,
			wantContains: []string{"theme-dark", "theme-light"},
		},
		{
			name:         "preserve palette-* class pattern",
			css:          `.palette-ocean { --primary: blue; } .palette-forest { --primary: green; }`,
			preserve:     []string{"palette-*"},
			wantKept:     true,
			wantContains: []string{"palette-ocean", "palette-forest"},
		},
		{
			name:         "preserve .dark and .light classes",
			css:          `.dark { color-scheme: dark; } .light { color-scheme: light; }`,
			preserve:     []string{"dark", "light"},
			wantKept:     true,
			wantContains: []string{".dark", ".light"},
		},
		{
			name:         "default patterns preserve theme classes",
			css:          `.dark { background: black; } .light { background: white; } .dark-mode { color: white; } .light-mode { color: black; }`,
			preserve:     DefaultPreservePatterns(),
			wantKept:     true,
			wantContains: []string{".dark", ".light", "dark-mode", "light-mode"},
		},
		{
			name:            "default preserve attributes preserve data-theme",
			css:             `[data-theme="dark"] body { background: #111; } .unused { margin: 0; }`,
			usedElements:    []string{"body"},
			preserveAttrs:   DefaultPreserveAttributes(),
			wantKept:        true,
			wantContains:    []string{"data-theme"},
			wantNotContains: []string{".unused"},
		},
		{
			name:          "combined theme attribute and class preservation",
			css:           `[data-theme="dark"] .theme-dark { display: block; } [data-theme="light"] .theme-light { display: block; }`,
			preserve:      []string{"theme-*"},
			preserveAttrs: []string{"data-theme"},
			wantKept:      true,
			wantContains:  []string{"data-theme", "theme-dark", "theme-light"},
		},
		{
			name:            "purge non-theme selectors",
			css:             `.unused-class { color: red; } [data-theme] .content { color: blue; }`,
			usedClasses:     []string{"content"},
			preserveAttrs:   []string{"data-theme"},
			wantKept:        true,
			wantContains:    []string{"data-theme", ".content"},
			wantNotContains: []string{".unused-class"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			used := NewUsedSelectors()
			for _, c := range tt.usedClasses {
				used.Classes[c] = true
			}
			for _, elem := range tt.usedElements {
				used.Elements[elem] = true
			}
			for _, attr := range tt.usedAttributes {
				used.Attributes[attr] = true
			}

			opts := PurgeOptions{
				Preserve:           tt.preserve,
				PreserveAttributes: tt.preserveAttrs,
			}

			output, stats := PurgeCSS(tt.css, used, opts)

			if tt.wantKept && stats.KeptRules == 0 {
				t.Errorf("expected rules to be kept, but none were kept")
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("output should contain %q, got: %s", want, output)
				}
			}

			for _, notWant := range tt.wantNotContains {
				if strings.Contains(output, notWant) {
					t.Errorf("output should NOT contain %q, got: %s", notWant, output)
				}
			}
		})
	}
}
