package conf

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// LoadFile loads the configuration from the given filename
func LoadFile(path string) *Config {
	if path == "" {
		path = "config.toml"
	}
	path, err := filepath.Abs(path)
	if err != nil {
		panic("Unknown path: " + err.Error())
	}
	cfg := Default()
	// read toml from file
	_, err = toml.DecodeFile(path, &cfg)
	if err != nil {
		panic("Error reading config file: " + err.Error())
	}
	return cfg
}
