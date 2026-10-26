package resourcehints

import (
	"fmt"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"
)

// Generator generates HTML resource hint tags.
type Generator struct{}

// NewGenerator creates a new Generator.
func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateHintTags generates HTML link tags for the given hints.
// Returns a string of newline-separated link tags.
func (g *Generator) GenerateHintTags(hints []SuggestedHint) string {
	if len(hints) == 0 {
		return ""
	}

	var tags []string

	// Sort hints for consistent output (preconnect first, then dns-prefetch, etc.)
	sortedHints := make([]SuggestedHint, len(hints))
	copy(sortedHints, hints)
	sort.Slice(sortedHints, func(i, j int) bool {
		// Priority: preconnect > dns-prefetch > preload > prefetch
		iPriority := hintTypePriority(sortedHints[i].HintTypes)
		jPriority := hintTypePriority(sortedHints[j].HintTypes)
		if iPriority != jPriority {
			return iPriority < jPriority
		}
		return sortedHints[i].Domain < sortedHints[j].Domain
	})

	for _, hint := range sortedHints {
		for _, hintType := range hint.HintTypes {
			tag := g.generateTag(hint, hintType)
			if tag != "" {
				tags = append(tags, tag)
			}
		}
	}

	if len(tags) == 0 {
		return ""
	}

	return strings.Join(tags, "\n")
}

// hintTypePriority returns a priority value for sorting hint types.
// Lower values = higher priority (should come first).
func hintTypePriority(types []HintType) int {
	for _, t := range types {
		switch t {
		case HintTypePreconnect:
			return 0
		case HintTypeDNSPrefetch:
			return 1
		case HintTypePreload:
			return 2
		case HintTypePrefetch:
			return 3
		}
	}
	return 4
}

// generateTag generates a single HTML link tag for a hint.
func (g *Generator) generateTag(hint SuggestedHint, hintType HintType) string {
	scheme := hint.Scheme
	if scheme == "" {
		scheme = "https"
	}

	href := fmt.Sprintf("%s://%s", scheme, hint.Domain)

	attrs := []string{
		fmt.Sprintf("rel=%q", string(hintType)),
		fmt.Sprintf("href=%q", href),
	}

	// Add crossorigin attribute if specified
	if hint.CrossOrigin != "" {
		if hint.CrossOrigin == "anonymous" {
			attrs = append(attrs, "crossorigin")
		} else {
			attrs = append(attrs, fmt.Sprintf("crossorigin=%q", hint.CrossOrigin))
		}
	}

	// Add "as" attribute for preload hints
	if hintType == HintTypePreload && hint.As != "" {
		attrs = append(attrs, fmt.Sprintf("as=%q", hint.As))
	}

	return fmt.Sprintf("<link %s>", strings.Join(attrs, " "))
}

// GenerateFromConfig generates hint tags from a ResourceHintsConfig.
// Combines manually configured domains with auto-detected ones.
func (g *Generator) GenerateFromConfig(config *models.ResourceHintsConfig, detectedDomains []DetectedDomain) string {
	var allHints []SuggestedHint

	// Add manually configured domains
	for _, d := range config.Domains {
		hint := SuggestedHint{
			Domain:      d.Domain,
			Scheme:      "https",
			CrossOrigin: d.CrossOrigin,
			As:          d.As,
		}

		// Convert string hint types to HintType
		for _, t := range d.HintTypes {
			switch strings.ToLower(t) {
			case "preconnect":
				hint.HintTypes = append(hint.HintTypes, HintTypePreconnect)
			case "dns-prefetch":
				hint.HintTypes = append(hint.HintTypes, HintTypeDNSPrefetch)
			case "preload":
				hint.HintTypes = append(hint.HintTypes, HintTypePreload)
			case "prefetch":
				hint.HintTypes = append(hint.HintTypes, HintTypePrefetch)
			}
		}

		if len(hint.HintTypes) > 0 {
			allHints = append(allHints, hint)
		}
	}

	// Add auto-detected domains if enabled
	if config.IsAutoDetectEnabled() && len(detectedDomains) > 0 {
		detector := NewDetector()
		detector.SetExcludeDomains(config.ExcludeDomains)

		// Filter out domains that are already manually configured
		configuredDomains := make(map[string]bool)
		for _, d := range config.Domains {
			configuredDomains[d.Domain] = true
		}

		var filteredDomains []DetectedDomain
		for _, d := range detectedDomains {
			if !configuredDomains[d.Domain] {
				filteredDomains = append(filteredDomains, d)
			}
		}

		suggestedHints := detector.SuggestHints(filteredDomains)
		allHints = append(allHints, suggestedHints...)
	}

	return g.GenerateHintTags(allHints)
}

// GenerateComment generates an HTML comment to wrap resource hints.
func GenerateComment(content string) string {
	if content == "" {
		return ""
	}
	return fmt.Sprintf("<!-- Auto-generated resource hints -->\n%s\n<!-- End resource hints -->", content)
}
