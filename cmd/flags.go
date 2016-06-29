package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/tucnak/climax"
)

// reuseable config flag
var configFlag = climax.Flag{
	Name:     "config",
	Short:    "c",
	Usage:    `--config="./config.toml"`,
	Help:     "Choose config file. Default is ./config.toml",
	Variable: true,
}

// get the config file from the config flag, or return "config.toml"
func getConfigPath(ctx climax.Context) string {
	if val, ok := ctx.Get("config"); ok && val != "" {
		return val
	}
	return "config.toml"
}

// debugFlag turns on debugging
var debugFlag = climax.Flag{
	Name:  "debug",
	Usage: `--debug`,
	Help:  "turn on debugging",
}

// handleDebugging checks for the debug flag and deals with it
func handleDebugging(ctx climax.Context) {
	if ctx.Is("debug") {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugln("debugging turned on")
	}
}
