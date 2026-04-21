package resolver

import (
	"errors"
	"os/exec"
	"regexp"
)

// Target holds the finalized, resolved instructions for the Execution layer.
type Target struct {
	RealName string // What the system actually calls it (e.g., "nodejs")
	Manager  string // "pacman", "paru", or "yay"
}

// IsValidName ensures no malicious characters are passed to the OS.
func IsValidName(name string) bool {
	// Arch packages only allow alphanumeric, hyphens, underscores, and dots.
	valid := regexp.MustCompile(`^[a-zA-Z0-9\-_\.]+$`)
	return valid.MatchString(name)
}

// Common aliases to make the Developer Experience (DX) flawless.
//not used anymore
var aliases = map[string]string{
	"node":   "nodejs",
	"golang": "go",
	"java":   "jre-openjdk",
	"java17": "jre17-openjdk",
	"java21": "jre21-openjdk",
}

// Resolve figures out exactly how to install the requested language.
func Resolve(pkgName string) (*Target, error) {
	// 1. Translate user intent using the alias map
	realName := pkgName
	if val, ok := Aliases[pkgName]; ok {
		realName = val
	}

	// 2. Check the official Arch repositories first (Fastest & Safest)
	if existsInPacman(realName) {
		return &Target{
			RealName: realName,
			Manager:  "pacman",
		}, nil
	}

	// 3. If not in main repos, hunt in the AUR using available helpers
	if helper, err := getAURHelper(); err == nil {
		if existsInAUR(helper, realName) {
			return &Target{
				RealName: realName,
				Manager:  helper,
			}, nil
		}
	}

	return nil, errors.New("package not found in official repositories or the AUR")
}

// existsInPacman safely queries the pacman sync database.
func existsInPacman(pkg string) bool {
	// SECURITY: We do not use bash -c. We pass arguments directly to the binary
	// to completely eliminate command injection vulnerabilities.
	// -Si prints info if it exists, exits with status 0.
	cmd := exec.Command("pacman", "-Si", pkg)
	err := cmd.Run()
	return err == nil
}

// existsInAUR safely queries the AUR via paru or yay.
func existsInAUR(helper, pkg string) bool {
	// -Si works identically in paru/yay to query the AUR without installing.
	cmd := exec.Command(helper, "-Si", pkg)
	err := cmd.Run()
	return err == nil
}

// getAURHelper dynamically detects what the CachyOS user has installed.
func getAURHelper() (string, error) {
	// Look for paru first (CachyOS default usually), then fallback to yay
	if _, err := exec.LookPath("paru"); err == nil {
		return "paru", nil
	}
	if _, err := exec.LookPath("yay"); err == nil {
		return "yay", nil
	}
	return "", errors.New("no AUR helper (paru/yay) found on system")
}