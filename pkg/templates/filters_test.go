package templates

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"

	"github.com/flosch/pongo2/v6"
)

func TestFilterRSSDate(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	date := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	post := &models.Post{Date: &date}
	ctx := NewContext(post, "", nil)

	result, err := engine.RenderString("{{ post.date | rss_date }}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}

	expected := "Mon, 15 Jan 2024 10:30:00 +0000"
	if result != expected {
		t.Errorf("rss_date: got %q, want %q", result, expected)
	}
}

func TestFilterAtomDate(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	date := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	post := &models.Post{Date: &date}
	ctx := NewContext(post, "", nil)

	result, err := engine.RenderString("{{ post.date | atom_date }}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}

	expected := "2024-01-15T10:30:00Z"
	if result != expected {
		t.Errorf("atom_date: got %q, want %q", result, expected)
	}
}

func TestFilterDateFormat(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	post := &models.Post{Date: &date}
	ctx := NewContext(post, "", nil)

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "default format",
			template: "{{ post.date | date_format }}",
			expected: "2024-01-15",
		},
		{
			name:     "custom format",
			template: "{{ post.date | date_format:\"January 2, 2006\" }}",
			expected: "January 15, 2024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.RenderString(tt.template, ctx)
			if err != nil {
				t.Fatalf("RenderString() error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFilterHumanDate(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	date := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	post := &models.Post{Date: &date}
	ctx := NewContext(post, "", nil)

	result, err := engine.RenderString("{{ post.date | human_date }}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}

	if result != "Jan 15, 2024" {
		t.Fatalf("human_date: got %q, want %q", result, "Jan 15, 2024")
	}
}

func TestTemplateTrees_UseHumanDateForVisibleHTMLDates(t *testing.T) {
	files := []string{
		"../../templates/components/post_byline.html",
		"../../templates/partials/card.html",
		"../../templates/partials/embed-card.html",
		"../../templates/partials/feed_preview.html",
		"../../templates/partials/cards/article-card.html",
		"../../templates/partials/cards/default-card.html",
		"../../templates/partials/cards/guide-card.html",
		"../../templates/partials/cards/inline-card.html",
		"../../templates/partials/cards/link-card.html",
		"../../templates/partials/cards/note-card.html",
		"../../templates/partials/cards/quote-card.html",
		"../../templates/partials/cards/simple-card.html",
		"../../templates/partials/cards/video-card.html",
		"../../pkg/themes/default/templates/card.html",
		"../../pkg/themes/default/templates/components/post_byline.html",
		"../../pkg/themes/default/templates/og-card.html",
		"../../pkg/themes/default/templates/partials/cards/article-card.html",
		"../../pkg/themes/default/templates/partials/cards/default-card.html",
		"../../pkg/themes/default/templates/partials/cards/guide-card.html",
		"../../pkg/themes/default/templates/partials/cards/inline-card.html",
		"../../pkg/themes/default/templates/partials/cards/link-card.html",
		"../../pkg/themes/default/templates/partials/cards/note-card.html",
		"../../pkg/themes/default/templates/partials/cards/quote-card.html",
		"../../pkg/themes/default/templates/partials/cards/simple-card.html",
		"../../pkg/themes/default/templates/partials/cards/video-card.html",
	}
	legacyPatterns := []string{
		`| date:"Jan 2, 2006"`,
		`| date:"January 2, 2006"`,
		`| date:"2006-01-02"`,
		`| date_format:"Jan 2, 2006"`,
		`| date_format:"January 2, 2006"`,
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("ReadFile(%q) error: %v", file, err)
		}
		text := string(content)
		if !strings.Contains(text, "| human_date") {
			t.Fatalf("template %q does not use human_date", file)
		}
		for _, pattern := range legacyPatterns {
			if strings.Contains(text, pattern) {
				t.Fatalf("template %q still contains legacy date pattern %q", file, pattern)
			}
		}
	}
}

func TestTemplateTrees_PreserveSizedMediaDimensions(t *testing.T) {
	tests := []struct {
		file     string
		expected []string
	}{
		{
			file:     "../../pkg/themes/default/templates/partials/cards/article-card.html",
			expected: []string{`width="220"`, `height="160"`},
		},
		{
			file:     "../../pkg/themes/default/templates/partials/cards/video-card.html",
			expected: []string{`width="1200"`, `height="675"`},
		},
		{
			file:     "../../templates/partials/cards/video-card.html",
			expected: []string{`width="1200"`, `height="675"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			content, err := os.ReadFile(tt.file)
			if err != nil {
				t.Fatalf("ReadFile(%q) error: %v", tt.file, err)
			}
			text := string(content)
			for _, needle := range tt.expected {
				if !strings.Contains(text, needle) {
					t.Fatalf("template %q does not contain %s", tt.file, needle)
				}
			}
		})
	}
}

func TestTemplateTrees_UseCinematicMediaSizingForDefaultCards(t *testing.T) {
	files := []string{
		"../../templates/partials/cards/default-card.html",
		"../../pkg/themes/default/templates/partials/cards/default-card.html",
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			content, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("ReadFile(%q) error: %v", file, err)
			}
			text := string(content)

			if !strings.Contains(text, `|with_size:"1200"`) {
				t.Fatalf("default card %q must use width-only source sizing", file)
			}
			if strings.Contains(text, `|with_size:"1200,500"`) {
				t.Fatalf("default card %q must not request a distorted fixed-size source", file)
			}
			if !strings.Contains(text, `width="1200" height="500"`) {
				t.Fatalf("default card %q must declare cinematic media dimensions", file)
			}
		})
	}
}

func TestFilterSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"Hello  World", "hello-world"},
		{"Hello World!", "hello-world"},
		{"  Spaces  ", "spaces"},
		{"UPPERCASE", "uppercase"},
		{"with_underscore", "with_underscore"},
		{"123 Numbers", "123-numbers"},
	}

	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			ctx := NewContext(nil, "", nil)
			ctx.Set("input", tt.input)
			result, err := engine.RenderString("{{ input | slugify }}", ctx)
			if err != nil {
				t.Fatalf("RenderString() error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("slugify(%q): got %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFilterTruncate(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "short string",
			template: "{{ \"Hello\" | truncate:10 }}",
			expected: "Hello",
		},
		{
			name:     "truncate",
			template: "{{ \"Hello World Today\" | truncate:12 }}",
			expected: "Hello World...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext(nil, "", nil)
			result, err := engine.RenderString(tt.template, ctx)
			if err != nil {
				t.Fatalf("RenderString() error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFilterTruncateWords(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "fewer words than limit",
			template: "{{ \"Hello World\" | truncatewords:5 }}",
			expected: "Hello World",
		},
		{
			name:     "truncate words",
			template: "{{ \"One Two Three Four Five Six\" | truncatewords:3 }}",
			expected: "One Two Three ...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext(nil, "", nil)
			result, err := engine.RenderString(tt.template, ctx)
			if err != nil {
				t.Fatalf("RenderString() error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFilterDefaultIfNone(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Test with nil post (no title)
	post := &models.Post{}
	ctx := NewContext(post, "", nil)

	result, err := engine.RenderString("{{ post.title | default_if_none:\"Untitled\" }}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}

	if result != "Untitled" {
		t.Errorf("got %q, want %q", result, "Untitled")
	}

	// Test with value set - use a new post to avoid cache issues
	title := "My Title"
	post2 := &models.Post{Title: &title}
	ctx = NewContext(post2, "", nil)

	result, err = engine.RenderString("{{ post.title | default_if_none:\"Untitled\" }}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}

	if result != "My Title" {
		t.Errorf("got %q, want %q", result, "My Title")
	}
}

func TestFilterSlidesReveal_SplitsOnHeadingsAndRules(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	ctx := NewContext(nil, "", nil)
	ctx.Set("html", `<h2>Intro</h2><p>Welcome</p><h3>Details</h3><p>More</p><hr/><h2>End</h2><p>Bye</p>`)

	result, err := engine.RenderString(`{{ html | slides_reveal }}`, ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}

	checks := []string{
		`<section><section><div class="slide-content"><h2>Intro</h2><p>Welcome</p></div></section><section><div class="slide-content"><h3>Details</h3><p>More</p></div></section></section>`,
		`<section><div class="slide-content"><h2>End</h2><p>Bye</p></div></section>`,
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Fatalf("slides_reveal output missing %q in %q", check, result)
		}
	}

	if strings.Contains(result, "<div class=&#34;slide-content&#34;>") {
		t.Fatalf("expected safe HTML output, got escaped content: %q", result)
	}
}

func TestFilterSlidesReveal_PreservesLeadingContentAsFirstSlide(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	ctx := NewContext(nil, "", nil)
	ctx.Set("html", `<p>Lead</p><h3>Deep Dive</h3><p>More</p>`)

	result, err := engine.RenderString(`{{ html | slides_reveal }}`, ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}

	if !strings.Contains(result, `<section><section><div class="slide-content"><p>Lead</p></div></section><section><div class="slide-content"><h3>Deep Dive</h3><p>More</p></div></section></section>`) {
		t.Fatalf("unexpected reveal output: %q", result)
	}
}

func TestFilterLength(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	tests := []struct {
		name     string
		template string
		ctx      Context
		expected string
	}{
		{
			name:     "string",
			template: "{{ \"hello\" | length }}",
			ctx:      NewContext(nil, "", nil),
			expected: "5",
		},
		{
			name:     "slice",
			template: "{{ post.tags | length }}",
			ctx:      NewContext(&models.Post{Tags: []string{"a", "b", "c"}}, "", nil),
			expected: "3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.RenderString(tt.template, tt.ctx)
			if err != nil {
				t.Fatalf("RenderString() error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFilterJoin(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	post := &models.Post{Tags: []string{"a", "b", "c"}}
	ctx := NewContext(post, "", nil)

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "with separator",
			template: "{{ post.tags | join:\", \" }}",
			expected: "a, b, c",
		},
		{
			name:     "with pipe separator",
			template: "{{ post.tags | join:\" | \" }}",
			expected: "a | b | c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.RenderString(tt.template, ctx)
			if err != nil {
				t.Fatalf("RenderString() error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFilterFirstLast(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	post := &models.Post{Tags: []string{"first", "middle", "last"}}
	ctx := NewContext(post, "", nil)

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "first",
			template: "{{ post.tags | first }}",
			expected: "first",
		},
		{
			name:     "last",
			template: "{{ post.tags | last }}",
			expected: "last",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.RenderString(tt.template, ctx)
			if err != nil {
				t.Fatalf("RenderString() error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFilterReverse(t *testing.T) {
	// pongo2 has a built-in reverse filter that works on slices
	// For string reversal, we need to use our custom approach
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Test slice reversal
	post := &models.Post{Tags: []string{"a", "b", "c"}}
	ctx := NewContext(post, "", nil)
	result, err := engine.RenderString("{% for t in post.tags | reverse %}{{ t }}{% endfor %}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}

	if result != "cba" {
		t.Errorf("reverse got %q, want %q", result, "cba")
	}
}

func TestFilterStripTags(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple tag",
			input:    "<p>Hello</p>",
			expected: "Hello",
		},
		{
			name:     "nested tags",
			input:    "<div><p>Hello <strong>World</strong></p></div>",
			expected: "Hello World",
		},
		{
			name:     "no tags",
			input:    "Plain text",
			expected: "Plain text",
		},
		{
			name:     "html entities nbsp",
			input:    "Hello&nbsp;World",
			expected: "Hello World",
		},
		{
			name:     "html entities amp - re-escaped for html safety",
			input:    "Tom &amp; Jerry",
			expected: "Tom &amp; Jerry", // pongo2 re-escapes & for HTML safety
		},
		{
			name:     "html entities lt gt - re-escaped for html safety",
			input:    "&lt;code&gt;",
			expected: "&lt;code&gt;", // pongo2 re-escapes < > for HTML safety
		},
		{
			name:     "html entities quot - re-escaped for html safety",
			input:    "He said &quot;hello&quot;",
			expected: "He said &quot;hello&quot;", // pongo2 re-escapes quotes
		},
		{
			name:     "html entities apos",
			input:    "It&#39;s a test",
			expected: "It&#39;s a test", // pongo2 re-escapes apostrophes
		},
		{
			name:     "multiple whitespace collapse",
			input:    "Hello    World",
			expected: "Hello World",
		},
		{
			name:     "whitespace with newlines",
			input:    "Hello\n\n  World",
			expected: "Hello World",
		},
		{
			name:     "leading trailing whitespace",
			input:    "  Hello World  ",
			expected: "Hello World",
		},
		{
			name:     "combined scenario - tags removed, nbsp decoded, whitespace collapsed",
			input:    "<p>Hello&nbsp;&nbsp;World</p>  <br/>  <span>Test</span>",
			expected: "Hello World Test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext(nil, "", nil)
			ctx.Set("input", tt.input)
			result, err := engine.RenderString("{{ input | striptags }}", ctx)
			if err != nil {
				t.Fatalf("RenderString() error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFilterLinebreaksBR(t *testing.T) {
	// Note: pongo2's built-in linebreaksbr filter escapes HTML and uses <br />
	// Our custom filter uses <br> but pongo2's takes precedence
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	ctx := NewContext(nil, "", nil)
	ctx.Set("text", "Line 1\nLine 2")
	result, err := engine.RenderString("{{ text | linebreaksbr }}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}

	// pongo2 escapes HTML and uses <br />
	expected := "Line 1&lt;br /&gt;Line 2"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestFilterURLEncode(t *testing.T) {
	// pongo2's urlencode uses + for spaces (form encoding)
	// Our custom filter uses %20 (URL encoding)
	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "spaces (uses + for form encoding)",
			template: "{{ \"hello world\" | urlencode }}",
			expected: "hello+world", // pongo2 default behavior
		},
		{
			name:     "ampersand",
			template: "{{ \"a&b\" | urlencode }}",
			expected: "a%26b",
		},
	}

	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext(nil, "", nil)
			result, err := engine.RenderString(tt.template, ctx)
			if err != nil {
				t.Fatalf("RenderString() error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFilterAbsoluteURL(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	config := &models.Config{URL: "https://example.com"}
	post := &models.Post{Href: "/my-post/"}
	ctx := NewContext(post, "", config)

	result, err := engine.RenderString("{{ post.href | absolute_url:config.url }}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}

	expected := "https://example.com/my-post/"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestFilterExternalURL(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		input string
		want  string
	}{
		{"https://example.com/path", "https://example.com/path"},
		{"http://example.com", "http://example.com"},
		{"javascript:alert(1)", ""},
		{"//example.com/path", ""},
		{"https://user@example.com/path", ""},
		{"https://example.com/\nother", ""},
	} {
		ctx := NewContext(nil, "", nil)
		ctx.Set("source", tc.input)
		got, err := engine.RenderString(`{{ source | external_url }}`, ctx)
		if err != nil {
			t.Fatalf("RenderString(%q): %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("external_url(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestFilterExcerpt(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	html := `<p>First paragraph with some text.</p>
<p>Second paragraph with more content.</p>
<p></p>
<p>Third paragraph after an empty one.</p>
<p>Fourth paragraph here.</p>
<p>Fifth paragraph at the end.</p>`

	ctx := NewContext(nil, "", nil)
	ctx.Set("html", html)

	result, err := engine.RenderString("{{ html | excerpt }}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}

	t.Logf("Result: %s", result)

	// Should include first 3 non-empty paragraphs
	if !contains(result, "First paragraph") {
		t.Errorf("Expected 'First paragraph' in result")
	}
	if !contains(result, "Second paragraph") {
		t.Errorf("Expected 'Second paragraph' in result")
	}
	if !contains(result, "Third paragraph") {
		t.Errorf("Expected 'Third paragraph' in result")
	}
}

func TestFilterExcerpt_StripsBrokenPreTags(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	html := `<p>Before code.</p><p><code>git push remote: Push to create is not enabled for users.</code></pre><p>After code.</p>`

	ctx := NewContext(nil, "", nil)
	ctx.Set("html", html)

	result, err := engine.RenderString("{{ html | excerpt:\"paragraphs=3,chars=500\" }}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}

	if strings.Contains(result, "</pre>") {
		t.Fatalf("excerpt should strip dangling </pre> tags, got %q", result)
	}
	if !strings.Contains(result, "Before code.") {
		t.Fatalf("excerpt should keep earlier paragraph text, got %q", result)
	}
	if !strings.Contains(result, "Push to create is not enabled for users.") {
		t.Fatalf("excerpt should preserve inline code text, got %q", result)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestFilterIsVideo(t *testing.T) {
	// Ensure filters are registered
	registerFilters()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"mp4", "photo.mp4", true},
		{"webm", "clip.webm", true},
		{"mov", "movie.mov", true},
		{"m4v", "video.m4v", true},
		{"ogv", "video.ogv", true},
		{"png", "photo.png", false},
		{"jpg", "photo.jpg", false},
		{"empty", "", false},
		{"mp4_uppercase", "VIDEO.MP4", true},
		{"no_extension", "filename", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := pongo2.AsValue(tt.input)
			result, err := filterIsVideo(in, nil)
			if err != nil {
				t.Fatalf("filterIsVideo() error: %v", err)
			}
			got := result.Bool()
			if got != tt.want {
				t.Errorf("filterIsVideo(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestFilterMediaURL(t *testing.T) {
	// Ensure filters are registered
	registerFilters()

	tests := []struct {
		name     string
		primary  string
		fallback string
		want     string
	}{
		{"primary_only", "img.png", "", "img.png"},
		{"fallback_only", "", "video.mp4", "video.mp4"},
		{"both_prefers_primary", "img.png", "video.mp4", "img.png"},
		{"both_empty", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := pongo2.AsValue(tt.primary)
			param := pongo2.AsValue(tt.fallback)
			result, err := filterMediaURL(in, param)
			if err != nil {
				t.Fatalf("filterMediaURL() error: %v", err)
			}
			got := result.String()
			if got != tt.want {
				t.Errorf("filterMediaURL(%q, %q) = %q, want %q", tt.primary, tt.fallback, got, tt.want)
			}
		})
	}
}

func TestFilterIsVideoWithQueryFragment(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	ctx := NewContext(nil, "", nil)
	result, err := engine.RenderString("{{ 'https://dropper.wayl.one/movie.mp4?token=1#watch' | is_video }}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}
	if result != "True" {
		t.Errorf("is_video should treat query/fragment URLs as video, got %q", result)
	}
}

func TestFilterWithSizeTrustedRelativeAndUntrusted(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	ctx := NewContext(nil, "", nil)

	trusted, err := engine.RenderString("{{ 'https://dropper.wayl.one/image.jpg' | with_size:\"1200,675\" }}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}
	if !strings.Contains(trusted, "w=1200") || !strings.Contains(trusted, "h=675") {
		t.Errorf("with_size should append params for trusted hosts, got %q", trusted)
	}

	trustedHTTP, err := engine.RenderString("{{ 'http://dropper.wayl.one/image.jpg' | with_size:\"1200,675\" }}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}
	if !strings.HasPrefix(trustedHTTP, "https://dropper.wayl.one/image.jpg?") || !strings.Contains(trustedHTTP, "w=1200") || !strings.Contains(trustedHTTP, "h=675") {
		t.Errorf("with_size should normalize trusted http URLs to https, got %q", trustedHTTP)
	}

	widthOnly, err := engine.RenderString("{{ 'https://dropper.wayl.one/image.jpg?h=675' | with_size:\"1200\" }}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}
	if !strings.Contains(widthOnly, "w=1200") || strings.Contains(widthOnly, "h=675") {
		t.Errorf("width-only with_size should remove fixed height, got %q", widthOnly)
	}

	relative, err := engine.RenderString("{{ '/media/image.png' | with_size:\"1200,675\" }}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}
	if !strings.HasPrefix(relative, "/media/image.png?") || !strings.Contains(relative, "w=1200") || !strings.Contains(relative, "h=675") {
		t.Errorf("with_size should decorate relative URLs, got %q", relative)
	}

	untrusted, err := engine.RenderString("{{ 'https://example.com/image.jpg' | with_size:\"1200,675\" }}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}
	if untrusted != "https://example.com/image.jpg" {
		t.Errorf("with_size should leave untrusted hosts untouched, got %q", untrusted)
	}
}

func TestFilterResponsiveSrcsetTrustedAndUntrusted(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("NewEngine() error: %v", err)
	}
	ctx := NewContext(nil, "", nil)

	got, err := engine.RenderString(`{{ "https://dropper.waylonwalker.com/photo.webp?token=abc&w=600"|responsive_srcset:"320,160,320" }}`, ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}
	want := "https://dropper.waylonwalker.com/photo.webp?token=abc&amp;w=160 160w, https://dropper.waylonwalker.com/photo.webp?token=abc&amp;w=320 320w"
	if got != want {
		t.Fatalf("responsive_srcset output = %q, want %q", got, want)
	}

	got, err = engine.RenderString(`{{ "https://example.com/photo.webp"|responsive_srcset:"160,320" }}`, ctx)
	if err != nil {
		t.Fatalf("RenderString() for untrusted URL error: %v", err)
	}
	if got != "" {
		t.Fatalf("untrusted responsive_srcset output = %q, want empty", got)
	}
}

func TestPhotoCardRendersResponsiveRemoteCandidatesAndCanonicalLightboxURL(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("NewEngine() error: %v", err)
	}

	render := func(image string) string {
		t.Helper()
		got, err := engine.RenderToString("partials/cards/photo-card.html", map[string]interface{}{
			"post": map[string]interface{}{
				"image": image,
				"href":  "/shots/test-photo/",
				"title": "Test photo",
			},
		})
		if err != nil {
			t.Fatalf("RenderToString() error: %v", err)
		}
		return got
	}

	trusted := render("https://dropper.waylonwalker.com/photo.webp?token=abc&w=600")
	for _, want := range []string{
		`width="600"`,
		`loading="lazy" decoding="async"`,
		`srcset="https://dropper.waylonwalker.com/photo.webp?token=abc&amp;w=160 160w`,
		`sizes="(max-width: 632px) calc(100vw - 2rem), 600px"`,
		`data-href="https://dropper.waylonwalker.com/photo.webp?token=abc"`,
	} {
		if !strings.Contains(trusted, want) {
			t.Errorf("trusted photo card missing %q in %s", want, trusted)
		}
	}

	untrusted := render("https://example.com/photo.webp")
	if strings.Contains(untrusted, "srcset=") || strings.Contains(untrusted, "data-href=") {
		t.Fatalf("untrusted photo card unexpectedly changed: %s", untrusted)
	}
	if !strings.Contains(untrusted, `width="600"`) {
		t.Fatalf("untrusted photo card must still reserve the same display slot: %s", untrusted)
	}
	if strings.Contains(trusted, `sizes="auto,`) {
		t.Fatalf("auto sizes activates Chromium's 300x150 containment placeholder without an intrinsic height: %s", trusted)
	}
}

func TestFilterMediaURLNormalizesTrustedHTTP(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	ctx := NewContext(nil, "", nil)

	primary, err := engine.RenderString("{{ 'http://dropper.wayl.one/image.jpg' | media_url }}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}
	if primary != "https://dropper.wayl.one/image.jpg" {
		t.Errorf("media_url should normalize trusted http URLs, got %q", primary)
	}

	fallback, err := engine.RenderString("{{ '' | media_url:'http://dropper.wayl.one/video.mp4' }}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}
	if fallback != "https://dropper.wayl.one/video.mp4" {
		t.Errorf("media_url should normalize trusted fallback URLs, got %q", fallback)
	}
}

func TestFilterVideoMIMEHandlesQueryFragment(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	ctx := NewContext(nil, "", nil)
	result, err := engine.RenderString("{{ 'https://dropper.wayl.one/movie.mp4?token=1#watch' | video_mime }}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}
	if result != "video/mp4" {
		t.Errorf("video_mime should ignore query/fragment, got %q", result)
	}
}

func TestFilterPosterURLAliasAndFallback(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	title := "Alias Test"
	post := &models.Post{
		Title: &title,
		Extra: map[string]interface{}{
			"poster_image": "http://dropper.wayl.one/poster.png",
			"poster":       "http://dropper.wayl.one/fallback.png",
			"thumbnail":    "http://dropper.wayl.one/thumb.png",
		},
	}
	ctx := NewContext(post, "", nil)
	result, err := engine.RenderString("{{ post | poster_url:'http://dropper.wayl.one/video.mp4' }}", ctx)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}
	if result != "https://dropper.wayl.one/poster.png" {
		t.Errorf("poster_url should honor poster_image first, got %q", result)
	}

	post2 := &models.Post{}
	ctx2 := NewContext(post2, "", nil)
	fallback, err := engine.RenderString("{{ post | poster_url:'http://dropper.wayl.one/video.mp4' }}", ctx2)
	if err != nil {
		t.Fatalf("RenderString() error: %v", err)
	}
	if fallback != "https://dropper.wayl.one/video.webp" {
		t.Errorf("poster_url should derive .webp when no alias exists, got %q", fallback)
	}
}

func TestFilterSummary(t *testing.T) {
	in := "<p>Recent Arch vulns.</p>\n<ol><li><p>AUR is community</p></li></ol>\n<pre><code>sudo pacman -Rns yay\necho \"== pkg\"</code></pre>\n<p>Keep &amp; watch <a href=\"/x\">links</a>.</p>"
	got, err := filterSummary(pongo2.AsValue(in), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := got.String()
	if strings.Contains(s, "\n") {
		t.Errorf("summary must be single-line, got %q", s)
	}
	if strings.Contains(s, "pacman") || strings.Contains(s, "==") {
		t.Errorf("code blocks must be dropped, got %q", s)
	}
	want := "Recent Arch vulns. AUR is community Keep &amp; watch links."
	if s != want {
		t.Errorf("got %q want %q", s, want)
	}

	short, err := filterSummary(pongo2.AsValue("<p>one two three four</p>"), pongo2.AsValue(7))
	if err != nil {
		t.Fatalf("unexpected truncation error: %v", err)
	}
	if short.String() != "one two…" {
		t.Errorf("truncation: got %q", short.String())
	}
}

func TestFilterSummaryHeadingAnchors(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "heading permalink does not leak hash",
			input: `<h2>What happened <a href="#what-happened" class="heading-anchor">#</a></h2><p>is my replica pool</p>`,
			want:  "What happened is my replica pool",
		},
		{
			name:  "single quoted attributes and multiple classes",
			input: `<h3>Intro <a class='visually-hidden heading-anchor' href='#intro'><span>#</span></a></h3><p>Body</p>`,
			want:  "Intro Body",
		},
		{
			name:  "several permalinks with other links preserved",
			input: `<h2>First<a class="heading-anchor" href="#first">#</a></h2><p>See <a href="/more">more</a>.</p><h3>Second<a href="#second" class="heading-anchor">#</a></h3>`,
			want:  "First See more. Second",
		},
		{
			name:  "other anchors and literal hashes remain",
			input: `<p><a class="not-heading-anchor" href="/x">#reference</a> and <a href="/tags">#tag</a></p>`,
			want:  "#reference and #tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filterSummary(pongo2.AsValue(tt.input), nil)
			if err != nil {
				t.Fatalf("filterSummary() error: %v", err)
			}
			if got.String() != tt.want {
				t.Errorf("filterSummary() = %q, want %q", got.String(), tt.want)
			}
		})
	}
}

func TestTemplateTrees_LinkPostsExposeSourceAttribution(t *testing.T) {
	postTemplates := []string{
		"../../templates/post.html",
		"../../pkg/themes/default/templates/post.html",
	}
	cardTemplates := []string{
		"../../templates/partials/cards/link-card.html",
		"../../pkg/themes/default/templates/partials/cards/link-card.html",
	}
	for _, path := range append(postTemplates, cardTemplates...) {
		t.Run(path, func(t *testing.T) {
			b, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			template := string(b)
			for _, expected := range []string{
				"post.link | external_url",
				"post.url | external_url",
				"if not link_source_url",
				"| external_url",
				"u-bookmark-of",
				`target="_blank"`,
				"noopener noreferrer",
			} {
				if !strings.Contains(template, expected) {
					t.Errorf("%s must contain %q", path, expected)
				}
			}
		})
	}
	for _, path := range postTemplates {
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), "{% if not post.private %}") {
			t.Errorf("%s must not leak a private source URL", path)
		}
	}
}

func TestLinkSourceTemplateFallbackToValidURL(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatal(err)
	}
	const sourceTemplate = `{% set link_source_url = post.link | external_url %}{% if not link_source_url %}{% set link_source_url = post.url | external_url %}{% endif %}{{ link_source_url }}`
	for _, tc := range []struct {
		name string
		link string
		url  string
		want string
	}{
		{"valid primary", "https://example.com/first", "https://example.com/second", "https://example.com/first"},
		{"invalid primary", "javascript:alert(1)", "https://example.com/fallback", "https://example.com/fallback"},
		{"missing primary", "", "https://example.com/fallback", "https://example.com/fallback"},
		{"both unsafe", "javascript:alert(1)", "data:text/html,bad", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			post := &models.Post{Extra: map[string]interface{}{
				"link": tc.link,
				"url":  tc.url,
			}}
			ctx := NewContext(post, "", nil)
			got, err := engine.RenderString(sourceTemplate, ctx)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("source = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestEditorialCardImageReducesOnlyYouTubeMaxresPosters(t *testing.T) {
	tcs := []struct {
		name, input string
		lead        bool
		want        string
	}{
		{
			name:  "secondary youtube jpeg",
			input: "https://i.ytimg.com/vi/abc123/maxresdefault.jpg",
			want:  "https://i.ytimg.com/vi/abc123/mqdefault.jpg",
		},
		{
			name:  "featured youtube jpeg",
			input: "https://i.ytimg.com/vi/abc123/maxresdefault.jpg",
			lead:  true,
			want:  "https://i.ytimg.com/vi/abc123/sddefault.jpg",
		},
		{
			name:  "youtube alias with query",
			input: "https://img.youtube.com/vi/abc123/maxresdefault.jpg?foo=bar",
			want:  "https://img.youtube.com/vi/abc123/mqdefault.jpg?foo=bar",
		},
		{
			name:  "other youtube thumbnail already sized",
			input: "https://i.ytimg.com/vi/abc123/hqdefault.jpg",
			want:  "https://i.ytimg.com/vi/abc123/hqdefault.jpg",
		},
		{
			name:  "unrelated site same filename",
			input: "https://photos.example.org/maxresdefault.jpg",
			want:  "https://photos.example.org/maxresdefault.jpg",
		},
		{
			name:  "similar hostname is not youtube",
			input: "https://i.ytimg.com.attacker.invalid/vi/a/maxresdefault.jpg",
			want:  "https://i.ytimg.com.attacker.invalid/vi/a/maxresdefault.jpg",
		},
		{
			name:  "dropper responsive URL retained",
			input: "https://dropper.waylonwalker.com/file/abc.webp?width=720",
			want:  "https://dropper.waylonwalker.com/file/abc.webp?width=720",
		},
		{
			name:  "only https is eligible",
			input: "http://i.ytimg.com/vi/abc123/maxresdefault.jpg",
			want:  "http://i.ytimg.com/vi/abc123/maxresdefault.jpg",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got, err := filterEditorialCardImage(pongo2.AsValue(tc.input), pongo2.AsValue(tc.lead))
			if err != nil {
				t.Fatal(err)
			}
			if got.String() != tc.want {
				t.Fatalf("poster = %q, want %q", got.String(), tc.want)
			}
		})
	}
}

func TestEditorialCardImageRegisteredWithPongo(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatal(err)
	}
	rendered, err := engine.RenderString(`{{ url | editorial_card_image:lead }}`,
		Context{Extra: map[string]interface{}{
			"url":  "https://i.ytimg.com/vi/id/maxresdefault.jpg",
			"lead": true,
		}})
	if err != nil {
		t.Fatalf("filter could not render: %v", err)
	}
	if rendered != "https://i.ytimg.com/vi/id/sddefault.jpg" {
		t.Fatalf("rendered poster = %q", rendered)
	}
}
