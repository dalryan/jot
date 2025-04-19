package jot

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Note represents a note in the jot system with metadata and content.
type Note struct {
	// ID is the unique identifier for the note.
	ID string `yaml:"id" json:"id"`
	// CreatedAt is the timestamp when the note was created.
	CreatedAt time.Time `yaml:"created_at" json:"created_at"`
	// UpdatedAt is the timestamp when the note was last updated.
	UpdatedAt time.Time `yaml:"updated_at" json:"updated_at"`
	// Tags is a list of tags associated with the note.
	Tags []string `yaml:"tags,omitempty" json:"tags"`
	// Links is a list of references to other notes or resources.
	Links []string `yaml:"links,omitempty" json:"links"`
	// Content is the main text content of the note.
	Content string `yaml:"-" json:"content"`
	// Context is the organizational context the note belongs to.
	Context string `yaml:"context,omitempty" json:"context,omitempty"`
}

// ToMarkdown converts a Note to a markdown string with YAML frontmatter.
// The frontmatter contains the note's metadata, and the content follows after.
// Returns the formatted markdown string and any error encountered during conversion.
func (n *Note) ToMarkdown() (string, error) {
	meta := struct {
		ID        string    `yaml:"id"`
		CreatedAt time.Time `yaml:"created_at"`
		UpdatedAt time.Time `yaml:"updated_at"`
		Tags      []string  `yaml:"tags,omitempty"`
		Links     []string  `yaml:"links,omitempty"`
		Context   string    `yaml:"context,omitempty"`
	}{
		ID:        n.ID,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
		Tags:      n.Tags,
		Links:     n.Links,
		Context:   n.Context,
	}

	var buf bytes.Buffer
	buf.WriteString("---\n")
	yml, err := yaml.Marshal(meta)
	if err != nil {
		return "", fmt.Errorf("failed to marshal note metadata to YAML for note ID '%s': %w", n.ID, err)
	}
	buf.Write(yml)
	buf.WriteString("---\n\n")
	buf.WriteString(n.Content)
	buf.WriteByte('\n')

	return buf.String(), nil
}

// SaveNote saves a note to the notes directory.
// It converts the note to markdown format and writes it to a file.
func SaveNote(cfg *Config, note *Note) error {
	if err := cfg.EnsureDirectories(); err != nil {
		return fmt.Errorf("failed to ensure directories exist for saving note ID '%s': %w", note.ID, err)
	}

	md, err := note.ToMarkdown()
	if err != nil {
		return fmt.Errorf("failed to convert note ID '%s' to markdown: %w", note.ID, err)
	}

	notePath := filepath.Join(cfg.NotesDir(), fmt.Sprintf("%s.md", note.ID))
	if err := os.WriteFile(notePath, []byte(md), 0644); err != nil {
		return fmt.Errorf("failed to write note ID '%s' to file path '%s': %w", note.ID, notePath, err)
	}
	return nil
}

// ParseNoteFile reads a markdown file with YAML frontmatter and converts it to a Note.
// It extracts metadata from the frontmatter and the content from the rest of the file.
// Returns the parsed Note and any error encountered during parsing.
func ParseNoteFile(path string) (*Note, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read note file at path '%s': %w", path, err)
	}

	parts := strings.SplitN(string(data), "---\n", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid frontmatter format in note file '%s': missing YAML delimiters", path)
	}

	metaPart := parts[1]
	content := strings.TrimSpace(parts[2])

	n := &Note{}
	if err := yaml.Unmarshal([]byte(metaPart), n); err != nil {
		return nil, fmt.Errorf("failed to parse YAML frontmatter in note file '%s': %w", path, err)
	}
	n.Content = content
	if n.Tags == nil {
		n.Tags = []string{}
	}
	if n.Links == nil {
		n.Links = []string{}
	}

	return n, nil
}

// UpdateTimestamp sets the UpdatedAt field of the note to the current time.
// This should be called whenever the note content is modified.
func (n *Note) UpdateTimestamp() {
	n.UpdatedAt = time.Now()
}
