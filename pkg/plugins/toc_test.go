package plugins

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// =============================================================================
// TocPlugin Tests
// =============================================================================

func TestTocPlugin_Name(t *testing.T) {
	p := NewTocPlugin()
	if p.Name() != "toc" {
		t.Errorf("expected name 'toc', got %q", p.Name())
	}
}

func TestTocPlugin_DefaultSettings(t *testing.T) {
	p := NewTocPlugin()

	if p.minLevel != 2 {
		t.Errorf("expected default minLevel 2, got %d", p.minLevel)
	}
	if p.maxLevel != 4 {
		t.Errorf("expected default maxLevel 4, got %d", p.maxLevel)
	}
}

func TestTocPlugin_Configure(t *testing.T) {
	p := NewTocPlugin()
	m := lifecycle.NewManager()
	config := m.Config()
	config.Extra = map[string]interface{}{
		"toc_min_level": 1,
		"toc_max_level": 6,
	}

	err := p.Configure(m)
	if err != nil {
		t.Fatalf("Configure error: %v", err)
	}

	if p.minLevel != 1 {
		t.Errorf("expected minLevel 1 after configuration, got %d", p.minLevel)
	}
	if p.maxLevel != 6 {
		t.Errorf("expected maxLevel 6 after configuration, got %d", p.maxLevel)
	}
}

func TestTocPlugin_ConfigureInvalidLevels(t *testing.T) {
	p := NewTocPlugin()
	m := lifecycle.NewManager()
	config := m.Config()
	config.Extra = map[string]interface{}{
		"toc_min_level": 0,  // Invalid: below 1
		"toc_max_level": 10, // Invalid: above 6
	}

	err := p.Configure(m)
	if err != nil {
		t.Fatalf("Configure error: %v", err)
	}

	// Should keep defaults for invalid values
	if p.minLevel != 2 {
		t.Errorf("expected minLevel to remain at default 2, got %d", p.minLevel)
	}
	if p.maxLevel != 4 {
		t.Errorf("expected maxLevel to remain at default 4, got %d", p.maxLevel)
	}
}

func TestTocPlugin_SetLevelRange(t *testing.T) {
	p := NewTocPlugin()

	p.SetLevelRange(1, 6)
	if p.minLevel != 1 {
		t.Errorf("expected minLevel 1, got %d", p.minLevel)
	}
	if p.maxLevel != 6 {
		t.Errorf("expected maxLevel 6, got %d", p.maxLevel)
	}

	// Test invalid ranges
	p.SetLevelRange(0, 7) // Both invalid
	if p.minLevel != 1 {
		t.Errorf("expected minLevel to remain 1, got %d", p.minLevel)
	}
	if p.maxLevel != 6 {
		t.Errorf("expected maxLevel to remain 6, got %d", p.maxLevel)
	}
}

func TestTocPlugin_GeneratesTOCFromMultipleHeadings(t *testing.T) {
	p := NewTocPlugin()
	p.SetLevelRange(1, 6)

	m := lifecycle.NewManager()
	post := &models.Post{
		Content: `# Heading 1
## Heading 2
### Heading 3
## Another H2`,
		Slug:  "test-post",
		Extra: make(map[string]interface{}),
	}

	m.SetPosts([]*models.Post{post})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	toc, ok := posts[0].Extra["toc"].([]*TocEntry)
	if !ok || toc == nil {
		t.Fatal("expected toc to be set in post Extra")
	}

	// Should have entries
	if len(toc) == 0 {
		t.Error("expected at least one TOC entry")
	}
}

func TestTocPlugin_RespectsMinLevel(t *testing.T) {
	p := NewTocPlugin()
	p.SetLevelRange(2, 4) // Exclude h1

	m := lifecycle.NewManager()
	post := &models.Post{
		Content: `# Heading 1 (should be excluded)
## Heading 2
### Heading 3`,
		Slug:  "test-post",
		Extra: make(map[string]interface{}),
	}

	m.SetPosts([]*models.Post{post})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	toc, ok := posts[0].Extra["toc"].([]*TocEntry)
	if !ok || toc == nil {
		t.Fatal("expected toc to be set in post Extra")
	}

	// Check that h1 is excluded
	for _, entry := range toc {
		if entry.Level == 1 {
			t.Errorf("h1 should be excluded when minLevel is 2, found entry: %+v", entry)
		}
	}
}

func TestTocPlugin_RespectsMaxLevel(t *testing.T) {
	p := NewTocPlugin()
	p.SetLevelRange(2, 3) // Exclude h4+

	m := lifecycle.NewManager()
	post := &models.Post{
		Content: `## Heading 2
### Heading 3
#### Heading 4 (should be excluded)
##### Heading 5 (should be excluded)`,
		Slug:  "test-post",
		Extra: make(map[string]interface{}),
	}

	m.SetPosts([]*models.Post{post})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	toc, ok := posts[0].Extra["toc"].([]*TocEntry)
	if !ok || toc == nil {
		t.Fatal("expected toc to be set in post Extra")
	}

	// Check all levels recursively
	var checkLevels func([]*TocEntry)
	checkLevels = func(entries []*TocEntry) {
		for _, entry := range entries {
			if entry.Level > 3 {
				t.Errorf("h%d should be excluded when maxLevel is 3, found entry: %+v", entry.Level, entry)
			}
			checkLevels(entry.Children)
		}
	}
	checkLevels(toc)
}

func TestTocPlugin_HandlesDuplicateHeadingIDs(t *testing.T) {
	p := NewTocPlugin()
	p.SetLevelRange(2, 4)

	m := lifecycle.NewManager()
	post := &models.Post{
		Content: `## Introduction
## Methods
## Introduction
## Introduction`,
		Slug:  "test-post",
		Extra: make(map[string]interface{}),
	}

	m.SetPosts([]*models.Post{post})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	toc, ok := posts[0].Extra["toc"].([]*TocEntry)
	if !ok || toc == nil {
		t.Fatal("expected toc to be set in post Extra")
	}

	// Collect all IDs
	ids := make(map[string]bool)
	var collectIDs func([]*TocEntry)
	collectIDs = func(entries []*TocEntry) {
		for _, entry := range entries {
			if ids[entry.ID] {
				t.Errorf("duplicate ID found: %s", entry.ID)
			}
			ids[entry.ID] = true
			collectIDs(entry.Children)
		}
	}
	collectIDs(toc)

	// Should have 4 unique IDs
	if len(ids) != 4 {
		t.Errorf("expected 4 unique IDs, got %d", len(ids))
	}
}

func TestTocPlugin_NoHeadings(t *testing.T) {
	p := NewTocPlugin()

	m := lifecycle.NewManager()
	post := &models.Post{
		Content: `Just some paragraph text without any headings.

Another paragraph here.`,
		Slug:  "test-post",
		Extra: make(map[string]interface{}),
	}

	m.SetPosts([]*models.Post{post})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	// TOC should not be set when there are no headings
	if _, ok := posts[0].Extra["toc"]; ok {
		t.Error("expected no toc when there are no headings")
	}
}

func TestTocPlugin_SkippedPost(t *testing.T) {
	p := NewTocPlugin()

	m := lifecycle.NewManager()
	post := &models.Post{
		Content: `## Heading 1
### Heading 2`,
		Slug:  "test-post",
		Skip:  true,
		Extra: make(map[string]interface{}),
	}

	m.SetPosts([]*models.Post{post})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	// TOC should not be set for skipped posts
	if _, ok := posts[0].Extra["toc"]; ok {
		t.Error("expected no toc for skipped post")
	}
}

func TestTocPlugin_EmptyContent(t *testing.T) {
	p := NewTocPlugin()

	m := lifecycle.NewManager()
	post := &models.Post{
		Content: "",
		Slug:    "test-post",
		Extra:   make(map[string]interface{}),
	}

	m.SetPosts([]*models.Post{post})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	if _, ok := posts[0].Extra["toc"]; ok {
		t.Error("expected no toc for empty content")
	}
}

func TestTocPlugin_HierarchicalStructure(t *testing.T) {
	p := NewTocPlugin()
	p.SetLevelRange(2, 4)

	m := lifecycle.NewManager()
	post := &models.Post{
		Content: `## Section 1
### Subsection 1.1
#### Sub-subsection 1.1.1
### Subsection 1.2
## Section 2
### Subsection 2.1`,
		Slug:  "test-post",
		Extra: make(map[string]interface{}),
	}

	m.SetPosts([]*models.Post{post})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	toc, ok := posts[0].Extra["toc"].([]*TocEntry)
	if !ok || toc == nil {
		t.Fatal("expected toc to be set in post Extra")
	}

	// Should have 2 root entries (Section 1 and Section 2)
	if len(toc) != 2 {
		t.Errorf("expected 2 root entries, got %d", len(toc))
	}

	// Section 1 should have 2 children
	if len(toc[0].Children) != 2 {
		t.Errorf("expected Section 1 to have 2 children, got %d", len(toc[0].Children))
	}

	// Subsection 1.1 should have 1 child
	if len(toc[0].Children[0].Children) != 1 {
		t.Errorf("expected Subsection 1.1 to have 1 child, got %d", len(toc[0].Children[0].Children))
	}

	// Section 2 should have 1 child
	if len(toc[1].Children) != 1 {
		t.Errorf("expected Section 2 to have 1 child, got %d", len(toc[1].Children))
	}
}

func TestTocPlugin_TocEntry(t *testing.T) {
	entry := &TocEntry{
		Level:    2,
		Text:     "Test Heading",
		ID:       "test-heading",
		Children: make([]*TocEntry, 0),
	}

	if entry.Level != 2 {
		t.Errorf("expected level 2, got %d", entry.Level)
	}
	if entry.Text != "Test Heading" {
		t.Errorf("expected text 'Test Heading', got %q", entry.Text)
	}
	if entry.ID != "test-heading" {
		t.Errorf("expected id 'test-heading', got %q", entry.ID)
	}
}

func TestTocPlugin_GenerateID(t *testing.T) {
	p := NewTocPlugin()

	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{"simple", "Hello World", "hello-world"},
		{"with special chars", "Hello, World!", "hello-world"},
		{"with numbers", "Section 123", "section-123"},
		{"uppercase", "UPPERCASE TEXT", "uppercase-text"},
		{"multiple spaces", "Hello   World", "hello-world"},
		{"leading/trailing spaces", "  Hello World  ", "hello-world"},
		{"hyphens preserved", "hello-world", "hello-world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idCounts := make(map[string]int)
			got := p.generateID(tt.text, idCounts)
			if got != tt.expected {
				t.Errorf("generateID(%q) = %q, want %q", tt.text, got, tt.expected)
			}
		})
	}
}

func TestTocPlugin_GenerateIDDuplicates(t *testing.T) {
	p := NewTocPlugin()
	idCounts := make(map[string]int)

	// First occurrence
	id1 := p.generateID("Introduction", idCounts)
	if id1 != "introduction" {
		t.Errorf("first ID should be 'introduction', got %q", id1)
	}

	// Second occurrence - should be different
	id2 := p.generateID("Introduction", idCounts)
	if id2 == id1 {
		t.Errorf("second ID should be different from first, both are %q", id1)
	}

	// Third occurrence - should be different from both
	id3 := p.generateID("Introduction", idCounts)
	if id3 == id1 || id3 == id2 {
		t.Errorf("third ID should be different from first two, got %q", id3)
	}
}

func TestTocPlugin_HeadingRegex(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantMatch bool
		wantLevel int
		wantText  string
	}{
		{"h1", "# Heading 1", true, 1, "Heading 1"},
		{"h2", "## Heading 2", true, 2, "Heading 2"},
		{"h3", "### Heading 3", true, 3, "Heading 3"},
		{"h4", "#### Heading 4", true, 4, "Heading 4"},
		{"h5", "##### Heading 5", true, 5, "Heading 5"},
		{"h6", "###### Heading 6", true, 6, "Heading 6"},
		{"h1 with trailing hashes", "# Heading #", true, 1, "Heading"},
		{"h2 with multiple trailing hashes", "## Heading ##", true, 2, "Heading"},
		{"not a heading", "Just text", false, 0, ""},
		{"no space after hash", "#NoSpace", false, 0, ""},
		{"too many hashes", "####### Too Many", false, 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := headingRegex.FindStringSubmatch(tt.input)
			gotMatch := matches != nil

			if gotMatch != tt.wantMatch {
				t.Errorf("regex match = %v, want %v", gotMatch, tt.wantMatch)
				return
			}

			if !tt.wantMatch {
				return
			}

			if len(matches) < 3 {
				t.Fatal("expected at least 3 groups in match")
			}

			gotLevel := len(matches[1])
			if gotLevel != tt.wantLevel {
				t.Errorf("level = %d, want %d", gotLevel, tt.wantLevel)
			}

			gotText := matches[2]
			// Note: the regex may include trailing spaces, so we compare trimmed
			if gotText != tt.wantText {
				t.Errorf("text = %q, want %q", gotText, tt.wantText)
			}
		})
	}
}

func TestTocPlugin_BuildHierarchy(t *testing.T) {
	p := NewTocPlugin()

	headings := []*TocEntry{
		{Level: 2, Text: "H2-1", ID: "h2-1", Children: make([]*TocEntry, 0)},
		{Level: 3, Text: "H3-1", ID: "h3-1", Children: make([]*TocEntry, 0)},
		{Level: 3, Text: "H3-2", ID: "h3-2", Children: make([]*TocEntry, 0)},
		{Level: 2, Text: "H2-2", ID: "h2-2", Children: make([]*TocEntry, 0)},
		{Level: 3, Text: "H3-3", ID: "h3-3", Children: make([]*TocEntry, 0)},
	}

	result := p.buildHierarchy(headings)

	// Should have 2 root entries
	if len(result) != 2 {
		t.Fatalf("expected 2 root entries, got %d", len(result))
	}

	// First root should have 2 children
	if len(result[0].Children) != 2 {
		t.Errorf("expected first root to have 2 children, got %d", len(result[0].Children))
	}

	// Second root should have 1 child
	if len(result[1].Children) != 1 {
		t.Errorf("expected second root to have 1 child, got %d", len(result[1].Children))
	}
}

func TestTocPlugin_ExtractTOC(t *testing.T) {
	p := NewTocPlugin()
	p.SetLevelRange(1, 6)

	content := `# Heading 1
## Heading 2
### Heading 3
## Another H2
# Another H1`

	toc := p.extractTOC(content)

	if len(toc) == 0 {
		t.Fatal("expected TOC entries")
	}

	// Verify structure
	if toc[0].Text != "Heading 1" {
		t.Errorf("expected first entry to be 'Heading 1', got %q", toc[0].Text)
	}
	if toc[0].Level != 1 {
		t.Errorf("expected first entry level to be 1, got %d", toc[0].Level)
	}
}

// Interface compliance tests
func TestTocPlugin_Interfaces(_ *testing.T) {
	p := NewTocPlugin()

	var _ lifecycle.Plugin = p
	var _ lifecycle.ConfigurePlugin = p
	var _ lifecycle.TransformPlugin = p
}

func TestCleanHeadingText(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain", "Plain heading", "Plain heading"},
		{"resolved wikilink html", `Posts tagged: <a href="/tags/launchdoc/" class="wikilink" data-title="Posts tagged: launchdoc">launchdoc</a>`, "Posts tagged: launchdoc"},
		{"markdown link", "Recent [TIL](/til/)", "Recent TIL"},
		{"wikilink alias", "See [[target|Alias Text]]", "See Alias Text"},
		{"wikilink plain", "See [[ daily-notes ]]", "See daily-notes"},
		{"inline markup", "**Bold** and `code` and ==mark==", "Bold and code and mark"},
		{"heading attrs", "Title {#custom-id .cls}", "Title"},
		{"entities", "Fish &amp; Chips", "Fish & Chips"},
		{"code span keeps underscores", "`update_meta`", "update_meta"},
		{"code span keeps markup", "`**not bold**` option", "**not bold** option"},
		{"intraword underscore", "snake_case and _emphasis_", "snake_case and emphasis"},
		{"unpaired underscore", "_headers File", "_headers File"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanHeadingText(tt.in); got != tt.want {
				t.Errorf("cleanHeadingText(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestTocPlugin_IDsMatchRenderedHeadings(t *testing.T) {
	content := "# Intro\n\n## Intro\n\n```toml\n## Not a heading\n# comment\n```\n\n" +
		"### `update_meta`\n\n### \"Quoted title\"\n\n## Snake_case and _emphasis_\n\n" +
		"## Custom {#my-id}\n\n## Intro\n\n### _headers File\n\n## Intro\n\n~~~\n## Also not\n~~~\n\n" +
		"## Pinned {#intro-9}\n\n## Café {#café}\n\n## Scoped {#a:b.c .wide}\n"

	p := NewTocPlugin()
	p.SetLevelRange(2, 4)
	var got []string
	var walk func([]*TocEntry)
	walk = func(entries []*TocEntry) {
		for _, e := range entries {
			got = append(got, e.ID)
			walk(e.Children)
		}
	}
	walk(p.extractTOC(content))

	md := createMarkdownRenderer("github", false, MarkdownExtensionConfig{TypographerEnabled: true})
	var buf bytes.Buffer
	if err := md.Convert([]byte(content), &buf); err != nil {
		t.Fatalf("convert: %v", err)
	}
	matches := regexp.MustCompile(`<h[234] id="([^"]+)"`).FindAllStringSubmatch(buf.String(), -1)
	want := make([]string, 0, len(matches))
	for _, m := range matches {
		want = append(want, m[1])
	}

	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("TOC ids %v do not match rendered heading ids %v", got, want)
	}
	for _, id := range []string{"intro-9", "café", "a:b.c"} {
		if !strings.Contains(strings.Join(want, ","), id) {
			t.Errorf("explicit id %q missing from rendered ids %v", id, want)
		}
	}
}
