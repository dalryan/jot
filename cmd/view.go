package cmd

import (
	"fmt"
	"github.com/dalryan/jot/internal/jot"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View a note by its ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		baseDir := cfg.StoragePath

		note, err := jot.FindNoteByID(baseDir, id)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		raw, _ := cmd.Flags().GetBool("raw")
		pretty, _ := cmd.Flags().GetBool("pretty")

		if raw {
			md, _ := note.ToMarkdown()
			fmt.Print(md)
			return
		}

		if pretty {
			renderPretty(note)
		} else {
			renderBasic(note)
		}
	},
}

func init() {
	viewCmd.Flags().Bool("raw", false, "Output raw markdown")
	viewCmd.Flags().Bool("pretty", false, "Render pretty output")

	rootCmd.AddCommand(viewCmd)
}

func renderBasic(n *jot.Note) {
	fmt.Printf("# Note: %s\n", n.ID)
	fmt.Printf("Created: %s\n", n.CreatedAt.Format("2006-01-02 15:04"))
	if len(n.Tags) > 0 {
		fmt.Printf("Tags:    %s\n", strings.Join(n.Tags, ", "))
	}
	if len(n.Links) > 0 {
		fmt.Printf("Links:   %s\n", strings.Join(n.Links, ", "))
	}
	fmt.Println("\n" + n.Content)
}

func renderPretty(n *jot.Note) {
	// minimal ANSI-styled render
	fmt.Printf("\033[1m%s\033[0m\n", firstLine(n.Content))
	fmt.Printf("📅 %s\n", n.CreatedAt.Format("Jan 2 2006, 3:04PM"))
	if len(n.Tags) > 0 {
		fmt.Printf("🏷️  %s\n", strings.Join(n.Tags, ", "))
	}
	if len(n.Links) > 0 {
		fmt.Printf("🔗 %s\n", strings.Join(n.Links, ", "))
	}
	fmt.Println("\n" + n.Content)
}
