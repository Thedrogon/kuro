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

var updateCmd = &cobra.Command{
	Use:     "update [language]",
	Aliases: []string{"up"},
	Short:   "Update a specific tracked language environment",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkgName := args[0]

		if !resolver.IsValidName(pkgName) {
			fmt.Printf("Error: '%s' contains invalid characters for a package name.\n", pkgName)
			os.Exit(1)
		}

		reg, err := state.Load()
		if err != nil {
			fmt.Println("Error loading state:", err)
			os.Exit(1)
		}

		lang, tracked := reg.Languages[pkgName]
		if !tracked {
			fmt.Printf("'%s' is not managed by kuro. Run 'kuro install %s' first.\n", pkgName, pkgName)
			os.Exit(1)
		}

		target, err := resolver.Resolve(pkgName)
		if err != nil {
			fmt.Println("Resolution failed:", err)
			os.Exit(1)
		}

		ui.PrintStep(fmt.Sprintf("Checking upstream for '%s' updates...", pkgName))
		latest, err := resolver.GetLatestVersion(target.Manager, target.RealName)
		if err != nil {
			fmt.Println("Failed to check upstream version:", err)
			os.Exit(1)
		}

		if latest == lang.Version || latest == "unknown" {
			ui.PrintSuccess(fmt.Sprintf("'%s' is already up-to-date (v%s).", pkgName, lang.Version))
			return
		}

		ui.PrintStep(fmt.Sprintf("Updating %s: %s -> %s via %s...", pkgName, lang.Version, latest, target.Manager))
		stream := &ui.LogStream{}
		
		err = runner.Install(target.Manager, target.RealName, stream)
		if err != nil {
			fmt.Println("\nUpdate failed:", err)
			os.Exit(1)
		}

		// Update the local registry
		lang.Version = latest
		reg.Languages[pkgName] = lang
		if err := state.Save(reg); err != nil {
			fmt.Println("Warning: Update succeeded but failed to save registry:", err)
		} else {
			ui.PrintSuccess(fmt.Sprintf("Successfully updated %s to v%s.", pkgName, latest))
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}