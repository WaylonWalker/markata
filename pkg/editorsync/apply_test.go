package editorsync

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"sync"
	"testing"
)

func newPrivateTestEngine(t *testing.T, db *sql.DB) *Engine {
	t.Helper()
	engine, err := NewEngine(db, func(_ context.Context, subject, site string) error {
		if subject == "" || site != "family-blog" {
			return ErrWriteForbidden
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return engine
}

func TestApplyPrivateMutationIsDurableAndIdempotent(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "sync.sqlite")
	db, err := OpenSQLite(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	m := testMutation()
	m.BaseRevision = 0
	engine := newPrivateTestEngine(t, db)
	first, err := engine.Apply(ctx, "subject-42", m)
	if err != nil {
		t.Fatal(err)
	}
	if first.Revision != 1 || first.Sequence != 1 || first.PayloadDigest != m.Digest() {
		t.Fatalf("bad first receipt: %+v", first)
	}
	replay, err := engine.Apply(ctx, "subject-42", m)
	if err != nil || replay != first {
		t.Fatalf("replay must return identical receipt: %+v, %v", replay, err)
	}
	changed := m
	changed.Source = []byte("# Different\n")
	if _, err := engine.Apply(ctx, "subject-42", changed); !errors.Is(err, ErrMutationIDReuse) {
		t.Fatalf("UUID reuse must reject changed payload: %v", err)
	}
	m2 := m
	m2.BaseRevision = 1
	m2.MutationID = "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafb"
	m2.Source = []byte("# Updated\n")
	second, err := engine.Apply(ctx, "subject-42", m2)
	if err != nil || second.Revision != 2 || second.Sequence != 2 {
		t.Fatalf("update failed: %+v, %v", second, err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = OpenSQLite(ctx, path)
	if err != nil {
		t.Fatalf("restart SQLite: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("close SQLite: %v", err)
		}
	}()
	engine = newPrivateTestEngine(t, db)
	replay, err = engine.Apply(ctx, "subject-42", m)
	if err != nil || replay != first {
		t.Fatalf("receipt lost after restart: %+v, %v", replay, err)
	}
	var revisions, changes, receipts int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM editor_document_versions").Scan(&revisions); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM editor_changes").Scan(&changes); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM editor_mutation_receipts").Scan(&receipts); err != nil {
		t.Fatal(err)
	}
	if revisions != 2 || changes != 2 || receipts != 2 {
		t.Fatalf("retry duplicated data: versions %d, changes %d, receipts %d", revisions, changes, receipts)
	}
	stale := m
	stale.MutationID = "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafc"
	stale.Source = []byte("# Stale\n")
	var conflict *ConflictError
	if _, err := engine.Apply(ctx, "subject-42", stale); !errors.As(err, &conflict) || conflict.CurrentRevision != 2 {
		t.Fatalf("expected recoverable conflict at revision 2: %v", err)
	}
}

func TestApplyPrivateMutationEnforcesOwnerAndSite(t *testing.T) {
	ctx := context.Background()
	db, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "sync.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("close SQLite: %v", err)
		}
	}()
	engine := newPrivateTestEngine(t, db)
	m := testMutation()
	m.BaseRevision = 0
	if _, err := engine.Apply(ctx, "forged-subject", m); !errors.Is(err, ErrWriteForbidden) {
		t.Fatalf("forged principal must fail: %v", err)
	}
	m.SiteID = "other-site"
	if _, err := engine.Apply(ctx, m.SubjectID, m); !errors.Is(err, ErrWriteForbidden) {
		t.Fatalf("ungranted site must fail: %v", err)
	}
	m.SiteID = "family-blog"
	m.Visibility = SharedDraft
	if _, err := engine.Apply(ctx, m.SubjectID, m); !errors.Is(err, ErrUnsupportedVisibility) {
		t.Fatalf("shared publishing not supported in private POC: %v", err)
	}
	m.Visibility = Private
	if _, err := engine.Apply(ctx, m.SubjectID, m); err != nil {
		t.Fatal(err)
	}
	bob := m
	bob.SubjectID = "bob"
	bob.BaseRevision = 1
	bob.MutationID = "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafb"
	if _, err := engine.Apply(ctx, "bob", bob); !errors.Is(err, ErrWriteForbidden) {
		t.Fatalf("private document not isolated: %v", err)
	}
}

func TestApplyRollbackLeavesNoAcknowledgedMutation(t *testing.T) {
	ctx := context.Background()
	db, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "sync.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("close SQLite: %v", err)
		}
	}()
	_, err = db.ExecContext(ctx, `CREATE TRIGGER block_sync_receipt BEFORE INSERT ON editor_mutation_receipts
		BEGIN SELECT RAISE(ABORT, 'simulated persistence failure'); END`)
	if err != nil {
		t.Fatal(err)
	}
	engine := newPrivateTestEngine(t, db)
	m := testMutation()
	m.BaseRevision = 0
	if _, err := engine.Apply(ctx, m.SubjectID, m); err == nil {
		t.Fatal("failure before durable receipt was incorrectly acknowledged")
	}
	for _, table := range []string{"editor_sites", "editor_documents", "editor_document_versions", "editor_changes", "editor_mutation_receipts"} {
		// table names are fixed trusted literals from this test, never user data.
		var count int
		if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 0 {
			t.Errorf("%s retained %d rows after rolled-back mutation", table, count)
		}
	}
}

func TestApplyConcurrentCreatePreservesSingleWriter(t *testing.T) {
	ctx := context.Background()
	db, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "sync.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("close SQLite: %v", err)
		}
	}()
	engine := newPrivateTestEngine(t, db)
	commands := []Mutation{testMutation(), testMutation()}
	commands[0].BaseRevision = 0
	commands[1].BaseRevision = 0
	commands[1].MutationID = "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafe"
	results := make([]error, len(commands))
	var group sync.WaitGroup
	for i := range commands {
		group.Add(1)
		go func(i int) {
			defer group.Done()
			_, results[i] = engine.Apply(ctx, commands[i].SubjectID, commands[i])
		}(i)
	}
	group.Wait()
	accepted, conflicts := 0, 0
	for _, err := range results {
		switch {
		case err == nil:
			accepted++
		case errors.Is(err, ErrRevisionConflict):
			conflicts++
		default:
			t.Fatalf("unexpected concurrent result: %v", err)
		}
	}
	if accepted != 1 || conflicts != 1 {
		t.Fatalf("expected 1 accepted + 1 conflict; got %d accepted, %d conflicts", accepted, conflicts)
	}
}

func TestApplyDeletionWithEmptySource(t *testing.T) {
	ctx := context.Background()
	db, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "delete.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("close SQLite: %v", err)
		}
	}()
	engine := newPrivateTestEngine(t, db)
	m := testMutation()
	m.BaseRevision = 0
	if _, err := engine.Apply(ctx, m.SubjectID, m); err != nil {
		t.Fatal(err)
	}
	m.MutationID = "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafb"
	m.BaseRevision = 1
	m.Deleted = true
	m.Source = nil
	if _, err := engine.Apply(ctx, m.SubjectID, m); err != nil {
		t.Fatalf("empty-body deletion must not write SQL NULL: %v", err)
	}
	var deleted, sourceLen, sourceNull int
	if err := db.QueryRowContext(ctx,
		"SELECT deleted, length(source), source IS NULL FROM editor_documents WHERE site_id=? AND document_id=?",
		m.SiteID, m.DocumentID,
	).Scan(&deleted, &sourceLen, &sourceNull); err != nil {
		t.Fatal(err)
	}
	if deleted != 1 || sourceLen != 0 || sourceNull != 0 {
		t.Fatalf("incorrect tombstone deleted=%d len=%d null=%d", deleted, sourceLen, sourceNull)
	}
}
