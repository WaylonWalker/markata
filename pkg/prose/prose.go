// Package prose provides deterministic editorial checks shared by CLI and LSP.
// Observations describe writing patterns, never claims about AI authorship.
package prose

import (
	"regexp"
	"sort"
	"strings"
)

// Issue uses zero-based line numbers and UTF-8 byte columns. LSP clients
// convert these positions to UTF-16. Check never modifies the document.
type Issue struct {
	Rule       string `json:"rule"`
	Message    string `json:"message"`
	Line       int    `json:"line"`
	Column     int    `json:"column"`
	EndLine    int    `json:"end_line"`
	EndColumn  int    `json:"end_column"`
	Suggestion string `json:"suggestion,omitempty"`
}

type Options struct {
	IgnoredRules   map[string]bool
	IgnoredPhrases map[string]map[string]bool
}

type rule struct {
	id         string
	message    string
	suggestion string
	re         *regexp.Regexp
}

var rules = []rule{
	{"prose/wordiness", "This introduction may be unnecessary.", "", regexp.MustCompile(`(?i)\b(?:it is important to note that|it's important to note that|it is worth noting that)\b`)},
	{"prose/cliche", "Consider replacing this familiar phrase with a specific description.", "", regexp.MustCompile(`(?i)\b(?:delve into|rich tapestry|in today's fast-paced world|in today's rapidly evolving landscape)\b`)},
	{"prose/intensifier", "Does this intensifier add useful meaning?", "", regexp.MustCompile(`(?i)\b(?:very unique|completely revolutionary)\b`)},
}

// Check returns non-blocking observations outside protected Markdown regions.
// A separate engine adapter will eventually provide spelling and grammar.
func Check(source string, options Options) []Issue {
	masked := maskNonProse(source)
	issues := make([]Issue, 0)
	for _, r := range rules {
		if options.IgnoredRules[r.id] {
			continue
		}
		for _, loc := range r.re.FindAllStringIndex(masked, -1) {
			matched := source[loc[0]:loc[1]]
			if options.IgnoredPhrases[r.id][normalizePhrase(matched)] {
				continue
			}
			sl, sc := position(source, loc[0])
			el, ec := position(source, loc[1])
			issues = append(issues, Issue{
				Rule:       r.id,
				Message:    r.message,
				Line:       sl,
				Column:     sc,
				EndLine:    el,
				EndColumn:  ec,
				Suggestion: suggestionFor(r, matched),
			})
		}
	}
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Line != issues[j].Line {
			return issues[i].Line < issues[j].Line
		}
		if issues[i].Column != issues[j].Column {
			return issues[i].Column < issues[j].Column
		}
		return issues[i].Rule < issues[j].Rule
	})
	return issues
}

func position(source string, offset int) (line, column int) {
	prefix := source[:offset]
	last := strings.LastIndexByte(prefix, '\n')
	return strings.Count(prefix, "\n"), offset - last - 1
}

var (
	inlineCode = regexp.MustCompile("`+[^\n`]*`+")
	urlPattern = regexp.MustCompile(`(?i)https?://[^\s<>)]+`)
	linkTarget = regexp.MustCompile(`\]\([^\n)]*\)`)
)

func maskRange(data []byte, start, end int) {
	for i := start; i < end; i++ {
		if data[i] != '\n' && data[i] != '\r' {
			data[i] = ' '
		}
	}
}

// maskNonProse replaces protected bytes with spaces to preserve all offsets.
func maskNonProse(source string) string {
	data := []byte(source)
	lines := strings.SplitAfter(source, "\n")
	offset := 0
	inFence := false
	fenceChar := byte(0)
	fenceLen := 0
	inFrontmatter := len(lines) > 0 && strings.TrimSpace(lines[0]) == "---"
	for index, line := range lines {
		trim := strings.TrimLeft(line, " \t")
		marker, count := fenceMarker(trim)
		switch {
		case inFrontmatter:
			maskRange(data, offset, offset+len(line))
			if index > 0 && (strings.TrimSpace(line) == "---" || strings.TrimSpace(line) == "...") {
				inFrontmatter = false
			}
		case inFence:
			maskRange(data, offset, offset+len(line))
			if marker == fenceChar && count >= fenceLen {
				inFence = false
			}
		case count >= 3:
			inFence, fenceChar, fenceLen = true, marker, count
			maskRange(data, offset, offset+len(line))
		case strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t"):
			maskRange(data, offset, offset+len(line))
		}
		offset += len(line)
	}
	for _, re := range []*regexp.Regexp{inlineCode, linkTarget, urlPattern} {
		for _, match := range re.FindAllIndex(data, -1) {
			maskRange(data, match[0], match[1])
		}
	}
	return string(data)
}

func fenceMarker(line string) (marker byte, count int) {
	if line == "" || (line[0] != '`' && line[0] != '~') {
		return 0, 0
	}
	marker = line[0]
	for count < len(line) && line[count] == marker {
		count++
	}
	return marker, count
}

// KnownRule reports whether an identifier names one of our built-in rules.
func KnownRule(id string) bool {
	for _, r := range rules {
		if r.id == id {
			return true
		}
	}
	return false
}

// suggestionFor returns only deterministic, explicit editorial replacements.
// Other observations remain explanations rather than automatic rewrites.
func suggestionFor(r rule, matched string) string {
	if r.id != "prose/intensifier" {
		return ""
	}
	parts := strings.SplitN(matched, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	result := parts[1]
	if parts[0] != strings.ToLower(parts[0]) {
		result = strings.ToUpper(result[:1]) + result[1:]
	}
	return result
}

// normalizePhrase consistently matches accepted phrases case-insensitively.
func normalizePhrase(phrase string) string {
	return strings.ToLower(strings.TrimSpace(phrase))
}

// ValidPhrase checks an exact known prose pattern. It is used by the browser
// preference endpoint to reject arbitrary strings and unbounded values.
func ValidPhrase(id, phrase string) bool {
	if len(phrase) == 0 || len(phrase) > 256 || strings.TrimSpace(phrase) != phrase {
		return false
	}
	for _, r := range rules {
		if r.id != id {
			continue
		}
		matches := r.re.FindStringIndex(phrase)
		return matches != nil && matches[0] == 0 && matches[1] == len(phrase)
	}
	return false
}
