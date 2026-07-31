// Package platform resolves where SimonDesktop stores its local data on
// disk, independent of the process's working directory.
package platform

import (
	"os"
	"path/filepath"
)

// Paths holds every directory/file SimonDesktop persists data under, all
// rooted at ~/Library/Application Support/SimonDesktop.
type Paths struct {
	Root           string
	DBPath         string
	ChatsDir       string
	KnowledgeDir   string
	AttachmentsDir string
}

// Resolve returns SimonDesktop's data paths and ensures every directory in
// them exists, creating them on first run.
func Resolve() (Paths, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Paths{}, err
	}

	root := filepath.Join(home, "Library", "Application Support", "SimonDesktop")
	p := Paths{
		Root:           root,
		DBPath:         filepath.Join(root, "app.db"),
		ChatsDir:       filepath.Join(root, "chats"),
		KnowledgeDir:   filepath.Join(root, "knowledge"),
		AttachmentsDir: filepath.Join(root, "attachments"),
	}

	for _, dir := range []string{p.Root, p.ChatsDir, p.KnowledgeDir, p.AttachmentsDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return Paths{}, err
		}
	}

	return p, nil
}
