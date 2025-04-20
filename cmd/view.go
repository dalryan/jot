package cmd

import (
	"bufio"
	"fmt"
	"github.com/dalryan/jot/internal/jot"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view [id]",
	Short: "View a note by its ID or from stdin",
	Run: func(cmd *cobra.Command, args []string) {
		var note *jot.Note
		var err error
		baseDir := cfg.StoragePath

		stat, err := os.Stdin.Stat()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to read stdin:", err)
			os.Exit(1)
		}

		if len(args) > 0 {
			id := args[0]
			note, err = jot.FindNoteByID(baseDir, id)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		} else if (stat.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				line := scanner.Text()
				if len(line) >= 8 {
					id := line[:8]
					note, err = jot.FindNoteByID(baseDir, id)
					if err != nil {
						fmt.Println("Error:", err)
						os.Exit(1)
					}
				} else {
					fmt.Println("Error: Input too short")
					os.Exit(1)
				}
			} else {
				fmt.Println("Error: No input provided")
				os.Exit(1)
			}
		} else {
			fmt.Println("Error: No ID provided")
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
