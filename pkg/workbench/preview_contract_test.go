package workbench

import (
	"encoding/json"
	"net/http"
	"testing"
)

// The server payload is intentionally additive so old preview clients remain
// compatible while a WASM transport can reject mismatched schemas.
func TestHandler_PreviewContractV1(t *testing.T) {
	h := newHandler(t)
	for _, tc := range []struct {
		name, body string
		wantDocument bool
	}{
		{"first frame", `{"markdown":"# First\n\nSecond paragraph","include_document":true}`, true},
		{"fragment", `{"markdown":"# First\n\nSecond paragraph","include_document":false}`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := request(h, tc.body)
			if w.Code != http.StatusOK {
				t.Fatalf("POST /render returned %d: %s", w.Code, w.Body.String())
			}
			var payload renderResponse
			if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
				t.Fatal(err)
			}
			if payload.SchemaVersion != 1 {
				t.Fatalf("preview schema version = %d, want 1", payload.SchemaVersion)
			}
			if (payload.Document != nil) != tc.wantDocument {
				t.Fatalf("initial document presence = %t, want %t", payload.Document != nil, tc.wantDocument)
			}
			if len(payload.SourceLines) != 2 || payload.SourceLines[0] != 0 || payload.SourceLines[1] != 2 {
				t.Errorf("source lines = %v, want [0, 2]", payload.SourceLines)
			}
			if payload.FeatureHeadHTML == "" && payload.BodyEndHTML == "" {
				// Both can be empty on a plain Markdown document; that's valid.
			}
		})
	}
}
