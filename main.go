package main

import (
	"github.com/bobisme/RestApiProject/cmd"
	"github.com/tucnak/climax"
)

func main() {
	cli := climax.New("rest-api")
	cli.AddCommand(cmd.GenerateConfig)
	cli.AddCommand(cmd.InitDB)
	cli.Run()
}
