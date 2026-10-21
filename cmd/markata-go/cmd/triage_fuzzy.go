package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
)

const (
	triageManualPostOption    = "Enter a Markdown path manually..."
	triageMaxPostCandidates   = 5000
	triageMaxPostNestedLevels = 10
)

var errTriageIntakeCanceled = errors.New("triage intake canceled")

// triagePicker lets tests exercise the question order without opening a
// terminal UI. The production picker uses the same fuzzyfinder as the
// installed-agent picker.
type triagePicker func(q triageQuestion, options []string) (string, error)

func triageFuzzyPick(q triageQuestion, options []string) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("no choices available for %s", q.key)
	}
	index, err := fuzzyfinder.Find(options, func(i int) string { return triageIntakeOptionLabel(q, options[i]) })
	if errors.Is(err, fuzzyfinder.ErrAbort) {
		return "", errTriageIntakeCanceled
	}
	if err != nil {
		return "", fmt.Errorf("triage %s selection: %w", q.key, err)
	}
	if index < 0 || index >= len(options) {
		return "", fmt.Errorf("triage %s selection returned invalid index", q.key)
	}
	return options[index], nil
}

func triageIntakeOptionLabel(q triageQuestion, option string) string {
	switch q.key {
	case "category":
		switch option {
		case "post":
			return "post — Debug or enhance an existing Markdown post"
		case "lint":
			return "lint — Find content errors in a post or the whole site"
		case "general":
			return "general — Troubleshoot Markata or its configuration"
		}
	case "scope":
		switch option {
		case "post":
			return "post — Lint a specific Markdown post"
		case "site":
			return "site — Lint the whole site"
		}
	case "action":
		switch option {
		case "debug":
			return "debug — Fix a broken post feature or renderer error"
		case "enhance":
			return "enhance — Add a chart, road, diagram, or other feature"
		case "review":
			return "review — Evaluate existing Markdown and structure"
		}
	case "feature":
		switch option {
		case "chart":
			return "chart — Chart data, syntax, axes, or rendering"
		case "road":
			return "road — Timeline, milestones, or static progress"
		case "mermaid":
			return "mermaid — Diagram syntax or rendering"
		case "markdown":
			return "markdown — Formatting, tables, and other Markdown"
		case "other":
			return "other — Another post feature"
		}
	}
	// Post paths remain visually identical to the values included in the
	// handoff, which makes substring/fuzzy matching predictable.
	return option
}

func triageIsMarkdown(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md" || ext == ".markdown"
}

// triageMarkdownCandidates returns *filenames only*. It neither parses site
// config nor opens Markdown contents. Restrict traversal to known content
// folders and refuse hidden, generated and dependency directories. Symlinks
// are never followed. Users with custom content locations can enter a path.
func triageIsRepositoryMeta(name string) bool {
	switch strings.ToLower(name) {
	case "readme.md", "agents.md", "contributing.md", "changelog.md",
		"code_of_conduct.md", "security.md", "license.md":
		return true
	}
	return false
}

func triageMarkdownCandidates(siteRoot string) []string {
	if siteRoot == "" {
		return nil
	}
	roots := []string{"posts", "pages", "docs", "content", "src/posts", "src/pages", "src/content"}
	excluded := map[string]bool{
		".git": true, ".markata": true, ".markata-cache": true,
		"node_modules": true, "vendor": true, "dist": true,
		"build": true, "output": true, "static": true,
		"__pycache__": true,
	}
	paths := make([]string, 0, 256)
	seen := make(map[string]bool)
	add := func(path string) {
		rel, err := filepath.Rel(siteRoot, path)
		if err != nil || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
			return
		}
		rel = filepath.ToSlash(rel)
		if len(rel) <= 512 && !seen[rel] {
			seen[rel] = true
			paths = append(paths, rel)
		}
	}

	// Include root-level Markdown posts without recursively walking the
	// entire repository.
	if entries, err := os.ReadDir(siteRoot); err == nil {
		for _, entry := range entries {
			if len(paths) == triageMaxPostCandidates {
				break
			}
			if entry.Type().IsRegular() && !strings.HasPrefix(entry.Name(), ".") &&
				!triageIsRepositoryMeta(entry.Name()) && triageIsMarkdown(entry.Name()) {
				add(filepath.Join(siteRoot, entry.Name()))
			}
		}
	}
	for _, folder := range roots {
		if len(paths) == triageMaxPostCandidates {
			break
		}
		base := filepath.Join(siteRoot, filepath.FromSlash(folder))
		info, err := os.Lstat(base)
		if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			continue
		}
		_ = filepath.WalkDir(base, func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return nil // Missing/unreadable content dirs should never block triage.
			}
			if len(paths) >= triageMaxPostCandidates {
				return filepath.SkipAll
			}
			if path != base && entry.IsDir() {
				rel, err := filepath.Rel(base, path)
				if err != nil {
					return filepath.SkipDir
				}
				if excluded[entry.Name()] || strings.HasPrefix(entry.Name(), ".") ||
					strings.Count(filepath.ToSlash(rel), "/") >= triageMaxPostNestedLevels {
					return filepath.SkipDir
				}
			}
			if !entry.IsDir() && entry.Type()&os.ModeSymlink == 0 &&
				!strings.HasPrefix(entry.Name(), ".") && triageIsMarkdown(entry.Name()) {
				add(path)
			}
			return nil
		})
	}
	sort.Strings(paths)
	return paths
}
