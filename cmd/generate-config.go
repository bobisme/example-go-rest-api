package cmd

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/bobisme/RestApiProject/conf"
	"github.com/tucnak/climax"
)

// GenerateConfig provide generate-config command to render the default config
// options to stdout
var GenerateConfig = climax.Command{
	Name:  "generate-config",
	Brief: "generate basic config",
	// Usage: cmdName + ` genconfig > config.toml`,
	Help: "This generates a toml configuration.\n",
	Examples: []climax.Example{
		{
			Usecase:     `> config.toml`,
			Description: `Writes defaults to config file: config.toml`,
		},
	},
	Handle: func(ctx climax.Context) int {
		defaults := conf.Default()
		encoder := toml.NewEncoder(os.Stdout)
		if err := encoder.Encode(defaults); err != nil {
			println("Could not encode defaults:", err)
			return 1
		}
		return 0
	},
}
