package main

import (
	"os"

	"libbeat/cmd"
	"libbeat/mock"
)

var RootCmd = cmd.GenRootCmd(mock.Name, mock.Version, mock.New)

func main() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
