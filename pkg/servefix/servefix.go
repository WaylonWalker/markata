// Package servefix plans reviewable source edits for serve diagnostics.
// Plans never modify files; callers must explicitly select edits to apply.
package servefix

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/lint"
)

var ErrStale = errors.New("source changed since fix preview")

var ambiguousSlashDate = regexp.MustCompile(`\b(\d{1,2})/(\d{1,2})/(?:\d{2}|\d{4})\b`)

const codeH1InContent = "h1-in-content"

// Safety describes how a proposed edit may be approved.
type Safety string

const (
	SafetySafe   Safety = "safe"
	SafetyReview Safety = "review"
	SafetyManual Safety = "manual"
)

// Edit is one independently selectable, source-positioned fix.
type Edit struct {
	ID          string `json:"id"`
	File        string `json:"file"`
	Code        string `json:"code"`
	Category    string `json:"category"`
	Message     string `json:"message"`
	Line        int    `json:"line"`
	Column      int    `json:"column"`
	EndLine     int    `json:"end_line"`
	EndColumn   int    `json:"end_column"`
	Before      string `json:"before"`
	After       string `json:"after"`
	Safety      Safety `json:"safety"`
	Explanation string `json:"explanation"`
	Start       int    `json:"-"`
	End         int    `json:"-"`
}

// Plan is a snapshot of one file and the fixes that can be reviewed.
type Plan struct {
	File   string `json:"file"`
	Digest string `json:"digest"`
	Edits  []Edit `json:"edits"`
}

// Selection supports one-by-one approval, a whole category, or all remaining fixes.
type Selection struct {
	IDs        []string `json:"ids,omitempty"`
	Categories []string `json:"categories,omitempty"`
	All        bool     `json:"all,omitempty"`
	SafeOnly   bool     `json:"safe_only,omitempty"`
}

// FileSelection freezes the selection and digest for one file in a batch.
type FileSelection struct {
	Path      string    `json:"path"`
	Digest    string    `json:"digest"`
	Selection Selection `json:"selection"`
}

// FilePreview contains one immutable file-level preview in a batch.
type FilePreview struct {
	Path      string    `json:"path"`
	Digest    string    `json:"digest"`
	Before    string    `json:"before"`
	After     string    `json:"after"`
	Edits     []Edit    `json:"edits"`
	Selection Selection `json:"selection"`
}

// FileResult reports the outcome for one file in a batch apply.
type FileResult struct {
	Path    string `json:"path"`
	Status  string `json:"status"`
	Reason  string `json:"reason,omitempty"`
	Applied []Edit `json:"applied,omitempty"`
}

// BatchResult reports applied, stale, and failed file outcomes.
type BatchResult struct {
	Files   []FileResult `json:"files"`
	Applied []Edit       `json:"applied"`
}

// PlanFile reads a Markdown source under root and proposes precise fixes.
// A diagnostics.Fixable flag alone is insufficient: fixes without an unambiguous
// text edit are intentionally omitted.
func PlanFile(root, path string) (Plan, error) {
	file, err := sourcePath(root, path)
	if err != nil {
		return Plan{}, err
	}
	content, err := os.ReadFile(file)
	if err != nil {
		return Plan{}, fmt.Errorf("read fix source %s: %w", path, err)
	}
	plan := Plan{File: file, Digest: digest(content), Edits: []Edit{}}
	// Diagnostics normalizes CRLF and bare CR before calculating ranges.
	// Until byte offsets are carried through that API, avoid misplaced edits.
	if bytes.IndexByte(content, '\r') >= 0 {
		return plan, nil
	}
	for _, issue := range diagnostics.Check(file, string(content), nil) {
		if !issue.Fixable && issue.Code != codeH1InContent && issue.Code != diagnostics.ReasonFrontmatterSuspiciousDelimiter && issue.Code != diagnostics.ReasonFrontmatterLeadingWhitespace && issue.Code != diagnostics.ReasonFrontmatterMalformedClosing {
			continue
		}
		edit, ok := candidate(content, issue)
		if !ok {
			continue
		}
		edit.ID = fmt.Sprintf("%s:%d:%d:%s", issue.Code, issue.Range.StartLine+1, issue.Range.StartCol+1, shortDigest([]byte(edit.Before)))
		edit.File = file
		plan.Edits = append(plan.Edits, edit)
	}
	return plan, nil
}

// Preview returns the selected content without changing the source file.
func Preview(plan Plan, content []byte, selection Selection) ([]byte, []Edit, error) {
	if digest(content) != plan.Digest {
		return nil, nil, ErrStale
	}
	ids := make(map[string]bool, len(selection.IDs))
	for _, id := range selection.IDs {
		ids[id] = true
	}
	categories := make(map[string]bool, len(selection.Categories))
	for _, category := range selection.Categories {
		categories[category] = true
	}
	selected := make([]Edit, 0, len(plan.Edits))
	for i := range plan.Edits {
		edit := plan.Edits[i]
		if (selection.All || ids[edit.ID] || categories[edit.Category]) && (!selection.SafeOnly || edit.Safety == SafetySafe) {
			selected = append(selected, edit)
			delete(ids, edit.ID)
		}
	}
	if len(ids) > 0 {
		return nil, nil, fmt.Errorf("selected fix is unavailable in this plan")
	}
	sort.Slice(selected, func(i, j int) bool { return selected[i].Start < selected[j].Start })
	for i := range selected {
		edit := selected[i]
		if edit.Start < 0 || edit.End < edit.Start || edit.End > len(content) || string(content[edit.Start:edit.End]) != edit.Before {
			return nil, nil, ErrStale
		}
		if i > 0 && edit.Start < selected[i-1].End {
			return nil, nil, errors.New("selected fixes overlap")
		}
	}
	result := append([]byte(nil), content...)
	for i := len(selected) - 1; i >= 0; i-- {
		edit := selected[i]
		result = append(append(append([]byte(nil), result[:edit.Start]...), edit.After...), result[edit.End:]...)
	}
	return result, selected, nil
}

// Apply validates a preview against the current file, then replaces that file
// atomically. Callers should trigger a rebuild after a successful edit.
func Apply(root string, plan Plan, selection Selection) ([]Edit, error) {
	file, err := sourcePath(root, plan.File)
	if err != nil {
		return nil, err
	}
	if file != plan.File {
		return nil, errors.New("fix plan file does not match source")
	}
	info, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	updated, applied, err := Preview(plan, content, selection)
	if err != nil || len(applied) == 0 {
		return applied, err
	}
	if err := replaceFile(file, content, updated, info.Mode().Perm()); err != nil {
		return nil, err
	}
	return applied, nil
}

// PreviewBatch freezes each requested file independently. Paths are kept in
// the caller's root-relative form so browser responses never expose absolute
// source locations.
func PreviewBatch(root string, requested []FileSelection) ([]FilePreview, error) {
	if len(requested) == 0 {
		return nil, errors.New("batch contains no files")
	}
	seen := make(map[string]struct{}, len(requested))
	previews := make([]FilePreview, 0, len(requested))
	for _, item := range requested {
		path := filepath.Clean(item.Path)
		if path == "." || filepath.IsAbs(path) {
			return nil, fmt.Errorf("invalid fix path %q", item.Path)
		}
		if _, exists := seen[path]; exists {
			return nil, fmt.Errorf("duplicate fix path %q", item.Path)
		}
		seen[path] = struct{}{}
		plan, err := PlanFile(root, path)
		if err != nil {
			return nil, err
		}
		content, err := os.ReadFile(plan.File)
		if err != nil {
			return nil, fmt.Errorf("read fix source %s: %w", path, err)
		}
		selection := item.Selection
		if len(selection.IDs) == 0 && len(selection.Categories) == 0 && !selection.All {
			selection.All = true
		}
		updated, edits, err := Preview(plan, content, selection)
		if err != nil {
			return nil, fmt.Errorf("preview %s: %w", path, err)
		}
		for i := range edits {
			edits[i].File = filepath.ToSlash(path)
		}
		previews = append(previews, FilePreview{Path: filepath.ToSlash(path), Digest: plan.Digest, Before: string(content), After: string(updated), Edits: edits, Selection: selection})
	}
	return previews, nil
}

// PreviewCurrentBatch returns only files that still match the frozen digests.
// It is intended for pre-apply integrations such as a filesystem watcher that
// needs to know the exact replacement content before the atomic write.
func PreviewCurrentBatch(root string, requested []FileSelection) []FilePreview {
	previews := make([]FilePreview, 0, len(requested))
	for _, item := range requested {
		plan, err := PlanFile(root, item.Path)
		if err != nil || plan.Digest != item.Digest {
			continue
		}
		content, err := os.ReadFile(plan.File)
		if err != nil {
			continue
		}
		updated, edits, err := Preview(plan, content, item.Selection)
		if err != nil || len(edits) == 0 {
			continue
		}
		for i := range edits {
			edits[i].File = filepath.ToSlash(filepath.Clean(item.Path))
		}
		previews = append(previews, FilePreview{Path: filepath.ToSlash(filepath.Clean(item.Path)), Digest: plan.Digest, Before: string(content), After: string(updated), Edits: edits, Selection: item.Selection})
	}
	return previews
}

// ApplyBatch validates all file digests before mutating any file, then applies
// each still-valid file independently. A stale or conflicted file is skipped;
// other valid files can proceed. Each replacement retains Apply's last-moment
// digest check to protect edits made during the batch itself.
func ApplyBatch(root string, previews []FileSelection) BatchResult {
	result := BatchResult{Files: make([]FileResult, 0, len(previews)), Applied: []Edit{}}
	if len(previews) == 0 {
		return result
	}
	type readyFile struct {
		index   int
		preview FileSelection
		plan    Plan
	}
	ready := make([]readyFile, 0, len(previews))
	outcomes := make([]FileResult, len(previews))
	seen := make(map[string]struct{}, len(previews))
	for index, preview := range previews {
		path := filepath.ToSlash(filepath.Clean(preview.Path))
		if _, exists := seen[path]; exists {
			outcomes[index] = FileResult{Path: path, Status: "skipped", Reason: "duplicate file in batch"}
			continue
		}
		seen[path] = struct{}{}
		plan, err := PlanFile(root, path)
		if err != nil {
			outcomes[index] = FileResult{Path: path, Status: "skipped", Reason: err.Error()}
			continue
		}
		if plan.Digest != preview.Digest {
			outcomes[index] = FileResult{Path: path, Status: "stale", Reason: ErrStale.Error()}
			continue
		}
		content, err := os.ReadFile(plan.File)
		if err != nil {
			outcomes[index] = FileResult{Path: path, Status: "skipped", Reason: err.Error()}
			continue
		}
		_, edits, err := Preview(plan, content, preview.Selection)
		if errors.Is(err, ErrStale) {
			outcomes[index] = FileResult{Path: path, Status: "stale", Reason: ErrStale.Error()}
			continue
		}
		if err != nil {
			outcomes[index] = FileResult{Path: path, Status: "skipped", Reason: err.Error()}
			continue
		}
		if len(edits) == 0 {
			outcomes[index] = FileResult{Path: path, Status: "unchanged"}
			continue
		}
		ready = append(ready, readyFile{index: index, preview: preview, plan: plan})
	}
	for i := range ready {
		item := ready[i]
		applied, err := Apply(root, item.plan, item.preview.Selection)
		if errors.Is(err, ErrStale) {
			outcomes[item.index] = FileResult{Path: item.preview.Path, Status: "stale", Reason: ErrStale.Error()}
			continue
		}
		if err != nil {
			outcomes[item.index] = FileResult{Path: item.preview.Path, Status: "failed", Reason: err.Error()}
			continue
		}
		for i := range applied {
			applied[i].File = filepath.ToSlash(filepath.Clean(item.preview.Path))
		}
		outcomes[item.index] = FileResult{Path: item.preview.Path, Status: "applied", Applied: applied}
		result.Applied = append(result.Applied, applied...)
	}
	result.Files = outcomes
	return result
}

//nolint:gocyclo // Each diagnostic code has a small independent safety-checked edit rule.
func candidate(content []byte, issue diagnostics.Issue) (Edit, bool) {
	line, offset, ok := lineAt(content, issue.Range.StartLine)
	if !ok {
		return Edit{}, false
	}
	edit := Edit{
		Code: issue.Code, Category: issue.Code, Message: issue.Message,
		Line: issue.Range.StartLine + 1, Column: issue.Range.StartCol + 1,
		EndLine: issue.Range.EndLine + 1, EndColumn: issue.Range.EndCol + 1,
		Safety: SafetySafe,
	}
	switch issue.Code {
	case codeH1InContent:
		if !strings.HasPrefix(line, "# ") {
			return Edit{}, false
		}
		edit.Start, edit.End = offset, offset+2
		edit.Before, edit.After = "# ", "## "
	case "protocol-less-url":
		start := issue.Range.StartCol
		if start < 0 || start+2 > len(line) || !strings.HasPrefix(line[start:], "//") {
			return Edit{}, false
		}
		edit.Start, edit.End = offset+start, offset+start+2
		edit.Before, edit.After = "//", "https://"
	case "admonition-fenced-code":
		edit.Safety = SafetyReview
		edit.Explanation = "Review the inserted blank line in the surrounding fenced code."
		end := offset + len(line)
		if end >= len(content) || content[end] != '\n' {
			return Edit{}, false
		}
		edit.Start, edit.End = end+1, end+1
		edit.Before, edit.After = "", "\n"
	case "invalid-date":
		// A date such as 01/02/2026 has two plausible interpretations.
		// Keep it in the inbox for manual review instead of assuming US order.
		if match := ambiguousSlashDate.FindStringSubmatch(line); len(match) == 3 {
			first, err := strconv.Atoi(match[1])
			if err != nil {
				return Edit{}, false
			}
			second, err := strconv.Atoi(match[2])
			if err != nil {
				return Edit{}, false
			}
			if first <= 12 && second <= 12 {
				return Edit{}, false
			}
		}
		fixed, changes := lint.NewDateTimeFixer(lint.DefaultDateTimeFixerConfig()).FixDateInContent(line)
		if len(changes) != 1 || fixed == line {
			return Edit{}, false
		}
		edit.Start, edit.End = offset, offset+len(line)
		edit.Before, edit.After = line, fixed
	case diagnostics.ReasonFrontmatterSuspiciousDelimiter,
		diagnostics.ReasonFrontmatterLeadingWhitespace,
		diagnostics.ReasonFrontmatterMalformedClosing:
		if strings.TrimSpace(line) != "---" && !onlyHyphens(strings.TrimSpace(line)) {
			return Edit{}, false
		}
		edit.Start, edit.End = offset, offset+len(line)
		edit.Before, edit.After = line, "---"
	default:
		// In particular, duplicate keys may have nested YAML data. The old
		// whole-file fixer is too destructive to approve a single occurrence.
		return Edit{}, false
	}
	if edit.Explanation == "" {
		edit.Explanation = safetyExplanation(edit.Code, edit.Safety)
	}
	return edit, true
}

func safetyExplanation(code string, safety Safety) string {
	switch safety {
	case SafetySafe:
		switch code {
		case codeH1InContent:
			return "The body heading becomes level two because the page title supplies the level-one heading."
		case "protocol-less-url":
			return "Add https:// to the protocol-relative URL."
		case "invalid-date":
			return "Normalize the unambiguous date using Markata's configured date fixer."
		default:
			return "Normalize the recognized frontmatter delimiter to three hyphens."
		}
	case SafetyReview:
		return "Review the surrounding Markdown before applying this edit."
	default:
		return "Open the source and correct this diagnostic manually."
	}
}

func onlyHyphens(value string) bool {
	if len(value) <= 3 {
		return false
	}
	for _, ch := range value {
		if ch != '-' {
			return false
		}
	}
	return true
}

func lineAt(content []byte, number int) (line string, offset int, ok bool) {
	if number < 0 {
		return "", 0, false
	}
	start := 0
	for i := 0; i < number; i++ {
		newline := bytes.IndexByte(content[start:], '\n')
		if newline < 0 {
			return "", 0, false
		}
		start += newline + 1
	}
	end := start
	for end < len(content) && content[end] != '\n' {
		end++
	}
	return string(content[start:end]), start, true
}

func sourcePath(root, path string) (string, error) {
	base, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", err
	}
	base, err = filepath.Abs(base)
	if err != nil {
		return "", err
	}
	file := path
	if !filepath.IsAbs(file) {
		file = filepath.Join(base, file)
	}
	file, err = filepath.Abs(file)
	if err != nil {
		return "", err
	}
	info, err := os.Lstat(file)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("fix source is not a regular file: %s", path)
	}
	resolved, err := filepath.EvalSymlinks(file)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(base, resolved)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("fix source is outside project root: %s", path)
	}
	return file, nil
}

func replaceFile(file string, original, updated []byte, mode os.FileMode) error {
	tmp, err := os.CreateTemp(filepath.Dir(file), ".markata-fix-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()
	if err := tmp.Chmod(mode); err != nil {
		return err
	}
	if _, err := tmp.Write(updated); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	current, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	if digest(current) != digest(original) {
		return ErrStale
	}
	return replacePath(tmp.Name(), file)
}

func digest(content []byte) string      { return fmt.Sprintf("%x", sha256.Sum256(content)) }
func shortDigest(content []byte) string { return digest(content)[:12] }
