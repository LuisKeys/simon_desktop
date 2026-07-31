package conversations

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Repository persists Conversation metadata in SQLite.
type Repository struct {
	db *sql.DB
}

// NewRepository returns a Repository backed by db. db's schema must already
// include the conversations table (see internal/persistence).
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a new conversation titled DefaultTitle and returns it.
func (r *Repository) Create(ctx context.Context) (Conversation, error) {
	now := time.Now().UTC()
	c := Conversation{
		ID:        uuid.NewString(),
		Title:     DefaultTitle,
		CreatedAt: now,
		UpdatedAt: now,
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO conversations (id, title, created_at, updated_at) VALUES (?, ?, ?, ?)`,
		c.ID, c.Title, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return Conversation{}, err
	}
	return c, nil
}

// List returns every conversation, most recently updated first.
func (r *Repository) List(ctx context.Context) ([]Conversation, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, created_at, updated_at FROM conversations ORDER BY updated_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []Conversation{}
	for rows.Next() {
		var c Conversation
		if err := rows.Scan(&c.ID, &c.Title, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// RecordMessage bumps conversationID's updated_at to now and, if it still
// has DefaultTitle, derives a title from the first user message.
func (r *Repository) RecordMessage(ctx context.Context, conversationID, userText string) error {
	var title string
	if err := r.db.QueryRowContext(ctx,
		`SELECT title FROM conversations WHERE id = ?`, conversationID,
	).Scan(&title); err != nil {
		return err
	}

	if title == DefaultTitle {
		title = deriveTitle(userText)
	}

	_, err := r.db.ExecContext(ctx,
		`UPDATE conversations SET title = ?, updated_at = ? WHERE id = ?`,
		title, time.Now().UTC(), conversationID,
	)
	return err
}

// deriveTitle builds a deterministic conversation title from a user
// message: newlines removed, truncated to 50 characters, with an ellipsis
// appended when truncated.
func deriveTitle(text string) string {
	flat := strings.Join(strings.Fields(strings.ReplaceAll(text, "\n", " ")), " ")
	runes := []rune(flat)
	if len(runes) <= 50 {
		return flat
	}
	return string(runes[:50]) + "…"
}
