package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thedrogon/kuro/resolver"
	"github.com/thedrogon/kuro/runner"
	"github.com/thedrogon/kuro/state"
	"github.com/thedrogon/kuro/ui"
)

var removeCmd = &cobra.Command{
	Use:     "remove [language]",
	Aliases: []string{"rm", "uninstall"},
	Short:   "Remove a language environment",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkgName := args[0]

		reg, err := state.Load()
		if err != nil {
			fmt.Println("Error loading state:", err)
			os.Exit(1)
		}

		// Check if we actually track it
		lang, exists := reg.Languages[pkgName]
		if !exists {
			fmt.Printf("'%s' is not tracked by kuro.\n", pkgName)
			os.Exit(1)
		}

		ui.PrintStep(fmt.Sprintf("Preparing to remove '%s'...", pkgName))
		target, err := resolver.Resolve(pkgName)
		if err != nil {
			fmt.Println("Resolution failed:", err)
			os.Exit(1)
		}

		stream := &ui.LogStream{}
		ui.PrintStep("Streaming uninstallation logs...")
		
		// runner.Remove needs to be added to runner/exec.go
		err = runner.Remove(target.Manager, target.RealName, stream)
		if err != nil {
			fmt.Println("\nRemoval failed:", err)
			os.Exit(1)
		}

		// Purge from TOML
		delete(reg.Languages, pkgName)
		if err := state.Save(reg); err != nil {
			fmt.Println("Warning: Removed successfully but failed to update registry:", err)
		} else {
			ui.PrintSuccess(fmt.Sprintf("Successfully removed %s.", pkgName))
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}