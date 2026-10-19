package diagnostics

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

const (
	diagnosticCodeDuplicateKey    = "duplicate-key"
	diagnosticCodeInvalidDate     = "invalid-date"
	diagnosticCodeH1InContent     = "h1-in-content"
	diagnosticCodeMissingAltText  = "missing-alt-text"
	diagnosticCodeProtocolLessURL = "protocol-less-url"
	diagnosticCodeAdmonitionFence = "admonition-fenced-code"
	diagnosticCodeBrokenWikilink  = "broken-wikilink"
	diagnosticCodeUnknownMention  = "unknown-mention"
)

// Severity indicates the severity of a diagnostic issue.
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

// Range represents a position range in a document.
type Range struct {
	StartLine int `json:"start_line"` // 0-based line number
	StartCol  int `json:"start_col"`  // 0-based column (character offset)
	EndLine   int `json:"end_line"`   // 0-based line number
	EndCol    int `json:"end_col"`    // 0-based column
}

// Issue represents a diagnostic issue found in a file.
type Issue struct {
	File     string   `json:"file"`     // File path
	Range    Range    `json:"range"`    // Position in the file
	Code     string   `json:"code"`     // Issue code (e.g., "duplicate-key", "broken-wikilink")
	Severity Severity `json:"severity"` // Severity level
	Message  string   `json:"message"`  // Human-readable message
	Fixable  bool     `json:"fixable"`  // Whether this issue can be automatically fixed
}

// Resolver provides lookup functionality for wikilinks and mentions.
// This interface allows the diagnostics package to check if references exist
// without directly depending on the LSP index or any specific implementation.
type Resolver interface {
	// ResolveSlug returns true if a post with the given slug exists.
	ResolveSlug(slug string) bool
	// ResolveHandle returns true if a mention handle exists in the blogroll.
	ResolveHandle(handle string) bool
}

// Check runs all diagnostic checks on the content and returns any issues found.
// The resolver is optional; if nil, wikilink and mention checks are skipped.
func Check(filePath, content string, resolver Resolver) []Issue {
	analysis := AnalyzeFrontmatter(filePath, content)
	issues := append([]Issue{}, analysis.Issues...)
	frontmatter := analysis.Inspection.Frontmatter
	body := analysis.Inspection.Body
	hasFrontmatter := analysis.Inspection.HasFrontmatter

	// Body checks
	issues = append(issues, checkImageLinks(filePath, body, hasFrontmatter, frontmatter)...)
	issues = append(issues, checkProtocollessURLs(filePath, content)...)
	issues = append(issues, checkH1Headings(filePath, body, hasFrontmatter, frontmatter)...)
	issues = append(issues, checkAdmonitionFencedCode(filePath, body, hasFrontmatter, frontmatter)...)

	// Reference checks (require resolver)
	if resolver != nil {
		issues = append(issues, checkWikilinks(filePath, body, hasFrontmatter, frontmatter, resolver)...)
		issues = append(issues, checkMentions(filePath, body, hasFrontmatter, frontmatter, resolver)...)
	}

	return issues
}

// checkDuplicateKeys finds duplicate YAML keys in frontmatter.
func checkDuplicateKeys(filePath, frontmatter string) []Issue {
	var issues []Issue
	seen := make(map[string]int) // key -> first line number
	scanner := bufio.NewScanner(strings.NewReader(frontmatter))
	lineNum := 1 // Start at 1 (after the opening ---)

	// Regex to match top-level YAML keys (not indented)
	keyRegex := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\s*:`)

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		if match := keyRegex.FindStringSubmatch(line); match != nil {
			key := match[1]
			if firstLine, exists := seen[key]; exists {
				issues = append(issues, Issue{
					File: filePath,
					Range: Range{
						StartLine: lineNum - 1, // 0-based
						StartCol:  0,
						EndLine:   lineNum - 1,
						EndCol:    len(line),
					},
					Code:     diagnosticCodeDuplicateKey,
					Severity: SeverityError,
					Message:  fmt.Sprintf("duplicate key '%s' (first occurrence at line %d)", key, firstLine),
					Fixable:  true,
				})
			} else {
				seen[key] = lineNum
			}
		}
	}

	return issues
}

// checkDateFormats validates date formats in frontmatter.
func checkDateFormats(filePath, frontmatter string) []Issue {
	var issues []Issue
	scanner := bufio.NewScanner(strings.NewReader(frontmatter))
	lineNum := 1

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		if match := frontmatterDateLineRegex.FindStringSubmatch(line); match != nil {
			key := strings.ToLower(match[1])
			if !dateKeyRegex.MatchString(key) {
				continue
			}
			value := strings.TrimSpace(match[2])
			if isNullDateValue(value) || isSupportedDateString(value) {
				continue
			}
			issues = append(issues, Issue{
				File: filePath,
				Range: Range{
					StartLine: lineNum - 1,
					StartCol:  0,
					EndLine:   lineNum - 1,
					EndCol:    len(line),
				},
				Code:     diagnosticCodeInvalidDate,
				Severity: SeverityWarning,
				Message:  fmt.Sprintf("invalid date format for '%s'", key),
				Fixable:  true,
			})
		}
	}

	return issues
}

func isNullDateValue(value string) bool {
	value = strings.ToLower(strings.Trim(strings.TrimSpace(value), "\"'"))
	switch value {
	case "", "null", "nil", "~":
		return true
	default:
		return false
	}
}

// checkImageLinks finds malformed image links.
func checkImageLinks(filePath, body string, hasFrontmatter bool, frontmatter string) []Issue {
	var issues []Issue

	// Regex for image links without alt text: ![](url)
	noAltRegex := regexp.MustCompile(`!\[\]\(([^)]+)\)`)

	// Calculate line offset for body
	lineOffset := 0
	if hasFrontmatter {
		lineOffset = frontmatterLineOffset(frontmatter)
	}

	scanner := bufio.NewScanner(strings.NewReader(body))
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		if matches := noAltRegex.FindAllStringSubmatchIndex(line, -1); matches != nil {
			for _, match := range matches {
				issues = append(issues, Issue{
					File: filePath,
					Range: Range{
						StartLine: lineNum + lineOffset - 1,
						StartCol:  match[0],
						EndLine:   lineNum + lineOffset - 1,
						EndCol:    match[1],
					},
					Code:     diagnosticCodeMissingAltText,
					Severity: SeverityWarning,
					Message:  "image link missing alt text",
					Fixable:  true,
				})
			}
		}
	}

	return issues
}

// checkProtocollessURLs finds protocol-less URLs.
func checkProtocollessURLs(filePath, content string) []Issue {
	var issues []Issue

	// Regex for protocol-less URLs: //example.com
	protocollessRegex := regexp.MustCompile(`[^:]//[a-zA-Z0-9][a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		if matches := protocollessRegex.FindAllStringIndex(line, -1); matches != nil {
			for _, match := range matches {
				issues = append(issues, Issue{
					File: filePath,
					Range: Range{
						StartLine: lineNum - 1,
						StartCol:  match[0] + 1, // +1 to skip the non-colon char
						EndLine:   lineNum - 1,
						EndCol:    match[1],
					},
					Code:     diagnosticCodeProtocolLessURL,
					Severity: SeverityWarning,
					Message:  "protocol-less URL found (should use https://)",
					Fixable:  true,
				})
			}
		}
	}

	return issues
}

// checkH1Headings finds H1 headings in markdown content.
func checkH1Headings(filePath, body string, hasFrontmatter bool, frontmatter string) []Issue {
	var issues []Issue

	lineOffset := 0
	if hasFrontmatter {
		// InspectFrontmatter returns only the YAML body, without the newline
		// surrounding the opening and closing delimiters.
		lineOffset = frontmatterLineOffset(frontmatter)
	}

	inCodeBlock := false
	scanner := bufio.NewScanner(strings.NewReader(body))
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "```") || strings.HasPrefix(trimmedLine, "~~~") {
			inCodeBlock = !inCodeBlock
			continue
		}

		if inCodeBlock {
			continue
		}

		if strings.HasPrefix(line, "# ") || line == "#" {
			issues = append(issues, Issue{
				File: filePath,
				Range: Range{
					StartLine: lineNum + lineOffset - 1,
					StartCol:  0,
					EndLine:   lineNum + lineOffset - 1,
					EndCol:    len(line),
				},
				Code:     diagnosticCodeH1InContent,
				Severity: SeverityWarning,
				Message:  "H1 heading found in content. Templates already add an H1 from frontmatter title. Use H2 (##) or deeper instead.",
				Fixable:  false,
			})
		}
	}

	return issues
}

// checkAdmonitionFencedCode detects fenced code blocks inside admonitions
// that don't have a blank line before them.
func checkAdmonitionFencedCode(filePath, body string, hasFrontmatter bool, frontmatter string) []Issue {
	var issues []Issue

	lineOffset := 0
	if hasFrontmatter {
		lineOffset = frontmatterLineOffset(frontmatter)
	}

	lines := strings.Split(body, "\n")

	admonitionRegex := regexp.MustCompile(`^(\s*)!!!\s+\w+`)
	fencedCodeRegex := regexp.MustCompile(`^\s*` + "```")

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		lineNum := i + 1 + lineOffset

		if match := admonitionRegex.FindStringSubmatch(line); match != nil {
			admonitionIndent := len(match[1])

			if i+1 < len(lines) {
				nextLine := lines[i+1]
				trimmedNext := strings.TrimLeft(nextLine, " \t")
				nextIndent := len(nextLine) - len(trimmedNext)

				if nextIndent > admonitionIndent && fencedCodeRegex.MatchString(nextLine) {
					issues = append(issues, Issue{
						File: filePath,
						Range: Range{
							StartLine: lineNum - 1,
							StartCol:  0,
							EndLine:   lineNum - 1,
							EndCol:    len(line),
						},
						Code:     diagnosticCodeAdmonitionFence,
						Severity: SeverityWarning,
						Message:  "fenced code block immediately follows admonition without blank line - this may not render correctly due to goldmark limitation",
						Fixable:  true,
					})
				}
			}
		}
	}

	return issues
}

// wikilinkRegex matches [[slug]] and [[slug|display text]] patterns.
var wikilinkRegex = regexp.MustCompile(`\[\[([^\]|]+)(?:\|([^\]]+))?\]\]`)

// checkWikilinks finds broken wikilinks in content.
func checkWikilinks(filePath, body string, hasFrontmatter bool, frontmatter string, resolver Resolver) []Issue {
	var issues []Issue

	lineOffset := 0
	if hasFrontmatter {
		lineOffset = frontmatterLineOffset(frontmatter)
	}

	lines := strings.Split(body, "\n")
	inCodeBlock := false
	codeBlockPattern := regexp.MustCompile("^```|^~~~")

	for lineNum, line := range lines {
		trimmed := strings.TrimSpace(line)
		if codeBlockPattern.MatchString(trimmed) {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			continue
		}

		matches := wikilinkRegex.FindAllStringSubmatchIndex(line, -1)
		for _, match := range matches {
			if len(match) < 4 {
				continue
			}

			slugStart := match[2]
			slugEnd := match[3]
			slug := strings.TrimSpace(line[slugStart:slugEnd])

			if !resolver.ResolveSlug(slug) {
				issues = append(issues, Issue{
					File: filePath,
					Range: Range{
						StartLine: lineNum + lineOffset,
						StartCol:  match[0],
						EndLine:   lineNum + lineOffset,
						EndCol:    match[1],
					},
					Code:     diagnosticCodeBrokenWikilink,
					Severity: SeverityWarning,
					Message:  fmt.Sprintf("broken wikilink: target post %q not found", slug),
					Fixable:  false,
				})
			}
		}
	}

	return issues
}

// mentionRegex matches @mentions and ~quote attributions.
var mentionRegex = regexp.MustCompile(`([@~])([a-zA-Z][a-zA-Z0-9_.-]*)`)

// resolveMentionForLint checks the complete handle first, allowing valid
// domain-style names, then retries without trailing sentence punctuation.
func resolveMentionForLint(resolver Resolver, handle string) (string, bool) {
	if resolver.ResolveHandle(handle) {
		return handle, true
	}
	stripped := strings.TrimRight(handle, ".,;:!?")
	if stripped == "" || stripped == handle {
		return handle, false
	}
	return stripped, resolver.ResolveHandle(stripped)
}

// checkMentions finds unknown mentions in content.
func checkMentions(filePath, body string, hasFrontmatter bool, frontmatter string, resolver Resolver) []Issue {
	var issues []Issue

	lineOffset := 0
	if hasFrontmatter {
		lineOffset = frontmatterLineOffset(frontmatter)
	}

	lines := strings.Split(body, "\n")
	inCodeBlock := false
	codeBlockPattern := regexp.MustCompile("^```|^~~~")

	for lineNum, line := range lines {
		trimmed := strings.TrimSpace(line)
		if codeBlockPattern.MatchString(trimmed) {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			continue
		}

		matches := mentionRegex.FindAllStringSubmatchIndex(line, -1)
		for _, match := range matches {
			if len(match) < 4 {
				continue
			}

			start := match[0]
			// Validate that @ is at a valid boundary
			if start > 0 {
				prevChar := line[start-1]
				if (prevChar >= 'a' && prevChar <= 'z') ||
					(prevChar >= 'A' && prevChar <= 'Z') ||
					(prevChar >= '0' && prevChar <= '9') || prevChar == '_' || prevChar == '@' || prevChar == '~' {
					continue // Skip email addresses and @@mentions
				}
			}

			marker := line[match[2]:match[3]]
			handleStart := match[4]
			handleEnd := match[5]
			handle := strings.ToLower(line[handleStart:handleEnd])
			handle, known := resolveMentionForLint(resolver, handle)
			if known {
				continue
			}

			issues = append(issues, Issue{
				File: filePath,
				Range: Range{
					StartLine: lineNum + lineOffset,
					StartCol:  match[0],
					EndLine:   lineNum + lineOffset,
					EndCol:    match[1],
				},
				Code:     diagnosticCodeUnknownMention,
				Severity: SeverityWarning,
				Message:  fmt.Sprintf("unknown mention: %s%s not found in blogroll; add a handle or alias to your blogroll", marker, handle),
				Fixable:  false,
			})
		}
	}

	return issues
}
