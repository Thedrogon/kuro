package runner

import (
	"errors"
	"io"
	"os"
	"strconv"
	"os/exec"
)

// Install executes the package manager strictly and securely.
// It streams all text output to logStream so the UI can render it in real-time.
func Install(manager, pkgName string, logStream io.Writer) error {
	var cmd *exec.Cmd

	// Construct the exact command sequence. 
	// We strictly use individual arguments to completely prevent command injection.
	switch manager {
	case "pacman":
		// Official repos require sudo.
		cmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", pkgName)
	case "paru", "yay":
		// AUR helpers handle root internally. Running them with sudo breaks them.
		cmd = exec.Command(manager, "-S", "--noconfirm", pkgName)
	default:
		return errors.New("unsupported package manager: " + manager)
	}

	// PERFORMANCE: Direct io.Writer routing. 
	// We don't buffer logs in memory; we stream them directly to the UI component.
	cmd.Stdout = logStream
	cmd.Stderr = logStream

	// CRITICAL: We bind standard input directly to the terminal.
	// If sudo requires a password, the terminal will still accept the user's keystrokes.
	cmd.Stdin = os.Stdin

	// Execute and wait for completion.
	return cmd.Run()
}

// Remove executes the uninstallation process securely.
func Remove(manager, pkgName string, logStream io.Writer) error {
	var cmd *exec.Cmd

	switch manager {
	case "pacman":
		// -Rns removes the package, configuration files, and unneeded dependencies
		cmd = exec.Command("sudo", "pacman", "-Rns", "--noconfirm", pkgName)
	case "paru", "yay":
		cmd = exec.Command(manager, "-Rns", "--noconfirm", pkgName)
	default:
		return errors.New("unsupported package manager: " + manager)
	}

	cmd.Stdout = logStream
	cmd.Stderr = logStream
	cmd.Stdin = os.Stdin

	return cmd.Run()
}