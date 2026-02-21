package config

import (
	"os"
	"path/filepath"
	"runtime"
)

func ConfigDir() string {
	if v := os.Getenv("NODEVEIL_CONFIG_DIR"); v != "" {
		return v
	}
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "windows":
		if appdata := os.Getenv("APPDATA"); appdata != "" {
			return filepath.Join(appdata, "Nodeveil")
		}
		return filepath.Join(home, "AppData", "Roaming", "Nodeveil")
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Nodeveil")
	default:
		return filepath.Join(home, ".config", "nodeveil")
	}
}

func DefaultDBPath() string { return filepath.Join(ConfigDir(), "nodeveil.db") }
