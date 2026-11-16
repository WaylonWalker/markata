package editorsync

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const maxPrivateSyncPage = 100

var ErrPrivateVersionNotFound = errors.New("private editor revision not found")

// PrivateManifestItem is metadata needed to lazily hydrate a single source
// revision; source bytes are intentionally not sent in a vault listing.
type PrivateManifestItem struct {
	DocumentID string
	Revision   int64
	Hash       string
}

type PrivateSnapshotPage struct {
	Items         []PrivateManifestItem
	NextCursor    string
	ChangesCursor string
}

type PrivateChange struct {
	DocumentID string
	Revision   int64
	Deleted    bool
}

type PrivateChangesPage struct {
	Items      []PrivateChange
	NextCursor string
}

func validateSyncPageSize(limit int) error {
	if limit < 1 || limit > maxPrivateSyncPage {
		return errors.New("private sync page size must be between 1 and 100")
	}
	return nil
}

func (r *PrivateReader) siteCheckpoint(ctx context.Context, site string) (int64, error) {
	var checkpoint int64
	err := r.db.QueryRowContext(ctx,
		"SELECT last_seq FROM editor_sites WHERE site_id=?", site,
	).Scan(&checkpoint)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("read private sync checkpoint: %w", err)
	}
	if checkpoint < 0 {
		return 0, ErrExpiredSyncCursor
	}
	return checkpoint, nil
}

// Snapshot uses one immutable high-watermark across all pages; newly saved
// revisions cannot perturb the manifest while the client hydrates a large site.
func (r *PrivateReader) Snapshot(ctx context.Context, subject, site, cursor string, limit int) (PrivateSnapshotPage, error) {
	if err := r.requireRead(ctx, subject, site); err != nil {
		return PrivateSnapshotPage{}, err
	}
	if err := validateSyncPageSize(limit); err != nil {
		return PrivateSnapshotPage{}, err
	}
	var state privateCursor
	if cursor == "" {
		checkpoint, err := r.siteCheckpoint(ctx, site)
		if err != nil {
			return PrivateSnapshotPage{}, err
		}
		state = privateCursor{Kind: "snapshot", Subject: subject, Site: site, Checkpoint: checkpoint}
	} else {
		var err error
		state, err = r.openCursor(cursor, "snapshot", subject, site)
		if err != nil {
			return PrivateSnapshotPage{}, err
		}
	}
	// If a backup was restored to an older database generation, the
	// previously issued manifest cannot be resumed faithfully.
	current, err := r.siteCheckpoint(ctx, site)
	if err != nil {
		return PrivateSnapshotPage{}, err
	}
	if current < state.Checkpoint {
		return PrivateSnapshotPage{}, ErrExpiredSyncCursor
	}
	const query = "WITH ranked AS (" +
		"SELECT document_id,revision,changed_seq,owner_subject,visibility,content_hash,deleted," +
		"ROW_NUMBER() OVER (PARTITION BY document_id ORDER BY changed_seq DESC) AS rank " +
		"FROM editor_document_versions WHERE site_id=? AND changed_seq<=?) " +
		"SELECT document_id,revision,content_hash FROM ranked " +
		"WHERE rank=1 AND owner_subject=? AND visibility='private' AND deleted=0 " +
		"AND document_id>? ORDER BY document_id LIMIT ?"
	rows, err := r.db.QueryContext(ctx, query, site, state.Checkpoint, subject, state.AfterID, limit+1)
	if err != nil {
		return PrivateSnapshotPage{}, fmt.Errorf("list private snapshot: %w", err)
	}
	defer func() { _ = rows.Close() }()
	out := PrivateSnapshotPage{Items: make([]PrivateManifestItem, 0, limit)}
	for rows.Next() {
		var item PrivateManifestItem
		if err := rows.Scan(&item.DocumentID, &item.Revision, &item.Hash); err != nil {
			return PrivateSnapshotPage{}, fmt.Errorf("scan private snapshot: %w", err)
		}
		out.Items = append(out.Items, item)
	}
	if err := rows.Err(); err != nil {
		return PrivateSnapshotPage{}, fmt.Errorf("iterate private snapshot: %w", err)
	}
	if len(out.Items) > limit {
		out.Items = out.Items[:limit]
		state.AfterID = out.Items[len(out.Items)-1].DocumentID
		out.NextCursor, err = r.sealCursor(state)
	} else {
		state.Kind = "changes"
		state.AfterSeq = state.Checkpoint
		state.AfterID = ""
		out.ChangesCursor, err = r.sealCursor(state)
	}
	if err != nil {
		return PrivateSnapshotPage{}, err
	}
	return out, nil
}

// Changes consumes a snapshot's change cursor and only returns owner-scoped
// private draft events. The returned encrypted cursor advances past invisible
// events, without revealing other subjects' document IDs or site activity.
func (r *PrivateReader) Changes(ctx context.Context, subject, site, cursor string, limit int) (PrivateChangesPage, error) {
	if err := r.requireRead(ctx, subject, site); err != nil {
		return PrivateChangesPage{}, err
	}
	if err := validateSyncPageSize(limit); err != nil {
		return PrivateChangesPage{}, err
	}
	state, err := r.openCursor(cursor, "changes", subject, site)
	if err != nil {
		return PrivateChangesPage{}, err
	}
	upper, err := r.siteCheckpoint(ctx, site)
	if err != nil {
		return PrivateChangesPage{}, err
	}
	if upper < state.AfterSeq {
		return PrivateChangesPage{}, ErrExpiredSyncCursor
	}
	const query = "SELECT seq,document_id,revision,deleted FROM editor_changes " +
		"WHERE site_id=? AND owner_subject=? AND visibility='private' AND seq>? AND seq<=? " +
		"ORDER BY seq LIMIT ?"
	rows, err := r.db.QueryContext(ctx, query, site, subject, state.AfterSeq, upper, limit+1)
	if err != nil {
		return PrivateChangesPage{}, fmt.Errorf("list private changes: %w", err)
	}
	defer func() { _ = rows.Close() }()
	out := PrivateChangesPage{Items: make([]PrivateChange, 0, limit)}
	last := state.AfterSeq
	for rows.Next() {
		var item PrivateChange
		var seq int64
		if err := rows.Scan(&seq, &item.DocumentID, &item.Revision, &item.Deleted); err != nil {
			return PrivateChangesPage{}, fmt.Errorf("scan private change: %w", err)
		}
		out.Items = append(out.Items, item)
		if len(out.Items) <= limit {
			last = seq
		}
	}
	if err := rows.Err(); err != nil {
		return PrivateChangesPage{}, fmt.Errorf("iterate private changes: %w", err)
	}
	if len(out.Items) > limit {
		out.Items = out.Items[:limit]
		state.AfterSeq = last
	} else {
		state.AfterSeq = upper
	}
	out.NextCursor, err = r.sealCursor(state)
	if err != nil {
		return PrivateChangesPage{}, err
	}
	return out, nil
}

// SourceAtRevision lazily hydrates a specific immutable private revision.
// Authorization is checked anew on every call, including after reconnect.
func (r *PrivateReader) SourceAtRevision(ctx context.Context, subject, site, document string, revision int64) ([]byte, error) {
	if err := r.requireRead(ctx, subject, site); err != nil {
		return nil, err
	}
	if !validSyncID(document, 256) || revision < 1 {
		return nil, ErrPrivateVersionNotFound
	}
	var source []byte
	err := r.db.QueryRowContext(ctx,
		"SELECT source FROM editor_document_versions WHERE site_id=? AND document_id=? "+
			"AND revision=? AND owner_subject=? AND visibility='private' AND deleted=0",
		site, document, revision, subject,
	).Scan(&source)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPrivateVersionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("hydrate private revision: %w", err)
	}
	return source, nil
}
