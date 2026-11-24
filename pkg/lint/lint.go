// Package lint provides markdown linting functionality for markata-go.
// It detects common issues in markdown files that can cause build failures.
//
// # Supported Checks
//
// The linter detects the following issues (via pkg/diagnostics):
//   - Duplicate YAML keys in frontmatter
//   - Invalid date formats (non-ISO 8601)
//   - Malformed image links (missing alt text)
//   - Protocol-less URLs (//example.com instead of https://example.com)
//   - H1 headings in content (templates add H1 from frontmatter title)
//   - Fenced code blocks in admonitions without blank line (goldmark limitation)
//
// All issues can be auto-fixed using the Fix function (except H1 headings).
package lint

import (
	"bufio"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/diagnostics"
)

// Pre-compiled regex patterns for lint fix operations.
var (
	keyRegex          = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\s*:`)
	noAltRegex        = regexp.MustCompile(`!\[\]\(([^)]+)\)`)
	protocollessRegex = regexp.MustCompile(`(\(|"|\s)//([a-zA-Z0-9][a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`)
	admonitionRegex   = regexp.MustCompile(`^(\s*)!!!\s+\w+`)
	fencedCodeRegex   = regexp.MustCompile(`^\s*` + "```")
)

// Issue represents a linting issue found in a file.
type Issue struct {
	File       string   // File path
	Line       int      // Line number (1-indexed)
	Column     int      // Column number (1-indexed, 0 if not applicable)
	Type       string   // Issue type (e.g., "duplicate-key", "invalid-date")
	Severity   Severity // Severity level
	Message    string   // Human-readable message
	Fixable    bool     // Whether this issue can be automatically fixed
	FixApplied bool     // Whether fix was applied
}

// Severity indicates the severity of a linting issue.
type Severity int

const (
	SeverityError Severity = iota
	SeverityWarning
	SeverityInfo
)

func (s Severity) String() string {
	switch s {
	case SeverityError:
		return "error"
	case SeverityWarning:
		return "warning"
	case SeverityInfo:
		return "info"
	default:
		return "unknown"
	}
}

// Result contains the linting results for a file.
type Result struct {
	File    string  // File path
	Issues  []Issue // Issues found
	Content string  // Original content
	Fixed   string  // Fixed content (same as Content if no fixes)
}

// HasErrors returns true if any issues are errors.
func (r *Result) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Severity == SeverityError {
			return true
		}
	}
	return false
}

// FixableCount returns the number of fixable issues.
func (r *Result) FixableCount() int {
	count := 0
	for _, issue := range r.Issues {
		if issue.Fixable {
			count++
		}
	}
	return count
}

// convertSeverity converts diagnostics.Severity to lint.Severity.
func convertSeverity(s diagnostics.Severity) Severity {
	switch s {
	case diagnostics.SeverityError:
		return SeverityError
	case diagnostics.SeverityWarning:
		return SeverityWarning
	case diagnostics.SeverityInfo:
		return SeverityInfo
	default:
		return SeverityWarning
	}
}

// convertIssue converts a diagnostics.Issue to a lint.Issue.
func convertIssue(di diagnostics.Issue) Issue {
	return Issue{
		File:     di.File,
		Line:     di.Range.StartLine + 1, // Convert 0-based to 1-based
		Column:   di.Range.StartCol + 1,  // Convert 0-based to 1-based
		Type:     di.Code,
		Severity: convertSeverity(di.Severity),
		Message:  di.Message,
		Fixable:  di.Fixable,
	}
}

// Lint analyzes content and returns any issues found.
func Lint(filePath, content string) *Result {
	result := &Result{
		File:    filePath,
		Content: content,
		Fixed:   content,
	}

	// Use shared diagnostics (without resolver - no wikilink/mention checks)
	diagIssues := diagnostics.Check(filePath, content, nil)

	for _, di := range diagIssues {
		result.Issues = append(result.Issues, convertIssue(di))
	}

	return result
}

// WithResolver analyzes content and returns any issues found,
// including wikilink and mention checks using the provided resolver.
func WithResolver(filePath, content string, resolver diagnostics.Resolver) *Result {
	result := &Result{
		File:    filePath,
		Content: content,
		Fixed:   content,
	}

	diagIssues := diagnostics.Check(filePath, content, resolver)

	for _, di := range diagIssues {
		result.Issues = append(result.Issues, convertIssue(di))
	}

	return result
}

// Fix applies automatic fixes to the content and returns the fixed content.
func Fix(filePath, content string) *Result {
	result := Lint(filePath, content)
	fixed := content

	// Apply fixes in order
	fixed = fixDuplicateKeys(fixed)
	fixed = fixDateFormats(fixed)
	fixed = fixImageLinks(fixed)
	fixed = fixProtocollessURLs(fixed)
	fixed = fixAdmonitionFencedCode(fixed)

	result.Fixed = fixed

	// Mark issues as fixed
	for i := range result.Issues {
		result.Issues[i].FixApplied = true
	}

	return result
}

// fixDuplicateKeys removes duplicate YAML keys, keeping the last occurrence.
func fixDuplicateKeys(content string) string {
	if !strings.HasPrefix(content, "---") {
		return content
	}

	parts := strings.SplitN(content[3:], "---", 2)
	if len(parts) < 2 {
		return content
	}

	frontmatter := parts[0]
	body := parts[1]

	// Track keys and their lines
	type keyLine struct {
		key  string
		line string
	}
	var lines []keyLine
	seen := make(map[string]int) // key -> index in lines

	scanner := bufio.NewScanner(strings.NewReader(frontmatter))
	var currentKey string
	var currentLines []string

	flushCurrent := func() {
		if currentKey != "" && len(currentLines) > 0 {
			combined := strings.Join(currentLines, "\n")
			if idx, exists := seen[currentKey]; exists {
				// Replace previous occurrence
				lines[idx] = keyLine{key: currentKey, line: combined}
			} else {
				seen[currentKey] = len(lines)
				lines = append(lines, keyLine{key: currentKey, line: combined})
			}
		}
		currentKey = ""
		currentLines = nil
	}

	for scanner.Scan() {
		line := scanner.Text()

		if match := keyRegex.FindStringSubmatch(line); match != nil {
			// New key found, flush previous
			flushCurrent()
			currentKey = match[1]
			currentLines = []string{line}
		} else if currentKey != "" && (strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") || line == "") {
			// Continuation of current key (indented or empty)
			currentLines = append(currentLines, line)
		} else if currentKey != "" {
			// Non-indented line that's not a key - might be a YAML array/list item
			currentLines = append(currentLines, line)
		}
	}
	flushCurrent()

	// Rebuild frontmatter
	fixedLines := make([]string, 0, len(lines))
	for _, kl := range lines {
		fixedLines = append(fixedLines, kl.line)
	}

	return "---\n" + strings.Join(fixedLines, "\n") + "\n---" + body
}

// fixDateFormats normalizes date formats to ISO 8601 using DateTimeFixer.
func fixDateFormats(content string) string {
	fixer := NewDateTimeFixer(DefaultDateTimeFixerConfig())
	fixed, _ := fixer.FixDateInContent(content)
	return fixed
}

// fixImageLinks adds placeholder alt text to images without it.
func fixImageLinks(content string) string {
	return noAltRegex.ReplaceAllString(content, "![image]($1)")
}

// fixProtocollessURLs adds https:// to protocol-less URLs.
func fixProtocollessURLs(content string) string {
	return protocollessRegex.ReplaceAllString(content, "${1}https://$2")
}

// fixAdmonitionFencedCode adds blank lines after admonition declarations
// that are immediately followed by fenced code blocks.
func fixAdmonitionFencedCode(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		result = append(result, line)

		// Check for admonition start
		if match := admonitionRegex.FindStringSubmatch(line); match != nil {
			admonitionIndent := len(match[1])

			// Check if next line is a fenced code block without blank line
			if i+1 < len(lines) {
				nextLine := lines[i+1]
				trimmedNext := strings.TrimLeft(nextLine, " \t")
				nextIndent := len(nextLine) - len(trimmedNext)

				// If next line is indented and starts with ``` (fenced code)
				if nextIndent > admonitionIndent && fencedCodeRegex.MatchString(nextLine) {
					// Add a blank line with proper indentation
					result = append(result, "")
				}
			}
		}
	}

	return strings.Join(result, "\n")
}
