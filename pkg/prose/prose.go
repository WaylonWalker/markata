// Package prose provides small, deterministic prose checks shared by the CLI and LSP.
// The checks are editorial suggestions, not claims about AI authorship.
package prose

import (
	"regexp"
	"sort"
	"strings"
)

// Issue uses zero-based line numbers and UTF-8 byte columns. Frontends may
// convert columns to UTF-16 for LSP; the source document is never rewritten.
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
	// IgnoredRules is keyed by stable rule ID (for example, "prose/wordiness").
	IgnoredRules map[string]bool
}

type rule struct {
	id, message, suggestion string
	re *regexp.Regexp
}

var rules = []rule{
	{"prose/wordiness", "This introduction may be unnecessary.", "", regexp.MustCompile(`(?i)\b(?:it is important to note that|it's important to note that|it is worth noting that)\b`)},
	{"prose/cliche", "Consider replacing this familiar phrase with a specific description.", "", regexp.MustCompile(`(?i)\b(?:delve into|rich tapestry|in today's fast-paced world|in today's rapidly evolving landscape)\b`)},
	{"prose/intensifier", "Does this intensifier add useful meaning?", "", regexp.MustCompile(`(?i)\b(?:very unique|completely revolutionary)\b`)},
}

// Check reports conservative, non-blocking style observations. It ignores
// fenced code, inline code, frontmatter, Markdown link targets and bare URLs.
// It does not attempt spelling/grammar correction (a future engine adapter).
func Check(source string, options Options) []Issue {
	masked := maskNonProse(source)
	issues := make([]Issue, 0)
	for _, r := range rules {
		if options.IgnoredRules[r.id] {
			continue
		}
		for _, loc := range r.re.FindAllStringIndex(masked, -1) {
			sl, sc := position(source, loc[0])
			el, ec := position(source, loc[1])
			issues = append(issues, Issue{
				Rule: r.id, Message: r.message, Line: sl, Column: sc,
				EndLine: el, EndColumn: ec, Suggestion: r.suggestion,
			})
		}
	}
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Line != issues[j].Line { return issues[i].Line < issues[j].Line }
		if issues[i].Column != issues[j].Column { return issues[i].Column < issues[j].Column }
		return issues[i].Rule < issues[j].Rule
	})
	return issues
}

func position(source string, offset int) (int, int) {
	prefix := source[:offset]
	last := strings.LastIndexByte(prefix, '\n')
	return strings.Count(prefix, "\n"), offset - last - 1
}

var (
	inlineCode = regexp.MustCompile("`+[^\n`]*`+")
	urlPattern = regexp.MustCompile(`(?i)https?://[^\s<>)]+`)
	linkTarget = regexp.MustCompile(`\]\([^\n)]*\)`)
)

// maskNonProse replaces excluded bytes with spaces without moving source
// offsets. This intentionally prefers missing a suggestion to linting code.
func maskNonProse(source string) string {
	data := []byte(source)
	mask := func(start, end int) {
		for i := start; i < end; i++ { if data[i] != '\n' && data[i] != '\r' { data[i] = ' ' } }
	}
	lines := strings.SplitAfter(source, "\n")
	offset := 0
	inFence := false
	fenceChar := byte(0)
	fenceLen := 0
	inFrontmatter := false
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" { inFrontmatter = true }
	for index, line := range lines {
		trim := strings.TrimLeft(line, " \t")
		if inFrontmatter {
			mask(offset, offset+len(line))
			if index > 0 && (strings.TrimSpace(line) == "---" || strings.TrimSpace(line) == "...") { inFrontmatter = false }
		} else {
			marker := byte(0)
			count := 0
			if len(trim) > 0 && (trim[0] == '`' || trim[0] == '~') {
				marker = trim[0]
				for count < len(trim) && trim[count] == marker { count++ }
			}
			if inFence {
				mask(offset, offset+len(line))
				if marker == fenceChar && count >= fenceLen { inFence = false }
			} else if count >= 3 {
				inFence, fenceChar, fenceLen = true, marker, count
				mask(offset, offset+len(line))
			} else if strings.HasPrefix(trim, "    ") || strings.HasPrefix(line, "\t") {
				mask(offset, offset+len(line))
			}
		}
		offset += len(line)
	}
	// Apply inline masks to the still-visible text, not to fenced regions.
	for _, re := range []*regexp.Regexp{inlineCode, linkTarget, urlPattern} {
		for _, match := range re.FindAllIndex(data, -1) { mask(match[0], match[1]) }
	}
	return string(data)
}
