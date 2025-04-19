package cmd

import (
	"fmt"
	"github.com/dalryan/jot/internal/jot"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Open or create today's daily note",
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
		editor := cfg.Editor

		if _, err := os.Stat(notePath); err == nil {
			err := jot.RunEditor(editor, notePath)
			if err != nil {
				fmt.Println("Error running editor:", err)
				return
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
				fmt.Printf("Warning: Failed to load template '%s': %v\n", templateName, err)
			}
		}

		tempPath := filepath.Join(os.TempDir(), "jot-"+id+".md")
		if err := jot.WriteTempMarkdown(note, tempPath); err != nil {
			fmt.Println("Error:", err)
			return
		}

		err = jot.RunEditor(editor, tempPath)
		if err != nil {
			fmt.Println("Error running editor:", err)
			return
		}

		noteFinal, err := jot.ParseNoteFile(tempPath)
		if err != nil {
			fmt.Println("Error parsing edited note:", err)
			return
		}

		if err := jot.SaveNote(cfg, noteFinal); err != nil {
			fmt.Println("Error saving note:", err)
			return
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
