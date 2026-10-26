package plugins

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

// BenchmarkRenderMarkdown_ColdStart measures markdown rendering with fresh goldmark instance.
func BenchmarkRenderMarkdown_ColdStart(b *testing.B) {
	posts := generateBenchmarkPosts(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		plugin := NewRenderMarkdownPlugin()
		m := lifecycle.NewManager()
		m.SetPosts(posts)

		if err := plugin.Render(m); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderMarkdown_HotCache measures markdown rendering reusing goldmark instance.
func BenchmarkRenderMarkdown_HotCache(b *testing.B) {
	posts := generateBenchmarkPosts(100)
	plugin := NewRenderMarkdownPlugin() // Reuse the same plugin

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := lifecycle.NewManager()
		// Create fresh posts each iteration (content must be re-rendered)
		freshPosts := make([]*models.Post, len(posts))
		for j, p := range posts {
			freshPosts[j] = &models.Post{
				Path:    p.Path,
				Slug:    p.Slug,
				Title:   p.Title,
				Content: p.Content,
			}
		}
		m.SetPosts(freshPosts)

		if err := plugin.Render(m); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderMarkdown_ContentSizes measures rendering with varying content sizes.
func BenchmarkRenderMarkdown_ContentSizes(b *testing.B) {
	sizes := []struct {
		name     string
		lines    int
		numPosts int
	}{
		{"small_10lines", 10, 100},
		{"medium_100lines", 100, 50},
		{"large_500lines", 500, 20},
		{"xlarge_1000lines", 1000, 10},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			posts := make([]*models.Post, size.numPosts)
			for i := 0; i < size.numPosts; i++ {
				title := fmt.Sprintf("Post %d", i)
				posts[i] = &models.Post{
					Path:    fmt.Sprintf("post-%d.md", i),
					Slug:    fmt.Sprintf("post-%d", i),
					Title:   &title,
					Content: generateMarkdownContent(size.lines),
				}
			}

			plugin := NewRenderMarkdownPlugin()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m := lifecycle.NewManager()
				freshPosts := clonePosts(posts)
				m.SetPosts(freshPosts)

				if err := plugin.Render(m); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkRenderMarkdown_SyntaxHighlighting measures code block rendering overhead.
func BenchmarkRenderMarkdown_SyntaxHighlighting(b *testing.B) {
	codeContent := `# Code Heavy Post

` + "```go\n" + `package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello, World!")
	os.Exit(0)
}
` + "```\n\n```python\n" + `
def hello():
    print("Hello, World!")

if __name__ == "__main__":
    hello()
` + "```\n\n```javascript\n" + `
function hello() {
    console.log("Hello, World!");
}

hello();
` + "```"

	posts := make([]*models.Post, 50)
	for i := 0; i < 50; i++ {
		title := fmt.Sprintf("Code Post %d", i)
		posts[i] = &models.Post{
			Path:    fmt.Sprintf("code-%d.md", i),
			Slug:    fmt.Sprintf("code-%d", i),
			Title:   &title,
			Content: codeContent,
		}
	}

	plugin := NewRenderMarkdownPlugin()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := lifecycle.NewManager()
		freshPosts := clonePosts(posts)
		m.SetPosts(freshPosts)

		if err := plugin.Render(m); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParseFrontmatter measures frontmatter parsing performance.
func BenchmarkParseFrontmatter(b *testing.B) {
	testCases := []struct {
		name    string
		content string
	}{
		{"simple", `---
title: "Simple Post"
date: 2024-01-15
published: true
---
Content here.`},
		{"with_tags", `---
title: "Post with Tags"
date: 2024-01-15
published: true
tags:
  - go
  - benchmark
  - performance
  - testing
---
Content here.`},
		{"complex", `---
title: "Complex Post with Many Fields"
description: "A very long description that goes on for quite a while to test parsing"
date: 2024-01-15T10:30:00Z
published: true
draft: false
template: custom.html
slug: custom-slug
tags:
  - go
  - benchmark
  - performance
extra_field1: value1
extra_field2: value2
nested:
  key: value
---
Content here.`},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, err := ParseFrontmatter(tc.content)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkGlob measures file discovery performance.
func BenchmarkGlob(b *testing.B) {
	testDir := setupBenchmarkDir(b, 500)
	defer os.RemoveAll(testDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		plugin := NewGlobPlugin()
		m := lifecycle.NewManager()
		m.Config().ContentDir = testDir
		m.Config().GlobPatterns = []string{"**/*.md"}

		if err := plugin.Configure(m); err != nil {
			b.Fatal(err)
		}
		if err := plugin.Glob(m); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGlob_WithGitignore measures glob with gitignore parsing.
func BenchmarkGlob_WithGitignore(b *testing.B) {
	testDir := setupBenchmarkDir(b, 500)
	defer os.RemoveAll(testDir)

	// Create .gitignore
	gitignoreContent := `
# Ignore output
output/
*.bak
.DS_Store
node_modules/
vendor/
.git/
`
	//nolint:gosec // G306: test files don't need restrictive permissions
	if err := os.WriteFile(filepath.Join(testDir, ".gitignore"), []byte(gitignoreContent), 0o644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		plugin := NewGlobPlugin()
		m := lifecycle.NewManager()
		m.Config().ContentDir = testDir
		m.Config().GlobPatterns = []string{"**/*.md"}

		if err := plugin.Configure(m); err != nil {
			b.Fatal(err)
		}
		if err := plugin.Glob(m); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkLoad measures file loading performance.
func BenchmarkLoad(b *testing.B) {
	testDir := setupBenchmarkDir(b, 100)
	defer os.RemoveAll(testDir)

	// Get file list
	plugin := NewGlobPlugin()
	m := lifecycle.NewManager()
	m.Config().ContentDir = testDir
	m.Config().GlobPatterns = []string{"**/*.md"}
	if err := plugin.Configure(m); err != nil {
		b.Fatal(err)
	}
	if err := plugin.Glob(m); err != nil {
		b.Fatal(err)
	}
	files := m.Files()

	loadPlugin := NewLoadPlugin()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := lifecycle.NewManager()
		m.Config().ContentDir = testDir
		m.SetFiles(files)

		if err := loadPlugin.Load(m); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkLoad_Concurrency measures load performance with varying concurrency.
func BenchmarkLoad_Concurrency(b *testing.B) {
	testDir := setupBenchmarkDir(b, 200)
	defer os.RemoveAll(testDir)

	// Get file list
	plugin := NewGlobPlugin()
	m := lifecycle.NewManager()
	m.Config().ContentDir = testDir
	m.Config().GlobPatterns = []string{"**/*.md"}
	if err := plugin.Configure(m); err != nil {
		b.Fatal(err)
	}
	if err := plugin.Glob(m); err != nil {
		b.Fatal(err)
	}
	files := m.Files()

	concurrencies := []int{1, 2, 4, 8, 16}
	loadPlugin := NewLoadPlugin()

	for _, conc := range concurrencies {
		b.Run(fmt.Sprintf("conc-%d", conc), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m := lifecycle.NewManager()
				m.Config().ContentDir = testDir
				m.SetConcurrency(conc)
				m.SetFiles(files)

				if err := loadPlugin.Load(m); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkTemplateEngine measures template rendering performance.
func BenchmarkTemplateEngine(b *testing.B) {
	// We need to test with actual embedded templates
	engine, err := templates.NewEngine("")
	if err != nil {
		b.Skip("Template engine not available")
	}

	post := &models.Post{
		Path:        "test.md",
		Slug:        "test-post",
		Content:     "Test content",
		ArticleHTML: "<p>Test content</p>",
		Published:   true,
	}
	title := "Test Post"
	post.Title = &title

	ctx := templates.NewContext(post, post.ArticleHTML, nil)

	b.Run("RenderString_Simple", func(b *testing.B) {
		tmpl := `<h1>{{ post.title }}</h1><div>{{ body }}</div>`
		for i := 0; i < b.N; i++ {
			_, err := engine.RenderString(tmpl, ctx)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("RenderString_Complex", func(b *testing.B) {
		tmpl := `<!DOCTYPE html>
<html>
<head><title>{{ post.title }}</title></head>
<body>
{% if post.title %}<h1>{{ post.title }}</h1>{% endif %}
{% if post.description %}<p>{{ post.description }}</p>{% endif %}
<article>{{ body }}</article>
{% if post.tags %}
<ul>{% for tag in post.tags %}<li>{{ tag }}</li>{% endfor %}</ul>
{% endif %}
</body>
</html>`
		for i := 0; i < b.N; i++ {
			_, err := engine.RenderString(tmpl, ctx)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkTemplateEngine_CacheHit measures template cache effectiveness.
func BenchmarkTemplateEngine_CacheHit(b *testing.B) {
	engine, err := templates.NewEngine("")
	if err != nil {
		b.Skip("Template engine not available")
	}

	// Check if post.html exists
	if !engine.TemplateExists("post.html") {
		b.Skip("post.html template not available")
	}

	post := &models.Post{
		Path:        "test.md",
		Slug:        "test-post",
		ArticleHTML: "<p>Test content</p>",
		Published:   true,
	}
	title := "Test Post"
	post.Title = &title
	ctx := templates.NewContext(post, post.ArticleHTML, nil)

	// First call to populate cache
	_, err = engine.Render("post.html", ctx)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Render("post.html", ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPublishHTML_Write measures file writing performance.
func BenchmarkPublishHTML_Write(b *testing.B) {
	posts := generateBenchmarkPosts(100)
	// Add HTML content
	for _, p := range posts {
		p.HTML = fmt.Sprintf("<html><body><h1>%s</h1><p>%s</p></body></html>", *p.Title, p.Content)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		outputDir, err := os.MkdirTemp("", "markata-bench-output-*")
		if err != nil {
			b.Fatal(err)
		}

		plugin := NewPublishHTMLPlugin()
		m := lifecycle.NewManager()
		m.Config().OutputDir = outputDir
		m.SetPosts(posts)

		if err := plugin.Write(m); err != nil {
			os.RemoveAll(outputDir)
			b.Fatal(err)
		}

		os.RemoveAll(outputDir)
	}
}

// Helper functions

func generateBenchmarkPosts(count int) []*models.Post {
	posts := make([]*models.Post, count)
	for i := 0; i < count; i++ {
		title := fmt.Sprintf("Benchmark Post %d", i)
		desc := fmt.Sprintf("Description for post %d", i)
		date := time.Now().AddDate(0, 0, -i)
		posts[i] = &models.Post{
			Path:        fmt.Sprintf("posts/post-%d.md", i),
			Slug:        fmt.Sprintf("post-%d", i),
			Title:       &title,
			Description: &desc,
			Date:        &date,
			Content:     generateMarkdownContent(50),
			Published:   i%2 == 0,
			Draft:       i%10 == 0,
			Tags:        []string{"benchmark", "test", fmt.Sprintf("tag-%d", i%5)},
		}
	}
	return posts
}

func generateMarkdownContent(lines int) string {
	var buf bytes.Buffer
	buf.WriteString("# Heading\n\n")
	for i := 0; i < lines; i++ {
		if i%10 == 0 {
			buf.WriteString(fmt.Sprintf("\n## Section %d\n\n", i/10))
		}
		if i%5 == 0 {
			buf.WriteString(fmt.Sprintf("- List item %d\n", i))
		} else {
			buf.WriteString(fmt.Sprintf("This is paragraph %d with some **bold** and *italic* text. ", i))
			buf.WriteString("Here's a [link](https://example.com) and `inline code`.\n\n")
		}
	}
	return buf.String()
}

func clonePosts(posts []*models.Post) []*models.Post {
	cloned := make([]*models.Post, len(posts))
	for i, p := range posts {
		cloned[i] = &models.Post{
			Path:        p.Path,
			Slug:        p.Slug,
			Title:       p.Title,
			Description: p.Description,
			Date:        p.Date,
			Content:     p.Content,
			Published:   p.Published,
			Draft:       p.Draft,
			Tags:        p.Tags,
		}
	}
	return cloned
}

func setupBenchmarkDir(b *testing.B, numFiles int) string {
	b.Helper()

	dir, err := os.MkdirTemp("", "markata-plugin-bench-*")
	if err != nil {
		b.Fatal(err)
	}

	subdirs := []string{"posts", "docs", "guides", "blog/2024", "blog/2023"}
	for _, subdir := range subdirs {
		if err := os.MkdirAll(filepath.Join(dir, subdir), 0o755); err != nil {
			b.Fatal(err)
		}
	}

	for i := 0; i < numFiles; i++ {
		subdir := subdirs[i%len(subdirs)]
		filename := fmt.Sprintf("post-%d.md", i)
		path := filepath.Join(dir, subdir, filename)
		content := generateBenchmarkFileContent(i)
		//nolint:gosec // G306: test files don't need restrictive permissions
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			b.Fatal(err)
		}
	}

	return dir
}

func generateBenchmarkFileContent(index int) string {
	date := time.Now().AddDate(0, 0, -index)
	return fmt.Sprintf(`---
title: "Benchmark Post %d"
description: "Test post for benchmarking"
date: %s
published: %t
tags:
  - benchmark
  - test
---

# Benchmark Post %d

This is content for benchmark post %d.

## Details

Some markdown content with **bold** and *italic* text.

- Item 1
- Item 2
- Item 3

`+"```go"+`
func example() {
    fmt.Println("Hello")
}
`+"```"+`
`, index, date.Format("2006-01-02"), index%2 == 0, index, index)
}
