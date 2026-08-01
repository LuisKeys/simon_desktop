package app

import (
	"database/sql"

	"simonpartner/internal/conversations"
	"simonpartner/internal/settings"
)

// Repositories groups every SQLite-backed repository AppService uses.
type Repositories struct {
	Conversations *conversations.Repository
	Settings      *settings.Repository
}

func newRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Conversations: conversations.NewRepository(db),
		Settings:      settings.NewRepository(db),
	}
}
