package diagnostics

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// FrontmatterProblem identifies a structural frontmatter problem.
type FrontmatterProblem string

const frontmatterDelimiter = "---"

const (
	yamlStringTag        = "!!str"
	frontmatterTagsField = "tags"
)

const (
	// FrontmatterProblemNone means that no structural problem was found.
	FrontmatterProblemNone FrontmatterProblem = ""
	// FrontmatterProblemMissingClosing means an exact opening has no close.
	FrontmatterProblemMissingClosing FrontmatterProblem = "missing_closing_delimiter"
	// FrontmatterProblemMalformedClosing means a close was found, but is not ---.
	FrontmatterProblemMalformedClosing FrontmatterProblem = "malformed_closing_delimiter"
)

// FrontmatterInspection contains structural frontmatter observations. Line
// numbers are one-based so they can be shown directly to authors.
type FrontmatterInspection struct {
	Frontmatter         string
	Body                string
	HasFrontmatter      bool
	OpeningDelimiter    string
	ClosingDelimiter    string
	OpeningLine         int
	ClosingLine         int
	SuspiciousDelimiter string
	LeadingWhitespace   bool
	Problem             FrontmatterProblem
}

// FrontmatterAnalysis combines structural inspection, YAML parsing, and
// frontmatter-specific diagnostics.
type FrontmatterAnalysis struct {
	Inspection FrontmatterInspection
	Valid      bool
	Issues     []Issue
}

var (
	frontmatterKeyRegex       = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_-]*\s*:`)
	frontmatterDateLineRegex  = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_-]*)\s*:\s*(.+)$`)
	yamlErrorLineRegex        = regexp.MustCompile(`(?i)\bline\s+(\d+)`)
	dateKeyRegex              = regexp.MustCompile(`^(date|publishdate|pubdate|published_date|created|modified|lastmod|updated|updated_at|last_modified)$`)
	diagnosticTimeRegex       = regexp.MustCompile(`(\d{1,2}):0*(\d{1,2}):0*(\d{1,2})`)
	diagnosticSingleHourRegex = regexp.MustCompile(`([ T])(\d):(\d{2})`)
)

// InspectFrontmatter finds a strict frontmatter block and detects likely
// delimiter mistakes without treating ordinary Markdown as metadata.
func InspectFrontmatter(content string) FrontmatterInspection {
	normalized := normalizeFrontmatterContent(content)
	lines := strings.Split(normalized, "\n")
	inspection := FrontmatterInspection{Body: normalized}
	if len(lines) == 0 {
		return inspection
	}

	firstLine := lines[0]
	if firstLine == frontmatterDelimiter {
		inspection.OpeningDelimiter = firstLine
		inspection.OpeningLine = 1
		if closeIndex := findClosingDelimiter(lines, 1); closeIndex >= 0 {
			inspection.HasFrontmatter = true
			inspection.ClosingDelimiter = lines[closeIndex]
			inspection.ClosingLine = closeIndex + 1
			inspection.Frontmatter = strings.Join(lines[1:closeIndex], "\n")
			inspection.Body = strings.Join(lines[closeIndex+1:], "\n")
			return inspection
		}

		if closeIndex := findMalformedClosingDelimiter(lines, 1); closeIndex >= 0 {
			inspection.ClosingDelimiter = lines[closeIndex]
			inspection.ClosingLine = closeIndex + 1
			inspection.Problem = FrontmatterProblemMalformedClosing
			return inspection
		}

		inspection.Problem = FrontmatterProblemMissingClosing
		return inspection
	}

	trimmedFirstLine := strings.TrimSpace(firstLine)
	closeIndex := findClosingDelimiter(lines, 1)
	if closeIndex < 0 || !looksLikeYAML(lines[1:closeIndex]) {
		return inspection
	}

	if trimmedFirstLine == frontmatterDelimiter && hasLeadingWhitespace(firstLine) {
		inspection.LeadingWhitespace = true
		inspection.OpeningDelimiter = firstLine
		inspection.OpeningLine = 1
		inspection.ClosingDelimiter = lines[closeIndex]
		inspection.ClosingLine = closeIndex + 1
		return inspection
	}

	if isLongHyphenDelimiter(trimmedFirstLine) {
		inspection.SuspiciousDelimiter = firstLine
		inspection.OpeningDelimiter = firstLine
		inspection.OpeningLine = 1
		inspection.ClosingDelimiter = lines[closeIndex]
		inspection.ClosingLine = closeIndex + 1
	}

	return inspection
}

func hasLeadingWhitespace(line string) bool {
	return strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t")
}

// AnalyzeFrontmatter validates frontmatter and returns stable source issues.
// It does not diagnose ordinary body content; callers can use Check for the
// complete Markdown diagnostics set.
func AnalyzeFrontmatter(filePath, content string) FrontmatterAnalysis {
	inspection := InspectFrontmatter(content)
	analysis := FrontmatterAnalysis{Inspection: inspection}

	if inspection.SuspiciousDelimiter != "" {
		analysis.Issues = append(analysis.Issues, Issue{
			File:     filePath,
			Range:    lineRange(inspection.OpeningLine, len(inspection.OpeningDelimiter)),
			Code:     ReasonFrontmatterSuspiciousDelimiter,
			Severity: SeverityWarning,
			Message: fmt.Sprintf(
				"found suspicious frontmatter opening delimiter %q; expected %q",
				inspection.SuspiciousDelimiter,
				frontmatterDelimiter,
			),
		})
	}
	if inspection.LeadingWhitespace {
		analysis.Issues = append(analysis.Issues, Issue{
			File:     filePath,
			Range:    lineRange(inspection.OpeningLine, len(inspection.OpeningDelimiter)),
			Code:     ReasonFrontmatterLeadingWhitespace,
			Severity: SeverityWarning,
			Message:  fmt.Sprintf("frontmatter opening delimiter must start at column 1; expected %q", frontmatterDelimiter),
		})
	}

	switch inspection.Problem {
	case FrontmatterProblemNone:
		// No structural error.
	case FrontmatterProblemMissingClosing:
		analysis.Issues = append(analysis.Issues, Issue{
			File:     filePath,
			Range:    lineRange(inspection.OpeningLine, len(inspection.OpeningDelimiter)),
			Code:     ReasonFrontmatterMissingClosing,
			Severity: SeverityError,
			Message:  fmt.Sprintf("frontmatter has no closing delimiter; expected %q", frontmatterDelimiter),
		})
	case FrontmatterProblemMalformedClosing:
		analysis.Issues = append(analysis.Issues, Issue{
			File:     filePath,
			Range:    lineRange(inspection.ClosingLine, len(inspection.ClosingDelimiter)),
			Code:     ReasonFrontmatterMalformedClosing,
			Severity: SeverityError,
			Message: fmt.Sprintf(
				"found malformed frontmatter closing delimiter %q; expected %q",
				inspection.ClosingDelimiter,
				frontmatterDelimiter,
			),
		})
	}

	if !inspection.HasFrontmatter {
		return analysis
	}
	if strings.TrimSpace(inspection.Frontmatter) == "" {
		analysis.Valid = true
		return analysis
	}

	analysis.Issues = append(analysis.Issues, checkDuplicateKeys(filePath, inspection.Frontmatter)...)
	analysis.Issues = append(analysis.Issues, checkDateFormats(filePath, inspection.Frontmatter)...)

	var metadata map[string]interface{}
	if err := yaml.Unmarshal([]byte(inspection.Frontmatter), &metadata); err != nil {
		analysis.Issues = append(analysis.Issues, frontmatterParseIssue(filePath, inspection, err))
		return analysis
	}

	var document yaml.Node
	if err := yaml.Unmarshal([]byte(inspection.Frontmatter), &document); err == nil {
		analysis.Issues = append(analysis.Issues, checkMetadataTypes(filePath, inspection, &document)...)
	}

	analysis.Valid = true
	return analysis
}

func normalizeFrontmatterContent(content string) string {
	content = strings.TrimPrefix(content, "\uFEFF")
	content = strings.ReplaceAll(content, "\r\n", "\n")
	return strings.ReplaceAll(content, "\r", "\n")
}

func findClosingDelimiter(lines []string, start int) int {
	for i := start; i < len(lines); i++ {
		if lines[i] == frontmatterDelimiter {
			return i
		}
	}
	return -1
}

func findMalformedClosingDelimiter(lines []string, start int) int {
	for i := start; i < len(lines); i++ {
		line := lines[i]
		if !isMalformedDelimiter(line) {
			continue
		}
		if i == start || looksLikeYAML(lines[start:i]) {
			return i
		}
	}
	return -1
}

func looksLikeYAML(lines []string) bool {
	for _, line := range lines {
		if frontmatterKeyRegex.MatchString(line) {
			return true
		}
	}
	return false
}

func isLongHyphenDelimiter(line string) bool {
	return len(line) > len(frontmatterDelimiter) && isOnlyHyphens(line)
}

func isMalformedDelimiter(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return false
	}
	if trimmed == frontmatterDelimiter {
		return line != frontmatterDelimiter
	}
	return strings.HasPrefix(trimmed, frontmatterDelimiter) || isOnlyHyphens(trimmed)
}

func isOnlyHyphens(line string) bool {
	if line == "" {
		return false
	}
	for _, char := range line {
		if char != '-' {
			return false
		}
	}
	return true
}

func lineRange(line, width int) Range {
	if line < 1 {
		line = 1
	}
	return Range{
		StartLine: line - 1,
		EndLine:   line - 1,
		EndCol:    width,
	}
}

func frontmatterLineOffset(frontmatter string) int {
	if frontmatter == "" {
		return 2
	}
	return strings.Count(frontmatter, "\n") + 3
}

func frontmatterParseIssue(filePath string, inspection FrontmatterInspection, err error) Issue {
	line := inspection.OpeningLine
	if matches := yamlErrorLineRegex.FindStringSubmatch(err.Error()); len(matches) == 2 {
		if relativeLine, parseErr := strconv.Atoi(matches[1]); parseErr == nil {
			line = inspection.OpeningLine + relativeLine
		}
	}
	return Issue{
		File:     filePath,
		Range:    lineRange(line, 0),
		Code:     ReasonFrontmatterParseError,
		Severity: SeverityError,
		Message:  "frontmatter YAML could not be parsed",
	}
}

func checkMetadataTypes(filePath string, inspection FrontmatterInspection, document *yaml.Node) []Issue {
	mapping := document
	if mapping.Kind == yaml.DocumentNode && len(mapping.Content) > 0 {
		mapping = mapping.Content[0]
	}
	if mapping.Kind != yaml.MappingNode {
		return nil
	}

	var issues []Issue
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		key := mapping.Content[index]
		value := mapping.Content[index+1]
		field := strings.ToLower(key.Value)
		expected, ok := frontmatterFieldType(field)
		if ok && !frontmatterFieldHasType(field, value) {
			issue := invalidTypeIssue(filePath, inspection, key, field, expected)
			if field == frontmatterTagsField && value.Kind == yaml.ScalarNode && value.Tag == yamlStringTag {
				issue.Message = `frontmatter field "tags" uses scalar shorthand; use a list of strings`
			}
			issues = append(issues, issue)
		}
	}
	return issues
}

func frontmatterFieldType(field string) (string, bool) {
	switch {
	case field == "title" || field == "description" || field == "slug" || field == "template":
		return "a string", true
	case field == "published" || field == "draft" || field == "private" || field == "skip":
		return "a boolean", true
	case field == frontmatterTagsField:
		return "a list of strings", true
	case dateKeyRegex.MatchString(field):
		return "a supported date", true
	case field == "templates":
		return "a mapping", true
	default:
		return "", false
	}
}

func frontmatterFieldHasType(field string, node *yaml.Node) bool {
	switch {
	case field == "title" || field == "description" || field == "slug" || field == "template":
		return isStringOrNull(node)
	case field == "published" || field == "draft" || field == "private" || field == "skip":
		return isBooleanNode(node)
	case field == frontmatterTagsField:
		return isStringListNode(node) || isNullNode(node)
	case dateKeyRegex.MatchString(field):
		return isDateNode(node)
	case field == "templates":
		return node.Kind == yaml.MappingNode || isNullNode(node)
	default:
		return true
	}
}

func invalidTypeIssue(filePath string, inspection FrontmatterInspection, key *yaml.Node, field, expected string) Issue {
	line := inspection.OpeningLine + key.Line
	startCol := key.Column - 1
	if startCol < 0 {
		startCol = 0
	}
	return Issue{
		File: filePath,
		Range: Range{
			StartLine: line - 1,
			StartCol:  startCol,
			EndLine:   line - 1,
			EndCol:    startCol + len(key.Value),
		},
		Code:     ReasonFrontmatterInvalidType,
		Severity: SeverityWarning,
		Message:  fmt.Sprintf("frontmatter field %q must be %s", field, expected),
	}
}

func isStringOrNull(node *yaml.Node) bool {
	return isNullNode(node) || (node.Kind == yaml.ScalarNode && node.Tag == yamlStringTag)
}

func isBooleanNode(node *yaml.Node) bool {
	if node.Kind != yaml.ScalarNode || node.Tag == "!!null" {
		return false
	}
	if node.Tag == "!!bool" {
		return true
	}
	switch strings.ToLower(strings.TrimSpace(node.Value)) {
	case "true", "false", "yes", "no", "on", "off":
		return node.Tag == yamlStringTag
	default:
		return false
	}
}

func isStringListNode(node *yaml.Node) bool {
	if node.Kind != yaml.SequenceNode {
		return false
	}
	for _, item := range node.Content {
		if item.Kind != yaml.ScalarNode || item.Tag != yamlStringTag {
			return false
		}
	}
	return true
}

func isDateNode(node *yaml.Node) bool {
	if isNullNode(node) {
		return true
	}
	if node.Kind != yaml.ScalarNode {
		return false
	}
	if node.Tag == "!!timestamp" {
		return true
	}
	return node.Tag == yamlStringTag
}

func isNullNode(node *yaml.Node) bool {
	return node.Kind == yaml.ScalarNode && node.Tag == "!!null"
}

func isSupportedDateString(value string) bool {
	value = normalizeDiagnosticDate(value)
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"2006/01/02 15:04:05",
		"2006/01/02 15:04",
		"2006/01/02",
		"01/02/2006 15:04:05",
		"01/02/2006 15:04",
		"01/02/2006",
		"02-01-2006",
		"January 2, 2006",
		"Jan 2, 2006",
		"2 January 2006",
		"2 Jan 2006",
	}
	for _, format := range formats {
		if _, err := time.Parse(format, value); err == nil {
			return true
		}
	}
	return false
}

func normalizeDiagnosticDate(value string) string {
	value = strings.Trim(strings.TrimSpace(value), "\"'")
	value = diagnosticTimeRegex.ReplaceAllStringFunc(value, func(match string) string {
		parts := diagnosticTimeRegex.FindStringSubmatch(match)
		if len(parts) != 4 {
			return match
		}
		hour, hourErr := strconv.Atoi(parts[1])
		minute, minuteErr := strconv.Atoi(parts[2])
		second, secondErr := strconv.Atoi(parts[3])
		if hourErr != nil || minuteErr != nil || secondErr != nil {
			return match
		}
		return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
	})
	return diagnosticSingleHourRegex.ReplaceAllString(value, "${1}0${2}:${3}")
}
