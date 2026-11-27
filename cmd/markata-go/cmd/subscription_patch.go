package cmd

import (
	"fmt"
	"maps"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/WaylonWalker/markata/pkg/models"
)

// patchSubscriptionBlock applies only the selected fields and checks that the
// resulting TOML differs semantically only by those fields. Unrecognized
// layouts (such as multiline values that cannot safely be replaced) fail
// before the caller writes anything.
func patchSubscriptionBlock(block string, fields map[string]string, names []string) (string, error) {
	headerEnd := strings.IndexByte(block, '\n')
	if headerEnd < 0 {
		return "", fmt.Errorf("subscription has no field body")
	}
	header, body := block[:headerEnd+1], block[headerEnd+1:]
	var before map[string]any
	if _, err := toml.Decode(body, &before); err != nil {
		return "", fmt.Errorf("read subscription fields: %w", err)
	}
	expected := maps.Clone(before)
	for _, name := range names {
		value, ok := fields[name]
		if !ok {
			return "", fmt.Errorf("unsupported subscription field %q", name)
		}
		var parsed map[string]any
		if _, err := toml.Decode(name+" = "+value, &parsed); err != nil {
			return "", fmt.Errorf("invalid subscription field %q: %w", name, err)
		}
		expected[name] = parsed[name]
		body = replaceSubscriptionField(body, name, value)
	}
	var after map[string]any
	if _, err := toml.Decode(body, &after); err != nil {
		return "", fmt.Errorf("cannot safely edit subscription; edit its TOML file directly: %w", err)
	}
	if !reflect.DeepEqual(expected, after) {
		return "", fmt.Errorf("cannot safely edit subscription without changing unrelated fields; edit its TOML file directly")
	}
	return header + body, nil
}

// Preserve indentation, trailing inline comments and original line endings
// when changing a single-line TOML value.
func replaceSubscriptionField(body, name, value string) string {
	pattern := regexp.MustCompile(`(?m)^([ \t]*` + regexp.QuoteMeta(name) + `[ \t]*=[ \t]*)([^\r\n]*)`)
	if !pattern.MatchString(body) {
		return strings.TrimRight(body, "\r\n") + "\n" + name + " = " + value + "\n"
	}
	return pattern.ReplaceAllStringFunc(body, func(line string) string {
		parts := pattern.FindStringSubmatch(line)
		return parts[1] + value + subscriptionInlineComment(parts[2])
	})
}

// subscriptionInlineComment locates a comment outside single or double
// quoted strings, respecting escapes in TOML basic strings.
func subscriptionInlineComment(value string) string {
	var quote byte
	escaped := false
	for i := 0; i < len(value); i++ {
		ch := value[i]
		switch {
		case escaped:
			escaped = false
		case quote == '"' && ch == '\\':
			escaped = true
		case quote != 0:
			if ch == quote {
				quote = 0
			}
		case ch == '"' || ch == '\'':
			quote = ch
		case ch == '#':
			start := i
			for start > 0 && (value[start-1] == ' ' || value[start-1] == '\t') {
				start--
			}
			return value[start:]
		}
	}
	return ""
}

func subscriptionUpdateMatches(feed *models.ExternalFeedConfig, req subscriptionRequest) bool {
	names := req.selectedFields
	if len(names) == 0 {
		names = []string{"title", "description", "category", "type", "tags", "handle", "active"}
	}
	actual := map[string]bool{
		"title":       feed.Title == strings.TrimSpace(req.Title),
		"description": feed.Description == strings.TrimSpace(req.Description),
		"category":    feed.Category == strings.TrimSpace(req.Category),
		"type":        string(feed.Type) == req.Type,
		"tags":        slices.Equal(feed.Tags, parseTagsBlogroll(req.Tags)),
		"handle":      feed.Handle == strings.TrimSpace(req.Handle),
		"active":      feed.IsActive() == req.Active,
	}
	for _, name := range names {
		if !actual[name] {
			return false
		}
	}
	return true
}
