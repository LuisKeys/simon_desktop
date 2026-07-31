package sessions

import (
	"context"

	"github.com/LuisKeys/simon/memory"
)

// ChatMessage is one user/assistant turn, shaped for the frontend.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// History reads conversationID's persisted memory (under chatsDir) and
// returns its user/assistant turns, in order. System and tool messages are
// omitted — the frontend only ever renders user/assistant turns.
func History(ctx context.Context, chatsDir, conversationID string) ([]ChatMessage, error) {
	mem := memory.NewJSONFileIn(chatsDir, conversationID+".json")
	defer mem.Close()

	msgs, err := mem.List(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]ChatMessage, 0, len(msgs))
	for _, m := range msgs {
		if m.Role != memory.RoleUser && m.Role != memory.RoleAssistant {
			continue
		}
		out = append(out, ChatMessage{Role: string(m.Role), Content: m.Content})
	}
	return out, nil
}
