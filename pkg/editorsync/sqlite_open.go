package editorsync

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	// Register the pure-Go SQLite database/sql driver for this opt-in package.
	_ "modernc.org/sqlite"
)

// OpenSQLite opens a trusted, on-disk editor sync database, sets connection
// pragmas, and applies idempotent migrations. Path must come from server
// configuration, never a browser-controlled URL or document path.
//
// POC uses a single connection to serialize writes and guarantee that
// connection-local foreign key and busy timeout settings are consistent.
// Scaling to multiple writer processes requires a separate WAL/retry review.
func OpenSQLite(ctx context.Context, path string) (*sql.DB, error) {
	if path == "" || path == ":memory:" || strings.ContainsRune(path, '\x00') {
		return nil, errors.New("a trusted SQLite file path is required")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve editor sync SQLite path: %w", err)
	}
	info, err := os.Stat(filepath.Dir(abs))
	if err != nil || !info.IsDir() {
		return nil, errors.New("editor sync SQLite parent directory must already exist")
	}
	// New database files must not inherit a permissive default creation mode:
	// they can contain unpublished private drafts.
	file, err := os.OpenFile(abs, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return nil, fmt.Errorf("create private editor sync SQLite file: %w", err)
	}
	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("close private editor sync SQLite file: %w", err)
	}
	params := url.Values{}
	params.Add("_pragma", "foreign_keys(1)")
	params.Add("_pragma", "busy_timeout(5000)")
	params.Add("_pragma", "journal_mode(WAL)")
	uriPath := filepath.ToSlash(abs)
	if runtime.GOOS == "windows" {
		uriPath = "/" + uriPath
	}
	uri := (&url.URL{Scheme: "file", Path: uriPath, RawQuery: params.Encode()}).String()
	db, err := sql.Open("sqlite", uri)
	if err != nil {
		return nil, fmt.Errorf("open SQLite: %w", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	if err := db.PingContext(ctx); err != nil {
		return nil, closeSQLiteOnError(db, fmt.Errorf("connect SQLite: %w", err))
	}
	var foreignKeys int
	if err := db.QueryRowContext(ctx, "PRAGMA foreign_keys").Scan(&foreignKeys); err != nil || foreignKeys != 1 {
		return nil, closeSQLiteOnError(db, errors.New("SQLite foreign key enforcement is not active"))
	}
	var journalMode string
	if err := db.QueryRowContext(ctx, "PRAGMA journal_mode").Scan(&journalMode); err != nil || !strings.EqualFold(journalMode, "wal") {
		return nil, closeSQLiteOnError(db, errors.New("SQLite WAL mode is not active"))
	}
	if err := Migrate(ctx, db); err != nil {
		return nil, closeSQLiteOnError(db, err)
	}
	return db, nil
}

// closeSQLiteOnError makes resource cleanup errors observable without losing the
// original database error. The caller must not use the closed handle.
func closeSQLiteOnError(db *sql.DB, cause error) error {
	if err := db.Close(); err != nil {
		return errors.Join(cause, fmt.Errorf("close editor sync SQLite: %w", err))
	}
	return cause
}
