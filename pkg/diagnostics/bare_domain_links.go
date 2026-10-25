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
	fenceChar := byte(0)
	fenceWidth := 0
	for lineNum, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimLeft(line, " ")
		indent := len(line) - len(trimmed)
		if indent <= 3 && len(trimmed) >= 3 && (trimmed[0] == '`' || trimmed[0] == '~') {
			width := leadingRun(trimmed, trimmed[0])
			if width >= 3 {
				if fenceChar == 0 {
					fenceChar, fenceWidth = trimmed[0], width
					continue
				}
				if trimmed[0] == fenceChar && width >= fenceWidth && strings.TrimSpace(trimmed[width:]) == "" {
					fenceChar, fenceWidth = 0, 0
				}
				continue
			}
		}
		if fenceChar != 0 || indent >= 4 {
			continue
		}
		visible := hideInlineCode(line)
		for i := 0; i < len(visible); i++ {
			if visible[i] != '[' || (i > 0 && (visible[i-1] == '!' || visible[i-1] == '\\')) {
				continue
			}
			depth, closeLabel := 1, -1
			for j := i + 1; j < len(visible); j++ {
				if visible[j] == '\\' {
					j++
					continue
				}
				if visible[j] == '[' {
					depth++
				} else if visible[j] == ']' {
					depth--
					if depth == 0 {
						closeLabel = j
						break
					}
				}
			}
			if closeLabel < 0 || closeLabel+1 >= len(visible) || visible[closeLabel+1] != '(' {
				continue
			}
			start := closeLabel + 2
			for start < len(visible) && (visible[start] == ' ' || visible[start] == '\t') {
				start++
			}
			end := start
			after := start
			if start < len(visible) && visible[start] == '<' {
				start++
				end = start
				for end < len(visible) && visible[end] != '>' {
					end++
				}
				if end == len(visible) {
					continue
				}
				after = end + 1
			} else {
				nesting := 0
				for end < len(visible) {
					if visible[end] == '\\' {
						end += 2
						continue
					}
					if visible[end] == '(' {
						nesting++
					} else if visible[end] == ')' {
						if nesting == 0 {
							break
						}
						nesting--
					} else if (visible[end] == ' ' || visible[end] == '\t') && nesting == 0 {
						break
					}
					end++
				}
				after = end
			}
			if end <= start || end > len(line) || !hasLinkClosingParen(visible, after) {
				continue
			}
			destination := line[start:end]
			if isLikelyBareDomain(destination) {
				issues = append(issues, Issue{
					File: filePath,
					Range: Range{StartLine: lineNum + lineOffset, StartCol: start, EndLine: lineNum + lineOffset, EndCol: end},
					Code: diagnosticCodeBareDomainLink,
					Severity: SeverityWarning,
					Message: fmt.Sprintf("Markdown link destination %q looks like a bare domain; use an explicit https:// URL if external", destination),
					Fixable: false,
				})
			}
			i = closeLabel
		}
	}
	return issues
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
		if label == "" || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for i := 0; i < len(label); i++ {
			c := label[i]
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-') {
				return false
			}
		}
	}
	tld := strings.ToLower(labels[len(labels)-1])
	if len(tld) < 2 || len(tld) > 24 {
		return false
	}
	for i := 0; i < len(tld); i++ {
		if tld[i] < 'a' || tld[i] > 'z' {
			return false
		}
	}
	// Avoid warnings for ordinary relative filenames.
	switch tld {
	case "md", "markdown", "html", "htm", "txt", "json", "yaml", "yml", "toml", "xml", "css", "js", "ts", "jsx", "tsx", "png", "jpg", "jpeg", "gif", "svg", "webp", "avif", "ico", "pdf", "zip", "gz", "tar", "woff", "woff2", "ttf", "mp4", "webm", "mov", "wasm", "go", "py", "sh", "rs", "lock", "csv", "ipynb":
		return false
	}
	return true
}
