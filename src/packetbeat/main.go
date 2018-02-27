package main

import (
	"os"

	"packetbeat/cmd"
)

var Name = "packetbeat"

//启动packetbeat
func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
