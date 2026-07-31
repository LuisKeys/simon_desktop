package app

import (
	"database/sql"

	"simonpartner/internal/conversations"
)

// Repositories groups every SQLite-backed repository AppService uses.
type Repositories struct {
	Conversations *conversations.Repository
}

func newRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Conversations: conversations.NewRepository(db),
	}
}
