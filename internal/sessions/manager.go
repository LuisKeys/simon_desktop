// Package sessions manages one simon.Session per conversation, created on
// demand and reused for the life of the process.
package sessions

import (
	"sync"

	"github.com/LuisKeys/simon"
)

// SystemPrompt is the system prompt every conversation's Session is created
// with.
const SystemPrompt = `You are Simon, a personal document assistant.

Answer using the available conversation and document context.

When the provided documents do not contain enough information, clearly say so.
Do not invent facts or claim that a document says something unsupported.
Keep responses direct and practical.`

// Manager creates and reuses one simon.Session per conversation ID.
// Sessions themselves are closed by simon.Runtime.Close, which closes every
// Session it created.
type Manager struct {
	rt *simon.Runtime

	mu       sync.Mutex
	sessions map[string]*simon.Session
}

// NewManager returns a Manager bound to rt.
func NewManager(rt *simon.Runtime) *Manager {
	return &Manager{rt: rt, sessions: make(map[string]*simon.Session)}
}

// GetOrCreate returns the existing Session for conversationID, creating one
// (with SystemPrompt) the first time it's requested.
func (m *Manager) GetOrCreate(conversationID string) (*simon.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if sess, ok := m.sessions[conversationID]; ok {
		return sess, nil
	}

	sess, err := m.rt.NewSession(conversationID, simon.WithSystemPrompt(SystemPrompt))
	if err != nil {
		return nil, err
	}
	m.sessions[conversationID] = sess
	return sess, nil
}

// Cancel cancels the active run on conversationID's session, if that
// session has been created and has a run in flight. A no-op otherwise.
func (m *Manager) Cancel(conversationID string) {
	m.mu.Lock()
	sess, ok := m.sessions[conversationID]
	m.mu.Unlock()
	if ok {
		sess.Cancel()
	}
}
