package editorsync

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSQLiteDurableOpenRestartAndDuplicateReceipt(t *testing.T) {
	file := filepath.Join(t.TempDir(), "sync.sqlite")
	ctx := context.Background()
	db, err := OpenSQLite(ctx, file)
	if err != nil {
		t.Fatal(err)
	}
	if err := Migrate(ctx, db); err != nil {
		t.Fatalf("idempotent migration: %v", err)
	}
	content := []byte("# First post\n")
	hash := sha256.Sum256(content)
	digest := hex.EncodeToString(hash[:])
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	steps := []struct {
		query string
		args  []any
	}{
		{"INSERT INTO editor_sites(site_id,last_seq) VALUES(?,?)", []any{"site-a", 1}},
		{`INSERT INTO editor_documents (site_id,document_id,owner_subject,visibility,
			revision,changed_seq,content_hash,source,deleted)
			VALUES (?,?,?,?,?,?,?,?,?)`, []any{"site-a", "document-1", "alice", "private", 1, 1, digest, content, 0}},
		{`INSERT INTO editor_document_versions (site_id,document_id,revision,changed_seq,
			owner_subject,visibility,content_hash,source,deleted)
			VALUES (?,?,?,?,?,?,?,?,?)`, []any{"site-a", "document-1", 1, 1, "alice", "private", digest, content, 0}},
		{`INSERT INTO editor_mutation_receipts (site_id,subject_id,mutation_id,payload_digest,
			document_id,base_revision,accepted_revision,changed_seq)
			VALUES (?,?,?,?,?,?,?,?)`, []any{"site-a", "alice", "mut-1", digest, "document-1", 0, 1, 1}},
		{`INSERT INTO editor_changes (site_id,seq,document_id,revision,owner_subject,visibility,deleted)
			VALUES (?,?,?,?,?,?,?)`, []any{"site-a", 1, "document-1", 1, "alice", "private", 0}},
	}
	for _, step := range steps {
		if _, err := tx.ExecContext(ctx, step.query, step.args...); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				t.Errorf("transaction rollback: %v", rollbackErr)
			}
			t.Fatalf("transaction insert: %v", err)
		}
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = OpenSQLite(ctx, file)
	if err != nil {
		t.Fatalf("reopen SQLite: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("close reopened SQLite: %v", err)
		}
	}()
	var stored []byte
	var revision, seq int64
	err = db.QueryRowContext(ctx,
		`SELECT d.source,d.revision,r.changed_seq FROM editor_documents d
		JOIN editor_mutation_receipts r ON d.site_id=r.site_id AND d.document_id=r.document_id
		WHERE r.site_id='site-a' AND r.subject_id='alice' AND r.mutation_id='mut-1'`,
	).Scan(&stored, &revision, &seq)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(stored, content) || revision != 1 || seq != 1 {
		t.Fatalf("unexpected durable revision=%d seq=%d source=%q", revision, seq, stored)
	}
	_, err = db.ExecContext(ctx,
		`INSERT INTO editor_mutation_receipts (site_id,subject_id,mutation_id,payload_digest,
		document_id,base_revision,accepted_revision,changed_seq)
		VALUES ('site-a','alice','mut-1',?,'document-1',0,1,1)`, digest)
	if err == nil {
		t.Fatal("reusing idempotency receipt key must fail in SQL")
	}
}

func TestSQLiteRejectsMissingParent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "sync.sqlite")
	if _, err := OpenSQLite(context.Background(), path); err == nil {
		t.Fatal("created an unexpected parent directory")
	}
}

func TestSQLiteNewDatabaseFileIsPrivate(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX permissions do not map directly to Windows ACLs")
	}
	file := filepath.Join(t.TempDir(), "private.sqlite")
	db, err := OpenSQLite(context.Background(), file)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0o077 != 0 {
		t.Fatalf("database permissions should be owner-only, got %04o", info.Mode().Perm())
	}
}
