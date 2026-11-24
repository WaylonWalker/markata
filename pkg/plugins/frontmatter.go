// Package plugins provides core plugins for the markata-go static site generator.
package plugins

import (
	"errors"
	"fmt"
	"strings"

	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"gopkg.in/yaml.v3"
)

var (
	// ErrInvalidFrontmatter indicates the frontmatter could not be parsed.
	ErrInvalidFrontmatter = errors.New("invalid frontmatter")

	// ErrRecoverableContent marks a source parsing error that can be isolated to
	// one content file. The load stage uses it to keep valid sibling files while
	// still propagating operational source errors.
	ErrRecoverableContent = errors.New("recoverable content error")
)

// ExtractFrontmatter splits content into frontmatter YAML string and body content.
// Returns:
//   - frontmatter: the raw YAML string between --- delimiters (empty if no frontmatter)
//   - body: the content after the frontmatter
//   - err: error if frontmatter is malformed
//
// Edge cases:
//   - No frontmatter (first line is not exactly ---): returns empty frontmatter, full content as body
//   - Empty frontmatter (---, then ---): returns empty frontmatter, content after second ---
//   - Unclosed frontmatter: returns error
func ExtractFrontmatter(content string) (frontmatter, body string, err error) {
	inspection := diagnostics.InspectFrontmatter(content)
	if inspection.Problem != "" {
		return "", "", fmt.Errorf("%w: invalid frontmatter delimiters", ErrInvalidFrontmatter)
	}
	if !inspection.HasFrontmatter {
		return "", inspection.Body, nil
	}
	return inspection.Frontmatter, inspection.Body, nil
}

// ParseFrontmatter parses content containing optional YAML frontmatter.
// Returns the parsed metadata as a map, the body content, and any error.
//
// The frontmatter must be delimited by --- at the start of the content.
//
// Example:
//
//	---
//	title: My Post
//	date: 2024-01-15
//	tags:
//	  - go
//	  - programming
//	---
//	# Content here
//
// Edge cases handled:
//   - No frontmatter: returns empty map, full content as body
//   - Empty frontmatter: returns empty map, content after delimiters
//   - Invalid YAML: returns error with context
func ParseFrontmatter(content string) (metadata map[string]interface{}, body string, err error) {
	var frontmatter string
	frontmatter, body, err = ExtractFrontmatter(content)
	if err != nil {
		return nil, "", err
	}

	// No frontmatter case
	if frontmatter == "" {
		return make(map[string]interface{}), body, nil
	}

	// Parse the YAML
	metadata = make(map[string]interface{})
	if err = yaml.Unmarshal([]byte(frontmatter), &metadata); err != nil {
		return nil, "", fmt.Errorf("%w: %w", ErrInvalidFrontmatter, err)
	}

	// Handle nil result from empty YAML
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	normalizeAliases(metadata)

	return metadata, body, nil
}

// ParseFrontmatterWithRaw is like ParseFrontmatter but also returns the raw frontmatter string.
// This is useful for computing input hashes for incremental builds.
func ParseFrontmatterWithRaw(content string) (metadata map[string]interface{}, body, rawFrontmatter string, err error) {
	rawFrontmatter, body, err = ExtractFrontmatter(content)
	if err != nil {
		return nil, "", "", err
	}

	// No frontmatter case
	if rawFrontmatter == "" {
		return make(map[string]interface{}), body, "", nil
	}

	// Parse the YAML
	metadata = make(map[string]interface{})
	if err = yaml.Unmarshal([]byte(rawFrontmatter), &metadata); err != nil {
		return nil, "", "", fmt.Errorf("%w: %w", ErrInvalidFrontmatter, err)
	}

	// Handle nil result from empty YAML
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	normalizeAliases(metadata)

	return metadata, body, rawFrontmatter, nil
}

// normalizeAliases normalizes frontmatter alias keys into "aliases".
// It merges values from: aliases, alias, handles, handle.
// Existing keys are preserved; "aliases" is populated if any aliases are found.
func normalizeAliases(metadata map[string]interface{}) {
	if metadata == nil {
		return
	}

	keys := []string{"aliases", "alias", "handles", "handle"}
	aliases := make([]string, 0)
	seen := make(map[string]bool)

	for _, key := range keys {
		value, ok := metadata[key]
		if !ok {
			continue
		}
		for _, alias := range normalizeAliasValue(value) {
			if alias == "" {
				continue
			}
			lower := strings.ToLower(alias)
			if seen[lower] {
				continue
			}
			seen[lower] = true
			aliases = append(aliases, alias)
		}
	}

	if len(aliases) == 0 {
		return
	}

	metadata["aliases"] = aliases
}

func normalizeAliasValue(value interface{}) []string {
	switch v := value.(type) {
	case string:
		if v == "" {
			return nil
		}
		return []string{v}
	case []string:
		return v
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if str, ok := item.(string); ok && str != "" {
				result = append(result, str)
			}
		}
		return result
	default:
		return nil
	}
}

// GetString extracts a string value from metadata, returning empty string if not found or wrong type.
func GetString(metadata map[string]interface{}, key string) string {
	if v, ok := metadata[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetBool extracts a boolean value from metadata, returning defaultVal if not found or wrong type.
// Handles YAML boolean values: true, false, yes, no, on, off.
func GetBool(metadata map[string]interface{}, key string, defaultVal bool) bool {
	v, ok := metadata[key]
	if !ok {
		return defaultVal
	}

	switch b := v.(type) {
	case bool:
		return b
	case string:
		switch strings.ToLower(b) {
		case BoolTrue, Yes, "on":
			return true
		case BoolFalse, No, Off:
			return false
		}
	}
	return defaultVal
}

// GetStringSlice extracts a string slice from metadata.
// Handles scalar strings as one-item lists, plus []interface{} (common from
// YAML) and []string.
func GetStringSlice(metadata map[string]interface{}, key string) []string {
	v, ok := metadata[key]
	if !ok {
		return nil
	}

	switch s := v.(type) {
	case string:
		if s == "" {
			return nil
		}
		return []string{s}
	case []string:
		return s
	case []interface{}:
		result := make([]string, 0, len(s))
		for _, item := range s {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		return result
	}
	return nil
}
