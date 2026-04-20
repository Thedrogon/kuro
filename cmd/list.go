package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thedrogon/kuro/state"
	"github.com/charmbracelet/lipgloss"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"}, // Allows typing 'kuro ls'
	Short:   "List all installed programming languages",
	Run: func(cmd *cobra.Command, args []string) {
		reg, err := state.Load()
		if err != nil {
			fmt.Println("Error loading state:", err)
			os.Exit(1)
		}

		if len(reg.Languages) == 0 {
			fmt.Println("No languages currently managed by kuro.")
			return
		}

		// Brutalist table formatting
		headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6")).Underline(true)
		rowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
		managerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

		fmt.Printf("%-15s %-15s %-15s\n", headerStyle.Render("LANGUAGE"), headerStyle.Render("VERSION"), headerStyle.Render("MANAGER"))
		
		for name, lang := range reg.Languages {
			fmt.Printf("%-15s %-15s %-15s\n", 
				rowStyle.Render(name), 
				rowStyle.Render(lang.Version), 
				managerStyle.Render(lang.Manager),
			)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}