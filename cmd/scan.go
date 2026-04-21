package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thedrogon/kuro/resolver"
	"github.com/thedrogon/kuro/state"
	"github.com/thedrogon/kuro/ui"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan the system for existing language installations and track them",
	Run: func(cmd *cobra.Command, args []string) {
		reg, err := state.Load()
		if err != nil {
			fmt.Println("Error loading state:", err)
			os.Exit(1)
		}

		ui.PrintStep("Scanning local system for unmanaged languages...")
		
		importedCount := 0

		// Iterate through our master list of supported languages
		for friendlyName, archName := range resolver.Aliases {
			
			// If Kuro already tracks it, skip it.
			if _, tracked := reg.Languages[friendlyName]; tracked {
				continue
			}

			// Ask the resolver if the system has this package
			version, manager, installed := resolver.CheckLocal(archName)
			
			if installed {
				// We found one! Add it to the registry.
				reg.Languages[friendlyName] = state.Language{
					Version: version,
					Manager: manager,
					Path:    "system", // Flagging it as pre-existing
				}
				
				ui.PrintSuccess(fmt.Sprintf("Discovered %s (v%s) via %s", friendlyName, version, manager))
				importedCount++
			}
		}

		if importedCount > 0 {
			if err := state.Save(reg); err != nil {
				fmt.Println("\nWarning: Scan successful but failed to save registry:", err)
			} else {
				fmt.Println() // Add a clean, empty line for visual breathing room
				msg := fmt.Sprintf("Successfully imported %d existing environments into Kuro.", importedCount)
				ui.PrintSuccess(msg)
			}
		} else {
			ui.PrintStep("No new unmanaged languages found. Your registry is perfectly synced.")
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}