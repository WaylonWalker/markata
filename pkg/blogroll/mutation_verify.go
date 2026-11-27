package blogroll

import (
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"

	"github.com/BurntSushi/toml"
)

// mutationDocument separates the feed array from every other parsed TOML
// value. Comparing these snapshots catches edits that cross a table boundary
// or accidentally modify a neighboring config section.
type mutationDocument struct {
	other   map[string]any
	feeds   []map[string]any
	present bool
	enabled any
}

func readMutationDocument(content string, changeEnabled bool) (mutationDocument, error) {
	var raw map[string]any
	if _, err := toml.Decode(content, &raw); err != nil {
		return mutationDocument{}, fmt.Errorf("invalid subscription TOML: %w", err)
	}
	doc := mutationDocument{other: raw}
	root := raw
	nested, hasNested := raw["markata-go"].(map[string]any)
	if hasNested {
		root = nested
	}
	section, ok := root["blogroll"].(map[string]any)
	if !ok {
		return doc, nil
	}
	doc.enabled = section["enabled"]
	value, present := section["feeds"]
	doc.present = present
	if present {
		feeds, err := mutationFeedTables(value)
		if err != nil {
			return mutationDocument{}, err
		}
		doc.feeds = feeds
	}
	delete(section, "feeds")
	if changeEnabled {
		delete(section, "enabled")
	}
	if len(section) == 0 {
		delete(root, "blogroll")
	}
	if hasNested && len(root) == 0 {
		delete(raw, "markata-go")
	}
	return doc, nil
}

func mutationFeedTables(value any) ([]map[string]any, error) {
	switch entries := value.(type) {
	case []map[string]any:
		return entries, nil
	case []any:
		result := make([]map[string]any, 0, len(entries))
		for _, entry := range entries {
			feed, ok := entry.(map[string]any)
			if !ok {
				return nil, errors.New("subscription array contains a non-table entry")
			}
			result = append(result, feed)
		}
		return result, nil
	default:
		return nil, errors.New("subscription configuration must contain an array of feeds")
	}
}

func verifyMutationDocument(original, proposed string, req TOMLMutation) error {
	changeEnabled := req.Action == MutationAdd && !req.Enabled
	before, err := readMutationDocument(original, changeEnabled)
	if err != nil {
		return err
	}
	after, err := readMutationDocument(proposed, changeEnabled)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(before.other, after.other) {
		return errors.New("subscription change would modify unrelated TOML values")
	}
	if changeEnabled && after.enabled != true {
		return errors.New("subscription change did not enable the blogroll")
	}
	expected, err := expectedMutationFeeds(before.feeds, req)
	if err != nil {
		return err
	}
	if !after.present || !slices.EqualFunc(expected, after.feeds, func(a, b map[string]any) bool {
		return reflect.DeepEqual(a, b)
	}) {
		return errors.New("subscription change would modify unexpected feed values; edit the owning TOML file directly")
	}
	return nil
}

func expectedMutationFeeds(feeds []map[string]any, req TOMLMutation) ([]map[string]any, error) {
	index := -1
	for i, feed := range feeds {
		if feed["url"] == req.URL {
			if index != -1 {
				return nil, errors.New("subscription URL is ambiguous")
			}
			index = i
		}
	}
	result := slices.Clone(feeds)
	if req.Action == MutationAdd {
		if index != -1 {
			return nil, errors.New("subscription already exists")
		}
		entry := map[string]any{"url": req.URL}
		if err := applyMutationValues(entry, req.Fields, mutationFields); err != nil {
			return nil, err
		}
		return append(result, entry), nil
	}
	if index == -1 {
		return nil, errors.New("subscription is not present in the owning TOML file")
	}
	if req.Action == MutationRemove {
		return slices.Delete(result, index, index+1), nil
	}
	entry := maps.Clone(result[index])
	if req.Action == MutationActive {
		entry[mutationFieldActive] = req.Active
	} else {
		names := req.SelectedFields
		if len(names) == 0 {
			names = mutationFields
		}
		if err := applyMutationValues(entry, req.Fields, names); err != nil {
			return nil, err
		}
	}
	result[index] = entry
	return result, nil
}

func applyMutationValues(entry map[string]any, fields map[string]string, names []string) error {
	for _, name := range names {
		value, ok := fields[name]
		if !ok || !slices.Contains(mutationFields, name) {
			return fmt.Errorf("unsupported subscription field %q", name)
		}
		var parsed map[string]any
		if _, err := toml.Decode(name+" = "+value, &parsed); err != nil {
			return fmt.Errorf("invalid subscription field %q: %w", name, err)
		}
		field, present := parsed[name]
		if len(parsed) != 1 || !present {
			return fmt.Errorf("invalid subscription field %q", name)
		}
		entry[name] = field
	}
	return nil
}
