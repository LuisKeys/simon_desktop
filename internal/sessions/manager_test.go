package sessions_test

import (
	"context"
	"testing"

	"github.com/LuisKeys/simon"
	"github.com/LuisKeys/simon/model"

	"simondesktop/internal/sessions"
)

func newTestRuntime(t *testing.T) *simon.Runtime {
	t.Helper()
	rt, err := simon.New(simon.WithModel(model.EchoModel{}))
	if err != nil {
		t.Fatalf("simon.New: %v", err)
	}
	t.Cleanup(func() { rt.Close() })
	return rt
}

func TestGetOrCreateReusesSessionPerConversation(t *testing.T) {
	mgr := sessions.NewManager(newTestRuntime(t))

	first, err := mgr.GetOrCreate("conversation-a")
	if err != nil {
		t.Fatalf("GetOrCreate: %v", err)
	}
	again, err := mgr.GetOrCreate("conversation-a")
	if err != nil {
		t.Fatalf("GetOrCreate (again): %v", err)
	}
	if first != again {
		t.Error("expected GetOrCreate to return the same *simon.Session for the same conversation ID")
	}

	other, err := mgr.GetOrCreate("conversation-b")
	if err != nil {
		t.Fatalf("GetOrCreate (other): %v", err)
	}
	if other == first {
		t.Error("expected a distinct *simon.Session for a different conversation ID")
	}
}

func TestManagerSessionRuns(t *testing.T) {
	mgr := sessions.NewManager(newTestRuntime(t))

	sess, err := mgr.GetOrCreate("conversation-a")
	if err != nil {
		t.Fatalf("GetOrCreate: %v", err)
	}

	resp, err := sess.Run(context.Background(), "hello")
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp.Text == "" {
		t.Error("expected EchoModel to return non-empty text")
	}
}
