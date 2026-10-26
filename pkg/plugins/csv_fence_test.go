package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestCSVFencePlugin_Name(t *testing.T) {
	p := NewCSVFencePlugin()
	if p.Name() != "csv_fence" {
		t.Errorf("expected name 'csv_fence', got %q", p.Name())
	}
}

func TestCSVFencePlugin_Priority(t *testing.T) {
	p := NewCSVFencePlugin()

	// Should have late priority for render stage
	if p.Priority(lifecycle.StageRender) != lifecycle.PriorityLate {
		t.Errorf("expected PriorityLate for render stage, got %d", p.Priority(lifecycle.StageRender))
	}

	// Should have default priority for other stages
	if p.Priority(lifecycle.StageTransform) != lifecycle.PriorityDefault {
		t.Errorf("expected PriorityDefault for transform stage, got %d", p.Priority(lifecycle.StageTransform))
	}
}

func TestCSVFencePlugin_Configure_Defaults(t *testing.T) {
	p := NewCSVFencePlugin()
	m := lifecycle.NewManager()

	if err := p.Configure(m); err != nil {
		t.Errorf("Configure returned error: %v", err)
	}

	// Verify defaults
	if !p.enabled {
		t.Error("expected enabled to be true by default")
	}
	if p.tableClass != "csv-table" {
		t.Errorf("expected tableClass 'csv-table', got %q", p.tableClass)
	}
	if !p.hasHeader {
		t.Error("expected hasHeader to be true by default")
	}
	if p.delimiter != ',' {
		t.Errorf("expected delimiter ',', got %c", p.delimiter)
	}
}

func TestCSVFencePlugin_Configure_CustomValues(t *testing.T) {
	p := NewCSVFencePlugin()
	m := lifecycle.NewManager()

	config := m.Config()
	config.Extra = map[string]interface{}{
		"csv_fence": map[string]interface{}{
			"enabled":     false,
			"table_class": "custom-table",
			"has_header":  false,
			"delimiter":   ";",
		},
	}

	if err := p.Configure(m); err != nil {
		t.Errorf("Configure returned error: %v", err)
	}

	if p.enabled {
		t.Error("expected enabled to be false")
	}
	if p.tableClass != "custom-table" {
		t.Errorf("expected tableClass 'custom-table', got %q", p.tableClass)
	}
	if p.hasHeader {
		t.Error("expected hasHeader to be false")
	}
	if p.delimiter != ';' {
		t.Errorf("expected delimiter ';', got %c", p.delimiter)
	}
}

func TestCSVFencePlugin_BasicTable(t *testing.T) {
	p := NewCSVFencePlugin()
	input := `<pre><code class="language-csv">Name,Age,City
Alice,30,New York
Bob,25,Los Angeles
Charlie,35,Chicago</code></pre>`

	post := &models.Post{ArticleHTML: input}
	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Check table structure
	if !strings.Contains(post.ArticleHTML, `<table class="csv-table">`) {
		t.Errorf("expected table with csv-table class, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "<thead>") {
		t.Error("expected thead element")
	}
	if !strings.Contains(post.ArticleHTML, "<tbody>") {
		t.Error("expected tbody element")
	}

	// Check header row
	if !strings.Contains(post.ArticleHTML, "<th>Name</th>") {
		t.Error("expected Name header")
	}
	if !strings.Contains(post.ArticleHTML, "<th>Age</th>") {
		t.Error("expected Age header")
	}
	if !strings.Contains(post.ArticleHTML, "<th>City</th>") {
		t.Error("expected City header")
	}

	// Check data rows
	if !strings.Contains(post.ArticleHTML, "<td>Alice</td>") {
		t.Error("expected Alice cell")
	}
	if !strings.Contains(post.ArticleHTML, "<td>30</td>") {
		t.Error("expected 30 cell")
	}
	if !strings.Contains(post.ArticleHTML, "<td>New York</td>") {
		t.Error("expected New York cell")
	}
}

func TestCSVFencePlugin_NoHeader(t *testing.T) {
	p := NewCSVFencePlugin()
	p.SetHasHeader(false)

	input := `<pre><code class="language-csv">1,2,3
4,5,6</code></pre>`

	post := &models.Post{ArticleHTML: input}
	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Should not have thead
	if strings.Contains(post.ArticleHTML, "<thead>") {
		t.Error("expected no thead element when hasHeader is false")
	}

	// Should have all rows in tbody
	if !strings.Contains(post.ArticleHTML, "<td>1</td>") {
		t.Error("expected first row in tbody")
	}
	if !strings.Contains(post.ArticleHTML, "<td>4</td>") {
		t.Error("expected second row in tbody")
	}
}

func TestCSVFencePlugin_CustomDelimiter(t *testing.T) {
	p := NewCSVFencePlugin()
	p.SetDelimiter(';')

	input := `<pre><code class="language-csv">A;B;C
1;2;3</code></pre>`

	post := &models.Post{ArticleHTML: input}
	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Check header cells
	if !strings.Contains(post.ArticleHTML, "<th>A</th>") {
		t.Error("expected A header")
	}
	if !strings.Contains(post.ArticleHTML, "<th>B</th>") {
		t.Error("expected B header")
	}

	// Check data cells
	if !strings.Contains(post.ArticleHTML, "<td>1</td>") {
		t.Error("expected 1 cell")
	}
	if !strings.Contains(post.ArticleHTML, "<td>2</td>") {
		t.Error("expected 2 cell")
	}
}

func TestCSVFencePlugin_CustomTableClass(t *testing.T) {
	p := NewCSVFencePlugin()
	p.SetTableClass("my-custom-class")

	input := `<pre><code class="language-csv">A,B
1,2</code></pre>`

	post := &models.Post{ArticleHTML: input}
	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<table class="my-custom-class">`) {
		t.Errorf("expected table with my-custom-class, got %q", post.ArticleHTML)
	}
}

func TestCSVFencePlugin_PerBlockOptions(t *testing.T) {
	p := NewCSVFencePlugin()

	// Per-block option should override global settings
	input := `<pre><code class="language-csv delimiter=&quot;;&quot; has_header=&quot;false&quot;">1;2;3
4;5;6</code></pre>`

	post := &models.Post{ArticleHTML: input}
	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Should not have thead (per-block has_header=false)
	if strings.Contains(post.ArticleHTML, "<thead>") {
		t.Error("expected no thead when per-block has_header=false")
	}

	// Should have correct cells (using semicolon delimiter)
	if !strings.Contains(post.ArticleHTML, "<td>1</td>") {
		t.Error("expected 1 cell")
	}
	if !strings.Contains(post.ArticleHTML, "<td>4</td>") {
		t.Error("expected 4 cell")
	}
}

func TestCSVFencePlugin_SkipPost(t *testing.T) {
	p := NewCSVFencePlugin()
	input := `<pre><code class="language-csv">A,B
1,2</code></pre>`

	post := &models.Post{
		ArticleHTML: input,
		Skip:        true,
	}

	originalHTML := post.ArticleHTML
	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if post.ArticleHTML != originalHTML {
		t.Error("expected no changes for skipped post")
	}
}

func TestCSVFencePlugin_EmptyContent(t *testing.T) {
	p := NewCSVFencePlugin()
	post := &models.Post{ArticleHTML: ""}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if post.ArticleHTML != "" {
		t.Error("expected empty ArticleHTML to remain empty")
	}
}

func TestCSVFencePlugin_EmptyCSV(t *testing.T) {
	p := NewCSVFencePlugin()
	input := `<pre><code class="language-csv"></code></pre>`

	post := &models.Post{ArticleHTML: input}
	originalHTML := post.ArticleHTML

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Empty CSV should return original block
	if post.ArticleHTML != originalHTML {
		t.Errorf("expected original block for empty CSV, got %q", post.ArticleHTML)
	}
}

func TestCSVFencePlugin_MultipleBlocks(t *testing.T) {
	p := NewCSVFencePlugin()
	input := `<p>First table:</p>
<pre><code class="language-csv">A,B
1,2</code></pre>
<p>Second table:</p>
<pre><code class="language-csv">C,D
3,4</code></pre>`

	post := &models.Post{ArticleHTML: input}
	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Both blocks should be converted
	if strings.Count(post.ArticleHTML, `<table class="csv-table">`) != 2 {
		t.Errorf("expected 2 tables, got %d", strings.Count(post.ArticleHTML, `<table class="csv-table">`))
	}

	// Check that surrounding content is preserved
	if !strings.Contains(post.ArticleHTML, "<p>First table:</p>") {
		t.Error("expected first paragraph to be preserved")
	}
	if !strings.Contains(post.ArticleHTML, "<p>Second table:</p>") {
		t.Error("expected second paragraph to be preserved")
	}
}

func TestCSVFencePlugin_HTMLEscaping(t *testing.T) {
	p := NewCSVFencePlugin()
	input := `<pre><code class="language-csv">Name,Description
Test,&lt;script&gt;alert('xss')&lt;/script&gt;</code></pre>`

	post := &models.Post{ArticleHTML: input}
	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// HTML entities should be properly handled
	// The content may have been unescaped from goldmark, then re-escaped
	// Just check it doesn't contain raw script tags
	if strings.Contains(post.ArticleHTML, "<script>") {
		t.Error("expected script tags to be escaped")
	}
}

func TestCSVFencePlugin_QuotedFields(t *testing.T) {
	p := NewCSVFencePlugin()
	input := `<pre><code class="language-csv">Name,Description
"Smith, John","Has a comma, in text"</code></pre>`

	post := &models.Post{ArticleHTML: input}
	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Quoted fields should be parsed correctly
	if !strings.Contains(post.ArticleHTML, "<td>Smith, John</td>") {
		t.Errorf("expected 'Smith, John' cell, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "<td>Has a comma, in text</td>") {
		t.Errorf("expected 'Has a comma, in text' cell, got %q", post.ArticleHTML)
	}
}

func TestCSVFencePlugin_VariableColumns(t *testing.T) {
	p := NewCSVFencePlugin()
	// CSV with variable number of columns
	input := `<pre><code class="language-csv">A,B,C
1,2
X,Y,Z,W</code></pre>`

	post := &models.Post{ArticleHTML: input}
	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Should handle variable columns gracefully
	if !strings.Contains(post.ArticleHTML, "<th>A</th>") {
		t.Error("expected A header")
	}
	if !strings.Contains(post.ArticleHTML, "<td>1</td>") {
		t.Error("expected 1 cell")
	}
	if !strings.Contains(post.ArticleHTML, "<td>W</td>") {
		t.Error("expected W cell from longer row")
	}
}

func TestCSVFencePlugin_Disabled(t *testing.T) {
	p := NewCSVFencePlugin()
	p.SetEnabled(false)

	m := lifecycle.NewManager()
	posts := []*models.Post{
		{ArticleHTML: `<pre><code class="language-csv">A,B
1,2</code></pre>`},
	}
	m.SetPosts(posts)

	err := p.Render(m)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	// Should not process when disabled
	if strings.Contains(m.Posts()[0].ArticleHTML, "<table") {
		t.Error("expected no table when plugin is disabled")
	}
}

func TestCSVFencePlugin_Render(t *testing.T) {
	p := NewCSVFencePlugin()
	m := lifecycle.NewManager()

	posts := []*models.Post{
		{ArticleHTML: `<pre><code class="language-csv">A,B
1,2</code></pre>`},
		{ArticleHTML: `<pre><code class="language-csv">C,D
3,4</code></pre>`},
		{ArticleHTML: `<pre><code class="language-csv">E,F</code></pre>`, Skip: true},
	}
	m.SetPosts(posts)

	err := p.Render(m)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	resultPosts := m.Posts()

	// First two should be processed
	if !strings.Contains(resultPosts[0].ArticleHTML, "<table") {
		t.Error("Post 1 not processed correctly")
	}
	if !strings.Contains(resultPosts[1].ArticleHTML, "<table") {
		t.Error("Post 2 not processed correctly")
	}

	// Skipped post should not be processed
	if strings.Contains(resultPosts[2].ArticleHTML, "<table") {
		t.Error("Skipped post should not be processed")
	}
}

func TestCSVFencePlugin_NonCSVBlocks(t *testing.T) {
	p := NewCSVFencePlugin()
	input := `<pre><code class="language-python">print("hello")</code></pre>
<pre><code class="language-csv">A,B
1,2</code></pre>
<pre><code class="language-javascript">console.log("world")</code></pre>`

	post := &models.Post{ArticleHTML: input}
	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Only CSV block should be converted
	if strings.Count(post.ArticleHTML, "<table") != 1 {
		t.Errorf("expected exactly 1 table, got %d", strings.Count(post.ArticleHTML, "<table"))
	}

	// Other code blocks should be preserved
	if !strings.Contains(post.ArticleHTML, `class="language-python"`) {
		t.Error("expected python code block to be preserved")
	}
	if !strings.Contains(post.ArticleHTML, `class="language-javascript"`) {
		t.Error("expected javascript code block to be preserved")
	}
}

func TestCSVFencePlugin_WhitespaceHandling(t *testing.T) {
	p := NewCSVFencePlugin()
	// CSV with extra whitespace
	input := `<pre><code class="language-csv">
Name,   Age,  City
  Alice  ,  30  ,  New York
</code></pre>`

	post := &models.Post{ArticleHTML: input}
	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Should have proper header (leading space in fields is trimmed)
	if !strings.Contains(post.ArticleHTML, "<th>Name</th>") {
		t.Error("expected Name header")
	}
}

func TestCSVFencePlugin_HeaderOnlyCSV(t *testing.T) {
	p := NewCSVFencePlugin()
	// CSV with only header row
	input := `<pre><code class="language-csv">Name,Age,City</code></pre>`

	post := &models.Post{ArticleHTML: input}
	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Should have thead but no tbody content
	if !strings.Contains(post.ArticleHTML, "<thead>") {
		t.Error("expected thead element")
	}
	if !strings.Contains(post.ArticleHTML, "<th>Name</th>") {
		t.Error("expected Name header")
	}
}

// Interface compliance tests

func TestCSVFencePlugin_Interfaces(_ *testing.T) {
	p := NewCSVFencePlugin()

	// Verify interface compliance
	var _ lifecycle.Plugin = p
	var _ lifecycle.ConfigurePlugin = p
	var _ lifecycle.RenderPlugin = p
	var _ lifecycle.PriorityPlugin = p
}
