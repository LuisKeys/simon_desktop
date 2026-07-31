package app

// Wails event names emitted for a conversation's run lifecycle. Stable,
// frontend-facing identifiers — do not rename without updating the
// frontend's EventsOn listeners.
const (
	EventRunStarted   = "chat.run.started"
	EventRunCompleted = "chat.run.completed"
	EventRunFailed    = "chat.run.failed"
	EventRunCancelled = "chat.run.cancelled"
)

// ChatRunEvent is the payload for chat.run.started and chat.run.cancelled.
type ChatRunEvent struct {
	ConversationID string `json:"conversationId"`
}

// ChatCompletedEvent is the payload for chat.run.completed.
type ChatCompletedEvent struct {
	ConversationID string `json:"conversationId"`
	Text           string `json:"text"`
	Model          string `json:"model"`
	Provider       string `json:"provider"`
	Steps          int    `json:"steps"`
}

// ChatFailedEvent is the payload for chat.run.failed.
type ChatFailedEvent struct {
	ConversationID string `json:"conversationId"`
	Message        string `json:"message"`
}
