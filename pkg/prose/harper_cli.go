package prose

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

// HarperCLIOptions configures the optional native Harper adapter. The default
// Markata prose checker works without an installed Harper executable.
type HarperCLIOptions struct {
	Binary  string
	Timeout time.Duration
	Options Options
}

type harperCLIFileResult struct {
	File  string          `json:"file"`
	Lints []harperCLILint `json:"lints"`
	Error string          `json:"error"`
}

type harperCLILint struct {
	Rule        string `json:"rule"`
	Kind        string `json:"kind"`
	Message     string `json:"message"`
	MatchedText string `json:"matched_text"`
	Span        struct {
		CharStart int `json:"char_start"`
		CharEnd   int `json:"char_end"`
	} `json:"span"`
	Suggestions []string `json:"suggestions"`
}

// ParseHarperCLIJSON converts Harper's character-based JSON spans to the shared
// Markata byte-offset contract. Malformed ranges or mismatched source content
// are rejected rather than producing an unsafe edit.
func ParseHarperCLIJSON(source string, data []byte, options Options) ([]Issue, error) {
	if !utf8.ValidString(source) {
		return nil, fmt.Errorf("Harper input must be valid UTF-8")
	}
	var files []harperCLIFileResult
	if err := json.Unmarshal(data, &files); err != nil {
		return nil, fmt.Errorf("decode Harper JSON: %w", err)
	}
	if len(files) != 1 {
		return nil, fmt.Errorf("expected one Harper file result, got %d", len(files))
	}
	if files[0].Error != "" {
		return nil, fmt.Errorf("Harper: %s", files[0].Error)
	}
	// Each entry is a byte offset at a Unicode character boundary; Harper's
	// char_start/char_end are Unicode scalar indices, not UTF-8 or UTF-16.
	byteOffsets := make([]int, 0, utf8.RuneCountInString(source)+1)
	for pos := range source {
		byteOffsets = append(byteOffsets, pos)
	}
	byteOffsets = append(byteOffsets, len(source))
	issues := make([]Issue, 0, len(files[0].Lints))
	for _, lint := range files[0].Lints {
		if lint.Rule == "" {
			return nil, fmt.Errorf("Harper returned a lint without a rule ID")
		}
		if lint.Span.CharStart < 0 || lint.Span.CharEnd <= lint.Span.CharStart ||
			lint.Span.CharEnd >= len(byteOffsets) {
			return nil, fmt.Errorf("Harper returned an invalid character range")
		}
		start := byteOffsets[lint.Span.CharStart]
		end := byteOffsets[lint.Span.CharEnd]
		if source[start:end] != lint.MatchedText {
			return nil, fmt.Errorf("Harper source mismatch for rule %s", lint.Rule)
		}
		ruleID := "harper/" + lint.Rule
		if options.IgnoredRules[ruleID] ||
			options.IgnoredPhrases[ruleID][normalizePhrase(lint.MatchedText)] {
			continue
		}
		// A Harper suggestion is a human-readable action, not raw replacement
		// text. Only the unambiguous ReplaceWith variant is safe as an edit.
		replacement := harperReplacement(lint.Suggestions)
		issue, err := NewIssue(source, "harper", ruleID, lint.Message, start, end, replacement)
		if err != nil {
			return nil, err
		}
		issues = append(issues, issue)
	}
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].StartByte != issues[j].StartByte {
			return issues[i].StartByte < issues[j].StartByte
		}
		return issues[i].Rule < issues[j].Rule
	})
	return issues, nil
}

func harperReplacement(suggestions []string) string {
	for _, suggestion := range suggestions {
		const prefix = "Replace with: “"
		if strings.HasPrefix(suggestion, prefix) && strings.HasSuffix(suggestion, "”") {
			return strings.TrimSuffix(strings.TrimPrefix(suggestion, prefix), "”")
		}
	}
	return ""
}

// CheckWithHarperCLI checks unsaved Markdown against an optional local
// executable without sending writing to a network service. No process is
// started unless the caller explicitly enables this adapter.
func CheckWithHarperCLI(ctx context.Context, source string, cfg HarperCLIOptions) ([]Issue, error) {
	binary := cfg.Binary
	if binary == "" {
		binary = "harper-cli"
	}
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 8 * time.Second
	}
	if !utf8.ValidString(source) {
		return nil, fmt.Errorf("Harper input must be valid UTF-8")
	}
	dir, err := os.MkdirTemp("", "markata-harper-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "draft.md")
	if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary, "lint", "--format", "json", "--quiet", path)
	stdout, err := cmd.Output()
	if ctx.Err() != nil {
		return nil, fmt.Errorf("Harper timed out or was cancelled: %w", ctx.Err())
	}
	if err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) || exit.ExitCode() != 1 {
			return nil, fmt.Errorf("run Harper: %w (install harper-cli or disable --harper)", err)
		}
	}
	issues, parseErr := ParseHarperCLIJSON(source, stdout, cfg.Options)
	if parseErr != nil {
		return nil, parseErr
	}
	if err != nil && len(issues) == 0 {
		return nil, fmt.Errorf("Harper exited unsuccessfully without findings: %w", err)
	}
	return issues, nil
}
