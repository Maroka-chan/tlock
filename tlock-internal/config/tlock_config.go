package config

import (
	"os"
	"path"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/kelindar/binary"
	"github.com/rs/zerolog/log"
)

// Default theme
var DEFAULT_THEME = "Catppuccin"

// Path to the config file
var CONFIG_DIR = path.Join(xdg.ConfigHome, "tlock")
var CONFIG_PATH = path.Join(xdg.ConfigHome, "tlock", "conf.bin")

// TLock config is the config which is overriden by tlock itself
type TLockConfig struct {
	// Current theme
	// Defaults to `Catppuccin`
	CurrentTheme string `yaml:"current_theme"`
}

// Returns the default config
func DefaultTLockConfig() TLockConfig {
	return TLockConfig{
		CurrentTheme: DEFAULT_THEME,
	}
}

// Loads the config from the file
func GetTLockConfig() TLockConfig {
	default_config := DefaultTLockConfig()

	// Read raw
	config_raw, err := os.ReadFile(CONFIG_PATH)

	// If error, just return the default config
	if err != nil {
		// Log
		log.Debug().Msg("[config] No config file found, returning the default config")

		return default_config
	}

	// Parse
	if err := binary.Unmarshal(config_raw, &default_config); err != nil {
		// Log
		log.Error().Err(err).Msg("[config] Failed to parse config, syntax error possibly?")

		// Return default config
		return default_config
	}

	// Return
	return default_config
}

// Writes the config
func (config TLockConfig) Write() {
	// Make directory
	os.MkdirAll(filepath.Dir(CONFIG_PATH), os.ModePerm)

	// Marshal
	data, _ := binary.Marshal(config)

	// Write
	file, err := os.Create(CONFIG_PATH)

	// If no error, write to file
	if err == nil {
		file.Write(data)
	} else {
		// Log
		log.Error().Err(err).Msg("[config] Failed to write to config")
	}
}
