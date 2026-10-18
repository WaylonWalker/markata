package builderadmin

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func TestStartupRevisionBuildReason(t *testing.T) {
	at := time.Date(2026, time.October, 16, 12, 0, 0, 0, time.UTC)
	old := BuildRecord{Kind: "build", Status: "success", BuilderCommit: revisionOldSHA, FinishedAt: at}
	newer := BuildRecord{Kind: "build", Status: "success", BuilderCommit: revisionNewSHA, FinishedAt: at.Add(time.Hour)}
	tests := []struct {
		name, current, published string
		builds                   []BuildRecord
		want                     bool
	}{
		{"upgrade", revisionNewSHA, "", []BuildRecord{old}, true},
		{"unchanged restart", revisionOldSHA, "", []BuildRecord{old}, false},
		{"latest successful, not older history", revisionOldSHA, "", []BuildRecord{old, newer}, true},
		{"matching latest build", revisionNewSHA, "", []BuildRecord{old, newer}, false},
		{"failed new revision retries", revisionNewSHA, "", []BuildRecord{old, {Kind: "build", Status: "failed", BuilderCommit: revisionNewSHA}}, true},
		{"rollback does not count", revisionNewSHA, "", []BuildRecord{old, {Kind: "rollback", Status: "success", BuilderCommit: revisionNewSHA}}, true},
		{"legacy state matches live artifact", revisionNewSHA, revisionNewSHA, []BuildRecord{{Kind: "build", Status: "success"}}, false},
		{"legacy state differs from live artifact", revisionNewSHA, revisionOldSHA, []BuildRecord{{Kind: "build", Status: "success"}}, true},
		{"first deployment establishes baseline", revisionNewSHA, "", nil, true},
		{"legacy unknown establishes baseline", revisionNewSHA, "", []BuildRecord{{Kind: "build", Status: "success"}}, true},
		{"unknown running SHA must not guess", "none", revisionOldSHA, []BuildRecord{old}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reason, got := startupRevisionBuildReason(tc.current, State{Builds: tc.builds}, tc.published)
			if got != tc.want {
				t.Fatalf("build needed = %t, want %t (reason: %q)", got, tc.want, reason)
			}
			if got && reason == "" {
				t.Fatal("expected an auditable trigger reason")
			}
		})
	}
}

func TestStartupRevisionBuildQueue(t *testing.T) {
	for _, tc := range []struct {
		name    string
		enabled bool
		pending []QueuedOperation
		builds  []BuildRecord
		want    int
	}{
		{"changed", true, nil, []BuildRecord{{Kind: "build", Status: "success", BuilderCommit: revisionOldSHA}}, 1},
		{"disabled", false, nil, []BuildRecord{{Kind: "build", Status: "success", BuilderCommit: revisionOldSHA}}, 0},
		{"unchanged", true, nil, []BuildRecord{{Kind: "build", Status: "success", BuilderCommit: revisionNewSHA}}, 0},
		{"ordinary build already pending", true, []QueuedOperation{{ID: "build-existing", Kind: "build"}}, []BuildRecord{{Kind: "build", Status: "success", BuilderCommit: revisionOldSHA}}, 1},
		{"webhook sync may finish without building", true, []QueuedOperation{{ID: "build-webhook", Kind: "build", SyncSource: true}}, []BuildRecord{{Kind: "build", Status: "success", BuilderCommit: revisionOldSHA}}, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			svc, err := New(Config{
				SourceDir:          dir,
				SiteDir:            dir,
				HistoryDir:         filepath.Join(dir, "history"),
				BuilderCommit:      revisionNewSHA,
				AutoBuildOnUpgrade: tc.enabled,
			})
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = svc.leaderLock.Close() })
			svc.state.Builds = tc.builds
			svc.state.Queue = tc.pending
			svc.resumeLeaderSession()
			if got := len(svc.state.Queue); got != tc.want {
				t.Fatalf("persisted queue length = %d, want %d: %+v", got, tc.want, svc.state.Queue)
			}
			if got := len(svc.queueCh); got != tc.want {
				t.Fatalf("worker queue length = %d, want %d", got, tc.want)
			}
			if tc.want > len(tc.pending) {
				last := svc.state.Queue[len(svc.state.Queue)-1]
				if last.TriggerType != "startup-revision" || last.Detail == "" {
					t.Fatalf("startup build missing trigger provenance: %+v", last)
				}
			}
		})
	}
}

func TestNewLeaderReloadsSuccessfulBuildHistory(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{
		SourceDir:          dir,
		SiteDir:            filepath.Join(dir, "site"),
		HistoryDir:         filepath.Join(dir, "history"),
		BuilderCommit:      revisionNewSHA,
		AutoBuildOnUpgrade: true,
	}
	standby, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = standby.leaderLock.Close() })
	// This build completes while the replacement pod is standing by.
	active, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	active.finishBuild(BuildRecord{ID: "finished-after-standby-start", Kind: "build", Status: "success", BuilderCommit: revisionNewSHA})
	_ = active.leaderLock.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	standby.tryBecomeLeader(ctx)
	defer standby.releaseLeadership()
	if !standby.isLeader() {
		t.Fatal("standby failed to take leadership")
	}
	state := standby.snapshotState()
	if len(state.Builds) != 1 || state.Builds[0].ID != "finished-after-standby-start" {
		t.Fatalf("leader retained stale state: %+v", state.Builds)
	}
	if len(state.Queue) != 0 {
		t.Fatalf("enqueued duplicate startup build after recent success: %+v", state.Queue)
	}
	found := false
	for _, job := range standby.controlRuntime.Snapshot().Jobs {
		if job.ID == "finished-after-standby-start" {
			found = true
		}
	}
	if !found {
		t.Fatal("leader control snapshot did not refresh from persisted history")
	}
}
