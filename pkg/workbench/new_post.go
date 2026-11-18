package workbench

import (
	"regexp"
	"strconv"
	"strings"
)

var newPostFrontmatterPattern = regexp.MustCompile(`(?s)^(---\r?\n)(.*?)(\r?\n---(?:\r?\n|$))`)

func sourceWithImage(source, imageURL string) string {
	quoted := "image: " + strconv.Quote(imageURL)
	match := newPostFrontmatterPattern.FindStringSubmatchIndex(source)
	if match == nil {
		return "---\n" + quoted + "\n---\n\n" + source
	}
	frontmatterStart, frontmatterEnd := match[4], match[5]
	lines := strings.Split(source[frontmatterStart:frontmatterEnd], "\n")
	found := false
	for i, line := range lines {
		lineEnding := ""
		if strings.HasSuffix(line, "\r") {
			lineEnding = "\r"
			line = strings.TrimSuffix(line, "\r")
		}
		if strings.HasPrefix(line, "image:") {
			lines[i] = quoted + lineEnding
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, quoted)
	}
	return source[:frontmatterStart] + strings.Join(lines, "\n") + source[frontmatterEnd:]
}
