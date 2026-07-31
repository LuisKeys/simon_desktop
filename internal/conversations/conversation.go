// Package conversations manages conversation metadata (title, timestamps).
// Message content is not stored here — it lives in Simon's own per-session
// memory.
package conversations

import "time"

// Conversation is a single conversation's metadata.
type Conversation struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// DefaultTitle is the title assigned to a conversation before its first
// message has been sent.
const DefaultTitle = "New conversation"
