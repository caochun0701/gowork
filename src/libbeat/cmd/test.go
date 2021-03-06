package cmd

import (
	"github.com/spf13/cobra"

	"libbeat/beat"
	"libbeat/cmd/test"
)

func genTestCmd(name, beatVersion string, beatCreator beat.Creator) *cobra.Command {
	exportCmd := &cobra.Command{
		Use:   "test",
		Short: "Test config",
	}

	exportCmd.AddCommand(test.GenTestConfigCmd(name, beatVersion, beatCreator))
	exportCmd.AddCommand(test.GenTestOutputCmd(name, beatVersion))

	return exportCmd
}
