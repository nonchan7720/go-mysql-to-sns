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

func runCommand() *cobra.Command {
	var (
		configFilePath string
	)

	cmd := cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			execute(ctx, configFilePath)
		},
	}
	flag := cmd.Flags()
	flag.StringVarP(&configFilePath, "config", "c", "config.yaml", "configuration file path")

	return &cmd
}

func execute(ctx context.Context, configFilePath string) {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	config, err := config.LoadConfig(configFilePath)
	if err != nil {
		panic(err)
	}
	binlog, err := mysql.NewBinlog(ctx, config)
	if err != nil {
		panic(err)
	}
	payload := make(chan interfaces.Payload)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		defer close(payload)
		return binlog.Run(ctx, payload)
	})
	eg.Go(func() error {
		for p := range payload {
			json, _ := p.ToJson()
			slog.Info(json)
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		panic(err)
	}
	stop()
	slog.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

}
