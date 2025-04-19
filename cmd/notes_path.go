package cmd

import (
	"fmt"
	"github.com/dalryan/jot/internal/jot"
	"os"

	"github.com/spf13/cobra"
)

var notesPathCmd = &cobra.Command{
	Use:   "notes-path",
	Short: "Print the path to the notes directory",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := jot.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}
		fmt.Println(cfg.NotesDir())
	},
}

func init() {
	rootCmd.AddCommand(notesPathCmd)
}
