package plugins

import (
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"
)

// dependencyIdentityKeys returns the lookup identities that can resolve a
// wikilink/embed target. Keep this aligned with PostIndex.LookupBySlug: exact
// case-insensitive identity first, then the slugified form.
func dependencyIdentityKeys(value string) []string {
	value = strings.TrimSpace(value)
	if fragment := strings.Index(value, "#"); fragment >= 0 {
		value = value[:fragment]
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	keys := make([]string, 0, 3)
	seen := make(map[string]bool, 3)
	add := func(key string) {
		key = strings.TrimSpace(key)
		if key == "" || seen[key] {
			return
		}
		seen[key] = true
		keys = append(keys, key)
	}

	add(value)
	add(strings.ToLower(value))
	add(models.Slugify(value))
	return keys
}

// postDependencyIdentities returns every identity whose appearance/change can
// satisfy an unresolved dependency: the canonical slug plus supported aliases.
func postDependencyIdentities(post *models.Post) []string {
	if post == nil {
		return nil
	}

	identities := dependencyIdentityKeys(post.Slug)
	seen := make(map[string]bool, len(identities)+4)
	for _, identity := range identities {
		seen[identity] = true
	}
	addAlias := func(alias string) {
		for _, identity := range dependencyIdentityKeys(alias) {
			if seen[identity] {
				continue
			}
			seen[identity] = true
			identities = append(identities, identity)
		}
	}

	if post.Extra == nil {
		return identities
	}
	switch aliases := post.Extra["aliases"].(type) {
	case []interface{}:
		for _, raw := range aliases {
			if alias, ok := raw.(string); ok {
				addAlias(alias)
			}
		}
	case []string:
		for _, alias := range aliases {
			addAlias(alias)
		}
	case string:
		addAlias(aliases)
	}
	return identities
}

// unresolvedLogicalDependencies finds wikilink/internal-embed target identities
// that remain unresolved after the normal transform plugins run. Resolved
// targets are already recorded in Post.Dependencies by those plugins.
func unresolvedLogicalDependencies(content string) []string {
	if content == "" {
		return nil
	}

	// Both wikilinks and internal embeds intentionally ignore fenced code.
	content = wikilinksCodeBlockRegex.ReplaceAllString(content, "")
	seen := make(map[string]bool)
	dependencies := make([]string, 0)
	addTarget := func(target string) {
		target = strings.TrimSpace(target)
		if target == "" || isExternalEmbedURL(target) {
			return
		}
		for _, identity := range dependencyIdentityKeys(target) {
			if seen[identity] {
				continue
			}
			seen[identity] = true
			dependencies = append(dependencies, identity)
		}
	}

	for _, groups := range wikilinkRegex.FindAllStringSubmatch(content, -1) {
		if len(groups) >= 2 {
			addTarget(groups[1])
		}
	}
	for _, groups := range internalEmbedRegex.FindAllStringSubmatch(content, -1) {
		if len(groups) >= 2 {
			addTarget(groups[1])
		}
	}
	return dependencies
}
