package editorsync

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrConflictNotFound = errors.New("private sync conflict not found")

// PrivateConflict is a stored, not-yet-applied mutation. Reading a conflict
// never grants permission to mutate or publish its underlying document.
type PrivateConflict struct {
	MutationID        string
	DocumentID        string
	StaleBaseRevision int64
	CurrentRevision   int64
	PayloadDigest     string
	Source            []byte
	Deleted           bool
}

// ReadPrivateConflict only returns a conflict originally recorded for the
// authenticated subject. Site write authority alone never permits reading a
// different subject's private conflict. No HTTP route is mounted here.
func (e *Engine) ReadPrivateConflict(ctx context.Context, principal, siteID, mutationID string) (PrivateConflict, error) {
	if e == nil || e.db == nil || e.authorize == nil || !validSyncID(principal, 256) ||
		!siteIDRE.MatchString(siteID) || !mutationIDRE.MatchString(mutationID) {
		return PrivateConflict{}, ErrWriteForbidden
	}
	if err := e.authorize(ctx, principal, siteID); err != nil {
		return PrivateConflict{}, ErrWriteForbidden
	}
	var record PrivateConflict
	var deleted int
	record.MutationID = mutationID
	err := e.db.QueryRowContext(ctx,
		`SELECT document_id, stale_base_revision, current_revision, payload_digest,
		proposed_source, proposed_deleted FROM editor_conflicts
		WHERE site_id=? AND subject_id=? AND mutation_id=?`,
		siteID, principal, mutationID,
	).Scan(&record.DocumentID, &record.StaleBaseRevision, &record.CurrentRevision,
		&record.PayloadDigest, &record.Source, &deleted)
	if errors.Is(err, sql.ErrNoRows) {
		return PrivateConflict{}, ErrConflictNotFound
	}
	if err != nil {
		return PrivateConflict{}, fmt.Errorf("read private editor sync conflict: %w", err)
	}
	record.Deleted = deleted == 1
	if len(record.PayloadDigest) != 64 {
		return PrivateConflict{}, ErrConflictNotFound
	}
	return record, nil
}
