package workbench

import (
	"sort"
	"strconv"
	"strings"
)

type sourceContainerSpan struct {
	line, width, endLine int
}

// markSourceContainers adds workbench-only source-line attributes to Markdown
// container openers. It preserves the source line count and ignores code fences.
func markSourceContainers(source string) (marked string, anchors []int) {
	lines := strings.Split(source, "\n")
	containers := sourceContainerSpans(lines)
	anchors = make([]int, 0, len(containers)*2)
	for _, item := range containers {
		attrs := make([]string, 0, 2)
		if !strings.Contains(lines[item.line], "data-markata-source-line=") {
			attrs = append(attrs, `data-markata-source-line="`+strconv.Itoa(item.line)+`"`)
		}
		anchors = append(anchors, item.line)
		if item.endLine >= 0 {
			anchors = append(anchors, item.endLine)
			if !strings.Contains(lines[item.line], "data-markata-source-end-line=") {
				attrs = append(attrs, `data-markata-source-end-line="`+strconv.Itoa(item.endLine)+`"`)
			}
		}
		if len(attrs) > 0 {
			lines[item.line] = appendContainerAttributes(lines[item.line], attrs)
		}
	}
	sort.Ints(anchors)
	anchors = compactLines(anchors)
	return strings.Join(lines, "\n"), anchors
}

func sourceContainerSpans(lines []string) []sourceContainerSpan {
	containers := make([]sourceContainerSpan, 0)
	stack := make([]int, 0)
	fence := byte(0)
	fenceWidth := 0

	for i, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		if len(line)-len(trimmed) > 3 || trimmed == "" {
			continue
		}

		if marker, width := markdownFence(trimmed); width >= 3 {
			if fence == 0 {
				fence, fenceWidth = marker, width
			} else if marker == fence && width >= fenceWidth && strings.TrimSpace(trimmed[width:]) == "" {
				fence, fenceWidth = 0, 0
			}
			continue
		}
		if fence != 0 {
			continue
		}

		if closeWidth := containerCloseWidth(trimmed); closeWidth >= 3 {
			stack = closeSourceContainers(containers, stack, closeWidth, i)
			continue
		}
		if !strings.HasPrefix(trimmed, ":::") {
			continue
		}

		width := 0
		for width < len(trimmed) && trimmed[width] == ':' {
			width++
		}
		if width < 3 || strings.TrimSpace(trimmed[width:]) == "" {
			continue
		}
		containers = append(containers, sourceContainerSpan{line: i, width: width, endLine: -1})
		stack = append(stack, len(containers)-1)
	}
	return containers
}

func closeSourceContainers(containers []sourceContainerSpan, stack []int, closeWidth, line int) []int {
	for len(stack) > 0 {
		index := stack[len(stack)-1]
		if closeWidth > containers[index].width {
			break
		}
		containers[index].endLine = line
		stack = stack[:len(stack)-1]
		if closeWidth == containers[index].width {
			break
		}
	}
	return stack
}

func appendContainerAttributes(line string, attributes []string) string {
	lineEnd := len(strings.TrimRight(line, " \t\r"))
	if lineEnd > 0 && line[lineEnd-1] == '}' {
		openBrace := strings.LastIndex(line[:lineEnd], "{")
		if openBrace >= 0 {
			return line[:lineEnd-1] + " " + strings.Join(attributes, " ") + line[lineEnd-1:]
		}
	}
	return line[:lineEnd] + " {" + strings.Join(attributes, " ") + "}" + line[lineEnd:]
}

func containerCloseWidth(line string) int {
	if strings.TrimSpace(strings.TrimLeft(line, ":")) != "" {
		return 0
	}
	width := 0
	for width < len(line) && line[width] == ':' {
		width++
	}
	return width
}

func compactLines(lines []int) []int {
	if len(lines) < 2 {
		return lines
	}
	write := 1
	for read := 1; read < len(lines); read++ {
		if lines[read] == lines[write-1] {
			continue
		}
		lines[write] = lines[read]
		write++
	}
	return lines[:write]
}

func markdownFence(line string) (marker byte, width int) {
	if line == "" || (line[0] != '`' && line[0] != '~') {
		return 0, 0
	}
	for width < len(line) && line[width] == line[0] {
		width++
	}
	return line[0], width
}
