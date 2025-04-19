package cmd

import (
	"fmt"
	"github.com/dalryan/jot/internal/jot"
	"os"
	"path/filepath"
	"strings"

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
		cfg, err := jot.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}

		templateDir := filepath.Join(cfg.StoragePath, "templates")
		entries, err := os.ReadDir(templateDir)
		if err != nil || len(entries) == 0 {
			fmt.Println("No templates found. Create one in ~/.jot/templates/")
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
		cfg, err := jot.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}

		templatePath := filepath.Join(cfg.StoragePath, "templates", name+".md")

		if _, err := os.Stat(templatePath); err == nil {
			fmt.Printf("Template '%s' already exists. Use 'jot templates edit %s'\n", name, name)
			return
		}

		templateDir := filepath.Join(cfg.StoragePath, "templates")
		err = os.MkdirAll(templateDir, 0755)
		if err != nil {
			fmt.Println("Failed to create templates directory:", err)
			return
		}

		templatePath = filepath.Join(templateDir, name+".md")

		editor := cfg.Editor

		err = jot.RunEditor(editor, templatePath)
		if err != nil {
			fmt.Println("Error running editor:", err)
			return
		}
	},
}

var editTemplateCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "Edit an existing template in your editor",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		cfg, err := jot.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}

		templatePath := filepath.Join(cfg.StoragePath, "templates", name+".md")

		if _, err := os.Stat(templatePath); err != nil {
			fmt.Printf("Template '%s' does not exist. Use 'jot templates new %s' to create it.\n", name, name)
			return
		}

		editor := cfg.Editor

		err = jot.RunEditor(editor, templatePath)
		if err != nil {
			fmt.Println("Error editing template:", err)
			return
		}
	},
}

func init() {
	templatesCmd.AddCommand(listTemplatesCmd)
	templatesCmd.AddCommand(newTemplateCmd)
	templatesCmd.AddCommand(editTemplateCmd)
	rootCmd.AddCommand(templatesCmd)
}
