// Package settings persists arbitrary key/value user preferences in the
// SQLite settings table (see internal/persistence).
package settings

import (
	"context"
	"database/sql"
)

// Repository persists key/value settings in SQLite.
type Repository struct {
	db *sql.DB
}

// NewRepository returns a Repository backed by db. db's schema must already
// include the settings table (see internal/persistence).
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Get returns key's stored value. found is false if key has never been set.
func (r *Repository) Get(ctx context.Context, key string) (value string, found bool, err error) {
	err = r.db.QueryRowContext(ctx, `SELECT value FROM settings WHERE key = ?`, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return value, true, nil
}

// Set stores value under key, overwriting any existing value.
func (r *Repository) Set(ctx context.Context, key, value string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		key, value,
	)
	return err
}
