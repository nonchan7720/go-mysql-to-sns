package cmd

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/interfaces"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/mysql"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func outboxRunCommand() *cobra.Command {
	var (
		configFilePath string
	)

	cmd := cobra.Command{
		Use: "outbox",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			executeOutbox(ctx, configFilePath)
		},
	}
	flag := cmd.Flags()
	flag.StringVarP(&configFilePath, "config", "c", "config.yaml", "configuration file path")

	return &cmd
}

func executeOutbox(ctx context.Context, configFilePath string) {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	config, err := config.LoadOutboxConfig(configFilePath)
	if err != nil {
		panic(err)
	}

	publisher, err := getPublisher(ctx, config.Publisher)
	if err != nil {
		panic(err)
	}

	binlog, err := mysql.NewOutboxPattern(ctx, config)
	if err != nil {
		panic(err)
	}
	defer binlog.Close()
	payload := make(chan interfaces.Outbox)
	savePoint := make(chan struct{})
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		defer close(payload)
		return binlog.Run(ctx, payload)
	})
	eg.Go(func() error {
		defer close(savePoint)
		for p := range payload {
			if err := publisher.PublishOutbox(ctx, p); err != nil {
				return err
			}
		}
		return nil
	})
	eg.Go(func() error {
		for range savePoint {
			if err := binlog.SavePosition(); err != nil {
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
}
