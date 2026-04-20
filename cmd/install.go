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

var installCmd = &cobra.Command{
	Use:   "install [language]",
	Short: "Install a new programming language environment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkgName := args[0]
		ui.PrintStep(fmt.Sprintf("Checking environment for '%s'...", pkgName))

		// 1. Load Local State
		reg, err := state.Load()
		if err != nil {
			fmt.Println("Error loading state:", err)
			os.Exit(1)
		}

		if _, exists := reg.Languages[pkgName]; exists {
			ui.PrintStep(fmt.Sprintf("'%s' is already installed and tracked by kuro.", pkgName))
			os.Exit(0)
		}

		// 2. Resolve the Package
		ui.PrintStep(fmt.Sprintf("Resolving '%s' upstream...", pkgName))
		target, err := resolver.Resolve(pkgName)
		if err != nil {
			fmt.Println("Resolution failed:", err)
			os.Exit(1)
		}
		ui.PrintStep(fmt.Sprintf("Target found: %s (via %s)", target.RealName, target.Manager))

		// 3. Execute with UI Streaming
		ui.PrintStep("Beginning installation stream...")
		
		stream := &ui.LogStream{}
		err = runner.Install(target.Manager, target.RealName, stream)
		if err != nil {
			fmt.Println("\nInstallation failed:", err)
			os.Exit(1)
		}

		// 4. Save to Registry on Success
		reg.Languages[pkgName] = state.Language{
			Version: "latest", // We will pull actual versions in a future pass
			Manager: target.Manager,
			Path:    "managed", 
		}
		
		if err := state.Save(reg); err != nil {
			fmt.Println("Warning: Installed successfully but failed to update registry:", err)
		} else {
			ui.PrintSuccess(fmt.Sprintf("Successfully installed and tracked %s.", pkgName))
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}