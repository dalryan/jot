package cmd

import (
	"fmt"
	"github.com/dalryan/jot/internal/jot"
	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Manage the active context",
}

var contextSetCmd = &cobra.Command{
	Use:   "set <name>",
	Short: "Set the current context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		baseDir := cfg.StoragePath
		err := jot.SetActiveContext(baseDir, name)
		if err != nil {
			fmt.Println("Error setting context:", err)
			return
		}
		fmt.Printf("Context set to: %s\n", name)
	},
}

var contextGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Show the current context",
	Run: func(cmd *cobra.Command, args []string) {
		baseDir := cfg.StoragePath
		ctx, err := jot.GetActiveContext(baseDir)
		if err != nil || ctx == "" {
			fmt.Println("No context is currently set.")
			return
		}
		fmt.Println("Current context:", ctx)
	},
}

var contextClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the current context",
	Run: func(cmd *cobra.Command, args []string) {
		baseDir := cfg.StoragePath
		err := jot.ClearContext(baseDir)
		if err != nil {
			fmt.Println("Error clearing context:", err)
			return
		}
		fmt.Println("Context cleared.")
	},
}

// init registers the context commands with the root command.
// This function sets up the command hierarchy for context management,
// adding set, get, and clear subcommands to the context command.
func init() {
	contextCmd.AddCommand(contextSetCmd)
	contextCmd.AddCommand(contextGetCmd)
	contextCmd.AddCommand(contextClearCmd)
	rootCmd.AddCommand(contextCmd)
}
