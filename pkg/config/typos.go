package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/suggest"
)

// likelyConfigTypo detects close spellings of known keys while preserving
// extension sections and fields that do not resemble built-in settings.
func likelyConfigTypo(wrapper map[string]any) error {
	if _, ok := wrapper["markata-go"]; !ok {
		for key := range wrapper {
			if matches := suggest.Closest(key, []string{"markata-go"}, 1); len(matches) > 0 {
				return fmt.Errorf("unknown config section %q; did you mean [markata-go]?", key)
			}
		}
		return nil
	}
	root, ok := wrapper["markata-go"].(map[string]any)
	if !ok {
		return nil
	}
	known := knownConfigKeyChildren()
	var walk func(map[string]any, string) error
	walk = func(values map[string]any, parent string) error {
		keys := make([]string, 0, len(values))
		for key := range values {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			path := key
			if parent != "" {
				path = parent + "." + key
			}
			valid := false
			for _, candidate := range known[parent] {
				if key == candidate {
					valid = true
					break
				}
			}
			if !valid {
				if matches := suggest.Closest(key, known[parent], 1); len(matches) > 0 {
					candidate := matches[0]
					// Root keys may belong to plugins. Only a close spelling
					// of a multiword built-in key is safe to diagnose here.
					if parent == "" && (!strings.Contains(key, "_") || !strings.Contains(candidate, "_")) {
						continue
					}
					if parent != "" {
						candidate = parent + "." + candidate
					}
					return fmt.Errorf("unknown config key %q; did you mean %q? Correct the key and rerun the command", path, candidate)
				}
				continue
			}
			if nested, ok := values[key].(map[string]any); ok && len(known[path]) > 0 {
				if err := walk(nested, path); err != nil {
					return err
				}
			}
		}
		return nil
	}
	if err := misplacedFeedFormats(root); err != nil {
		return err
	}
	return walk(root, "")
}

// misplacedFeedFormats catches known format toggles silently ignored at feed
// level. Unknown extension fields are deliberately left alone.
func misplacedFeedFormats(root map[string]any) error {
	feeds, ok := root["feeds"].([]any)
	if !ok {
		return nil
	}
	formatKeys := []string{"html", "simple_html", "rss", "atom", "json", "markdown", "text", "sitemap"}
	for i, value := range feeds {
		feed, ok := value.(map[string]any)
		if !ok {
			continue
		}
		for _, key := range formatKeys {
			if _, exists := feed[key]; !exists {
				continue
			}
			location := fmt.Sprintf("feeds[%d]", i)
			if slug, ok := feed["slug"].(string); ok && slug != "" {
				location += fmt.Sprintf(" (slug %q)", slug)
			}
			return fmt.Errorf("misplaced feed format key %q in %s: move it to feeds[%d].formats.%s", key, location, i, key)
		}
	}
	return nil
}

func knownConfigKeyChildren() map[string][]string {
	known := map[string][]string{"": {"include", "tailwind", "css_purge"}}
	fields := settingFieldTemplates()
	for fieldIndex := range fields {
		parts := strings.Split(fields[fieldIndex].Key, ".")
		for i, part := range parts {
			parent := strings.Join(parts[:i], ".")
			known[parent] = append(known[parent], part)
		}
	}
	return known
}
