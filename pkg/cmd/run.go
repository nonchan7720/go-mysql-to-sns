package cmd

import (
	"context"
	"errors"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/backend/aws"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/config"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/interfaces"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/mysql"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/service"
	backend "github.com/nonchan7720/go-storage-to-messenger/pkg/service/aws"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/service/healthcheck"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type runArgs struct {
	mainArgs
}

func runCommand() *cobra.Command {
	var (
		args runArgs
	)

	cmd := cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, _ []string) {
			ctx := cmd.Context()
			execute(ctx, &args)
		},
	}
	flag := cmd.Flags()
	args.setpflag(flag)

	return &cmd
}

func execute(ctx context.Context, args *runArgs) {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	config, err := config.LoadConfig(args.configFilePath)
	if err != nil {
		panic(err)
	}

	if err := executeValidation(config); err != nil {
		panic(err)
	}

	publisher, err := getPublisher(ctx, config.Publisher)
	if err != nil {
		panic(err)
	}

	binlog, err := mysql.NewBinlog(ctx, config)
	if err != nil {
		panic(err)
	}
	defer binlog.Close()
	var healthCheckServer service.HealthCheck = healthcheck.New(binlog, args.healthCheckAddr)
	payload := make(chan interfaces.Payload)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return healthCheckServer.Start(ctx)
	})
	eg.Go(func() error {
		defer close(payload)
		return binlog.Run(ctx, payload)
	})
	eg.Go(func() error {
		for p := range payload {
			if err := publisher.PublishBinlog(ctx, p); err != nil {
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

func getPublisher(ctx context.Context, conf *config.Publisher) (interfaces.Publisher, error) {
	noSelectedErr := errors.New("Set the Publisher.")
	var publisher interfaces.BackendPublisher
	if conf == nil {
		return nil, noSelectedErr
	}
	if conf.IsAWS() {
		switch {
		case conf.AWS.IsSNS():
			if client, err := aws.NewSNSClient(ctx, conf.AWS); err != nil {
				return nil, err
			} else {
				publisher = backend.NewAWSSNS(ctx, client, conf.AWS)
			}
		case conf.AWS.IsSQS():
			if client, err := aws.NewSQSClient(ctx, conf.AWS); err != nil {
				return nil, err
			} else {
				publisher = backend.NewAWSSQS(ctx, client, conf.AWS)
			}
		}
	}

	if publisher == nil {
		return nil, noSelectedErr
	}
	return service.New(publisher), nil
}

func executeValidation(conf *config.Config) error {
	return config.ValidateStruct(conf,
		validation.Field(&conf.Saver, validation.NotNil),
	)
}
