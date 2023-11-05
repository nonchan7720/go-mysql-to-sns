package logging

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/rollbar/rollbar-go"
	slogrollbar "github.com/samber/slog-rollbar/v2"
)

type RollbarConfig struct {
	LogLevel string  `yaml:"logLevel" default:"warn"`
	Token    string  `yaml:"token"`
	Backend  Backend `yaml:"backend"`

	client *rollbar.Client
	Client *http.Client
}

func (c *RollbarConfig) Init(env, version, serverRoot string) {
	client := rollbar.NewAsync(c.Token, env, version, "", serverRoot)
	if c.Client != nil {
		client.SetHTTPClient(c.Client)
	}
	c.client = client
}

func (c *RollbarConfig) Close() {
	if c.client != nil {
		_ = c.client.Close()
	}
}

func (c *RollbarConfig) getLevel() slog.Level {
	switch strings.ToLower(c.LogLevel) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "err", "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func NewRollbarHandler(conf *RollbarConfig) Handle {
	option := slogrollbar.Option{
		Level:  conf.getLevel(),
		Client: conf.client,
	}
	return NewErrorTracking(NewAsyncHandler(option.NewRollbarHandler()))
}
