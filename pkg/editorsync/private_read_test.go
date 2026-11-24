package editorsync

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func newPrivateReadFixture(t *testing.T) (*Engine, *PrivateReader) {
	t.Helper()
	db, err := OpenSQLite(context.Background(), filepath.Join(t.TempDir(), "private-read.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close private read test database: %v", err)
		}
	})
	engine := newPrivateTestEngine(t, db)
	reader, err := NewPrivateReader(db, func(_ context.Context, subject, site string) error {
		if site != "family-blog" || (subject != "subject-42" && subject != "bob") {
			return ErrReadForbidden
		}
		return nil
	}, []byte("0123456789abcdef0123456789abcdef"), time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	return engine, reader
}

func savePrivateTestDocument(t *testing.T, engine *Engine, subject, doc, id string, base uint64, body string) {
	t.Helper()
	m := testMutation()
	m.SubjectID = subject
	m.DocumentID = doc
	m.MutationID = id
	m.BaseRevision = base
	m.Source = []byte(body)
	if _, err := engine.Apply(context.Background(), subject, m); err != nil {
		t.Fatalf("save %s: %v", doc, err)
	}
}

func TestPrivateSnapshotPagesStayAtOneHighWatermark(t *testing.T) {
	ctx := context.Background()
	engine, reader := newPrivateReadFixture(t)
	savePrivateTestDocument(t, engine, "subject-42", "post-a",
		"8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafa", 0, "# A")
	savePrivateTestDocument(t, engine, "subject-42", "post-b",
		"8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafb", 0, "# B")
	savePrivateTestDocument(t, engine, "bob", "post-bob",
		"8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafc", 0, "# Bob")
	first, err := reader.Snapshot(ctx, "subject-42", "family-blog", "", 1)
	if err != nil || len(first.Items) != 1 || first.Items[0].DocumentID != "post-a" ||
		first.NextCursor == "" || first.ChangesCursor != "" {
		t.Fatalf("unexpected first page: %+v, %v", first, err)
	}
	savePrivateTestDocument(t, engine, "subject-42", "post-b",
		"8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafd", 1, "# B updated")
	savePrivateTestDocument(t, engine, "subject-42", "post-c",
		"8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafe", 0, "# C new")
	second, err := reader.Snapshot(ctx, "subject-42", "family-blog", first.NextCursor, 1)
	if err != nil || len(second.Items) != 1 || second.Items[0].DocumentID != "post-b" ||
		second.Items[0].Revision != 1 || second.NextCursor != "" || second.ChangesCursor == "" {
		t.Fatalf("snapshot must not include later revisions: %+v, %v", second, err)
	}
	oldSource, err := reader.SourceAtRevision(ctx, "subject-42", "family-blog", "post-b", 1)
	if err != nil || string(oldSource) != "# B" {
		t.Fatalf("cannot hydrate old snapshot revision: %q, %v", oldSource, err)
	}
	if _, err := reader.SourceAtRevision(ctx, "bob", "family-blog", "post-b", 1); !errors.Is(err, ErrPrivateVersionNotFound) {
		t.Fatalf("other subject accessed a private source: %v", err)
	}
	changes, err := reader.Changes(ctx, "subject-42", "family-blog", second.ChangesCursor, 1)
	if err != nil || len(changes.Items) != 1 || changes.Items[0].DocumentID != "post-b" ||
		changes.Items[0].Revision != 2 {
		t.Fatalf("missed update made during snapshot: %+v, %v", changes, err)
	}
	changes, err = reader.Changes(ctx, "subject-42", "family-blog", changes.NextCursor, 1)
	if err != nil || len(changes.Items) != 1 || changes.Items[0].DocumentID != "post-c" {
		t.Fatalf("missed new document after snapshot: %+v, %v", changes, err)
	}
	changes, err = reader.Changes(ctx, "subject-42", "family-blog", changes.NextCursor, 1)
	if err != nil || len(changes.Items) != 0 {
		t.Fatalf("duplicate change events: %+v, %v", changes, err)
	}
}

func TestPrivateCursorsFailClosed(t *testing.T) {
	ctx := context.Background()
	engine, reader := newPrivateReadFixture(t)
	savePrivateTestDocument(t, engine, "subject-42", "post-a",
		"8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafa", 0, "# A")
	savePrivateTestDocument(t, engine, "subject-42", "post-b",
		"8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafb", 0, "# B")
	page, err := reader.Snapshot(ctx, "subject-42", "family-blog", "", 1)
	if err != nil || page.NextCursor == "" {
		t.Fatalf("expected snapshot token: %+v, %v", page, err)
	}
	token := page.NextCursor
	if strings.Contains(token, "subject-42") || strings.Contains(token, "family-blog") ||
		strings.Contains(token, "post-a") {
		t.Fatal("snapshot cursor leaked its plaintext metadata")
	}
	if _, err := reader.Snapshot(ctx, "bob", "family-blog", token, 1); !errors.Is(err, ErrInvalidSyncCursor) {
		t.Fatalf("cross-subject cursor accepted: %v", err)
	}
	if _, err := reader.Changes(ctx, "subject-42", "family-blog", token, 1); !errors.Is(err, ErrInvalidSyncCursor) {
		t.Fatalf("wrong-purpose cursor accepted: %v", err)
	}
	tampered := token[:len(token)-3] + "abc"
	if _, err := reader.Snapshot(ctx, "subject-42", "family-blog", tampered, 1); !errors.Is(err, ErrInvalidSyncCursor) {
		t.Fatalf("tampered cursor accepted: %v", err)
	}
	if _, err := reader.Snapshot(ctx, "subject-42", "other-blog", token, 1); !errors.Is(err, ErrReadForbidden) {
		t.Fatalf("ungranted site accepted: %v", err)
	}
	originalAuthorize := reader.authorize
	reader.authorize = func(context.Context, string, string) error { return nil }
	if _, err := reader.Snapshot(ctx, "subject-42", "other-blog", token, 1); !errors.Is(err, ErrInvalidSyncCursor) {
		t.Fatalf("cross-site cursor accepted after site permission granted: %v", err)
	}
	reader.authorize = originalAuthorize
	frozen := time.Unix(1_800_000_000, 0)
	reader.now = func() time.Time { return frozen }
	issued, err := reader.Snapshot(ctx, "subject-42", "family-blog", "", 1)
	if err != nil {
		t.Fatal(err)
	}
	reader.now = func() time.Time { return frozen.Add(2 * time.Hour) }
	if _, err := reader.Snapshot(ctx, "subject-42", "family-blog", issued.NextCursor, 1); !errors.Is(err, ErrExpiredSyncCursor) {
		t.Fatalf("expired cursor accepted: %v", err)
	}
}

func TestPrivateChangesIncludeTombstonesAndSkipOtherPeople(t *testing.T) {
	ctx := context.Background()
	engine, reader := newPrivateReadFixture(t)
	snapshot, err := reader.Snapshot(ctx, "subject-42", "family-blog", "", 10)
	if err != nil || snapshot.ChangesCursor == "" {
		t.Fatalf("empty snapshot not resumable: %+v, %v", snapshot, err)
	}
	savePrivateTestDocument(t, engine, "bob", "other-person",
		"8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafa", 0, "hidden")
	savePrivateTestDocument(t, engine, "subject-42", "post-delete",
		"8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafb", 0, "to delete")
	m := testMutation()
	m.MutationID = "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafc"
	m.DocumentID = "post-delete"
	m.BaseRevision = 1
	m.Deleted = true
	m.Source = nil
	if _, err := engine.Apply(ctx, m.SubjectID, m); err != nil {
		t.Fatal(err)
	}
	page, err := reader.Changes(ctx, "subject-42", "family-blog", snapshot.ChangesCursor, 1)
	if err != nil || len(page.Items) != 1 || page.Items[0].Deleted {
		t.Fatalf("missing private initial change: %+v, %v", page, err)
	}
	next, err := reader.Changes(ctx, "subject-42", "family-blog", page.NextCursor, 1)
	if err != nil || len(next.Items) != 1 || !next.Items[0].Deleted || next.Items[0].Revision != 2 {
		t.Fatalf("missing private tombstone: %+v, %v", next, err)
	}
	none, err := reader.Changes(ctx, "subject-42", "family-blog", next.NextCursor, 1)
	if err != nil || len(none.Items) != 0 {
		t.Fatalf("change cursor did not advance: %+v, %v", none, err)
	}
}

func TestPrivateSnapshotRechecksRevokedReadAccess(t *testing.T) {
	ctx := context.Background()
	engine, reader := newPrivateReadFixture(t)
	savePrivateTestDocument(t, engine, "subject-42", "post-a",
		"8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafa", 0, "# A")
	savePrivateTestDocument(t, engine, "subject-42", "post-b",
		"8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafb", 0, "# B")
	page, err := reader.Snapshot(ctx, "subject-42", "family-blog", "", 1)
	if err != nil || page.NextCursor == "" {
		t.Fatalf("expected paginated checkpoint: %+v, %v", page, err)
	}
	reader.authorize = func(context.Context, string, string) error { return ErrReadForbidden }
	if _, err := reader.Snapshot(ctx, "subject-42", "family-blog", page.NextCursor, 1); !errors.Is(err, ErrReadForbidden) {
		t.Fatalf("revoked subject got later snapshot page: %v", err)
	}
	if _, err := reader.SourceAtRevision(ctx, "subject-42", "family-blog", "post-a", 1); !errors.Is(err, ErrReadForbidden) {
		t.Fatalf("revoked subject hydrated private source: %v", err)
	}
}

func TestPrivateSnapshotRejectsRestoredOlderCheckpoint(t *testing.T) {
	ctx := context.Background()
	engine, reader := newPrivateReadFixture(t)
	savePrivateTestDocument(t, engine, "subject-42", "post-a",
		"8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafa", 0, "# A")
	savePrivateTestDocument(t, engine, "subject-42", "post-b",
		"8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafb", 0, "# B")
	page, err := reader.Snapshot(ctx, "subject-42", "family-blog", "", 1)
	if err != nil || page.NextCursor == "" {
		t.Fatalf("expected cursor before restored checkpoint: %+v, %v", page, err)
	}
	if _, err := engine.db.ExecContext(ctx,
		"UPDATE editor_sites SET last_seq=0 WHERE site_id='family-blog'",
	); err != nil {
		t.Fatal(err)
	}
	if _, err := reader.Snapshot(ctx, "subject-42", "family-blog", page.NextCursor, 1); !errors.Is(err, ErrExpiredSyncCursor) {
		t.Fatalf("restored old checkpoint should require new snapshot: %v", err)
	}
}
