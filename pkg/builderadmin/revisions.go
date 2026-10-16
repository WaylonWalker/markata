package builderadmin

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/WaylonWalker/markata-go/pkg/diagnostics"
	"github.com/WaylonWalker/markata-go/pkg/sourcegit"
)

const (
	revisionUnknown = "unknown"
	revisionBuilt   = "built"
	revisionUnbuilt = "unbuilt"
)

// RevisionStatus compares the running Builder executable with successful build
// provenance. Unknown is intentional for legacy records without reliable SHA.
type RevisionStatus struct {
	BuilderVersion    string    `json:"builder_version,omitempty"`
	BuilderCommit     string    `json:"builder_commit,omitempty"`
	SourceCommit      string    `json:"source_commit,omitempty"`
	LastBuiltCommit   string    `json:"last_built_commit,omitempty"`
	LastBuiltSource   string    `json:"last_built_source,omitempty"`
	LastBuiltAt       time.Time `json:"last_built_at,omitempty"`
	LastBuildID       string    `json:"last_build_id,omitempty"`
	LiveBuilderCommit string    `json:"live_builder_commit,omitempty"`
	LiveSourceCommit  string    `json:"live_source_commit,omitempty"`
	SourceChanged     bool      `json:"source_changed"`
	Status            string    `json:"status"`
}

// reliableBuilderRevision must not treat dev, none, branch names or version tags
// as comparable commits. Git (and Go VCS stamps) use hexadecimal SHA values.
func reliableBuilderRevision(raw string) string {
	raw = strings.TrimSpace(raw)
	if len(raw) < 7 || len(raw) > 64 {
		return ""
	}
	for _, c := range raw {
		if c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F' {
			continue
		}
		return ""
	}
	return strings.ToLower(raw)
}

// compareBuilderRevisions requires positive evidence before saying a newly
// deployed revision has not been built. Rollbacks do not count as builds.
func compareBuilderRevisions(current string, state State, publishedCommit string) RevisionStatus {
	result := RevisionStatus{BuilderCommit: reliableBuilderRevision(current), Status: revisionUnknown}
	var latest *BuildRecord
	for i := range state.Builds {
		record := &state.Builds[i]
		if record.Kind == "rollback" || record.RollbackRelease != "" || record.Status != "success" {
			continue
		}
		if latest == nil || record.FinishedAt.After(latest.FinishedAt) {
			latest = record
		}
		if result.BuilderCommit != "" && reliableBuilderRevision(record.BuilderCommit) == result.BuilderCommit {
			result.Status = revisionBuilt
		}
	}
	if latest != nil {
		result.LastBuiltCommit = reliableBuilderRevision(latest.BuilderCommit)
		result.LastBuiltSource = reliableBuilderRevision(latest.SourceCommit)
		result.LastBuiltAt = latest.FinishedAt
		result.LastBuildID = latest.ID
	}
	// The current published release can be older than the latest successful
	// build following a rollback. Its artifact is still build provenance.
	result.LiveBuilderCommit = reliableBuilderRevision(publishedCommit)
	if result.BuilderCommit == "" {
		return result
	}
	if result.LiveBuilderCommit == result.BuilderCommit {
		result.Status = revisionBuilt
	}
	if result.Status == revisionUnknown && (result.LastBuiltCommit != "" || result.LiveBuilderCommit != "") {
		result.Status = revisionUnbuilt
	}
	return result
}

// startupRevisionBuildReason compares against the latest successful real
// build, not any historical build of the same revision. Rollback is not a
// build. The live diagnostics artifact backs up legacy history records.
func startupRevisionBuildReason(current string, state State, publishedCommit string) (string, bool) {
	current = reliableBuilderRevision(current)
	if current == "" {
		return "", false // Mutable tags and dev/none are not revision identities.
	}
	var latest *BuildRecord
	for i := range state.Builds {
		record := &state.Builds[i]
		if record.Status != "success" || record.Kind == "rollback" || record.RollbackRelease != "" {
			continue
		}
		if latest == nil || record.FinishedAt.After(latest.FinishedAt) {
			latest = record
		}
	}
	previous := ""
	if latest != nil {
		previous = reliableBuilderRevision(latest.BuilderCommit)
	}
	if previous == "" {
		previous = reliableBuilderRevision(publishedCommit)
	}
	if previous == current {
		return "", false
	}
	if previous == "" {
		return "No previous successful build with reliable commit provenance; building " + current, true
	}
	return "Builder revision changed from " + previous + " to " + current, true
}

func (s *Service) revisionStatus(state State) RevisionStatus {
	builder, source := s.liveReleaseRevisions()
	result := compareBuilderRevisions(s.cfg.BuilderCommit, state, builder)
	result.BuilderVersion = strings.TrimSpace(s.cfg.BuilderVersion)
	result.LiveSourceCommit = source

	// HEAD is deliberately read without scanning every file. A dirty checkout
	// may still contain unbuilt content even if its Git HEAD is unchanged.
	ctx, cancel := context.WithTimeout(context.Background(), 450*time.Millisecond)
	defer cancel()
	if head, err := sourcegit.Head(ctx, s.cfg.SourceDir); err == nil {
		result.SourceCommit = reliableBuilderRevision(head)
	}
	result.SourceChanged = result.SourceCommit != "" && result.LiveSourceCommit != "" && result.SourceCommit != result.LiveSourceCommit
	return result
}

func (s *Service) liveReleaseRevisions() (builder, source string) {
	path := filepath.Join(s.cfg.SiteDir, "current", diagnostics.DefaultArtifactPath)
	info, err := os.Stat(path)
	if err != nil || info.Size() > 2<<20 {
		return "", ""
	}
	payload, err := os.ReadFile(path)
	if err != nil {
		return "", ""
	}
	artifact, err := diagnostics.ParseArtifact(payload)
	if err != nil {
		return "", ""
	}
	builder = reliableBuilderRevision(artifact.Generator.Commit)
	if artifact.Source != nil {
		source = reliableBuilderRevision(artifact.Source.Commit)
	}
	return builder, source
}
