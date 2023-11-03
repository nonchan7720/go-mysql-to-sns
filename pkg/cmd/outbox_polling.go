package cmd

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/nonchan7720/go-storage-to-messenger/pkg/config"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/mysql"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/mysql/client"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/service"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/service/healthcheck"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type outboxPollingArgs struct {
	mainArgs
}

func outboxPollingCommand() *cobra.Command {
	var (
		args outboxPollingArgs
	)

	cmd := cobra.Command{
		Use: "polling",
		Run: func(cmd *cobra.Command, _ []string) {
			ctx := cmd.Context()
			executePolling(ctx, &args)
		},
	}
	flag := cmd.Flags()
	args.setpflag(flag)
	return &cmd
}

func executePolling(parent context.Context, args *outboxPollingArgs) {
	ctx, stop := signal.NotifyContext(parent, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	config, err := config.LoadOutboxPollingConfig(args.configFilePath)
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
	var healthCheckServer service.HealthCheck = healthcheck.New(poller, args.healthCheckAddr)
	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return healthCheckServer.Start(ctx)
	})
	group.Go(func() error {
		return poller.Start(ctx)
	})

	if err := group.Wait(); err != nil && err != context.Canceled {
		slog.Error(err.Error())
	}
	stop()
	slog.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	<-ctx.Done()
}
