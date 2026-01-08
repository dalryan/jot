package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dalryan/jot/internal/jot"

	"github.com/spf13/cobra"
)

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Open or create today's daily note",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cfg.EnsureDirectories(); err != nil {
			fmt.Fprintln(os.Stderr, "Error ensuring directories exist:", err)
			os.Exit(1)
		}

		today := time.Now().Format("2006-01-02")
		id := "today-" + time.Now().Format("20060102")
		title := "Journal for " + today

		contextFlag, _ := cmd.Flags().GetString("context")
		templateName, _ := cmd.Flags().GetString("template")

		context := contextFlag
		if context == "" {
			context = "journal"
		}

		notePath := filepath.Join(cfg.NotesDir(), id+".md")

		// If today's note already exists, just open it
		if _, err := os.Stat(notePath); err == nil {
			if err := jot.RunEditor(cfg.Editor, notePath); err != nil {
				fmt.Fprintln(os.Stderr, "Error running editor:", err)
				os.Exit(1)
			}
			return
		}

		note := &jot.Note{
			ID:        id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Context:   context,
			Content:   "# " + title + "\n\n",
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
			fmt.Fprintln(os.Stderr, "Error running editor:", err)
			os.Exit(1)
		}

		noteFinal, err := jot.ParseNoteFile(tempPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing edited note:", err)
			os.Exit(1)
		}

		if err := jot.SaveNote(cfg, noteFinal); err != nil {
			fmt.Fprintln(os.Stderr, "Error saving note:", err)
			os.Exit(1)
		}

		fmt.Printf("Journal saved: %s\n", noteFinal.ID)
	},
}

// init sets up the today command and its flags.
// This function registers the today command with the root command and
// defines the available flags for context and template selection.
func init() {
	todayCmd.Flags().String("context", "", "Context for the note (default: journal)")
	todayCmd.Flags().String("template", "", "Template name (e.g. 'daily')")
	rootCmd.AddCommand(todayCmd)
}
