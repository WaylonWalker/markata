package servecontrol

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/WaylonWalker/markata/pkg/servefix"
)

const (
	maxFixBatchFiles  = 256
	fixPreviewTTL     = 5 * time.Minute
	maxFixPreviewKeys = 1024
)

var fixPreviewAuthorizations = struct {
	sync.Mutex
	entries map[string]time.Time
}{entries: make(map[string]time.Time)}

type batchFixPreviewRequest struct {
	Files []servefix.FileSelection `json:"files"`
}

type batchFixApplyRequest struct {
	Files []servefix.FileSelection `json:"files"`
}

type batchFixApplyResponse struct {
	Result           servefix.BatchResult `json:"result"`
	RebuildRequested bool                 `json:"rebuild_requested"`
	RebuildError     string               `json:"rebuild_error,omitempty"`
}

func batchFixPreviewHandler(sourceRoot string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !sameOriginLocalRequest(r, requestScheme(r)) {
			http.Error(w, "fix preview requires a same-origin local request", http.StatusForbidden)
			return
		}
		var request batchFixPreviewRequest
		if err := decodeJSONLimit(w, r, &request, 1<<20); err != nil {
			http.Error(w, "invalid fix request", http.StatusBadRequest)
			return
		}
		if sourceRoot == "" {
			http.Error(w, "source fixes are unavailable", http.StatusNotFound)
			return
		}
		if len(request.Files) == 0 || len(request.Files) > maxFixBatchFiles {
			http.Error(w, fmt.Sprintf("batch must contain between 1 and %d files", maxFixBatchFiles), http.StatusBadRequest)
			return
		}
		previews, err := servefix.PreviewBatch(sourceRoot, request.Files)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if err := authorizeFixPreview(sourceRoot, r, previews); err != nil {
			http.Error(w, "could not authorize fix preview", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		if err := json.NewEncoder(w).Encode(map[string]any{"files": previews}); err != nil {
			return
		}
	}
}

func batchFixApplyHandler(runtime *Runtime, sourceRoot string, afterApply func([]servefix.FilePreview)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !sameOriginLocalRequest(r, requestScheme(r)) {
			http.Error(w, "fix apply requires a same-origin local request", http.StatusForbidden)
			return
		}
		var request batchFixApplyRequest
		if err := decodeJSONLimit(w, r, &request, 1<<20); err != nil {
			http.Error(w, "invalid fix request", http.StatusBadRequest)
			return
		}
		if sourceRoot == "" {
			http.Error(w, "source fixes are unavailable", http.StatusNotFound)
			return
		}
		if len(request.Files) == 0 || len(request.Files) > maxFixBatchFiles {
			http.Error(w, fmt.Sprintf("batch must contain between 1 and %d files", maxFixBatchFiles), http.StatusBadRequest)
			return
		}
		if !consumeFixPreviewAuthorization(sourceRoot, r, request.Files) {
			http.Error(w, "matching fix preview is required before apply", http.StatusConflict)
			return
		}

		result := servefix.ApplyBatch(sourceRoot, request.Files)
		if afterApply != nil {
			if expectedWrites := appliedFixWrites(sourceRoot, result); len(expectedWrites) > 0 {
				afterApply(expectedWrites)
			}
		}
		response := batchFixApplyResponse{Result: result}
		if len(result.Applied) > 0 {
			if err := runtime.Trigger(ActionRequest{Kind: "build"}); err != nil {
				response.RebuildError = err.Error()
			} else {
				response.RebuildRequested = true
			}
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			return
		}
	}
}

func authorizeFixPreview(sourceRoot string, r *http.Request, previews []servefix.FilePreview) error {
	files := make([]servefix.FileSelection, 0, len(previews))
	for i := range previews {
		preview := previews[i]
		files = append(files, servefix.FileSelection{Path: preview.Path, Digest: preview.Digest, Selection: preview.Selection})
	}
	key, err := fixPreviewAuthorizationKey(sourceRoot, r, files)
	if err != nil {
		return err
	}
	now := time.Now()
	fixPreviewAuthorizations.Lock()
	pruneFixPreviewAuthorizationsLocked(now)
	if len(fixPreviewAuthorizations.entries) >= maxFixPreviewKeys {
		clear(fixPreviewAuthorizations.entries)
	}
	fixPreviewAuthorizations.entries[key] = now.Add(fixPreviewTTL)
	fixPreviewAuthorizations.Unlock()
	return nil
}

func consumeFixPreviewAuthorization(sourceRoot string, r *http.Request, files []servefix.FileSelection) bool {
	key, err := fixPreviewAuthorizationKey(sourceRoot, r, files)
	if err != nil {
		return false
	}
	now := time.Now()
	fixPreviewAuthorizations.Lock()
	defer fixPreviewAuthorizations.Unlock()
	pruneFixPreviewAuthorizationsLocked(now)
	expires, ok := fixPreviewAuthorizations.entries[key]
	if !ok || now.After(expires) {
		delete(fixPreviewAuthorizations.entries, key)
		return false
	}
	delete(fixPreviewAuthorizations.entries, key)
	return true
}

func pruneFixPreviewAuthorizationsLocked(now time.Time) {
	for key, expires := range fixPreviewAuthorizations.entries {
		if now.After(expires) {
			delete(fixPreviewAuthorizations.entries, key)
		}
	}
}

func fixPreviewAuthorizationKey(sourceRoot string, r *http.Request, files []servefix.FileSelection) (string, error) {
	root, err := filepath.Abs(sourceRoot)
	if err != nil {
		return "", err
	}
	normalized := make([]servefix.FileSelection, len(files))
	copy(normalized, files)
	for i := range normalized {
		normalized[i].Path = filepath.ToSlash(filepath.Clean(normalized[i].Path))
	}
	payload, err := json.Marshal(normalized)
	if err != nil {
		return "", err
	}
	client := r.RemoteAddr
	if host, _, splitErr := net.SplitHostPort(r.RemoteAddr); splitErr == nil {
		client = host
	}
	digest := sha256.New()
	_, _ = digest.Write([]byte(root))
	_, _ = digest.Write([]byte{0})
	_, _ = digest.Write([]byte(client))
	_, _ = digest.Write([]byte{0})
	_, _ = digest.Write([]byte(r.Host))
	_, _ = digest.Write([]byte{0})
	_, _ = digest.Write(payload)
	return fmt.Sprintf("%x", digest.Sum(nil)), nil
}

func appliedFixWrites(sourceRoot string, result servefix.BatchResult) []servefix.FilePreview {
	previews := make([]servefix.FilePreview, 0, len(result.Files))
	for i := range result.Files {
		file := result.Files[i]
		if file.Status != "applied" {
			continue
		}
		path := filepath.Join(sourceRoot, filepath.FromSlash(file.Path))
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		// The watcher only needs the exact bytes Markata successfully wrote.
		// Setting Before and After to the replacement prevents a later user undo
		// from being mistaken for the Control Center's own filesystem event.
		previews = append(previews, servefix.FilePreview{
			Path: file.Path, Before: string(content), After: string(content), Edits: append([]servefix.Edit(nil), file.Applied...),
		})
	}
	return previews
}

func decodeJSONLimit(w http.ResponseWriter, r *http.Request, dst any, limit int64) error {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		return fmt.Errorf("expected application/json")
	}
	r.Body = http.MaxBytesReader(w, r.Body, limit)
	return json.NewDecoder(r.Body).Decode(dst)
}
