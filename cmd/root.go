package cmd

import (
	"fmt"

	"github.com/dalryan/jot/internal/jot"
	"github.com/spf13/cobra"
)

var cfg *jot.Config

var rootCmd = &cobra.Command{
	Use:   "jot",
	Short: "jot is a minimalist CLI note-taking app",
	Long:  `Jot lets you quickly create, view, and organize plain-text notes from the command line.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = jot.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
