 package resolver

import (
	"bytes"
	"os/exec"
	"strings"
)

// CheckLocal queries the local pacman database to see if a package is installed.
// It returns the installed version, the manager (pacman or aur helper), and a boolean if found.
func CheckLocal(pkgName string) (version string, manager string, installed bool) {
	// 1. pacman -Q <pkg> checks local installations.
	cmd := exec.Command("pacman", "-Q", pkgName)
	var out bytes.Buffer
	cmd.Stdout = &out

	// If it fails, the package is not on the system.
	if err := cmd.Run(); err != nil {
		return "", "", false
	}

	// Output format is: "package-name 1.2.3-1\n"
	parts := strings.Fields(out.String())
	if len(parts) >= 2 {
		version = parts[1]
	} else {
		version = "unknown"
	}

	// 2. Determine if it was installed via main repos or AUR.
	// pacman -Qm lists local packages that are "foreign" (from the AUR).
	manager = "pacman"
	cmdAur := exec.Command("pacman", "-Qm", pkgName)
	if err := cmdAur.Run(); err == nil {
		// It's from the AUR. Let's figure out what helper they use.
		if helper, err := getAURHelper(); err == nil {
			manager = helper
		} else {
			manager = "aur" // generic fallback
		}
	}

	return version, manager, true
}