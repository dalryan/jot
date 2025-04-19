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
		cfg, err := jot.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}

		if err := cfg.EnsureDirectories(); err != nil {
			fmt.Println("Error ensuring directories exist:", err)
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
			context = cfg.DefaultContext
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
				fmt.Printf("Warning: Failed to load template '%s': %v\n", templateName, err)
			}
		}

		tempPath := filepath.Join(os.TempDir(), "jot-"+id+".md")
		if err := jot.WriteTempMarkdown(note, tempPath); err != nil {
			fmt.Println("Error:", err)
			return
		}

		err = jot.RunEditor(cfg.Editor, tempPath)
		if err != nil {
			fmt.Println("Failed to open editor:", err)
			return
		}

		noteFinal, err := jot.ParseNoteFile(tempPath)
		if err != nil {
			fmt.Println("Failed to read updated note:", err)
			return
		}

		err = jot.SaveNote(cfg, noteFinal)
		if err != nil {
			fmt.Println("Failed to save note:", err)
			return
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
