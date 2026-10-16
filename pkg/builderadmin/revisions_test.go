package builderadmin

import (
	"strings"
	"testing"
	"time"
)

const (
	revisionOldSHA = "0123456789abcdef0123456789abcdef01234567"
	revisionNewSHA = "abcdef0123456789abcdef0123456789abcdef01"
)

func TestBuilderRevisionComparison(t *testing.T) {
	at := time.Date(2026, time.October, 15, 14, 0, 0, 0, time.UTC)
	prior := BuildRecord{ID: "build-one", Kind: "build", Status: "success", BuilderCommit: revisionOldSHA, FinishedAt: at}
	current := BuildRecord{ID: "build-two", Kind: "build", Status: "success", BuilderCommit: revisionNewSHA, FinishedAt: at.Add(time.Hour)}
	tests := []struct {
		name, revision, published, want string
		builds                          []BuildRecord
	}{
		{"new not built", revisionNewSHA, revisionOldSHA, revisionUnbuilt, []BuildRecord{prior}},
		{"current built", revisionNewSHA, revisionOldSHA, revisionBuilt, []BuildRecord{current, prior}},
		{"built live artifact", revisionNewSHA, revisionNewSHA, revisionBuilt, []BuildRecord{prior}},
		{"unknown legacy build", revisionNewSHA, "", revisionUnknown, []BuildRecord{{Status: "success"}}},
		{"unknown builder", "none", revisionOldSHA, revisionUnknown, []BuildRecord{prior}},
		{"failed build does not count", revisionNewSHA, revisionOldSHA, revisionUnbuilt, []BuildRecord{prior, {Kind: "build", Status: "failed", BuilderCommit: revisionNewSHA}}},
		{"rollback is not a build", revisionNewSHA, revisionOldSHA, revisionUnbuilt, []BuildRecord{prior, {Kind: "rollback", Status: "success", BuilderCommit: revisionNewSHA}}},
		{"same historical commit was built", revisionOldSHA, revisionNewSHA, revisionBuilt, []BuildRecord{prior, current}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			state := State{Builds: tc.builds}
			result := compareBuilderRevisions(tc.revision, state, tc.published)
			if result.Status != tc.want {
				t.Fatalf("status = %q, want %q; result %+v", result.Status, tc.want, result)
			}
		})
	}
}

func TestRevisionStringValidation(t *testing.T) {
	for _, raw := range []string{"", "none", "dev", "revision-main", "123", "cafegggf"} {
		if value := reliableBuilderRevision(raw); value != "" {
			t.Errorf("invalid revision %q accepted as %q", raw, value)
		}
	}
	if got := reliableBuilderRevision(" ABCDEF0123456 "); got != "abcdef0123456" {
		t.Fatalf("normalized SHA = %q", got)
	}
}

func TestBuilderIndexShowsRevisionComparison(t *testing.T) {
	for _, text := range []string{
		`id="revision-banner"`,
		`id="revision-current"`,
		`id="revision-last-built"`,
		`id="revision-source"`,
		`renderRevision(payload.revision, state);`,
		`New Builder revision has not built yet`,
	} {
		if !strings.Contains(indexHTML, text) {
			t.Errorf("Builder dashboard missing revision marker %q", text)
		}
	}
	for _, text := range []string{
		`id="report-builder-revision"`,
		`id="report-source-revision"`,
	} {
		if !strings.Contains(buildDetailHTML, text) {
			t.Errorf("build details missing revision marker %q", text)
		}
	}
}
