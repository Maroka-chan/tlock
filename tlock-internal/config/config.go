package config

import (
	"os"
	"path"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// Default theme
var DEFAULT_THEME = "Catppuccin"

// Path to the config file
var CONFIG_DIR = path.Join(xdg.ConfigHome, "tlock")
var CONFIG_PATH = path.Join(xdg.ConfigHome, "tlock", "tlock.yaml")

// Represents theme config
type Config struct {
	// Current theme
	// Defaults to `Catppuccin`
	CurrentTheme string `yaml:"current_theme"`

	/// Enable icons or not
	EnableIcon bool `yaml:"enable_icon"`
}

// Returns the default config
func DefaultConfig() Config {
	return Config{
		CurrentTheme: DEFAULT_THEME,
		EnableIcon:   false,
	}
}

// Loads the config from the file
func GetConfig() Config {
	default_config := DefaultConfig()

	// Read raw
	config_raw, err := os.ReadFile(CONFIG_PATH)

	// If error, just return the default config
	if err != nil {
		// Log
		log.Debug().Msg("[config] No config file found, returning the default config")

		return default_config
	}

	// Parse
	if err := yaml.Unmarshal(config_raw, &default_config); err != nil {
		// Log
		log.Error().Err(err).Msg("[config] Failed to parse config, syntax error possibly?")

		// Return default config
		return default_config
	}

	// Return
	return default_config
}

// Writes the config
func (config Config) Write() {
	// Make directory
	os.MkdirAll(filepath.Dir(CONFIG_PATH), os.ModePerm)

	// Marshal
	data, _ := yaml.Marshal(config)

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
