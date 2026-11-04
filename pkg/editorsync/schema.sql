-- Hosted editor draft revision schema, v2 (SQLite).
-- This is deliberately not a public API. All authorization is enforced in Go.
-- Per-site counters sequence edits; browser-facing cursors must be opaque and scoped.
CREATE TABLE IF NOT EXISTS editor_sites (
  site_id TEXT PRIMARY KEY NOT NULL CHECK (length(site_id) BETWEEN 1 AND 64),
  last_seq INTEGER NOT NULL DEFAULT 0 CHECK (last_seq >= 0)
);

CREATE TABLE IF NOT EXISTS editor_documents (
  site_id TEXT NOT NULL,
  document_id TEXT NOT NULL,
  owner_subject TEXT NOT NULL CHECK (length(owner_subject) BETWEEN 1 AND 256),
  visibility TEXT NOT NULL CHECK (visibility IN ('private', 'shared_draft', 'published_source')),
  revision INTEGER NOT NULL CHECK (revision >= 1),
  changed_seq INTEGER NOT NULL CHECK (changed_seq >= 1),
  content_hash TEXT NOT NULL CHECK (length(content_hash) = 64),
  source BLOB NOT NULL,
  deleted INTEGER NOT NULL DEFAULT 0 CHECK (deleted IN (0, 1)),
  PRIMARY KEY (site_id, document_id),
  FOREIGN KEY (site_id) REFERENCES editor_sites(site_id)
);

-- Immutable revisions are required for consistent paginated snapshots while new
-- mutations arrive. An index page must identify the exact revision to hydrate.
CREATE TABLE IF NOT EXISTS editor_document_versions (
  site_id TEXT NOT NULL,
  document_id TEXT NOT NULL,
  revision INTEGER NOT NULL CHECK (revision >= 1),
  changed_seq INTEGER NOT NULL CHECK (changed_seq >= 1),
  owner_subject TEXT NOT NULL,
  visibility TEXT NOT NULL CHECK (visibility IN ('private', 'shared_draft', 'published_source')),
  content_hash TEXT NOT NULL CHECK (length(content_hash) = 64),
  source BLOB NOT NULL,
  deleted INTEGER NOT NULL DEFAULT 0 CHECK (deleted IN (0, 1)),
  PRIMARY KEY (site_id, document_id, revision),
  UNIQUE (site_id, changed_seq),
  FOREIGN KEY (site_id, document_id) REFERENCES editor_documents(site_id, document_id)
);

-- Persistent receipts survive lost HTTP acknowledgements and process restarts.
-- Idempotency identity is (site, subject, mutation UUID); a reused UUID with
-- a different payload_digest must fail, never silently replay another write.
CREATE TABLE IF NOT EXISTS editor_mutation_receipts (
  site_id TEXT NOT NULL,
  subject_id TEXT NOT NULL,
  mutation_id TEXT NOT NULL,
  payload_digest TEXT NOT NULL CHECK (length(payload_digest) = 64),
  document_id TEXT NOT NULL,
  base_revision INTEGER NOT NULL CHECK (base_revision >= 0),
  accepted_revision INTEGER NOT NULL CHECK (accepted_revision >= 1),
  changed_seq INTEGER NOT NULL CHECK (changed_seq >= 1),
  PRIMARY KEY (site_id, subject_id, mutation_id),
  UNIQUE (site_id, changed_seq),
  FOREIGN KEY (site_id, document_id, accepted_revision)
    REFERENCES editor_document_versions(site_id, document_id, revision)
);

-- The internal sequence is NOT itself an authorization token or a browser cursor.
-- Change feeds must filter owner/visibility and hide unauthorized metadata.
CREATE TABLE IF NOT EXISTS editor_changes (
  site_id TEXT NOT NULL,
  seq INTEGER NOT NULL CHECK (seq >= 1),
  document_id TEXT NOT NULL,
  revision INTEGER NOT NULL CHECK (revision >= 1),
  owner_subject TEXT NOT NULL,
  visibility TEXT NOT NULL CHECK (visibility IN ('private', 'shared_draft', 'published_source')),
  deleted INTEGER NOT NULL CHECK (deleted IN (0, 1)),
  PRIMARY KEY (site_id, seq),
  FOREIGN KEY (site_id, document_id, revision)
    REFERENCES editor_document_versions(site_id, document_id, revision)
);

-- Conflicting edits are retained for review rather than overwritten.
-- Old v1 databases are extended with payload_digest and proposed_deleted by Migrate;
-- empty legacy digests must not be mistaken for verified mutation receipts.
CREATE TABLE IF NOT EXISTS editor_conflicts (
  site_id TEXT NOT NULL,
  subject_id TEXT NOT NULL,
  mutation_id TEXT NOT NULL,
  document_id TEXT NOT NULL,
  stale_base_revision INTEGER NOT NULL CHECK (stale_base_revision >= 0),
  current_revision INTEGER NOT NULL CHECK (current_revision >= 1),
  proposed_hash TEXT NOT NULL CHECK (length(proposed_hash) = 64),
  proposed_source BLOB NOT NULL,
  payload_digest TEXT NOT NULL DEFAULT '' CHECK (length(payload_digest) IN (0, 64)),
  proposed_deleted INTEGER NOT NULL DEFAULT 0 CHECK (proposed_deleted IN (0, 1)),
  PRIMARY KEY (site_id, subject_id, mutation_id),
  FOREIGN KEY (site_id, document_id) REFERENCES editor_documents(site_id, document_id)
);

CREATE INDEX IF NOT EXISTS idx_editor_versions_snapshot
  ON editor_document_versions(site_id, changed_seq, document_id);
CREATE INDEX IF NOT EXISTS idx_editor_changes_visibility
  ON editor_changes(site_id, seq, visibility, owner_subject);
CREATE INDEX IF NOT EXISTS idx_editor_conflicts_owner
  ON editor_conflicts(site_id, subject_id, document_id);
