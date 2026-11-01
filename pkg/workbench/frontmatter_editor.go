package workbench

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"gopkg.in/yaml.v3"
)

// Frontmatter operations are in-memory transformations, not workspace writes.
// The existing editor change/autosave path remains the only writer.
type frontmatterRequest struct {
	Source string `json:"source"`
	Field  string `json:"field"`
	Value  any    `json:"value"`
	Remove bool   `json:"remove"`
}

type frontmatterProperty struct {
	Name  string `json:"name"`
	Kind  string `json:"kind"`
	Value any    `json:"value"`
}

var frontmatterFieldName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_-]*$`)

func frontmatterDocument(source string) (diagnostics.FrontmatterInspection, *yaml.Node, error) {
	inspection := diagnostics.InspectFrontmatter(source)
	if inspection.Problem != "" || inspection.LeadingWhitespace || inspection.SuspiciousDelimiter != "" {
		return inspection, nil, fmt.Errorf("frontmatter delimiters are malformed; fix the YAML in source mode")
	}
	document := &yaml.Node{}
	if !inspection.HasFrontmatter || strings.TrimSpace(inspection.Frontmatter) == "" {
		return inspection, document, nil
	}
	if err := yaml.Unmarshal([]byte(inspection.Frontmatter), document); err != nil {
		return inspection, nil, fmt.Errorf("invalid YAML frontmatter: %w", err)
	}
	if len(document.Content) == 0 { // comment-only YAML is a valid empty mapping
		return inspection, document, nil
	}
	if len(document.Content) != 1 || document.Content[0].Kind != yaml.MappingNode {
		return inspection, nil, fmt.Errorf("frontmatter must be a YAML object (mapping)")
	}
	// Decode also rejects duplicate keys; the syntax tree alone does not.
	var values map[string]any
	if err := document.Decode(&values); err != nil {
		return inspection, nil, fmt.Errorf("invalid YAML frontmatter: %w", err)
	}
	for i := 0; i < len(document.Content[0].Content); i += 2 {
		if document.Content[0].Content[i].Kind != yaml.ScalarNode {
			return inspection, nil, fmt.Errorf("frontmatter property names must be strings")
		}
	}
	return inspection, document, nil
}

func frontmatterProperties(document *yaml.Node) ([]frontmatterProperty, error) {
	fields := []frontmatterProperty{}
	if len(document.Content) == 0 {
		return fields, nil
	}
	content := document.Content[0].Content
	for i := 0; i+1 < len(content); i += 2 {
		key, value := content[i], content[i+1]
		var decoded any
		if value.Kind == yaml.ScalarNode && value.Tag == "!!timestamp" {
			// Keep author-supplied date formatting rather than exposing a UTC JSON timestamp.
			decoded = value.Value
		} else if err := value.Decode(&decoded); err != nil {
			return nil, err
		}
		kind := "string"
		switch value.Kind {
		case yaml.SequenceNode:
			kind = "list"
		case yaml.MappingNode:
			kind = "object"
		case yaml.AliasNode:
			kind = "alias"
		case yaml.ScalarNode:
			switch value.Tag {
			case "!!bool":
				kind = "boolean"
			case "!!int", "!!float":
				kind = "number"
			case "!!null":
				kind = "null"
			}
		}
		fields = append(fields, frontmatterProperty{Name: key.Value, Value: decoded, Kind: kind})
	}
	return fields, nil
}

func marshalFrontmatterField(name string, value any) (string, error) {
	var out strings.Builder
	encoder := yaml.NewEncoder(&out)
	encoder.SetIndent(2)
	if err := encoder.Encode(map[string]any{name: value}); err != nil {
		return "", err
	}
	if err := encoder.Close(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func replaceFrontmatterProperty(source, name string, value any, remove bool) (string, error) {
	if !frontmatterFieldName.MatchString(name) || len(name) > 128 {
		return "", fmt.Errorf("invalid frontmatter property name")
	}
	inspection, document, err := frontmatterDocument(source)
	if err != nil {
		return "", err
	}
	if !inspection.HasFrontmatter {
		if remove {
			return source, nil
		}
		encoded, err := marshalFrontmatterField(name, value)
		if err != nil {
			return "", err
		}
		newline := "
"
		if strings.Contains(source, "
") {
			newline = "
"
		}
		next := "---" + newline + strings.ReplaceAll(string(encoded), "
", newline) + "---" + newline + source
		if _, _, err := frontmatterDocument(next); err != nil {
			return "", err
		}
		return next, nil
	}

	// Keep all unchanged YAML and all body bytes exactly as authored, including
	// comments, quoting, key order, and mixed line endings. Re-encode only the
	// property the user changed.
	lines := strings.SplitAfter(source, "
")
	closeIndex := inspection.ClosingLine - 1
	content := []*yaml.Node{}
	if len(document.Content) > 0 {
		content = document.Content[0].Content
	}
	start, end := -1, -1
	for i := 0; i+1 < len(content); i += 2 {
		if content[i].Value != name {
			continue
		}
		start, end = content[i].Line, closeIndex
		if i+2 < len(content) {
			end = content[i+2].Line
		}
		break
	}
	var encoded string
	if !remove {
		raw, err := marshalFrontmatterField(name, value)
		if err != nil {
			return "", fmt.Errorf("cannot encode property: %w", err)
		}
		encoded = raw
		if strings.Contains(source, "
") {
			encoded = strings.ReplaceAll(encoded, "
", "
")
		}
	}
	if start == -1 {
		if remove {
			return source, nil
		}
		start, end = closeIndex, closeIndex
	} else {
		// Preserve standalone comments and blank lines following a property.
		for end > start+1 {
			last := strings.TrimSpace(lines[end-1])
			if last != "" && !strings.HasPrefix(lines[end-1], "#") {
				break
			}
			end--
		}
	}
	next := strings.Join(lines[:start], "") + encoded + strings.Join(lines[end:], "")
	if _, _, err := frontmatterDocument(next); err != nil {
		return "", fmt.Errorf("edited property would produce invalid YAML: %w", err)
	}
	return next, nil
}

func (h *Handler) serveFrontmatterEditor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || !sameOriginWorkspaceMutation(r) {
		fail(w, http.StatusForbidden, "same-origin POST required")
		return
	}
	var request frontmatterRequest
	decoder := json.NewDecoder(io.LimitReader(r.Body, 2*maxSource+4096))
	if err := decoder.Decode(&request); err != nil {
		fail(w, http.StatusBadRequest, "invalid editor request")
		return
	}
	if len(request.Source) > maxSource {
		fail(w, http.StatusRequestEntityTooLarge, "document exceeds 200,000 bytes")
		return
	}
	if r.URL.Path == "/workspace/frontmatter/inspect" {
		inspection, document, err := frontmatterDocument(request.Source)
		if err != nil {
			fail(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		fields, err := frontmatterProperties(document)
		if err != nil {
			fail(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"has_frontmatter": inspection.HasFrontmatter,
			"fields": fields,
		})
		return
	}
	if r.URL.Path == "/workspace/frontmatter/edit" {
		source, err := replaceFrontmatterProperty(request.Source, request.Field, request.Value, request.Remove)
		if err != nil {
			fail(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"source": source})
		return
	}
	fail(w, http.StatusNotFound, "unknown frontmatter operation")
}
