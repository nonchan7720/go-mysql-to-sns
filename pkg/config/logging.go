package config

import (
	"log/slog"
	"os"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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

var (
	logHandler = validation.NewStringRuleWithError(
		func(value string) bool {
			v := strings.ToLower(value)
			return v == "json" || v == "text"
		},
		validation.NewError("validation_is_log_handle", "must be a value with json or text"),
	)
)

func (l Logging) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Handler, logHandler),
	)
}
