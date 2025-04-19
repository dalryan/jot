package jot

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// LoadTemplate loads a template from the templates directory and applies the given data.
// It returns the processed template content as a string.
func LoadTemplate(cfg *Config, name string, data map[string]string) (string, error) {
	if err := cfg.EnsureDirectories(); err != nil {
		return "", fmt.Errorf("failed to ensure directories exist for template '%s': %w", name, err)
	}

	path := filepath.Join(cfg.TemplatesDir(), name+".md")
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read template file '%s' at path '%s': %w", name, path, err)
	}

	tmpl, err := template.New(name).Parse(string(raw))
	if err != nil {
		return "", fmt.Errorf("failed to parse template '%s': %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template '%s' with provided data: %w", name, err)
	}
	return buf.String(), nil
}
