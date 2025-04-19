package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dalryan/jot/internal/jot"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	listCmd.Flags().StringSlice("tag", nil, "Filter notes by tag(s)")
	listCmd.Flags().String("context", "", "Override or set the context filter")
	listCmd.Flags().Bool("json", false, "Output notes as JSON")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List existing notes",
	Run: func(cmd *cobra.Command, args []string) {
		baseDir := cfg.StoragePath
		filterTags, _ := cmd.Flags().GetStringSlice("tag")
		filterContext, _ := cmd.Flags().GetString("context")

		if filterContext == "" {
			ctx, err := jot.GetActiveContext(baseDir)
			if err == nil {
				filterContext = ctx
			}
		}

		notes, err := jot.LoadAllNotes(baseDir)
		if err != nil {
			fmt.Println("Error loading notes:", err)
			return
		}

		outputJSON, _ := cmd.Flags().GetBool("json")

		if outputJSON {
			if notes == nil {
				notes = []*jot.Note{}
			}
			err = json.NewEncoder(os.Stdout).Encode(notes)
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
				return
			}
		} else {
			for _, n := range notes {
				if !jot.HasAllTags(n, filterTags) {
					continue
				}
				if filterContext != "" && n.Context != filterContext {
					continue
				}

				summary := fmt.Sprintf(
					"%-8s  %s  %-20s  %s",
					n.ID[:8],
					n.CreatedAt.Format("2006-01-02"),
					fmt.Sprintf("[%s]", joinStrings(n.Tags, ",")),
					firstLine(n.Content),
				)
				fmt.Println(summary)
			}
		}
	}}

func joinStrings(ss []string, sep string) string {
	return strings.Join(ss, sep)
}

func firstLine(s string) string {
	if idx := strings.IndexByte(s, '\n'); idx != -1 {
		return s[:idx]
	}
	return s
}

func init() {
	rootCmd.AddCommand(listCmd)
}
