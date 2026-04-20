package resolver

import (
	"bytes"
	"os/exec"
	"strings"
)

// GetLatestVersion safely asks the package manager for the current upstream version.
func GetLatestVersion(manager, pkg string) (string, error) {
	cmd := exec.Command(manager, "-Si", pkg)
	
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Parse the output looking for the "Version" line
	for _, line := range strings.Split(out.String(), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "Version") {
			// e.g., "Version         : 1.22.0-1" -> "1.22.0-1"
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	return "unknown", nil
}