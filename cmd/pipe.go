package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/dalryan/jot/internal/jot"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var pipeCmd = &cobra.Command{
	Use:   "pipe",
	Short: "Parse note file paths from stdin and display summaries",
	Run: func(cmd *cobra.Command, args []string) {
		tagFilter, _ := cmd.Flags().GetStringSlice("tag")
		contextFilter, _ := cmd.Flags().GetString("context")

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			path := scanner.Text()
			absPath, _ := filepath.Abs(path)

			note, err := jot.ParseNoteFile(absPath)
			if err != nil {
				continue
			}
			if !jot.HasAllTags(note, tagFilter) {
				continue
			}
			if contextFilter != "" && note.Context != contextFilter {
				continue
			}

			outputJSON, _ := cmd.Flags().GetBool("json")

			if outputJSON {
				if note == nil {
					note = &jot.Note{}
				}
				err = json.NewEncoder(os.Stdout).Encode(note)
				if err != nil {
					fmt.Println("Error encoding JSON:", err)
					return
				}
			} else {
				fmt.Printf("🧠 %s  %s  [%s]  %s\n",
					note.ID[:8],
					note.CreatedAt.Format("2006-01-02"),
					jot.JoinTags(note.Tags),
					jot.FirstLine(note.Content),
				)
			}
		}
	},
}

func init() {
	pipeCmd.Flags().StringSlice("tag", nil, "Filter by tag(s)")
	pipeCmd.Flags().String("context", "", "Filter by context")
	pipeCmd.Flags().Bool("json", false, "Output notes as JSON")
	rootCmd.AddCommand(pipeCmd)
}
