package app

import (
	"database/sql"

	"simondesktop/internal/conversations"
	"simondesktop/internal/settings"
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
