package jot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FindNoteByID locates and loads a note by its ID or ID prefix.
// It searches the notes directory for a file with a name starting with the given ID.
// Returns the parsed Note if found, or an error if the note doesn't exist or can't be parsed.
func FindNoteByID(baseDir, id string) (*Note, error) {
	noteDir := filepath.Join(baseDir, "notes")

	files, err := os.ReadDir(noteDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read notes directory at path '%s': %w", noteDir, err)
	}

	var match string
	for _, f := range files {
		if strings.HasPrefix(f.Name(), id) && strings.HasSuffix(f.Name(), ".md") {
			match = filepath.Join(noteDir, f.Name())
			break
		}
	}

	if match == "" {
		return nil, fmt.Errorf("note with ID or ID prefix '%s' not found in directory '%s'", id, noteDir)
	}

	note, err := ParseNoteFile(match)
	if err != nil {
		return nil, fmt.Errorf("failed to parse note with ID prefix '%s' at path '%s': %w", id, match, err)
	}
	return note, nil
}

// ResolveNotePath finds the full file path of a note by its ID or ID prefix.
// Unlike FindNoteByID, this function only returns the path to the note file, not the parsed note.
// Returns the full path if found, or an error if the note doesn't exist.
func ResolveNotePath(baseDir, id string) (string, error) {
	noteDir := filepath.Join(baseDir, "notes")

	files, err := os.ReadDir(noteDir)
	if err != nil {
		return "", fmt.Errorf("failed to read notes directory at path '%s': %w", noteDir, err)
	}

	for _, f := range files {
		if strings.HasPrefix(f.Name(), id) && strings.HasSuffix(f.Name(), ".md") {
			return filepath.Join(noteDir, f.Name()), nil
		}
	}

	return "", fmt.Errorf("note with ID or ID prefix '%s' not found in directory '%s'", id, noteDir)
}
