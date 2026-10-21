package fontpacks

import (
	"errors"
	"io"
	"slices"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/html"
)

// Coverage is an immutable, sorted set of distinct visible runes. Its zero
// value represents no text. Copies can safely be reused by multiple resolvers.
type Coverage struct {
	runes []rune
}

// Signature returns the canonical visible-rune set as a string, not a hash.
func (c Coverage) Signature() string { return string(c.runes) }

// CollectCoverage tokenizes one continuous HTML stream, including across reader
// boundaries in an io.MultiReader. It matches VisibleText's text and whitespace
// semantics, excluding replacement runes, and surfaces non-EOF reader errors.
func CollectCoverage(reader io.Reader) (Coverage, error) {
	tokenizer := html.NewTokenizer(reader)
	seen := make(map[rune]struct{}, 256)
	skipElement := ""
	for {
		switch tokenizer.Next() {
		case html.TextToken:
			if skipElement == "" {
				for _, r := range string(tokenizer.Text()) {
					if r != utf8.RuneError {
						// HTML collapses tabs, line breaks, form feeds, and
						// carriage returns to ordinary spaces during rendering.
						// Subset fonts need U+0020, not separate control glyphs.
						if r == '\t' || r == '\n' || r == '\f' || r == '\r' {
							r = ' '
						}
						seen[r] = struct{}{}
					}
				}
				// VisibleText appends a space after every visible text token.
				seen[' '] = struct{}{}
			}
		case html.StartTagToken:
			name, _ := tokenizer.TagName()
			if string(name) == "script" || string(name) == "style" {
				skipElement = string(name)
			}
		case html.EndTagToken:
			name, _ := tokenizer.TagName()
			if string(name) == skipElement {
				skipElement = ""
			}
		case html.SelfClosingTagToken, html.CommentToken, html.DoctypeToken:
		case html.ErrorToken:
			if !errors.Is(tokenizer.Err(), io.EOF) {
				return Coverage{}, tokenizer.Err()
			}
			runes := make([]rune, 0, len(seen))
			for r := range seen {
				runes = append(runes, r)
			}
			slices.Sort(runes)
			return Coverage{runes: runes}, nil
		}
	}
}

func coverageFromString(renderedHTML string) Coverage {
	// A strings.Reader cannot produce a non-EOF reader error.
	coverage, _ := CollectCoverage(strings.NewReader(renderedHTML)) //nolint:errcheck // strings.Reader produces only EOF, handled by CollectCoverage.
	return coverage
}

func baseTiers(pack FontPack) map[string]map[string]bool {
	result := map[string]map[string]bool{}
	for _, role := range pack.Roles {
		if role.Source != "" {
			if result[role.Source] == nil {
				result[role.Source] = map[string]bool{}
			}
			result[role.Source][role.Tier] = true
		}
	}
	return result
}

func requiredTiers(pack FontPack, coverage Coverage, profiles map[string]compiledSubsetProfile, manifests map[string]Manifest) map[string]map[string]bool {
	result := baseTiers(pack)
	if manifests != nil {
		result = tiersForManifests(result, manifests)
	}
	for source, tiers := range result {
		if !tiers["full"] {
			for _, r := range coverage.runes {
				if tiers["prose-extra"] && inAnyCompiledProfile(profiles, tiers, r) {
					continue // the partition already covers this extended character
				}
				if compiledProfileContains(profiles["latin-ext"], r) {
					// A source without a latin-ext tier can use the browser's
					// normal family fallback for these characters. Loading its
					// unrestricted full face would make that large file eligible
					// for every character and defeat the smaller base tier.
					if manifest, ok := manifests[source]; ok {
						if _, hasExtended := manifest.Tiers["latin-ext"]; !hasExtended {
							continue
						}
					}
					tiers["latin-ext"] = true
					continue
				}
				if !inAnyCompiledProfile(profiles, tiers, r) {
					// Characters outside the bundled subsets use the next family
					// in the CSS font stack (typically the system fallback). Do not
					// add an unrestricted full face just to cover an exceptional
					// character; it can otherwise replace the core face for ASCII.
					continue
				}
			}
		}
		if tiers["full"] {
			result[source] = map[string]bool{"full": true}
		}
	}
	return result
}

func tiersForManifests(requested map[string]map[string]bool, manifests map[string]Manifest) map[string]map[string]bool {
	result := make(map[string]map[string]bool, len(requested))
	for source, tiers := range requested {
		manifest := manifests[source]
		selected := make(map[string]bool)
		for tier := range tiers {
			_, basic := manifest.Tiers["prose-basic"]
			_, extra := manifest.Tiers["prose-extra"]
			_, display := manifest.Tiers["display-core"]
			if basic && extra && (tier == "prose-core" || (tier == "display-core" && !display)) {
				selected["prose-basic"] = true
				selected["prose-extra"] = true
				continue
			}
			if _, ok := manifest.Tiers[tier]; ok {
				selected[tier] = true
			} else if tier == "display-core" {
				if _, ok := manifest.Tiers["prose-core"]; ok {
					// prose-core includes the display-core repertoire and remains a
					// small subset. Prefer it over an unrestricted full face.
					selected["prose-core"] = true
				} else if _, ok := manifest.Tiers["full"]; ok {
					selected["full"] = true
				} else {
					selected[tier] = true // retain the useful missing-tier diagnostic
				}
			} else if _, ok := manifest.Tiers["full"]; ok {
				selected["full"] = true
			} else {
				selected[tier] = true // retain the useful missing-tier diagnostic
			}
		}
		if selected["full"] {
			selected = map[string]bool{"full": true}
		}
		result[source] = selected
	}
	return result
}
