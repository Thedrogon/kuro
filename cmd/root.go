package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/thedrogon/kuro/state"
)

// Minimalist styling for our base CLI text
var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")) // Bright White
	errorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))  // Red
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kuro",
	Short: "A brutalist language environment manager",
	Long: titleStyle.Render("KURO") + `
A lightning-fast, local meta-package manager for programming languages.
Tracks your environment, standardizes installations, and stays out of your way.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {

	defer state.ReleaseLock()
	state.AcquireLock()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error: %v", err)))
		return err
	}
	return nil
}