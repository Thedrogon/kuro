package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thedrogon/kuro/resolver"
	"github.com/thedrogon/kuro/runner"
	"github.com/thedrogon/kuro/state"
	"github.com/thedrogon/kuro/ui"
	"github.com/charmbracelet/lipgloss"
)

var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"outdated"},
	Short:   "Check for updates and optionally install them",
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

		ui.PrintStep("Pinging Arch repositories for updates...")

		// Styles
		highlight := lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
		dim := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
		warning := lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // Yellow

		// Track what actually needs updating
		type updateItem struct {
			Name    string
			Current string
			Latest  string
			Target  *resolver.Target
		}
		var updates []updateItem

		for name, lang := range reg.Languages {
			target, err := resolver.Resolve(name)
			if err != nil {
				continue // Skip if resolution fails
			}

			latest, err := resolver.GetLatestVersion(target.Manager, target.RealName)
			if err != nil || latest == "unknown" {
				continue
			}

			// If the version string doesn't match, we assume an update is available
			if lang.Version != latest {
				updates = append(updates, updateItem{
					Name:    name,
					Current: lang.Version,
					Latest:  latest,
					Target:  target,
				})
			}
		}

		if len(updates) == 0 {
			ui.PrintSuccess("All language environments are up to date.")
			return
		}

		// Print the updates list
		fmt.Printf("\n%s\n", warning.Render("Updates Available:"))
		for i, u := range updates {
			fmt.Printf("[%d] %-10s %s -> %s\n", i+1, highlight.Render(u.Name), dim.Render(u.Current), highlight.Render(u.Latest))
		}

		// THE PROMPT LOGIC
		fmt.Printf("\n%s\n> ", dim.Render("Enter numbers to update (e.g., '1,3'), '0' for all, or press Enter to cancel."))
		
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Edge Case Handled: User just pressed Enter
		if input == "" {
			fmt.Println("Update cancelled. Exiting.")
			os.Exit(0)
		}

		// TODO: Add the loop here to parse the numbers (e.g., "1,2" or "0") 
		// and run runner.Install() for each selected package, then update state.Save()
		fmt.Printf("Proceeding to update selection: %s\n", input)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}