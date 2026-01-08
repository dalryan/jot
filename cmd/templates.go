package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dalryan/jot/internal/jot"
	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Manage note templates",
}

var listTemplatesCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	Run: func(cmd *cobra.Command, args []string) {
		entries, err := os.ReadDir(cfg.TemplatesDir())
		if err != nil || len(entries) == 0 {
			fmt.Println("No templates found. Create one with 'jot templates new <name>'")
			return
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				name := strings.TrimSuffix(entry.Name(), ".md")
				fmt.Println(name)
			}
		}
	},
}

var newTemplateCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new template and open it in your editor",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		templatePath := filepath.Join(cfg.TemplatesDir(), name+".md")

		if _, err := os.Stat(templatePath); err == nil {
			fmt.Fprintf(os.Stderr, "Template '%s' already exists. Use 'jot templates edit %s'\n", name, name)
			os.Exit(1)
		}

		if err := os.MkdirAll(cfg.TemplatesDir(), 0755); err != nil {
			fmt.Fprintln(os.Stderr, "Error creating templates directory:", err)
			os.Exit(1)
		}

		if err := jot.RunEditor(cfg.Editor, templatePath); err != nil {
			fmt.Fprintln(os.Stderr, "Error running editor:", err)
			os.Exit(1)
		}
	},
}

var editTemplateCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "Edit an existing template in your editor",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		templatePath := filepath.Join(cfg.TemplatesDir(), name+".md")

		if _, err := os.Stat(templatePath); err != nil {
			fmt.Fprintf(os.Stderr, "Template '%s' does not exist. Use 'jot templates new %s' to create it.\n", name, name)
			os.Exit(1)
		}

		if err := jot.RunEditor(cfg.Editor, templatePath); err != nil {
			fmt.Fprintln(os.Stderr, "Error running editor:", err)
			os.Exit(1)
		}
	},
}

func init() {
	templatesCmd.AddCommand(listTemplatesCmd)
	templatesCmd.AddCommand(newTemplateCmd)
	templatesCmd.AddCommand(editTemplateCmd)
	rootCmd.AddCommand(templatesCmd)
}
