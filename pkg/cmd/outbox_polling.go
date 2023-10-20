package cmd

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/mysql"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/mysql/client"
	"github.com/spf13/cobra"
)

func outboxPollingCommand() *cobra.Command {
	var (
		configFilePath string
	)

	cmd := cobra.Command{
		Use: "polling",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			executePolling(ctx, configFilePath)
		},
	}
	flag := cmd.Flags()
	flag.StringVarP(&configFilePath, "config", "c", "config.yaml", "configuration file path")

	return &cmd
}

func executePolling(ctx context.Context, configFilePath string) {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	config, err := config.LoadOutboxPollingConfig(configFilePath)
	if err != nil {
		panic(err)
	}

	publisher, err := getPublisher(ctx, config.Publisher)
	if err != nil {
		panic(err)
	}

	poller, err := mysql.NewOutboxPolling(ctx, config, publisher, client.RunInTransaction)
	if err != nil {
		panic(err)
	}

	if err := poller.Start(ctx); err != nil && err != context.Canceled {
		slog.Error(err.Error())
	}
	stop()
	slog.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	<-ctx.Done()
}
