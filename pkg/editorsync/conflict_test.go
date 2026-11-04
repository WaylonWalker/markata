package editorsync

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
)

func TestPrivateConflictPersistsAcrossRestart(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "conflicts.sqlite")
	db, err := OpenSQLite(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	engine := newPrivateTestEngine(t, db)
	base := testMutation()
	base.BaseRevision = 0
	if _, err := engine.Apply(ctx, base.SubjectID, base); err != nil {
		t.Fatal(err)
	}
	stale := base
	stale.MutationID = "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafb"
	stale.Source = []byte("# Unsynced offline work\n")
	stale.BaseRevision = 0
	var conflict *ConflictError
	if _, err := engine.Apply(ctx, stale.SubjectID, stale); !errors.As(err, &conflict) ||
		!conflict.Preserved || conflict.CurrentRevision != 1 {
		t.Fatalf("stale work must persist before conflict: %v", err)
	}
	inspect := func(db *sql.DB) {
		t.Helper()
		var proposed, current []byte
		var baseRevision, currentRevision, deleted, changes int
		var digest string
		err := db.QueryRowContext(ctx, `SELECT proposed_source, stale_base_revision, current_revision,
			payload_digest, proposed_deleted FROM editor_conflicts
			WHERE site_id=? AND subject_id=? AND mutation_id=?`,
			stale.SiteID, stale.SubjectID, stale.MutationID,
		).Scan(&proposed, &baseRevision, &currentRevision, &digest, &deleted)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(proposed, stale.Source) || digest != stale.Digest() ||
			baseRevision != 0 || currentRevision != 1 || deleted != 0 {
			t.Fatal("stale source or intent not persisted exactly")
		}
		if err := db.QueryRowContext(ctx, `SELECT source FROM editor_documents
			WHERE site_id=? AND document_id=?`, base.SiteID, base.DocumentID).Scan(&current); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(current, base.Source) {
			t.Fatal("conflict overwrote live server document")
		}
		if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM editor_changes").Scan(&changes); err != nil {
			t.Fatal(err)
		}
		if changes != 1 {
			t.Fatalf("conflict must not create public change event, got %d", changes)
		}
	}
	inspect(db)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = OpenSQLite(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("close SQLite: %v", err)
		}
	}()
	inspect(db)
	engine = newPrivateTestEngine(t, db)
	if _, err := engine.Apply(ctx, stale.SubjectID, stale); !errors.As(err, &conflict) ||
		!conflict.Preserved || conflict.CurrentRevision != 1 {
		t.Fatalf("retry must report durable original conflict: %v", err)
	}
	changed := stale
	changed.Source = []byte("# Attempted overwrite of conflict\n")
	if _, err := engine.Apply(ctx, changed.SubjectID, changed); !errors.Is(err, ErrMutationIDReuse) {
		t.Fatalf("mutation UUID cannot replace conflict contents: %v", err)
	}
	inspect(db)
}

func TestPrivateDeletionConflictPreservesIntent(t *testing.T) {
	ctx := context.Background()
	db, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "delete-conflict.sqlite"))
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
	m.Deleted = true
	m.Source = nil
	_, err = engine.Apply(ctx, m.SubjectID, m)
	var conflict *ConflictError
	if !errors.As(err, &conflict) || !conflict.Preserved {
		t.Fatalf("delete conflict must be durable: %v", err)
	}
	var deleted int
	var digest string
	if err := db.QueryRowContext(ctx, `SELECT proposed_deleted, payload_digest FROM editor_conflicts
		WHERE site_id=? AND subject_id=? AND mutation_id=?`,
		m.SiteID, m.SubjectID, m.MutationID).Scan(&deleted, &digest); err != nil {
		t.Fatal(err)
	}
	if deleted != 1 || digest != m.Digest() {
		t.Fatal("delete conflict lost deletion intent")
	}
}

func TestConflictingAttemptCannotReadOtherPrivateDocuments(t *testing.T) {
	ctx := context.Background()
	db, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "isolation.sqlite"))
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
	other := m
	other.SubjectID = "bob"
	other.MutationID = "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafb"
	_, err = engine.Apply(ctx, other.SubjectID, other)
	if !errors.Is(err, ErrWriteForbidden) {
		t.Fatalf("cross-owner attempt accepted: %v", err)
	}
	var copies int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM editor_conflicts").Scan(&copies); err != nil {
		t.Fatal(err)
	}
	if copies != 0 {
		t.Fatal("cross-owner attempt created a conflict copy")
	}
}

func TestConflictWriteFailureIsNotAcknowledged(t *testing.T) {
	ctx := context.Background()
	db, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "failure.sqlite"))
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
	if _, err := db.ExecContext(ctx, `CREATE TRIGGER block_conflicts BEFORE INSERT
		ON editor_conflicts BEGIN SELECT RAISE(ABORT, 'simulated full disk'); END`); err != nil {
		t.Fatal(err)
	}
	m.MutationID = "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafb"
	_, err = engine.Apply(ctx, m.SubjectID, m)
	var conflict *ConflictError
	if err == nil || (errors.As(err, &conflict) && conflict.Preserved) {
		t.Fatalf("failed conflict insert incorrectly reported durability: %v", err)
	}
}

func TestConflictSchemaUpgradesLegacyFile(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "legacy.sqlite")
	db, err := OpenSQLite(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, "DROP TABLE editor_conflicts"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `CREATE TABLE editor_conflicts (
		site_id TEXT NOT NULL,
		subject_id TEXT NOT NULL,
		mutation_id TEXT NOT NULL,
		document_id TEXT NOT NULL,
		stale_base_revision INTEGER NOT NULL,
		current_revision INTEGER NOT NULL,
		proposed_hash TEXT NOT NULL,
		proposed_source BLOB NOT NULL,
		PRIMARY KEY (site_id, subject_id, mutation_id)
	)`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx,
		`INSERT INTO editor_conflicts (site_id, subject_id, mutation_id, document_id,
		stale_base_revision, current_revision, proposed_hash, proposed_source)
		VALUES (?,?,?,?,?,?,?,?)`,
		"legacy-site", "alice", "legacy-id", "legacy-doc", 0, 1,
		"0000000000000000000000000000000000000000000000000000000000000000",
		[]byte("# Preserve legacy bytes"),
	); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = OpenSQLite(ctx, path)
	if err != nil {
		t.Fatalf("upgrade v1 conflict table: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("close SQLite: %v", err)
		}
	}()
	for _, column := range []string{"payload_digest", "proposed_deleted"} {
		var count int
		if err := db.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM pragma_table_info('editor_conflicts') WHERE name=?", column,
		).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("migration failed to add %s", column)
		}
	}
	var legacyBytes []byte
	if err := db.QueryRowContext(ctx,
		"SELECT proposed_source FROM editor_conflicts WHERE mutation_id='legacy-id'",
	).Scan(&legacyBytes); err != nil {
		t.Fatalf("legacy conflict was lost during migration: %v", err)
	}
	if string(legacyBytes) != "# Preserve legacy bytes" {
		t.Fatalf("legacy conflict contents changed: %q", legacyBytes)
	}
	if err := Migrate(ctx, db); err != nil {
		t.Fatalf("repeat migration: %v", err)
	}
}

func TestReadPrivateConflictIsOwnerScoped(t *testing.T) {
	ctx := context.Background()
	db, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "read-conflicts.sqlite"))
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
	m.Source = []byte("# Keep this offline change")
	_, err = engine.Apply(ctx, m.SubjectID, m)
	var conflict *ConflictError
	if !errors.As(err, &conflict) || !conflict.Preserved {
		t.Fatalf("expected preserved conflict: %v", err)
	}
	record, err := engine.ReadPrivateConflict(ctx, m.SubjectID, m.SiteID, m.MutationID)
	if err != nil {
		t.Fatal(err)
	}
	if record.DocumentID != m.DocumentID || !bytes.Equal(record.Source, m.Source) ||
		record.PayloadDigest != m.Digest() || record.Deleted || record.CurrentRevision != 1 {
		t.Fatalf("incorrect private conflict record: %+v", record)
	}
	if _, err := engine.ReadPrivateConflict(ctx, "bob", m.SiteID, m.MutationID); !errors.Is(err, ErrConflictNotFound) {
		t.Fatalf("other site writer must not access owner's private conflict: %v", err)
	}
	if _, err := engine.ReadPrivateConflict(ctx, m.SubjectID, "other-site", m.MutationID); !errors.Is(err, ErrWriteForbidden) {
		t.Fatalf("cross-site conflict access must be forbidden: %v", err)
	}
	if _, err := engine.ReadPrivateConflict(ctx, m.SubjectID, m.SiteID, "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafc"); !errors.Is(err, ErrConflictNotFound) {
		t.Fatalf("unknown conflict should not reveal metadata: %v", err)
	}
}

func TestOrphanConflictIsDurableAndPrivateAfterRestart(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "orphan.sqlite")
	db, err := OpenSQLite(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	engine := newPrivateTestEngine(t, db)
	m := testMutation()
	m.BaseRevision = 99
	m.Source = []byte("# Offline copy absent on server")
	_, err = engine.Apply(ctx, m.SubjectID, m)
	var conflict *ConflictError
	if !errors.As(err, &conflict) || !conflict.Preserved || conflict.CurrentRevision != 0 {
		t.Fatalf("orphaned draft must be persisted before conflict response: %v", err)
	}
	record, err := engine.ReadPrivateConflict(ctx, m.SubjectID, m.SiteID, m.MutationID)
	if err != nil {
		t.Fatal(err)
	}
	if record.DocumentID != m.DocumentID || !bytes.Equal(record.Source, m.Source) ||
		record.StaleBaseRevision != 99 || record.PayloadDigest != m.Digest() {
		t.Fatalf("orphan conflict lost proposal metadata: %+v", record)
	}
	if _, err := engine.ReadPrivateConflict(ctx, "bob", m.SiteID, m.MutationID); !errors.Is(err, ErrConflictNotFound) {
		t.Fatalf("orphan conflict leaked to other site writer: %v", err)
	}
	var conflicts, documents, changes int
	if err := db.QueryRowContext(ctx,
		`SELECT (SELECT COUNT(*) FROM editor_orphan_conflicts),
		(SELECT COUNT(*) FROM editor_documents),
		(SELECT COUNT(*) FROM editor_changes)`,
	).Scan(&conflicts, &documents, &changes); err != nil {
		t.Fatal(err)
	}
	if conflicts != 1 || documents != 0 || changes != 0 {
		t.Fatalf("orphan wrongly published: conflicts=%d documents=%d changes=%d", conflicts, documents, changes)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = OpenSQLite(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("close SQLite: %v", err)
		}
	}()
	engine = newPrivateTestEngine(t, db)
	_, err = engine.Apply(ctx, m.SubjectID, m)
	if !errors.As(err, &conflict) || !conflict.Preserved || conflict.CurrentRevision != 0 {
		t.Fatalf("orphan retry must replay original conflict: %v", err)
	}
	changed := m
	changed.Source = []byte("# Unauthorized mutation-ID reuse")
	if _, err := engine.Apply(ctx, changed.SubjectID, changed); !errors.Is(err, ErrMutationIDReuse) {
		t.Fatalf("UUID reused with changed intent must fail: %v", err)
	}
	deletion := m
	deletion.MutationID = "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafd"
	deletion.Deleted = true
	deletion.Source = nil
	_, err = engine.Apply(ctx, deletion.SubjectID, deletion)
	if !errors.As(err, &conflict) || !conflict.Preserved {
		t.Fatalf("orphaned delete must persist empty source: %v", err)
	}
	var deleted, sourceNull, sourceLen int
	if err := db.QueryRowContext(ctx,
		`SELECT proposed_deleted, proposed_source IS NULL, length(proposed_source)
		FROM editor_orphan_conflicts WHERE mutation_id=?`, deletion.MutationID,
	).Scan(&deleted, &sourceNull, &sourceLen); err != nil {
		t.Fatal(err)
	}
	if deleted != 1 || sourceNull != 0 || sourceLen != 0 {
		t.Fatalf("orphan delete intent lost: deleted=%d null=%d len=%d", deleted, sourceNull, sourceLen)
	}
}

func TestOrphanConflictFailedInsertNeverAcknowledged(t *testing.T) {
	ctx := context.Background()
	db, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "orphan-failure.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("close SQLite: %v", err)
		}
	}()
	if _, err := db.ExecContext(ctx, `CREATE TRIGGER fail_orphan BEFORE INSERT
		ON editor_orphan_conflicts BEGIN SELECT RAISE(ABORT, 'simulated disk failure'); END`); err != nil {
		t.Fatal(err)
	}
	engine := newPrivateTestEngine(t, db)
	m := testMutation()
	m.BaseRevision = 17
	_, err = engine.Apply(ctx, m.SubjectID, m)
	var conflict *ConflictError
	if err == nil || (errors.As(err, &conflict) && conflict.Preserved) {
		t.Fatalf("failed orphan write incorrectly claimed durability: %v", err)
	}
	var count int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM editor_orphan_conflicts").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("failed orphan mutation left an acknowledged conflict")
	}
}
