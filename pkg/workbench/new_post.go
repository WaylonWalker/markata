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
	lines := strings.Split(source[match[2]:match[3]], "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, "image:") {
			lines[i] = quoted
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, quoted)
	}
	return source[:match[2]] + strings.Join(lines, "\n") + source[match[3]:]
}
