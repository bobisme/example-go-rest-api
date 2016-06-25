package conf

// Config holds the configuration for the whole app
type Config struct {
	// DBPath is the path to the sqlite3 database file
	DBPath string `toml:"db_path"`
}

// Default returns a configuration with default values
func Default() *Config {
	return &Config{
		DBPath: "database.sqlite3",
	}
}
