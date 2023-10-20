package cmd

import "github.com/spf13/cobra"

func rootCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:           "mysql-to-sns",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(runCommand())
	cmd.AddCommand(outboxCommand())

	return &cmd
}
