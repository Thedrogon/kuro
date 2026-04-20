package cmd

import (
	"os"
	"time"
	"github.com/pelletier/go-toml/v2"
)

// The overall state file structure
type Registry struct {
	LastUpdated time.Time             `toml:"last_updated"`
	Languages   map[string]Language   `toml:"languages"`
}

// What we actually track for each language
type Language struct {
	Version string `toml:"version"`
	Manager string `toml:"manager"` // e.g., "pacman", "paru", "binary"
	Path    string `toml:"path"`    // e.g., "/usr/bin/java"
}

// Function to load the state instantly
func LoadRegistry(path string) (*Registry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var reg Registry
	err = toml.Unmarshal(data, &reg)
	return &reg, err
}