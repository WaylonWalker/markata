package prose

import (
	"fmt"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// ContractVersion is the stable schema version for cross-runtime prose results.
// Readers may ignore unknown fields; new fields must not change byte semantics.
const ContractVersion = 1

// TextEdit describes one replacement using half-open UTF-8 byte offsets.
// Applying an edit is always opt-in and must revalidate the document version.
type TextEdit struct {
	StartByte   int    `json:"start_byte"`
	EndByte     int    `json:"end_byte"`
	Replacement string `json:"replacement"`
}

// UTF16Position is the line/column representation required by LSP and
// CodeMirror. A UTF-16 column is not necessarily a byte or rune column.
type UTF16Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// UTF16Span represents a half-open range in LSP coordinates.
type UTF16Span struct {
	Start UTF16Position `json:"start"`
	End   UTF16Position `json:"end"`
}

// PointAtByteOffset converts a source byte position into both UTF-8 and UTF-16
// line/column positions. Invalid UTF-8 boundaries are rejected rather than
// allowing a fix to split an encoded character.
func PointAtByteOffset(source string, offset int) (line, byteColumn, utf16Column int, err error) {
	if offset < 0 || offset > len(source) {
		return 0, 0, 0, fmt.Errorf("source offset %d out of bounds", offset)
	}
	if !utf8.ValidString(source) {
		return 0, 0, 0, fmt.Errorf("source contains invalid UTF-8")
	}
	if offset < len(source) && !utf8.RuneStart(source[offset]) {
		return 0, 0, 0, fmt.Errorf("offset %d splits UTF-8 code point", offset)
	}
	prefix := source[:offset]
	line = strings.Count(prefix, "\n")
	newline := strings.LastIndexByte(prefix, '\n')
	start := newline + 1
	byteColumn = offset - start
	for _, r := range source[start:offset] {
		utf16Column += utf16.RuneLen(r)
	}
	return line, byteColumn, utf16Column, nil
}

// NewIssue constructs a diagnostic with stable source offsets and a single
// optional, explicit replacement. The Markata CLI, LSP and browser should all
// consume these same byte ranges before converting to frontend coordinates.
func NewIssue(source, provider, rule, message string, start, end int, replacement string) (Issue, error) {
	if start < 0 || end < start || end > len(source) {
		return Issue{}, fmt.Errorf("invalid prose range [%d,%d) in %d bytes", start, end, len(source))
	}
	line, column, _, err := PointAtByteOffset(source, start)
	if err != nil {
		return Issue{}, err
	}
	endLine, endColumn, _, err := PointAtByteOffset(source, end)
	if err != nil {
		return Issue{}, err
	}
	issue := Issue{
		ContractVersion: ContractVersion,
		Provider:        provider,
		Severity:        "hint",
		Rule:            rule,
		Message:         message,
		StartByte:       start,
		EndByte:         end,
		Line:            line,
		Column:          column,
		EndLine:         endLine,
		EndColumn:       endColumn,
		Suggestion:      replacement,
	}
	if replacement != "" {
		issue.Edits = []TextEdit{{StartByte: start, EndByte: end, Replacement: replacement}}
	}
	return issue, nil
}

// UTF16Range converts the authoritative byte offsets for a diagnostic to an
// LSP-compatible range. Failing rather than guessing protects Unicode edits.
func UTF16Range(source string, issue Issue) (UTF16Span, error) {
	sl, _, sc, err := PointAtByteOffset(source, issue.StartByte)
	if err != nil {
		return UTF16Span{}, err
	}
	el, _, ec, err := PointAtByteOffset(source, issue.EndByte)
	if err != nil {
		return UTF16Span{}, err
	}
	return UTF16Span{
		Start: UTF16Position{Line: sl, Character: sc},
		End:   UTF16Position{Line: el, Character: ec},
	}, nil
}

// OriginalText retrieves the exact source spelling for a diagnostic.
func OriginalText(source string, issue Issue) (string, error) {
	if issue.StartByte < 0 || issue.EndByte <= issue.StartByte || issue.EndByte > len(source) {
		return "", fmt.Errorf("invalid source range")
	}
	if _, _, _, err := PointAtByteOffset(source, issue.StartByte); err != nil {
		return "", err
	}
	if _, _, _, err := PointAtByteOffset(source, issue.EndByte); err != nil {
		return "", err
	}
	return source[issue.StartByte:issue.EndByte], nil
}
