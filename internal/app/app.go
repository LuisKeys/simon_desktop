// Package app wires SimonDesktop's backend services together and exposes
// the facade Wails binds to the frontend.
package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/LuisKeys/simon"
	"github.com/LuisKeys/simon/memory"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"simonpartner/internal/conversations"
	"simonpartner/internal/persistence"
	"simonpartner/internal/platform"
	"simonpartner/internal/sessions"
)

// BaseWindowWidth and BaseWindowHeight are the window's dimensions at 100%
// UI scale; main.go uses these as the initial window size, and
// SetWindowScale multiplies them by the user's chosen scale factor.
const (
	BaseWindowWidth  = 1200
	BaseWindowHeight = 800
)

// AppService is the facade Wails binds to the frontend. It owns the single
// simon.Runtime for the process's lifetime, the SQLite metadata store, and
// every path SimonDesktop persists data under.
type AppService struct {
	ctx            context.Context
	paths          platform.Paths
	db             *sql.DB
	repositories   *Repositories
	runtime        *simon.Runtime
	sessionManager *sessions.Manager
}

// New creates an AppService. Startup must be called (by Wails' OnStartup)
// before any other method is used.
func New() *AppService {
	return &AppService{}
}

// Startup resolves SimonDesktop's local data directories and creates the
// process's single simon.Runtime. Called by Wails once, before the
// frontend is shown.
func (a *AppService) Startup(ctx context.Context) {
	a.ctx = ctx

	paths, err := platform.Resolve()
	if err != nil {
		log.Printf("simondesktop: could not resolve application support directory: %v", err)
		return
	}
	a.paths = paths

	db, err := persistence.Open(paths.DBPath)
	if err != nil {
		log.Printf("simondesktop: could not open the metadata database: %v", err)
		return
	}
	a.db = db
	a.repositories = newRepositories(db)

	memoryFactory := memory.FactoryFunc(func(_ context.Context, sessionID string) (memory.Memory, error) {
		return memory.NewJSONFileIn(paths.ChatsDir, sessionID+".json"), nil
	})

	rt, err := simon.New(simon.WithMemoryFactory(memoryFactory))
	if err != nil {
		log.Printf("simondesktop: could not initialize Simon: %v", err)
		return
	}
	a.runtime = rt
	a.sessionManager = sessions.NewManager(rt)
}

// Shutdown closes the Runtime (and every Session it created) and the
// metadata database cleanly. Called by Wails once, as the application is
// closing.
func (a *AppService) Shutdown(ctx context.Context) {
	if a.runtime != nil {
		if err := a.runtime.Close(); err != nil {
			log.Printf("simondesktop: error closing Simon runtime: %v", err)
		}
	}
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			log.Printf("simondesktop: error closing database: %v", err)
		}
	}
}

// CreateConversation creates a new, untitled conversation.
func (a *AppService) CreateConversation() (conversations.Conversation, error) {
	if a.repositories == nil {
		return conversations.Conversation{}, fmt.Errorf("simondesktop: not initialized")
	}
	return a.repositories.Conversations.Create(a.ctx)
}

// ListConversations returns every conversation, most recently updated
// first.
func (a *AppService) ListConversations() ([]conversations.Conversation, error) {
	if a.repositories == nil {
		return nil, fmt.Errorf("simondesktop: not initialized")
	}
	return a.repositories.Conversations.List(a.ctx)
}

// GetConversationMessages returns conversationID's user/assistant history.
func (a *AppService) GetConversationMessages(conversationID string) ([]sessions.ChatMessage, error) {
	if a.paths.ChatsDir == "" {
		return nil, errors.New("Could not open the conversation.")
	}
	msgs, err := sessions.History(a.ctx, a.paths.ChatsDir, conversationID)
	if err != nil {
		log.Printf("simondesktop: could not read conversation %s: %v", conversationID, err)
		return nil, errors.New("Could not open the conversation.")
	}
	return msgs, nil
}

// SendMessage starts a run on conversationID's session with text, returning
// once the run has started (not once it finishes). Progress and outcome
// are reported asynchronously via the chat.run.* Wails events.
func (a *AppService) SendMessage(conversationID string, text string) error {
	if a.runtime == nil || a.sessionManager == nil {
		return errors.New("Could not initialize Simon.")
	}

	sess, err := a.sessionManager.GetOrCreate(conversationID)
	if err != nil {
		return errors.New(friendlyError(err))
	}

	events, err := sess.Stream(context.Background(), text)
	if err != nil {
		return errors.New(friendlyError(err))
	}

	if a.repositories != nil {
		if err := a.repositories.Conversations.RecordMessage(a.ctx, conversationID, text); err != nil {
			log.Printf("simondesktop: could not update conversation metadata: %v", err)
		}
	}

	go a.consumeRun(conversationID, events)
	return nil
}

// GetSetting returns the stored value for key, or "" if it has never been
// set.
func (a *AppService) GetSetting(key string) (string, error) {
	if a.repositories == nil {
		return "", fmt.Errorf("simondesktop: not initialized")
	}
	value, _, err := a.repositories.Settings.Get(a.ctx, key)
	return value, err
}

// SetSetting stores value under key, overwriting any existing value.
func (a *AppService) SetSetting(key string, value string) error {
	if a.repositories == nil {
		return fmt.Errorf("simondesktop: not initialized")
	}
	return a.repositories.Settings.Set(a.ctx, key, value)
}

// SetWindowScale resizes the window to BaseWindowWidth/BaseWindowHeight
// scaled by factor, so the window grows along with the UI's font size.
func (a *AppService) SetWindowScale(factor float64) error {
	if factor <= 0 {
		return fmt.Errorf("simondesktop: invalid scale factor")
	}
	width := int(math.Round(BaseWindowWidth * factor))
	height := int(math.Round(BaseWindowHeight * factor))
	wailsruntime.WindowSetSize(a.ctx, width, height)
	return nil
}

// CancelRun cancels conversationID's active run, if any.
func (a *AppService) CancelRun(conversationID string) error {
	if a.sessionManager == nil {
		return errors.New("Could not initialize Simon.")
	}
	a.sessionManager.Cancel(conversationID)
	return nil
}

// consumeRun drains a Session.Stream event channel, translating each
// simon.Event into a stable chat.run.* Wails event for the frontend.
func (a *AppService) consumeRun(conversationID string, events <-chan simon.Event) {
	for ev := range events {
		switch ev.Type {
		case simon.EventRunStarted:
			wailsruntime.EventsEmit(a.ctx, EventRunStarted, ChatRunEvent{ConversationID: conversationID})

		case simon.EventRunCompleted:
			resp, _ := ev.Data.(simon.Response)
			wailsruntime.EventsEmit(a.ctx, EventRunCompleted, ChatCompletedEvent{
				ConversationID: conversationID,
				Text:           resp.Text,
				Model:          resp.Model,
				Provider:       resp.Provider,
				Steps:          resp.Steps,
			})

		case simon.EventRunFailed:
			message := "The assistant could not complete this request."
			if data, ok := ev.Data.(map[string]any); ok {
				if raw, ok := data["error"].(string); ok {
					message = friendlyRunFailureMessage(raw)
				}
			}
			wailsruntime.EventsEmit(a.ctx, EventRunFailed, ChatFailedEvent{ConversationID: conversationID, Message: message})

		case simon.EventRunCancelled:
			wailsruntime.EventsEmit(a.ctx, EventRunCancelled, ChatRunEvent{ConversationID: conversationID})
		}
	}
}
