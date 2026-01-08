package cmd

import (
	"fmt"
	"github.com/dalryan/jot/internal/jot"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [title]",
	Short: "Create a new note in your editor",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := cfg.EnsureDirectories(); err != nil {
			fmt.Fprintln(os.Stderr, "Error ensuring directories exist:", err)
			os.Exit(1)
		}

		title := ""
		if len(args) > 0 {
			title = args[0]
		}

		tags, _ := cmd.Flags().GetStringSlice("tag")
		links, _ := cmd.Flags().GetStringSlice("link")
		templateName, _ := cmd.Flags().GetString("template")
		explicitContext, _ := cmd.Flags().GetString("context")

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
			Content:   "",
		}

		if title != "" {
			note.Content = "# " + title + "\n\n"
		}

		if templateName != "" {
			content, err := jot.LoadTemplate(cfg, templateName, map[string]string{
				"date":    time.Now().Format("2006-01-02"),
				"context": context,
				"title":   title,
			})
			if err == nil {
				note.Content += content
			} else {
				fmt.Fprintf(os.Stderr, "Warning: failed to load template '%s': %v\n", templateName, err)
			}
		}

		tempPath := filepath.Join(os.TempDir(), "jot-"+id+".md")
		if err := jot.WriteTempMarkdown(note, tempPath); err != nil {
			fmt.Fprintln(os.Stderr, "Error writing temp file:", err)
			os.Exit(1)
		}

		if err := jot.RunEditor(cfg.Editor, tempPath); err != nil {
			fmt.Fprintln(os.Stderr, "Error opening editor:", err)
			os.Exit(1)
		}

		noteFinal, err := jot.ParseNoteFile(tempPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading edited note:", err)
			os.Exit(1)
		}

		if err := jot.SaveNote(cfg, noteFinal); err != nil {
			fmt.Fprintln(os.Stderr, "Error saving note:", err)
			os.Exit(1)
		}

		fmt.Printf("Note saved: %s\n", noteFinal.ID)
	},
}

// init sets up the new command and its flags.
// This function registers the new command with the root command and
// defines the available flags for tags, links, context, and template selection.
func init() {
	newCmd.Flags().StringSlice("tag", nil, "Tags for the note")
	newCmd.Flags().StringSlice("link", nil, "Links to other notes")
	newCmd.Flags().String("context", "", "Context for the note")
	newCmd.Flags().String("template", "", "Use a template (from templates directory)")
	rootCmd.AddCommand(newCmd)
}
