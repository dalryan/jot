package cmd

import (
	"fmt"
	"github.com/dalryan/jot/internal/jot"
	"github.com/google/uuid"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var quickCmd = &cobra.Command{
	Use:   "quick [message]",
	Short: "Capture a quick, timestamped note",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var message string

		if len(args) > 0 {
			message = strings.Join(args, " ")
		} else {
			stat, err := os.Stdin.Stat()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to read stdin:", err)
				os.Exit(1)
			}

			if (stat.Mode() & os.ModeCharDevice) == 0 {
				input, err := io.ReadAll(os.Stdin)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
					os.Exit(1)
				}
				message = strings.TrimSpace(string(input))
			}
		}

		if message == "" {
			fmt.Fprintln(os.Stderr, "Error: no note content provided (use args or pipe)")
			os.Exit(1)
		}
		tags, _ := cmd.Flags().GetStringSlice("tag")
		links, _ := cmd.Flags().GetStringSlice("link")
		explicitContext, _ := cmd.Flags().GetString("context")

		if err := cfg.EnsureDirectories(); err != nil {
			fmt.Println("Error ensuring directories exist:", err)
			os.Exit(1)
		}

		context := explicitContext
		if context == "" {
			activeContext, err := jot.GetActiveContext(cfg.StoragePath)
			if err == nil && activeContext != "" {
				context = activeContext
			} else {
				context = cfg.DefaultContext
			}
		}

		id := uuid.New().String()[:8]
		now := time.Now()

		note := &jot.Note{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
			Context:   context,
			Tags:      tags,
			Links:     links,
			Content:   message,
		}

		err := jot.SaveNote(cfg, note)
		if err != nil {
			fmt.Println("Failed to save note:", err)
			return
		}

		fmt.Printf("Quick note saved: %s\n", note.ID)
	},
}

// init sets up the flags for the quick command.
// This function defines the available flags for the quick command,
// including tags, links, and context override.
func init() {
	quickCmd.Flags().StringSlice("tag", nil, "Tags for the note (comma-separated or repeat)")
	quickCmd.Flags().StringSlice("link", nil, "Links to other notes")
	quickCmd.Flags().String("context", "", "Override the active context")
	rootCmd.AddCommand(quickCmd)
}
