package cmd

import (
	"fmt"
	"github.com/dalryan/jot/internal/jot"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit a note by ID",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		baseDir := cfg.StoragePath

		notePath, err := jot.ResolveNotePath(baseDir, id)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		editor := cfg.Editor

		c := exec.Command(editor, notePath)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		if err := c.Run(); err != nil {
			fmt.Printf("Error opening editor: %v\n", err)
			return
		}

		note, err := jot.ParseNoteFile(notePath)
		if err != nil {
			fmt.Printf("Warning: could not parse edited note: %v\n", err)
			return
		}

		note.UpdateTimestamp()
		if err := jot.SaveNote(cfg, note); err != nil {
			fmt.Printf("Warning: could not update timestamp: %v\n", err)
		}

		fmt.Printf("Updated note %s\n", note.ID)
	},
}

// init registers the edit command with the root command.
// This function is automatically called by Go when the package is initialized.
func init() {
	rootCmd.AddCommand(editCmd)
}
