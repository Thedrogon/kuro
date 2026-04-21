package state

import (
	"errors"
	"os"
	"path/filepath"
)

func getLockPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "state", "kuro", "kuro.lock"), nil
}

// AcquireLock attempts to create a lockfile. If it exists, another Kuro instance is running.
func AcquireLock() error {
	path, err := getLockPath()
	if err != nil {
		return err
	}

	// Ensure dir exists
	os.MkdirAll(filepath.Dir(path), 0755)

	// O_EXCL ensures this call fails if the file already exists. This is thread-safe.
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			return errors.New("kuro is already running in another process")
		}
		return err
	}
	file.Close()
	return nil
}

// ReleaseLock removes the lockfile when Kuro finishes.
func ReleaseLock() {
	if path, err := getLockPath(); err == nil {
		os.Remove(path)
	}
}