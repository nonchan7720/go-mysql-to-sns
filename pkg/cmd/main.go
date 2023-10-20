package cmd

import (
	"log/slog"
	"os"
)

func Execute() {
	exitCode := 0
	cmd := rootCommand()
	if err := cmd.Execute(); err != nil {
		slog.Error(err.Error())
		exitCode = 1
	}
	os.Exit(exitCode)
}
