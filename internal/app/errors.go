package app

import (
	"context"
	"errors"
	"strings"

	"github.com/LuisKeys/simon/pkg/simonerr"
)

// friendlyError turns a Go error returned directly from the Simon SDK into
// one of the simple messages the UI shows the user.
func friendlyError(err error) string {
	switch {
	case err == nil:
		return ""
	case errors.Is(err, context.Canceled), errors.Is(err, simonerr.ErrRunCancelled):
		return "The request was cancelled."
	case errors.Is(err, simonerr.ErrProvider):
		return "No model provider is configured."
	case errors.Is(err, simonerr.ErrRuntimeClosed), errors.Is(err, simonerr.ErrSessionClosed):
		return "Could not open the conversation."
	case errors.Is(err, simonerr.ErrSessionBusy):
		return "A response is already in progress for this conversation."
	default:
		return err.Error()
	}
}

// friendlyRunFailureMessage turns the error message string carried by a
// run.failed event into one of the simple UI messages. The SDK only
// exposes the error's string form on this event (see session.go's finish),
// so this is best-effort keyword matching rather than errors.Is.
func friendlyRunFailureMessage(raw string) string {
	lower := strings.ToLower(raw)
	switch {
	case strings.Contains(lower, "provider"):
		return "No model provider is configured."
	case strings.Contains(lower, "cancel"):
		return "The request was cancelled."
	default:
		return raw
	}
}
