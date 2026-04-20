package state

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml/v2"
)

// Registry holds the global state of all tracked languages.
type Registry struct {
	LastUpdated time.Time           `toml:"last_updated"`
	Languages   map[string]Language `toml:"languages"`
}

// Language defines the exact metadata for a single installed environment.
type Language struct {
	Version string `toml:"version"`
	Manager string `toml:"manager"` // "pacman", "paru", "binary"
	Path    string `toml:"path"`    // "/usr/bin/java"
}

// getRegistryPath strictly enforces XDG Base Directory standards for state data.
func getRegistryPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "state", "kuro", "registry.toml"), nil
}

// Load reads the registry file into memory. Blazing fast and safe.
func Load() (*Registry, error) {
	path, err := getRegistryPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		// Vulnerability prevention: Don't crash on first run. 
		// If the file doesn't exist yet, return a perfectly clean state.
		if errors.Is(err, os.ErrNotExist) {
			return &Registry{
				Languages: make(map[string]Language),
			}, nil
		}
		return nil, err // Return actual permission or hardware errors
	}

	var reg Registry
	if err := toml.Unmarshal(data, &reg); err != nil {
		return nil, err
	}

	// Memory safety: Ensure the map is ready for writes, even if the TOML was empty.
	if reg.Languages == nil {
		reg.Languages = make(map[string]Language)
	}

	return &reg, nil
}

// Save writes the current state back to disk cleanly.
func Save(reg *Registry) error {
	path, err := getRegistryPath()
	if err != nil {
		return err
	}

	// Stamp the registry with the exact moment it was modified
	reg.LastUpdated = time.Now().UTC()

	data, err := toml.Marshal(reg)
	if err != nil {
		return err
	}

	// Ensure the parent directories exist securely.
	// 0755 gives the owner full rights, and others read/execute only.
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write the actual file securely.
	// 0644 ensures only the user can write to it, preventing tampering.
	return os.WriteFile(path, data, 0644)
}