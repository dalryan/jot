package jot

import (
	"os"
	"path/filepath"
	"strings"
)

// GetActiveContext retrieves the currently active context from the context file.
// It takes the base directory path as input and returns the active context name.
// If the context file doesn't exist or can't be read, it returns an empty string.
func GetActiveContext(baseDir string) (string, error) {
	data, err := os.ReadFile(filepath.Join(baseDir, "context"))
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(string(data)), nil
}

// SetActiveContext sets the specified context as the active one by writing it to the context file.
// It takes the base directory path and the context name as input.
// Returns an error if the file cannot be written.
func SetActiveContext(baseDir, name string) error {
	return os.WriteFile(filepath.Join(baseDir, "context"), []byte(name+"\n"), 0644)
}

// ClearContext removes the context file, effectively clearing the active context.
// It takes the base directory path as input.
// Returns an error if the file cannot be removed.
func ClearContext(baseDir string) error {
	return os.Remove(filepath.Join(baseDir, "context"))
}
