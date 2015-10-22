package main

import (
	"github.com/FoxComm/libs/logger"
	"github.com/FoxComm/FoxComm/router/registry"
	"github.com/FoxComm/vulcand/vctl/command"
	"os"
)

func main() {
	r, err := registry.GetRegistry()
	if err != nil {
		logger.Error("Error: %s\n", err)
		return
	}
	cmd := command.NewCommand(r)
	if err := cmd.Run(os.Args); err != nil {
		logger.Error("Error: %s\n", err)
	}
}
