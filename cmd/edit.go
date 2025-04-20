package cmd

import (
	"bufio"
	"fmt"
	"github.com/dalryan/jot/internal/jot"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var editCmd = &cobra.Command{
	Use:   "edit [id]",
	Short: "Edit a note by ID or from stdin",
	Run: func(cmd *cobra.Command, args []string) {
		var id string
		baseDir := cfg.StoragePath

		stat, err := os.Stdin.Stat()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to read stdin:", err)
			os.Exit(1)
		}

		if len(args) > 0 {
			id = args[0]
		} else if (stat.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				line := scanner.Text()
				if len(line) >= 8 {
					id = line[:8]
				} else {
					fmt.Println("Error: Input too short")
					os.Exit(1)
				}
			} else {
				fmt.Println("Error: No input provided")
				os.Exit(1)
			}
		} else {
			fmt.Println("Error: No ID provided")
			os.Exit(1)
		}

		notePath, err := jot.ResolveNotePath(baseDir, id)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		editor := cfg.Editor

		c := exec.Command(editor, notePath)

		// If stdin is a pipe, open /dev/tty to ensure the editor gets input from the terminal
		// This prevents the "Input is not from a terminal" warning when piping to jot edit
		// might be a vim specific issue. However, that's what I tested with.
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			tty, err := os.Open("/dev/tty")
			if err != nil {
				c.Stdin = os.Stdin
			} else {
				c.Stdin = tty
				defer func(tty *os.File) {
					err := tty.Close()
					if err != nil {
						fmt.Printf("Error closing tty: %v\n", err)
					}
				}(tty)
			}
		} else {
			c.Stdin = os.Stdin
		}

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
