//go:build ignore

// Package main generates deterministic benchmark posts for performance testing.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Deterministic content segments for reproducible benchmarks
var (
	loremParagraphs = []string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.",
		"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident.",
		"Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis.",
		"Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt.",
		"Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt.",
	}

	codeSnippets = map[string]string{
		"go": `package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}`,
		"python": `import asyncio
from dataclasses import dataclass
from typing import List, Optional

@dataclass
class Post:
    title: str
    content: str
    tags: List[str]
    published: bool = False

async def fetch_posts(limit: Optional[int] = None) -> List[Post]:
    """Fetch posts from the database asynchronously."""
    await asyncio.sleep(0.1)  # Simulate I/O
    posts = [
        Post("Hello", "World", ["intro"], True),
        Post("Second", "Content", ["update"], True),
    ]
    return posts[:limit] if limit else posts`,
		"javascript": `import { useState, useEffect } from 'react';

const useFetch = (url) => {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch(url);
        if (!response.ok) throw new Error('Network response was not ok');
        const json = await response.json();
        setData(json);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [url]);

  return { data, loading, error };
};

export default useFetch;`,
		"rust": `use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::RwLock;

#[derive(Debug, Clone)]
pub struct Cache<T> {
    data: Arc<RwLock<HashMap<String, T>>>,
}

impl<T: Clone> Cache<T> {
    pub fn new() -> Self {
        Self {
            data: Arc::new(RwLock::new(HashMap::new())),
        }
    }

    pub async fn get(&self, key: &str) -> Option<T> {
        let data = self.data.read().await;
        data.get(key).cloned()
    }

    pub async fn set(&self, key: String, value: T) {
        let mut data = self.data.write().await;
        data.insert(key, value);
    }
}`,
		"sql": `-- Create tables for blog application
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    content TEXT,
    published BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_posts_slug ON posts(slug);
CREATE INDEX idx_posts_published ON posts(published);

-- Insert sample data
INSERT INTO posts (title, slug, content, published) VALUES
    ('Hello World', 'hello-world', 'Welcome to my blog!', TRUE),
    ('Second Post', 'second-post', 'More content here.', TRUE);`,
	}

	tags = []string{
		"go", "python", "javascript", "rust", "sql",
		"tutorial", "guide", "reference", "howto",
		"performance", "benchmark", "testing", "ci",
		"web", "api", "database", "cache",
	}
)

func main() {
	baseDir := "benchmarks/site/posts"

	// Generate blog posts (60 posts)
	for i := 1; i <= 60; i++ {
		month := ((i - 1) % 12) + 1
		day := ((i - 1) % 28) + 1
		path := filepath.Join(baseDir, "blog", "2024", fmt.Sprintf("%02d", month), fmt.Sprintf("post-%03d.md", i))
		content := generateBlogPost(i, 2024, month, day)
		writeFile(path, content)
	}

	// Generate documentation guides (40 posts)
	for i := 1; i <= 40; i++ {
		path := filepath.Join(baseDir, "docs", "guides", fmt.Sprintf("guide-%03d.md", i))
		content := generateGuide(i)
		writeFile(path, content)
	}

	fmt.Println("Generated 100 benchmark posts")
}

func generateBlogPost(index, year, month, day int) string {
	title := fmt.Sprintf("Blog Post %d: %s Deep Dive", index, selectTag(index))
	slug := fmt.Sprintf("blog/%d/%02d/post-%03d", year, month, index)
	date := fmt.Sprintf("%d-%02d-%02d", year, month, day)
	postTags := selectTags(index, 3)
	published := index%10 != 0 // 90% published

	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %q\n", title))
	sb.WriteString(fmt.Sprintf("slug: %s\n", slug))
	sb.WriteString(fmt.Sprintf("date: %s\n", date))
	sb.WriteString(fmt.Sprintf("published: %t\n", published))
	sb.WriteString(fmt.Sprintf("description: \"A comprehensive guide to %s with practical examples and best practices.\"\n", selectTag(index)))
	sb.WriteString("tags:\n")
	for _, tag := range postTags {
		sb.WriteString(fmt.Sprintf("  - %s\n", tag))
	}
	sb.WriteString("---\n\n")

	// Content
	sb.WriteString(fmt.Sprintf("# %s\n\n", title))
	sb.WriteString(selectParagraph(index) + "\n\n")

	// Add code examples
	sb.WriteString("## Code Example\n\n")
	lang, code := selectCode(index)
	sb.WriteString(fmt.Sprintf("Here's an example in %s:\n\n", lang))
	sb.WriteString(fmt.Sprintf("```%s\n%s\n```\n\n", lang, code))

	// More content
	sb.WriteString("## Key Points\n\n")
	for j := 0; j < 5; j++ {
		sb.WriteString(fmt.Sprintf("- Point %d: %s\n", j+1, selectParagraph(index + j)[:50]+"..."))
	}
	sb.WriteString("\n")

	sb.WriteString("## Summary\n\n")
	sb.WriteString(selectParagraph(index+2) + "\n")

	return sb.String()
}

func generateGuide(index int) string {
	topics := []string{"Configuration", "Templates", "Plugins", "Themes", "Feeds", "Deployment", "Performance", "Testing"}
	topic := topics[index%len(topics)]
	title := fmt.Sprintf("Guide %d: %s", index, topic)
	slug := fmt.Sprintf("docs/guides/guide-%03d", index)
	date := fmt.Sprintf("2024-%02d-%02d", ((index-1)%12)+1, ((index-1)%28)+1)
	postTags := selectTags(index+100, 2)
	postTags = append(postTags, "documentation")

	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %q\n", title))
	sb.WriteString(fmt.Sprintf("slug: %s\n", slug))
	sb.WriteString(fmt.Sprintf("date: %s\n", date))
	sb.WriteString("published: true\n")
	sb.WriteString(fmt.Sprintf("description: \"Complete guide to %s in markata-go.\"\n", topic))
	sb.WriteString("tags:\n")
	for _, tag := range postTags {
		sb.WriteString(fmt.Sprintf("  - %s\n", tag))
	}
	sb.WriteString("---\n\n")

	sb.WriteString(fmt.Sprintf("# %s Guide\n\n", topic))
	sb.WriteString(selectParagraph(index) + "\n\n")

	// Multiple code examples for documentation
	sb.WriteString("## Getting Started\n\n")
	sb.WriteString(selectParagraph(index+1) + "\n\n")

	sb.WriteString("### Installation\n\n")
	sb.WriteString("```bash\ngo install github.com/WaylonWalker/markata/cmd/markata-go@latest\n```\n\n")

	sb.WriteString("### Configuration\n\n")
	sb.WriteString("```toml\n[markata-go]\ntitle = \"My Site\"\noutput_dir = \"public\"\n```\n\n")

	// Add language-specific example
	lang, code := selectCode(index)
	sb.WriteString(fmt.Sprintf("### %s Example\n\n", strings.Title(lang)))
	sb.WriteString(fmt.Sprintf("```%s\n%s\n```\n\n", lang, code))

	sb.WriteString("## Advanced Usage\n\n")
	sb.WriteString(selectParagraph(index+2) + "\n\n")

	// Checklist
	sb.WriteString("## Checklist\n\n")
	sb.WriteString("- [x] Initial setup\n")
	sb.WriteString("- [x] Basic configuration\n")
	sb.WriteString("- [ ] Advanced features\n")
	sb.WriteString("- [ ] Production deployment\n\n")

	// Table
	sb.WriteString("## Reference\n\n")
	sb.WriteString("| Option | Type | Default | Description |\n")
	sb.WriteString("|--------|------|---------|-------------|\n")
	sb.WriteString("| enabled | bool | true | Enable the feature |\n")
	sb.WriteString("| timeout | int | 30 | Timeout in seconds |\n")
	sb.WriteString("| verbose | bool | false | Enable verbose output |\n\n")

	sb.WriteString("## See Also\n\n")
	sb.WriteString("- [Configuration Guide](/docs/guides/configuration/)\n")
	sb.WriteString("- [Plugin Development](/docs/guides/plugins/)\n")

	return sb.String()
}

func selectParagraph(index int) string {
	return loremParagraphs[index%len(loremParagraphs)]
}

func selectTag(index int) string {
	return tags[index%len(tags)]
}

func selectTags(index, count int) []string {
	result := make([]string, count)
	for i := 0; i < count; i++ {
		result[i] = tags[(index+i)%len(tags)]
	}
	return result
}

func selectCode(index int) (string, string) {
	languages := []string{"go", "python", "javascript", "rust", "sql"}
	lang := languages[index%len(languages)]
	return lang, codeSnippets[lang]
}

func writeFile(path, content string) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", dir, err)
		os.Exit(1)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file %s: %v\n", path, err)
		os.Exit(1)
	}
}
