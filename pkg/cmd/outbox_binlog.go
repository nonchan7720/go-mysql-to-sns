package cmd

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/nonchan7720/go-storage-to-messenger/pkg/config"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/interfaces"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/mysql"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/service"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/service/healthcheck"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type outboxBinlogArgs struct {
	mainArgs
}

func outboxBinlogCommand() *cobra.Command {
	var (
		args outboxBinlogArgs
	)

	cmd := cobra.Command{
		Use: "binlog",
		Run: func(cmd *cobra.Command, _ []string) {
			ctx := cmd.Context()
			executeOutboxBinlog(ctx, &args)
		},
	}
	flag := cmd.Flags()
	args.setpflag(flag)

	return &cmd
}

func executeOutboxBinlog(ctx context.Context, args *outboxBinlogArgs) {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	config, err := config.LoadOutboxConfig(args.configFilePath)
	if err != nil {
		panic(err)
	}
	defer config.Logging.Close()

	if err := executeValidation(&config.Config); err != nil {
		panic(err)
	}

	publisher, err := getPublisher(ctx, config.Publisher)
	if err != nil {
		panic(err)
	}

	outboxBinlog, err := mysql.NewOutboxPattern(ctx, config)
	if err != nil {
		panic(err)
	}
	defer outboxBinlog.Close()
	payload := make(chan interfaces.BinlogOutbox)
	savePoint := make(chan struct{})
	eg, ctx := errgroup.WithContext(ctx)
	var healthCheckServer service.HealthCheck = healthcheck.New(outboxBinlog, args.healthCheckAddr)
	eg.Go(func() error {
		return healthCheckServer.Start(ctx)
	})
	eg.Go(func() error {
		defer close(payload)
		return outboxBinlog.Run(ctx, payload)
	})
	eg.Go(func() error {
		defer close(savePoint)
		for p := range payload {
			if err := publisher.PublishOutbox(ctx, p.Producer, p.Outbox); err != nil {
				return err
			}
			savePoint <- struct{}{}
		}
		return nil
	})
	eg.Go(func() error {
		for range savePoint {
			if err := outboxBinlog.SavePosition(); err != nil {
				return err
			}
		}
		return nil
	})
	if err := eg.Wait(); err != nil && err != context.Canceled {
		slog.Error(err.Error())
	}
	stop()
	slog.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	<-ctx.Done()
}
