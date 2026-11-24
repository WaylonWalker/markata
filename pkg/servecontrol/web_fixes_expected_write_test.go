package servecontrol

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/servefix"
)

func TestBatchFixApplyRequiresMatchingPreview(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "post.md")
	original := "# Heading\n[link](//example.com)\n"
	if err := os.WriteFile(path, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	plan, err := servefix.PlanFile(root, "post.md")
	if err != nil {
		t.Fatal(err)
	}
	selection := servefix.Selection{All: true, SafeOnly: true}
	applyBody, err := json.Marshal(batchFixApplyRequest{Files: []servefix.FileSelection{{
		Path: "post.md", Digest: plan.Digest, Selection: selection,
	}}})
	if err != nil {
		t.Fatal(err)
	}

	runtime := NewRuntime()
	runtime.SetActionHandler(func(ActionRequest) error { return nil })
	handler := NewWebHandlerWithSourceRoot(runtime, root)

	response := serveLocalPost(handler, "/_markata/api/fixes/batch/apply", applyBody)
	if response.Code != http.StatusConflict {
		t.Fatalf("direct apply status=%d body=%s", response.Code, response.Body.String())
	}
	if got, err := os.ReadFile(path); err != nil || string(got) != original {
		t.Fatalf("direct apply changed source: %q err=%v", got, err)
	}

	previewBody, err := json.Marshal(batchFixPreviewRequest{Files: []servefix.FileSelection{{
		Path: "post.md", Selection: selection,
	}}})
	if err != nil {
		t.Fatal(err)
	}
	response = serveLocalPost(handler, "/_markata/api/fixes/batch/preview", previewBody)
	if response.Code != http.StatusOK {
		t.Fatalf("preview status=%d body=%s", response.Code, response.Body.String())
	}

	mismatchBody, err := json.Marshal(batchFixApplyRequest{Files: []servefix.FileSelection{{
		Path: "post.md", Digest: plan.Digest, Selection: servefix.Selection{All: true},
	}}})
	if err != nil {
		t.Fatal(err)
	}
	response = serveLocalPost(handler, "/_markata/api/fixes/batch/apply", mismatchBody)
	if response.Code != http.StatusConflict {
		t.Fatalf("mismatched apply status=%d body=%s", response.Code, response.Body.String())
	}

	response = serveLocalPost(handler, "/_markata/api/fixes/batch/apply", applyBody)
	if response.Code != http.StatusOK {
		t.Fatalf("previewed apply status=%d body=%s", response.Code, response.Body.String())
	}
	response = serveLocalPost(handler, "/_markata/api/fixes/batch/apply", applyBody)
	if response.Code != http.StatusConflict {
		t.Fatalf("reused preview status=%d body=%s", response.Code, response.Body.String())
	}
}

func TestBatchFixApplyHookExposesOnlyAppliedReplacementBytes(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "post.md")
	original := "# Heading\n"
	if err := os.WriteFile(path, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}

	plan, err := servefix.PlanFile(root, "post.md")
	if err != nil {
		t.Fatal(err)
	}
	selection := servefix.Selection{All: true, SafeOnly: true}

	var captured []servefix.FilePreview
	runtime := NewRuntime()
	runtime.SetActionHandler(func(ActionRequest) error { return nil })
	handler := NewWebHandlerWithFixHooks(runtime, root, func(previews []servefix.FilePreview) {
		captured = append([]servefix.FilePreview(nil), previews...)
	})

	previewBody, err := json.Marshal(batchFixPreviewRequest{Files: []servefix.FileSelection{{
		Path: "post.md", Selection: selection,
	}}})
	if err != nil {
		t.Fatal(err)
	}
	response := serveLocalPost(handler, "/_markata/api/fixes/batch/preview", previewBody)
	if response.Code != http.StatusOK {
		t.Fatalf("preview status=%d body=%s", response.Code, response.Body.String())
	}

	applyBody, err := json.Marshal(batchFixApplyRequest{Files: []servefix.FileSelection{{
		Path: "post.md", Digest: plan.Digest, Selection: selection,
	}}})
	if err != nil {
		t.Fatal(err)
	}
	response = serveLocalPost(handler, "/_markata/api/fixes/batch/apply", applyBody)
	if response.Code != http.StatusOK {
		t.Fatalf("apply status=%d body=%s", response.Code, response.Body.String())
	}
	if len(captured) != 1 {
		t.Fatalf("afterApply previews=%d, want 1", len(captured))
	}
	if captured[0].Before != captured[0].After {
		t.Fatalf("watcher hook exposed pre-fix bytes; before=%q after=%q", captured[0].Before, captured[0].After)
	}
	if captured[0].After == original {
		t.Fatal("watcher hook did not expose the replacement content")
	}

	updated, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(updated) != captured[0].After {
		t.Fatalf("applied content=%q, watcher expected=%q", string(updated), captured[0].After)
	}
}
