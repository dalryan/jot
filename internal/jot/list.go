package jot

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// LoadAllNotes loads all notes from the notes directory.
// It recursively walks through the directory and parses all markdown files.
// Returns a slice of all successfully parsed notes and any error encountered during directory traversal.
// Note that parsing errors for individual files are logged to stderr but don't stop the process.
func LoadAllNotes(baseDir string) ([]*Note, error) {
	noteDir := filepath.Join(baseDir, "notes")
	var notes []*Note

	err := filepath.WalkDir(noteDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		n, err := ParseNoteFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse note file at path '%s': %v\n", path, err)
			return nil
		}

		notes = append(notes, n)
		return nil
	})

	if err != nil {
		return notes, fmt.Errorf("failed to walk notes directory at path '%s': %w", noteDir, err)
	}

	return notes, nil
}
