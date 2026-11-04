package editorsync

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
)

// MVP supports private drafts only. Shared/published content and conflict-copy
// replication require separate policy and permission-filtered feed work.
var ErrWriteForbidden = errors.New("editor sync write forbidden")
var ErrRevisionConflict = errors.New("editor sync revision conflict")
var ErrMutationIDReuse = errors.New("editor sync mutation ID reused with different payload")
var ErrUnsupportedVisibility = errors.New("editor sync visibility not yet supported")

// SiteWriteAccess checks the authenticated principal's authority to write a
// site. Document-level private ownership is enforced separately by Apply.
type SiteWriteAccess func(context.Context, string, string) error

// Engine requires both a SQLite connection and explicit authorization.
type Engine struct {
	db        *sql.DB
	authorize SiteWriteAccess
}

func NewEngine(db *sql.DB, authorize SiteWriteAccess) (*Engine, error) {
	if db == nil || authorize == nil {
		return nil, errors.New("SQLite connection and site authorizer are required")
	}
	return &Engine{db: db, authorize: authorize}, nil
}

// Receipt proves a mutation was durably accepted by SQLite. It does NOT imply
// a Git commit, site build, or public publication.
type Receipt struct {
	MutationID    string
	Revision      int64
	Sequence      int64
	PayloadDigest string
}

// ConflictError exposes the observed server revision without discarding the
// device's unsent edit. The browser must retain its outbox entry on this error.
type ConflictError struct {
	CurrentRevision int64
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("%v (current revision %d)", ErrRevisionConflict, e.CurrentRevision)
}
func (e *ConflictError) Unwrap() error { return ErrRevisionConflict }

// Apply validates the caller, transactionally writes exactly one private
// document revision and a receipt, and only then acknowledges it. A retry with
// the same mutation ID/payload returns the original durable receipt. Editing
// shared/public source, persisted conflict copies and transport are later work.
func (e *Engine) Apply(ctx context.Context, principal string, m Mutation) (receipt Receipt, resultErr error) {
	if e == nil || e.db == nil || e.authorize == nil || principal == "" || m.SubjectID != principal {
		return Receipt{}, ErrWriteForbidden
	}
	if err := m.Validate(); err != nil {
		return Receipt{}, err
	}
	if m.Visibility != Private {
		return Receipt{}, ErrUnsupportedVisibility
	}
	if m.BaseRevision > math.MaxInt64 {
		return Receipt{}, errors.New("invalid base revision")
	}
	if err := e.authorize(ctx, principal, m.SiteID); err != nil {
		return Receipt{}, ErrWriteForbidden
	}
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return Receipt{}, fmt.Errorf("start sync mutation: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			resultErr = errors.Join(resultErr, fmt.Errorf("rollback sync mutation: %w", rollbackErr))
		}
		if resultErr != nil {
			receipt = Receipt{}
		}
	}()
	digest := m.Digest()
	if oldReceipt, exists, err := findSyncReceipt(ctx, tx, m, digest); err != nil {
		return Receipt{}, err
	} else if exists {
		return oldReceipt, nil
	}
	old, exists, err := currentSyncDocument(ctx, tx, m.SiteID, m.DocumentID)
	if err != nil {
		return Receipt{}, err
	}
	if exists && (old.Owner != principal || old.Visibility != Private) {
		return Receipt{}, ErrWriteForbidden
	}
	if uint64(old.Revision) != m.BaseRevision {
		return Receipt{}, &ConflictError{CurrentRevision: old.Revision}
	}
	sequence, err := nextSyncSequence(ctx, tx, m.SiteID)
	if err != nil {
		return Receipt{}, err
	}
	sum := sha256.Sum256(m.Source)
	receipt = Receipt{MutationID: m.MutationID, Revision: old.Revision + 1, Sequence: sequence, PayloadDigest: digest}
	if err := writeCurrentSyncDocument(ctx, tx, m, old, exists, receipt, hex.EncodeToString(sum[:])); err != nil {
		return Receipt{}, err
	}
	if err := appendSyncAudit(ctx, tx, m, receipt, hex.EncodeToString(sum[:])); err != nil {
		return Receipt{}, err
	}
	if err := tx.Commit(); err != nil {
		return Receipt{}, fmt.Errorf("commit sync mutation: %w", err)
	}
	return receipt, nil
}

type currentDocument struct {
	Revision   int64
	Owner      string
	Visibility Visibility
}

func findSyncReceipt(ctx context.Context, tx *sql.Tx, m Mutation, digest string) (Receipt, bool, error) {
	var oldDigest string
	receipt := Receipt{MutationID: m.MutationID, PayloadDigest: digest}
	err := tx.QueryRowContext(ctx,
		`SELECT payload_digest, accepted_revision, changed_seq
		FROM editor_mutation_receipts
		WHERE site_id=? AND subject_id=? AND mutation_id=?`,
		m.SiteID, m.SubjectID, m.MutationID,
	).Scan(&oldDigest, &receipt.Revision, &receipt.Sequence)
	if errors.Is(err, sql.ErrNoRows) {
		return Receipt{}, false, nil
	}
	if err != nil {
		return Receipt{}, false, fmt.Errorf("query sync receipt: %w", err)
	}
	if oldDigest != digest {
		return Receipt{}, false, ErrMutationIDReuse
	}
	return receipt, true, nil
}

func currentSyncDocument(ctx context.Context, tx *sql.Tx, site, document string) (currentDocument, bool, error) {
	var old currentDocument
	err := tx.QueryRowContext(ctx,
		`SELECT revision, owner_subject, visibility FROM editor_documents
		WHERE site_id=? AND document_id=?`, site, document,
	).Scan(&old.Revision, &old.Owner, &old.Visibility)
	if errors.Is(err, sql.ErrNoRows) {
		return currentDocument{}, false, nil
	}
	if err != nil {
		return currentDocument{}, false, fmt.Errorf("read current sync document: %w", err)
	}
	return old, true, nil
}

func nextSyncSequence(ctx context.Context, tx *sql.Tx, site string) (int64, error) {
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO editor_sites(site_id, last_seq) VALUES(?,0)
		ON CONFLICT(site_id) DO NOTHING`, site,
	); err != nil {
		return 0, fmt.Errorf("ensure sync site: %w", err)
	}
	if _, err := tx.ExecContext(ctx,
		"UPDATE editor_sites SET last_seq=last_seq+1 WHERE site_id=?", site,
	); err != nil {
		return 0, fmt.Errorf("advance sync sequence: %w", err)
	}
	var seq int64
	if err := tx.QueryRowContext(ctx,
		"SELECT last_seq FROM editor_sites WHERE site_id=?", site,
	).Scan(&seq); err != nil {
		return 0, fmt.Errorf("read sync sequence: %w", err)
	}
	return seq, nil
}

func writeCurrentSyncDocument(ctx context.Context, tx *sql.Tx, m Mutation, old currentDocument, exists bool, receipt Receipt, hash string) error {
	deleted := 0
	if m.Deleted {
		deleted = 1
	}
	if !exists {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO editor_documents
			(site_id,document_id,owner_subject,visibility,revision,changed_seq,content_hash,source,deleted)
			VALUES (?,?,?,?,?,?,?,?,?)`,
			m.SiteID, m.DocumentID, m.SubjectID, string(Private), receipt.Revision,
			receipt.Sequence, hash, m.Source, deleted,
		); err != nil {
			return fmt.Errorf("create sync document: %w", err)
		}
		return nil
	}
	result, err := tx.ExecContext(ctx,
		`UPDATE editor_documents SET revision=?, changed_seq=?, content_hash=?, source=?, deleted=?
		WHERE site_id=? AND document_id=? AND revision=? AND owner_subject=? AND visibility='private'`,
		receipt.Revision, receipt.Sequence, hash, m.Source, deleted,
		m.SiteID, m.DocumentID, old.Revision, m.SubjectID,
	)
	if err != nil {
		return fmt.Errorf("update sync document: %w", err)
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("verify sync compare-and-swap: %w", err)
	}
	if changed != 1 {
		return &ConflictError{CurrentRevision: old.Revision}
	}
	return nil
}

func appendSyncAudit(ctx context.Context, tx *sql.Tx, m Mutation, receipt Receipt, hash string) error {
	deleted := 0
	if m.Deleted {
		deleted = 1
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO editor_document_versions
		(site_id,document_id,revision,changed_seq,owner_subject,visibility,content_hash,source,deleted)
		VALUES (?,?,?,?,?,?,?,?,?)`,
		m.SiteID, m.DocumentID, receipt.Revision, receipt.Sequence, m.SubjectID, string(Private),
		hash, m.Source, deleted,
	); err != nil {
		return fmt.Errorf("append immutable revision: %w", err)
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO editor_changes
		(site_id,seq,document_id,revision,owner_subject,visibility,deleted)
		VALUES (?,?,?,?,?,?,?)`,
		m.SiteID, receipt.Sequence, m.DocumentID, receipt.Revision, m.SubjectID, string(Private), deleted,
	); err != nil {
		return fmt.Errorf("append sync change: %w", err)
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO editor_mutation_receipts
		(site_id,subject_id,mutation_id,payload_digest,document_id,
		base_revision,accepted_revision,changed_seq) VALUES (?,?,?,?,?,?,?,?)`,
		m.SiteID, m.SubjectID, m.MutationID, receipt.PayloadDigest, m.DocumentID, m.BaseRevision,
		receipt.Revision, receipt.Sequence,
	); err != nil {
		return fmt.Errorf("persist sync receipt: %w", err)
	}
	return nil
}
