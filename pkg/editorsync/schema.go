// Package editorsync defines the durable editor mutation contract and SQLite
// schema. Transport, browser persistence and a concrete SQLite driver are future
// slices; this package cannot claim to sync drafts on its own.
package editorsync

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"embed"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"regexp"
	"strings"
)

//go:embed schema.sql
var migrations embed.FS

var (
	siteIDRE     = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,62}[a-z0-9])?$`)
	mutationIDRE = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
)

// Visibility determines server-controlled draft availability, never browser
// authorization. Document/media ACLs must still be checked by a caller.
type Visibility string

const (
	Private         Visibility = "private"
	SharedDraft     Visibility = "shared_draft"
	PublishedSource Visibility = "published_source"
)

// Mutation is a proposed change. Only a transactionally committed receipt is
// acknowledgement; computing a digest must never advance the client outbox.
type Mutation struct {
	SiteID       string
	DocumentID   string
	SubjectID    string
	MutationID   string
	BaseRevision uint64
	Visibility   Visibility
	Source       []byte
	Deleted      bool
}

// Validate checks structural constraints only; the caller still needs an
// authenticated subject, write permission and document-level authorization.
func (m Mutation) Validate() error {
	if !siteIDRE.MatchString(m.SiteID) || len(m.SiteID) > 64 {
		return errors.New("invalid site ID")
	}
	if !validSyncID(m.DocumentID, 256) {
		return errors.New("invalid document ID")
	}
	if !validSyncID(m.SubjectID, 256) {
		return errors.New("invalid subject ID")
	}
	if !mutationIDRE.MatchString(m.MutationID) {
		return errors.New("invalid mutation UUID")
	}
	switch m.Visibility {
	case Private, SharedDraft, PublishedSource:
	default:
		return errors.New("invalid document visibility")
	}
	if len(m.Source) > 2<<20 {
		return errors.New("source exceeds 2 MiB POC limit")
	}
	return nil
}

func validSyncID(value string, limit int) bool {
	return value != "" && len(value) <= limit &&
		strings.TrimSpace(value) == value && !strings.ContainsAny(value, "\r\n\x00")
}

// Digest is a deterministic, length-framed SHA-256 over the entire intent,
// including the base revision, visibility and deletion state. This prevents a
// retried mutation UUID from being reused with a different payload.
func (m Mutation) Digest() string {
	h := sha256.New()
	for _, part := range []string{m.SiteID, m.SubjectID, m.DocumentID, m.MutationID, string(m.Visibility)} {
		writeDigestField(h, []byte(part))
	}
	var revision [8]byte
	binary.BigEndian.PutUint64(revision[:], m.BaseRevision)
	writeDigestField(h, revision[:])
	if m.Deleted {
		writeDigestField(h, []byte{1})
	} else {
		writeDigestField(h, []byte{0})
	}
	writeDigestField(h, m.Source)
	return hex.EncodeToString(h.Sum(nil))
}

func writeDigestField(h hash.Hash, data []byte) {
	var length [8]byte
	binary.BigEndian.PutUint64(length[:], uint64(len(data)))
	_, _ = h.Write(length[:]) // SHA-256 hash.Hash.Write always returns len(data), nil.
	_, _ = h.Write(data)
}

// Migrate applies the schema in a transaction to a *caller-provided SQLite*
// database/sql connection. It intentionally does not register a driver, open a
// production database, or configure site permissions. Openers must enable
// SQLite foreign_keys on every connection, WAL, bounded busy_timeout and
// backups. A driver-backed integration test is required before enabling sync.
func Migrate(ctx context.Context, db *sql.DB) error {
	if db == nil {
		return errors.New("SQLite database is required")
	}
	source, err := migrations.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("read editor sync schema: %w", err)
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin editor sync migration: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var statement strings.Builder
	for _, line := range strings.Split(string(source), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "--") || trimmed == "" {
			continue
		}
		statement.WriteString(line)
		statement.WriteByte('\n')
		if strings.HasSuffix(trimmed, ";") {
			if _, err := tx.ExecContext(ctx, statement.String()); err != nil {
				return fmt.Errorf("migrate editor sync schema: %w", err)
			}
			statement.Reset()
		}
	}
	if strings.TrimSpace(statement.String()) != "" {
		return errors.New("unterminated editor sync migration statement")
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit editor sync migration: %w", err)
	}
	return nil
}
