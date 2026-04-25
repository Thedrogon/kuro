package runner

import (
	"errors"
	"io"
	"os"
	"os/exec"
)

// Install executes the package manager strictly and securely.
func Install(manager, pkgName string, logStream io.Writer) error {
	// PHASE 1: DOWNLOAD ONLY (Atomicity)
	// -Sw tells pacman/paru to sync and download the tarball, but NOT install it.
	var dlCmd *exec.Cmd
	switch manager {
	case "pacman":
		dlCmd = exec.Command("sudo", "pacman", "-Sw", "--noconfirm", pkgName)
	case "paru", "yay":
		dlCmd = exec.Command(manager, "-Sw", "--noconfirm", pkgName)
	default:
		return errors.New("unsupported package manager: " + manager)
	}

	dlCmd.Stdout = logStream
	dlCmd.Stderr = logStream
	dlCmd.Stdin = os.Stdin

	// If the download fails (network drop, bad signature), we halt. The system is untouched.
	if err := dlCmd.Run(); err != nil {
		return errors.New("phase 1 (download) failed: " + err.Error())
	}

	// PHASE 2: INSTALL FROM CACHE
	var instCmd *exec.Cmd
	switch manager {
	case "pacman":
		instCmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", pkgName)
	case "paru", "yay":
		instCmd = exec.Command(manager, "-S", "--noconfirm", pkgName)
	}

	instCmd.Stdout = logStream
	instCmd.Stderr = logStream
	instCmd.Stdin = os.Stdin

	// Since the package is already downloaded and verified, this executes entirely offline.
	if err := instCmd.Run(); err != nil {
		return errors.New("phase 2 (install) failed: " + err.Error())
	}

	return nil
}