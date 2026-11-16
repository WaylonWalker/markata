package editorhub

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const orphanConflictID = "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafc"

func privateConflictRequest(subject, mutationID string) *http.Request {
	r := httptest.NewRequestWithContext(context.Background(), http.MethodGet,
		privateHTTPOrigin+"/api/v1/sites/family-blog/conflicts/"+mutationID, http.NoBody)
	if subject != "" {
		r = r.WithContext(context.WithValue(r.Context(), verifiedSubjectKey{}, subject))
	}
	return r
}

func TestHTTPConflictRecoveryIsOwnerPrivate(t *testing.T) {
	handler, _ := newHTTPMutationFixture(t)
	first := httptest.NewRecorder()
	handler.ServeHTTP(first, mutationRequest("alice", acceptedPrivateBody()))
	if first.Code != http.StatusOK {
		t.Fatalf("baseline document not saved: %d: %s", first.Code, first.Body.String())
	}
	stale := strings.Replace(acceptedPrivateBody(), privateMutationID, otherMutationID, 1)
	rejected := httptest.NewRecorder()
	handler.ServeHTTP(rejected, mutationRequest("alice", stale))
	if rejected.Code != http.StatusConflict {
		t.Fatalf("stale edit did not conflict: %d", rejected.Code)
	}
	owner := httptest.NewRecorder()
	handler.ServeHTTP(owner, privateConflictRequest("alice", otherMutationID))
	if owner.Code != http.StatusOK {
		t.Fatalf("original author cannot recover rejected edit: %d: %s", owner.Code, owner.Body.String())
	}
	if owner.Header().Get("Cache-Control") != "private, no-store" {
		t.Fatal("private conflict allowed caching")
	}
	var record struct {
		MutationID        string `json:"mutation_id"`
		DocumentID        string `json:"document_id"`
		StaleBaseRevision int64  `json:"stale_base_revision"`
		CurrentRevision   int64  `json:"current_revision"`
		Digest            string `json:"digest"`
		SourceBase64      string `json:"source_base64"`
		Deleted           bool   `json:"deleted"`
	}
	if err := json.Unmarshal(owner.Body.Bytes(), &record); err != nil {
		t.Fatal(err)
	}
	recovered, err := base64.StdEncoding.DecodeString(record.SourceBase64)
	if err != nil || string(recovered) != "# Hello from offline" ||
		record.MutationID != otherMutationID || record.DocumentID != "post-001" ||
		record.StaleBaseRevision != 0 || record.CurrentRevision != 1 ||
		len(record.Digest) != 64 || record.Deleted {
		t.Fatalf("rejected draft not recovered verbatim: %+v, %v", record, err)
	}
	for _, person := range []struct {
		subject string
		want    int
	}{
		{"eve", http.StatusNotFound},  // Has site write permission, not the rejected source.
		{"bob", http.StatusForbidden}, // Site read access alone is insufficient.
		{"", http.StatusUnauthorized},
		{"stranger", http.StatusForbidden},
	} {
		other := httptest.NewRecorder()
		handler.ServeHTTP(other, privateConflictRequest(person.subject, otherMutationID))
		if other.Code != person.want {
			t.Fatalf("%q read owner conflict: got %d, want %d", person.subject, other.Code, person.want)
		}
		if strings.Contains(other.Body.String(), "Hello from offline") {
			t.Fatal("private rejected draft leaked in unauthorized response")
		}
	}
}

func TestHTTPRecoversOrphanedPrivateMutationAndRejectsCrossSite(t *testing.T) {
	handler, _ := newHTTPMutationFixture(t)
	orphanBody := `{"document_id":"deleted-post","mutation_id":"` + orphanConflictID +
		`","base_revision":17,"source":"# Offline orphan content","deleted":false}`
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, mutationRequest("alice", orphanBody))
	if response.Code != http.StatusConflict {
		t.Fatalf("orphan mutation did not yield durable conflict: %d: %s", response.Code, response.Body.String())
	}
	for _, tc := range []struct {
		name   string
		change func(*http.Request)
		status int
	}{
		{"owner", func(*http.Request) {}, http.StatusOK},
		{"cross-site metadata", func(r *http.Request) { r.Header.Set("Sec-Fetch-Site", "cross-site") }, http.StatusForbidden},
		{"cross-site origin", func(r *http.Request) { r.Header.Set("Origin", "https://evil.example") }, http.StatusForbidden},
		{"duplicate origins", func(r *http.Request) {
			r.Header.Add("Origin", privateHTTPOrigin)
			r.Header.Add("Origin", privateHTTPOrigin)
		}, http.StatusForbidden},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := privateConflictRequest("alice", orphanConflictID)
			tc.change(r)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)
			if w.Code != tc.status {
				t.Fatalf("unexpected status %d, want %d: %s", w.Code, tc.status, w.Body.String())
			}
			if tc.status == http.StatusOK {
				var payload struct {
					Source string `json:"source_base64"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
					t.Fatal(err)
				}
				bytes, err := base64.StdEncoding.DecodeString(payload.Source)
				if err != nil || string(bytes) != "# Offline orphan content" {
					t.Fatalf("lost orphan content %q: %v", bytes, err)
				}
			}
		})
	}
	forged := privateConflictRequest("", orphanConflictID)
	forged.Header.Set("X-Hlab-User-Id", "alice")
	forged.Header.Set("X-Hlab-Auth-Method", "session")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, forged)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("forged identity authorized orphan conflict: %d", w.Code)
	}
}
