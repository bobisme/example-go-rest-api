package conf

// Config holds the configuration for the whole app
type Config struct {
	// DBPath is the path to the sqlite3 database file
	DBPath string `toml:"db_path"`
	// Port is the API Server HTTP port
	Port int `toml:"port"`
	// ReleaseMode is true if this is to be run in production
	ReleaseMode bool
}

// Default returns a configuration with default values
func Default() *Config {
	return &Config{
		Port:   8080,
		DBPath: "database.sqlite3",
	}
}
