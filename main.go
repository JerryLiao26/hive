package main

import (
	"os"

	"github.com/JerryLiao26/hive/cli"
	"github.com/JerryLiao26/hive/config"
	"github.com/JerryLiao26/hive/helper"
)

func str2comm(str string) helper.Command {
	return helper.Command(str)
}

func checkSupport(comm helper.Command) int {
	for i := 0; i < len(helper.SupportedCommands); i++ {
		if helper.SupportedCommands[i] == comm {
			return i
		}
	}
	return 0
}

func main() {
	if len(os.Args) >= 2 {
		// Pre-load
		flag := config.LoadConf()

		// First time
		if !flag {
			cli.FirstHandler()
		}

		// Get command
		comm := str2comm(os.Args[1])
		// Handler for command
		cli.SupportedCommandHandlers[checkSupport(comm)]()
	} else {
		cli.HelpHandler()
	}
}
