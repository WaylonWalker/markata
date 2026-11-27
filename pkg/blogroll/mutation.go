package blogroll

import (
	"errors"
	"fmt"
	"maps"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

// MutationAction identifies the intended change to an external subscription.
type MutationAction string

const (
	MutationAdd    MutationAction = "add"
	MutationUpdate MutationAction = "update"
	MutationRemove MutationAction = "remove"
	MutationActive MutationAction = "set-active"
)

// TOMLMutation is a frontend-neutral request. Fields contain TOML-encoded
// values produced by subscription metadata validation. Only the selected
// field names are changed for an update. The caller supplies the effective
// owner's content and whether the composed subscription list is enabled.
type TOMLMutation struct {
	Action         MutationAction
	URL            string
	Fields         map[string]string
	SelectedFields []string
	Enabled        bool
	Active         bool
	RemoveLast     bool
}

// PlanTOMLMutation produces a proposed configuration file without writing,
// loading workspace state, refreshing feeds, or contacting the network.
// The same result can be previewed in an editor, discarded or written by a
// transaction adapter that verifies the effective composed configuration.
func PlanTOMLMutation(original string, req TOMLMutation) (string, error) {
	u, err := url.Parse(req.URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" || u.User != nil {
		return "", errors.New("expected an absolute HTTP(S) feed URL without credentials")
	}
	var proposed string
	switch req.Action {
	case MutationAdd:
		proposed, err = planAddSubscription(original, req)
	case MutationUpdate, MutationActive, MutationRemove:
		proposed, err = planExistingSubscription(original, req)
	default:
		return "", errors.New("invalid subscription action")
	}
	if err != nil {
		return "", err
	}
	if err := verifyMutationDocument(original, proposed, req); err != nil {
		return "", err
	}
	return proposed, nil
}

func blogrollTablePrefix(content string) string {
	if strings.Contains(content, "[blogroll]") || strings.Contains(content, "[[blogroll.feeds]]") {
		return "blogroll"
	}
	return "markata-go.blogroll"
}

const mutationFieldActive = "active"

var mutationFields = []string{"title", "description", "category", "type", "tags", "handle", mutationFieldActive}

func planAddSubscription(original string, req TOMLMutation) (string, error) {
	content := ClearEmptySubscriptionList(original)
	if !req.Enabled {
		content = EnableSubscriptions(content)
	}
	var block strings.Builder
	block.WriteString("\n\n[[" + blogrollTablePrefix(content) + ".feeds]]\n")
	block.WriteString("url = " + strconv.Quote(req.URL) + "\n")
	for _, name := range mutationFields {
		value, ok := req.Fields[name]
		if !ok {
			return "", fmt.Errorf("missing subscription field %q", name)
		}
		block.WriteString(name + " = " + value + "\n")
	}
	return strings.TrimRight(content, "\r\n") + block.String(), nil
}

func planExistingSubscription(content string, req TOMLMutation) (string, error) {
	selected, err := UniqueSubscriptionBlock(content, req.URL)
	if err != nil {
		return "", err
	}
	replacement := ""
	switch req.Action {
	case MutationUpdate:
		names := req.SelectedFields
		if len(names) == 0 {
			names = mutationFields
		}
		replacement, err = PatchSubscriptionBlock(content[selected.Start:selected.End], req.Fields, names)
	case MutationActive:
		replacement, err = PatchSubscriptionBlock(content[selected.Start:selected.End], map[string]string{mutationFieldActive: strconv.FormatBool(req.Active)}, []string{mutationFieldActive})
	}
	if err != nil {
		return "", err
	}
	proposed := content[:selected.Start] + replacement + content[selected.End:]
	if req.Action == MutationRemove && req.RemoveLast {
		proposed = EmptySubscriptionList(proposed, blogrollTablePrefix(content))
	}
	return proposed, nil
}

var feedHeader = regexp.MustCompile(`(?m)^[ \t]*\[\[(?:markata-go\.)?blogroll\.feeds\]\][ \t]*(?:\r?\n|$)`)
var anyTableHeader = regexp.MustCompile(`(?m)^[ \t]*\[`)
var feedURLLine = regexp.MustCompile(`(?m)^[ \t]*url[ \t]*=[ \t]*(?:"([^"]+)"|'([^']+)')[ \t]*(?:#.*)?$`)
var blogrollTable = regexp.MustCompile(`(?m)^[ \t]*\[(?:markata-go\.)?blogroll\][ \t]*(?:\r?\n|$)`)

// SubscriptionBlock identifies the byte range of a feed table.
type SubscriptionBlock struct {
	Start, End int
	URL        string
}

// SubscriptionBlocks locates feed table boundaries in a TOML document.
func SubscriptionBlocks(content string) []SubscriptionBlock {
	headers := feedHeader.FindAllStringIndex(content, -1)
	result := make([]SubscriptionBlock, 0, len(headers))
	for _, header := range headers {
		end := len(content)
		if following := anyTableHeader.FindStringIndex(content[header[1]:]); following != nil {
			end = header[1] + following[0]
		}
		url := ""
		if match := feedURLLine.FindStringSubmatch(content[header[1]:end]); match != nil {
			url = match[1]
			if url == "" {
				url = match[2]
			}
		}
		result = append(result, SubscriptionBlock{header[0], end, url})
	}
	return result
}

// UniqueSubscriptionBlock refuses missing or ambiguous feed table matches.
func UniqueSubscriptionBlock(content, address string) (SubscriptionBlock, error) {
	var selected SubscriptionBlock
	found := 0
	for _, block := range SubscriptionBlocks(content) {
		if block.URL == address {
			selected = block
			found++
		}
	}
	if found != 1 {
		return selected, fmt.Errorf("cannot safely identify subscription in its TOML file (%d matching blocks)", found)
	}
	return selected, nil
}

// SetSubscriptionField sets a simple field in a table body.
func SetSubscriptionField(body, name, value string) string {
	return replaceSubscriptionField(body, name, value)
}

// EnableSubscriptions enables the blogroll table in a TOML document.
func EnableSubscriptions(content string) string {
	if loc := blogrollTable.FindStringIndex(content); loc != nil {
		end := len(content)
		if next := anyTableHeader.FindStringIndex(content[loc[1]:]); next != nil {
			end = loc[1] + next[0]
		}
		body := SetSubscriptionField(content[loc[1]:end], "enabled", "true")
		return content[:loc[1]] + body + content[end:]
	}
	return strings.TrimRight(content, "\r\n") + "\n\n[" + blogrollTablePrefix(content) + "]\nenabled = true\n"
}

// EmptySubscriptionList preserves an explicit empty array at the effective
// owner so deleting the final feed cannot reveal a shadowed lower layer.
func EmptySubscriptionList(content, prefix string) string {
	if loc := blogrollTable.FindStringIndex(content); loc != nil {
		return content[:loc[1]] + "feeds = []\n" + content[loc[1]:]
	}
	return strings.TrimRight(content, "\r\n") + "\n\n[" + prefix + "]\nfeeds = []\n"
}

// ClearEmptySubscriptionList removes the empty-array sentinel before appending
// a TOML array-of-tables, which cannot coexist with an inline array.
func ClearEmptySubscriptionList(content string) string {
	loc := blogrollTable.FindStringIndex(content)
	if loc == nil {
		return content
	}
	end := len(content)
	if next := anyTableHeader.FindStringIndex(content[loc[1]:]); next != nil {
		end = loc[1] + next[0]
	}
	empty := regexp.MustCompile(`(?m)^[ \t]*feeds[ \t]*=[ \t]*\[[ \t]*\][ \t]*(#[^\r\n]*)?(?:\r?\n|$)`)
	body := empty.ReplaceAllString(content[loc[1]:end], "$1\n")
	return content[:loc[1]] + body + content[end:]
}

// PatchSubscriptionBlock applies only the selected fields and checks that the
// resulting TOML differs semantically only by those fields. Unrecognized
// layouts (such as multiline values that cannot safely be replaced) fail
// before the caller writes anything.
func PatchSubscriptionBlock(block string, fields map[string]string, names []string) (string, error) {
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
