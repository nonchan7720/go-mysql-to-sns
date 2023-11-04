package config

import (
	"log/slog"
	"os"
	"strings"
)

type LoggingHandle string

const (
	JsonHandler = LoggingHandle("json")
	TextHandler = LoggingHandle("text")
)

type Logging struct {
	Level   slog.Level    `yaml:"level"`
	Handler LoggingHandle `yaml:"handler" default:"text"`
}

func (l *Logging) SetDefaults() {
	l.SetUpSlog()
}

func (l *Logging) SetUpSlog() {
	var h slog.Handler
	if strings.EqualFold(string(l.Handler), "json") {
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: l.Level,
		})
	} else {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: l.Level,
		})
	}
	log := slog.New(h)
	slog.SetDefault(log)
}
