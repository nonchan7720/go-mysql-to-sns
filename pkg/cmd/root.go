package cmd

import "github.com/spf13/cobra"

func rootCommand() *cobra.Command {
	cmd := cobra.Command{
		Use: "mysql-to-sns",
	}

	cmd.AddCommand(runCommand())

	return &cmd
}
