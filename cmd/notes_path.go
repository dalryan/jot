package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var notesPathCmd = &cobra.Command{
	Use:   "notes-path",
	Short: "Print the path to the notes directory",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cfg.NotesDir())
	},
}

func init() {
	rootCmd.AddCommand(notesPathCmd)
}
