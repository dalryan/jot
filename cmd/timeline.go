package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dalryan/jot/internal/jot"
	"os"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

var timelineCmd = &cobra.Command{
	Use:   "timeline",
	Short: "Show notes in reverse chronological order",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := jot.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}
		baseDir := cfg.StoragePath

		tagFilter, _ := cmd.Flags().GetStringSlice("tag")
		contextFilter, _ := cmd.Flags().GetString("context")
		sinceStr, _ := cmd.Flags().GetString("since")
		beforeStr, _ := cmd.Flags().GetString("before")
		limit, _ := cmd.Flags().GetInt("limit")

		// should we crash if the flags error?

		notes, err := jot.LoadAllNotes(baseDir)
		if err != nil {
			fmt.Println("Error loading notes:", err)
			return
		}

		var since, before time.Time
		if sinceStr != "" {
			since, _ = parseTime(sinceStr)
		}
		if beforeStr != "" {
			before, _ = parseTime(beforeStr)
		}

		var filtered []*jot.Note
		for _, n := range notes {
			if !jot.HasAllTags(n, tagFilter) {
				continue
			}
			if contextFilter != "" && n.Context != contextFilter {
				continue
			}
			if !since.IsZero() && n.CreatedAt.Before(since) {
				continue
			}
			if !before.IsZero() && n.CreatedAt.After(before) {
				continue
			}
			filtered = append(filtered, n)
		}

		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
		})

		if limit > 0 && len(filtered) > limit {
			filtered = filtered[:limit]
		}

		outputJSON, _ := cmd.Flags().GetBool("json")

		if outputJSON {
			if filtered == nil {
				filtered = []*jot.Note{}
			}
			err = json.NewEncoder(os.Stdout).Encode(filtered)
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
				return 
			}
		} else {
			for _, n := range filtered {
				fmt.Printf("[%s]  %-8s  %-12s  %s\n",
					n.CreatedAt.Format("2006-01-02 15:04"),
					n.ID[:8],
					n.Context,
					jot.FirstLine(n.Content),
				)
			}
		}
	},
}

func parseTime(input string) (time.Time, error) {
	if d, err := time.Parse("2006-01-02", input); err == nil {
		return d, nil
	}
	if dur, err := time.ParseDuration(input); err == nil {
		return time.Now().Add(-dur), nil
	}
	return time.Time{}, fmt.Errorf("invalid time format: %s", input)
}

func init() {
	timelineCmd.Flags().StringSlice("tag", nil, "Filter by tag(s)")
	timelineCmd.Flags().String("context", "", "Filter by context")
	timelineCmd.Flags().String("since", "", "Only notes after (e.g. '7d' or '2025-04-01')")
	timelineCmd.Flags().String("before", "", "Only notes before a date")
	timelineCmd.Flags().Int("limit", 0, "Limit number of results")
	timelineCmd.Flags().Bool("json", false, "Output notes as JSON")
	rootCmd.AddCommand(timelineCmd)
}
