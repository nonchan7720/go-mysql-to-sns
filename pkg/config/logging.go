package config

import (
	"log/slog"
	"os"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/logging"
)

type Logging struct {
	Level   slog.Level             `yaml:"level"`
	Handler logging.LoggingHandle  `yaml:"handler" default:"text"`
	Sentry  *logging.SentryConfig  `yaml:"sentry"`
	Rollbar *logging.RollbarConfig `yaml:"rollbar"`

	closer func()
}

func (l *Logging) SetupLog() {
	var (
		h slog.Handler
	)
	switch strings.ToLower(string(l.Handler)) {
	case "sentry":
		h = logging.NewSentryHandler(l.Sentry, DefaultConfig().App.TrackingEnv)
		h = logging.NewHandler(setUpLog(string(l.Sentry.Backend.Handler), l.Sentry.Backend.Level), h)
	case "rollbar":
		l.Rollbar.Init(DefaultConfig().App.TrackingEnv, "v1", DefaultConfig().App.ServiceName)
		l.closer = l.Rollbar.Close
		h = logging.NewRollbarHandler(l.Rollbar)
		h = logging.NewHandler(setUpLog(string(l.Rollbar.Backend.Handler), l.Rollbar.Backend.Level), h)
	case "json":
		h = logging.NewJSONHandler(
			logging.WithWriter(os.Stdout),
			logging.WithLevel(l.Level),
		)
	default:
		h = logging.NewTextHandler(
			logging.WithWriter(os.Stdout),
			logging.WithLevel(l.Level),
		)
	}
	log := slog.New(h)
	slog.SetDefault(log)
}

func (l Logging) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Handler, validation.NotIn("json", "text", "sentry", "rollbar")),
		validation.Field(&l.Sentry, validation.When(l.Handler == "sentry", validation.NotNil)),
		validation.Field(&l.Rollbar, validation.When(l.Handler == "rollbar", validation.NotNil)),
	)
}

func (l Logging) Close() {
	if l.closer != nil {
		l.closer()
	}
}

func setUpLog(handlerLogName string, level slog.Level) slog.Handler {
	switch strings.ToLower(handlerLogName) {
	case "json":
		return logging.NewJSONHandler(
			logging.WithWriter(os.Stdout),
			logging.WithLevel(level),
		)
	default:
		return logging.NewTextHandler(
			logging.WithWriter(os.Stdout),
			logging.WithLevel(level),
		)
	}
}
