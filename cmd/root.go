package cmd

import (
	"fmt"
	"github.com/dalryan/jot/internal/jot"
	"os"

	"github.com/spf13/cobra"
)

var cfg *jot.Config

var rootCmd = &cobra.Command{
	Use:   "jot",
	Short: "jot is a minimalist CLI note-taking app",
	Long:  `Jot lets you quickly create, view, and organize plain-text notes from the command line.`,
}

func Execute() error {
	var err error
	cfg, err = jot.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load config:", err)
		os.Exit(1)
	}
	err = rootCmd.Execute()
	if err != nil {
		fmt.Println("Failed to execute:", err)
	}
	return err
}

func init() {
	rootCmd.AddCommand(quickCmd)
	rootCmd.AddCommand(listCmd)
}
