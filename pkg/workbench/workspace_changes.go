package workbench

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// workspaceChange is intentionally a read-only Git observation. A candidate
// Markdown path is not authorization to publish or discard that path.
type workspaceChange struct {
	Path            string `json:"path"`
	PreviousPath    string `json:"previous_path,omitempty"`
	Status          string `json:"status"`
	Kind            string `json:"kind"`
	Staged          bool   `json:"staged"`
	Unstaged        bool   `json:"unstaged"`
	Untracked       bool   `json:"untracked"`
	Workspace       bool   `json:"workspace"`
	Markdown        bool   `json:"markdown"`
}

type workspaceChanges struct {
	Head      string            `json:"head"`
	Entries   []workspaceChange `json:"entries"`
	Total     int               `json:"total"`
	Staged    int               `json:"staged"`
	Unstaged  int               `json:"unstaged"`
	Untracked int               `json:"untracked"`
}

// parseWorkspaceChanges decodes -z porcelain v1, including the reversed
// destination/source fields used for rename and copy records. It does not
// split on whitespace: filenames may include spaces, tabs and newlines.
func parseWorkspaceChanges(status string, workspacePrefix string) ([]workspaceChange, error) {
	entries := make([]workspaceChange, 0)
	for status != "" {
		record, remainder, ok := strings.Cut(status, "\x00")
		if !ok || len(record) < 4 || record[2] != ' ' {
			return nil, fmt.Errorf("malformed Git status record")
		}
		code, name := record[:2], record[3:]
		if name == "" {
			return nil, fmt.Errorf("empty Git status path")
		}
		status = remainder
		entry := workspaceChange{
			Path: name, Status: code,
			Staged: code[0] != ' ' && code[0] != '?',
			Unstaged: code[1] != ' ' && code[1] != '?',
			Untracked: code == "??",
		}
		if strings.ContainsAny(code, "RC") {
			entry.PreviousPath, status, ok = strings.Cut(status, "\x00")
			if !ok || entry.PreviousPath == "" {
				return nil, fmt.Errorf("malformed Git rename/copy source")
			}
		}
		switch {
		case entry.Untracked:
			entry.Kind = "untracked"
		case strings.Contains(code, "U") || code == "AA" || code == "DD":
			entry.Kind = "conflict"
		case strings.Contains(code, "R"):
			entry.Kind = "renamed"
		case strings.Contains(code, "C"):
			entry.Kind = "copied"
		case strings.Contains(code, "D"):
			entry.Kind = "deleted"
		case strings.Contains(code, "A"):
			entry.Kind = "added"
		default:
			entry.Kind = "modified"
		}
		// Git status paths are repository-relative, even when the editor
		// workspace is nested inside that repository. Do not omit changes
		// outside the selected site; flag their scope for the UI instead.
		prefix := filepath.ToSlash(workspacePrefix)
		entry.Workspace = prefix == "." || name == prefix || strings.HasPrefix(name, prefix+"/")
		entry.Markdown = strings.HasSuffix(strings.ToLower(name), ".md")
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
	return entries, nil
}

func (h *Handler) serveWorkspaceChanges(w http.ResponseWriter, r *http.Request) {
	// Only the existing loopback-only /workspace handler mounts this endpoint.
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		fail(w, http.StatusMethodNotAllowed, "GET required")
		return
	}
	workspace, err := filepath.Abs(h.root)
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "invalid workspace root")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	root, err := gitOutput(ctx, workspace, "rev-parse", "--show-toplevel")
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "workspace is not in a Git repository")
		return
	}
	root = strings.TrimSpace(root)
	prefix, err := filepath.Rel(root, workspace)
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "cannot locate Git workspace")
		return
	}
	head, err := gitOutput(ctx, root, "rev-parse", "--verify", "HEAD")
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "repository has no committed baseline")
		return
	}
	porcelain, err := gitOutput(ctx, root, "status", "--porcelain=v1", "-z", "--untracked-files=all")
	if err != nil {
		fail(w, http.StatusServiceUnavailable, "cannot enumerate Git changes")
		return
	}
	entries, err := parseWorkspaceChanges(porcelain, prefix)
	if err != nil {
		fail(w, http.StatusServiceUnavailable, "invalid Git change listing")
		return
	}
	result := workspaceChanges{
		Head: strings.TrimSpace(head), Entries: entries, Total: len(entries),
	}
	for _, entry := range entries {
		if entry.Staged {
			result.Staged++
		}
		if entry.Unstaged {
			result.Unstaged++
		}
		if entry.Untracked {
			result.Untracked++
		}
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		return
	}
}
