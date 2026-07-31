// Package persistence owns SimonDesktop's SQLite metadata store (the
// conversations, documents, and settings tables). Message content itself
// lives in Simon's own memory, not here.
package persistence

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS conversations (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS documents (
    id TEXT PRIMARY KEY,
    original_path TEXT NOT NULL,
    managed_path TEXT NOT NULL,
    display_name TEXT NOT NULL,
    status TEXT NOT NULL,
    error_message TEXT,
    created_at DATETIME NOT NULL,
    indexed_at DATETIME
);

CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
`

// Open opens (creating if necessary) the SQLite database at path and
// applies the schema.
func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
