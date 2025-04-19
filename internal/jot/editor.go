package jot

import (
	"os"
	"os/exec"
)

// WriteTempMarkdown converts a Note to markdown format and writes it to a temporary file.
// It takes a Note and a file path as input and returns an error if the operation fails.
func WriteTempMarkdown(n *Note, path string) error {
	md, err := n.ToMarkdown()
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(md), 0644)
}

// RunEditor launches an external editor to edit a file.
// It takes the editor command and file path as input and returns an error if the operation fails.
// The function connects the editor's standard input, output, and error to the current process.
func RunEditor(editor, path string) error {
	c := exec.Command(editor, path)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
