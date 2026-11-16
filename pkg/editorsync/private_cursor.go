package editorsync

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

var ErrReadForbidden = errors.New("editor sync read forbidden")
var ErrInvalidSyncCursor = errors.New("invalid editor sync cursor")
var ErrExpiredSyncCursor = errors.New("editor sync cursor expired")

// SiteReadAccess verifies a stable authenticated subject has read permission
// for the site. This is independent from the private-document owner check.
type SiteReadAccess func(context.Context, string, string) error

// PrivateReader exposes a private-draft-only read protocol. It does not mount
// HTTP endpoints or grant access to shared or published documents.
type PrivateReader struct {
	db        *sql.DB
	authorize SiteReadAccess
	aead      cipher.AEAD
	ttl       time.Duration
	now       func() time.Time
}

type privateCursor struct {
	Version    int    `json:"v"`
	Kind       string `json:"kind"`
	Site       string `json:"site"`
	Subject    string `json:"subject"`
	Checkpoint int64  `json:"checkpoint"`
	AfterID    string `json:"after_id,omitempty"`
	AfterSeq   int64  `json:"after_seq,omitempty"`
	Expires    int64  `json:"expires"`
}

// NewPrivateReader accepts a server-held 32-byte key. The same key must be
// retained across server restarts while outstanding cursors remain valid.
// Cursors use authenticated encryption, not readable or merely signed JSON.
func NewPrivateReader(db *sql.DB, authorize SiteReadAccess, key []byte, ttl time.Duration) (*PrivateReader, error) {
	if db == nil || authorize == nil || len(key) != 32 || ttl <= 0 || ttl > 24*time.Hour {
		return nil, errors.New("private sync reader requires SQLite, authorization, 32-byte key and bounded TTL")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("init editor cursor cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("init editor cursor AEAD: %w", err)
	}
	return &PrivateReader{db: db, authorize: authorize, aead: aead, ttl: ttl, now: time.Now}, nil
}

func (r *PrivateReader) requireRead(ctx context.Context, subject, site string) error {
	if r == nil || r.db == nil || r.authorize == nil || r.aead == nil ||
		!validSyncID(subject, 256) || !siteIDRE.MatchString(site) {
		return ErrReadForbidden
	}
	if err := r.authorize(ctx, subject, site); err != nil {
		return ErrReadForbidden
	}
	return nil
}

func (r *PrivateReader) sealCursor(cursor privateCursor) (string, error) {
	cursor.Version = 1
	cursor.Expires = r.now().Add(r.ttl).Unix()
	plain, err := json.Marshal(cursor)
	if err != nil {
		return "", fmt.Errorf("encode private sync cursor: %w", err)
	}
	nonce := make([]byte, r.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate private sync cursor nonce: %w", err)
	}
	ciphertext := r.aead.Seal(nonce, nonce, plain, []byte("markata-private-sync-v1"))
	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

func (r *PrivateReader) openCursor(raw, purpose, subject, site string) (privateCursor, error) {
	if raw == "" || len(raw) > 4096 {
		return privateCursor{}, ErrInvalidSyncCursor
	}
	payload, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil || len(payload) < r.aead.NonceSize()+r.aead.Overhead() {
		return privateCursor{}, ErrInvalidSyncCursor
	}
	nonce := payload[:r.aead.NonceSize()]
	encrypted := payload[r.aead.NonceSize():]
	plain, err := r.aead.Open(nil, nonce, encrypted, []byte("markata-private-sync-v1"))
	if err != nil {
		return privateCursor{}, ErrInvalidSyncCursor
	}
	var cursor privateCursor
	if err := json.Unmarshal(plain, &cursor); err != nil || cursor.Version != 1 ||
		cursor.Kind != purpose || cursor.Subject != subject || cursor.Site != site ||
		cursor.Checkpoint < 0 || cursor.AfterSeq < 0 {
		return privateCursor{}, ErrInvalidSyncCursor
	}
	if cursor.Expires <= r.now().Unix() {
		return privateCursor{}, ErrExpiredSyncCursor
	}
	return cursor, nil
}
