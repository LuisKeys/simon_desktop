package conversations_test

import (
	"context"
	"path/filepath"
	"testing"

	"simonpartner/internal/conversations"
	"simonpartner/internal/persistence"
)

func TestCreateAndList(t *testing.T) {
	dir := t.TempDir()
	db, err := persistence.Open(filepath.Join(dir, "app.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer db.Close()

	repo := conversations.NewRepository(db)
	ctx := context.Background()

	created, err := repo.Create(ctx)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.Title != conversations.DefaultTitle {
		t.Errorf("Title = %q, want %q", created.Title, conversations.DefaultTitle)
	}

	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 1 || list[0].ID != created.ID {
		t.Fatalf("List = %+v, want one conversation with ID %q", list, created.ID)
	}
}

func TestListOrdersByUpdatedAtDesc(t *testing.T) {
	dir := t.TempDir()
	db, err := persistence.Open(filepath.Join(dir, "app.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer db.Close()

	repo := conversations.NewRepository(db)
	ctx := context.Background()

	first, err := repo.Create(ctx)
	if err != nil {
		t.Fatalf("Create first: %v", err)
	}
	second, err := repo.Create(ctx)
	if err != nil {
		t.Fatalf("Create second: %v", err)
	}

	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("List returned %d conversations, want 2", len(list))
	}
	if list[0].ID != second.ID || list[1].ID != first.ID {
		t.Errorf("List order = [%s, %s], want most-recently-created first: [%s, %s]", list[0].ID, list[1].ID, second.ID, first.ID)
	}
}

func TestRecordMessageDerivesTitleOnce(t *testing.T) {
	dir := t.TempDir()
	db, err := persistence.Open(filepath.Join(dir, "app.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer db.Close()

	repo := conversations.NewRepository(db)
	ctx := context.Background()

	created, err := repo.Create(ctx)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	longText := "This is a much longer first message that should be truncated to fifty characters for the title"
	if err := repo.RecordMessage(ctx, created.ID, longText); err != nil {
		t.Fatalf("RecordMessage: %v", err)
	}

	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("List = %+v, want one conversation", list)
	}
	got := list[0].Title
	wantPrefix := string([]rune(longText)[:50]) + "…"
	if got != wantPrefix {
		t.Errorf("Title = %q, want %q", got, wantPrefix)
	}

	// A second message must not overwrite the derived title.
	if err := repo.RecordMessage(ctx, created.ID, "second message"); err != nil {
		t.Fatalf("RecordMessage (second): %v", err)
	}
	list, err = repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if list[0].Title != got {
		t.Errorf("Title after second message = %q, want unchanged %q", list[0].Title, got)
	}
}
