package jot

import "strings"

// HasAllTags checks if a note contains all the specified tags.
// It returns true if the note has all the tags in the provided list, or if the list is empty.
// Returns false if any tag is missing from the note.
func HasAllTags(note *Note, tags []string) bool {
	if len(tags) == 0 {
		return true
	}

	tagSet := make(map[string]bool)
	for _, t := range note.Tags {
		tagSet[t] = true
	}

	for _, ft := range tags {
		if !tagSet[ft] {
			return false
		}
	}
	return true
}

// JoinTags combines a slice of tags into a single comma-separated string.
// This is useful for displaying tags in a human-readable format.
func JoinTags(tags []string) string {
	return strings.Join(tags, ", ")
}

// FirstLine extracts the first line from a multi-line string.
// If the string doesn't contain a newline character, the entire string is returned.
func FirstLine(s string) string {
	if idx := strings.IndexByte(s, '\n'); idx != -1 {
		return s[:idx]
	}
	return s
}
