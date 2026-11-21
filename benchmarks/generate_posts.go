//go:build ignore

// Package main generates deterministic benchmark posts for performance testing.
package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

var (
	titleAdjectives = []string{"Fast", "Practical", "Random", "Resilient", "Modern", "Focused"}
	titleTopics     = []string{"Rendering", "Caching", "Publishing", "Markdown", "Templates", "Feeds"}
	sentenceParts   = []string{
		"Randomized benchmark content keeps content generation reproducible",
		"Markdown rendering remains faithful across varied post lengths",
		"Small fixture changes should expose build and cache behavior",
		"Progressive paragraphs make scaling effects easier to compare",
		"Repeated warm builds should restore stable cache results",
		"Generated examples exercise headings emphasis lists and code",
	}
	tags = []string{"benchmark", "markdown", "performance", "testing", "go", "random"}
)

func main() {
	baseDir := "benchmarks/site/posts"
	if err := os.RemoveAll(baseDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing old benchmark posts: %v\n", err)
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(20261120))
	for i := 1; i <= 100; i++ {
		path := filepath.Join(baseDir, fmt.Sprintf("post-%03d.md", i))
		writeFile(path, generatePost(rng, i))
	}

	fmt.Println("Generated 100 benchmark posts")
}

func generatePost(rng *rand.Rand, index int) string {
	title := fmt.Sprintf("%s %s Benchmark %03d",
		titleAdjectives[rng.Intn(len(titleAdjectives))],
		titleTopics[rng.Intn(len(titleTopics))],
		index)
	date := fmt.Sprintf("2024-%02d-%02d", ((index-1)%12)+1, ((index-1)%28)+1)

	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %q\n", title))
	sb.WriteString(fmt.Sprintf("slug: %q\n", fmt.Sprintf("post-%03d", index)))
	sb.WriteString(fmt.Sprintf("date: %s\n", date))
	sb.WriteString("published: true\n")
	sb.WriteString("description: \"A progressive benchmark post with varied Markdown content\"\n")
	sb.WriteString("tags:\n")
	for i := 0; i < 3; i++ {
		sb.WriteString(fmt.Sprintf("  - %s\n", tags[(index+i)%len(tags)]))
	}
	sb.WriteString("---\n\n")

	sb.WriteString(fmt.Sprintf("# %s\n\n", title))

	paragraphCount := 1 + (index-1)*2/99
	for paragraph := 0; paragraph < paragraphCount; paragraph++ {
		sentenceCount := 1 + (index-1)/50
		if index == 1 {
			sentenceCount = 1
		}
		sentences := make([]string, 0, sentenceCount)
		for sentence := 0; sentence < sentenceCount; sentence++ {
			part := sentenceParts[rng.Intn(len(sentenceParts))]
			sentences = append(sentences, fmt.Sprintf("%s for benchmark item %03d.", part, index))
		}
		sb.WriteString(strings.Join(sentences, " "))
		sb.WriteString("\n\n")
	}

	sb.WriteString("## Markdown samples\n\n")
	sb.WriteString(fmt.Sprintf("This paragraph has **bold text**, *italic text*, and a [project link](/). Item %d uses `inline code`.\n\n", index))
	sb.WriteString("- Randomized fixture item\n")
	sb.WriteString("- Reproducible fixed-seed content\n")
	sb.WriteString("- Progressive paragraph coverage\n\n")
	sb.WriteString("> Benchmark Markdown should remain readable and structurally varied.\n\n")
	sb.WriteString("```go\n")
	sb.WriteString(fmt.Sprintf("var benchmarkItem = %d\n", index))
	sb.WriteString("```\n")

	return sb.String()
}

func writeFile(path, content string) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory %v: %v\n", dir, err)
		os.Exit(1)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file %v: %v\n", path, err)
		os.Exit(1)
	}
}
