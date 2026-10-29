package diagnostics

import (
	"fmt"
	"strings"
)

const diagnosticCodeBareDomainLink = "bare-domain-link"

// checkBareDomainLinks warns when inline Markdown links look like external
// hostnames but will resolve relative to the current page. Never auto-fix:
// dotted paths can be valid local files.
func checkBareDomainLinks(filePath, body string, hasFrontmatter bool, frontmatter string) []Issue {
	var issues []Issue
	lineOffset := 0
	if hasFrontmatter {
		lineOffset = frontmatterLineOffset(frontmatter)
	}
	var fenceChar byte
	fenceWidth := 0
	for lineNum, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimLeft(line, " ")
		indent := len(line) - len(trimmed)
		var skip bool
		fenceChar, fenceWidth, skip = advanceBareDomainFence(trimmed, indent, fenceChar, fenceWidth)
		if skip || indent >= 4 {
			continue
		}
		issues = append(issues, scanBareDomainLine(filePath, line, lineNum+lineOffset)...)
	}
	return issues
}

// advanceBareDomainFence tracks Markdown fenced blocks. The final bool
// indicates whether this line is inside a fence or is a fence delimiter.
func advanceBareDomainFence(trimmed string, indent int, char byte, minWidth int) (nextChar byte, nextWidth int, skip bool) {
	if indent > 3 || len(trimmed) < 3 {
		return char, minWidth, char != 0
	}
	marker := trimmed[0]
	if marker != 96 && marker != '~' {
		return char, minWidth, char != 0
	}
	width := leadingRun(trimmed, marker)
	if width < 3 {
		return char, minWidth, char != 0
	}
	if char == 0 {
		return marker, width, true
	}
	if marker == char && width >= minWidth && strings.TrimSpace(trimmed[width:]) == "" {
		return 0, 0, true
	}
	return char, minWidth, true
}

func scanBareDomainLine(filePath, line string, lineNum int) []Issue {
	var issues []Issue
	visible := hideInlineCode(line)
	for i := 0; i < len(visible); i++ {
		if visible[i] != '[' || (i > 0 && (visible[i-1] == '!' || visible[i-1] == '\\')) {
			continue
		}
		closeLabel := findMarkdownLabelClose(visible, i)
		if closeLabel < 0 {
			continue
		}
		start, end, ok := findMarkdownDestination(visible, closeLabel)
		if !ok {
			continue
		}
		destination := line[start:end]
		if isLikelyBareDomain(destination) {
			issues = append(issues, Issue{
				File:     filePath,
				Range:    Range{StartLine: lineNum, StartCol: start, EndLine: lineNum, EndCol: end},
				Code:     diagnosticCodeBareDomainLink,
				Severity: SeverityWarning,
				Message:  fmt.Sprintf("Markdown link destination %q looks like a bare domain; use an explicit https:// URL if external", destination),
				Fixable:  false,
			})
		}
		i = closeLabel
	}
	return issues
}

func findMarkdownLabelClose(line string, start int) int {
	depth := 1
	for i := start + 1; i < len(line); i++ {
		switch line[i] {
		case '\\':
			i++
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func findMarkdownDestination(line string, closeLabel int) (startIndex, endIndex int, ok bool) {
	if closeLabel+1 >= len(line) || line[closeLabel+1] != '(' {
		return 0, 0, false
	}
	start := closeLabel + 2
	for start < len(line) && (line[start] == ' ' || line[start] == '\t') {
		start++
	}
	var end, after int
	switch {
	case start < len(line) && line[start] == '<':
		start++
		end = start
		for end < len(line) && line[end] != '>' {
			end++
		}
		if end == len(line) {
			return 0, 0, false
		}
		after = end + 1
	default:
		end = scanMarkdownDestinationEnd(line, start)
		after = end
	}
	if end <= start || !hasLinkClosingParen(line, after) {
		return 0, 0, false
	}
	return start, end, true
}

func scanMarkdownDestinationEnd(line string, start int) int {
	nesting := 0
	for i := start; i < len(line); i++ {
		switch line[i] {
		case '\\':
			i++
		case '(':
			nesting++
		case ')':
			if nesting == 0 {
				return i
			}
			nesting--
		case ' ', '\t':
			if nesting == 0 {
				return i
			}
		}
	}
	return len(line)
}

func leadingRun(s string, char byte) int {
	n := 0
	for n < len(s) && s[n] == char {
		n++
	}
	return n
}

// Mask code spans without changing byte positions for source diagnostics.
func hideInlineCode(line string) string {
	visible := []byte(line)
	for i := 0; i < len(line); {
		if line[i] == '\\' && i+1 < len(line) {
			i += 2
			continue
		}
		if line[i] != '`' {
			i++
			continue
		}
		width := leadingRun(line[i:], '`')
		closing := -1
		for j := i + width; j < len(line); {
			if line[j] != '`' {
				j++
				continue
			}
			run := leadingRun(line[j:], '`')
			if run == width {
				closing = j + run
				break
			}
			j += run
		}
		if closing < 0 {
			i += width
			continue
		}
		for j := i; j < closing; j++ {
			visible[j] = ' '
		}
		i = closing
	}
	return string(visible)
}

// Accept a normal closing parenthesis, optionally following a quoted title.
func hasLinkClosingParen(line string, pos int) bool {
	for pos < len(line) && (line[pos] == ' ' || line[pos] == '\t') {
		pos++
	}
	if pos >= len(line) {
		return false
	}
	if line[pos] == ')' {
		return true
	}
	quote := line[pos]
	if quote != '"' && quote != '\'' {
		return false
	}
	pos++
	for pos < len(line) && line[pos] != quote {
		if line[pos] == '\\' {
			pos++
		}
		pos++
	}
	if pos >= len(line) {
		return false
	}
	pos++
	for pos < len(line) && (line[pos] == ' ' || line[pos] == '\t') {
		pos++
	}
	return pos < len(line) && line[pos] == ')'
}

func isLikelyBareDomain(destination string) bool {
	if destination == "" || strings.ContainsAny(destination, ":\\@%") ||
		strings.HasPrefix(destination, ".") || strings.HasPrefix(destination, "/") ||
		strings.HasPrefix(destination, "#") || strings.HasPrefix(destination, "?") {
		return false
	}
	host := destination
	if at := strings.IndexAny(host, "/?#"); at >= 0 {
		host = host[:at]
	}
	labels := strings.Split(host, ".")
	if len(labels) < 2 {
		return false
	}
	for _, label := range labels {
		if !validBareDomainLabel(label) {
			return false
		}
	}
	tld := strings.ToLower(labels[len(labels)-1])
	if len(tld) < 2 || len(tld) > 24 || isCommonLocalExtension(tld) {
		return false
	}
	for i := 0; i < len(tld); i++ {
		if tld[i] < 'a' || tld[i] > 'z' {
			return false
		}
	}
	return true
}

func validBareDomainLabel(label string) bool {
	if label == "" || label[0] == '-' || label[len(label)-1] == '-' {
		return false
	}
	for i := 0; i < len(label); i++ {
		c := label[i]
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-') {
			return false
		}
	}
	return true
}

func isCommonLocalExtension(tld string) bool {
	switch tld {
	case "md", "markdown", "html", "htm", "txt", "json", "yaml", "yml", "toml", "xml", "css", "js", "ts", "jsx", "tsx", "png", "jpg", "jpeg", "gif", "svg", "webp", "avif", "ico", "pdf", "zip", "gz", "tar", "woff", "woff2", "ttf", "mp4", "webm", "mov", "wasm", "go", "py", "sh", "rs", "lock", "csv", "ipynb":
		return true
	}
	return false
}
