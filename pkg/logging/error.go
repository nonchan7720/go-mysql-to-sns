package logging

import (
	"fmt"
	"log/slog"

	"github.com/pkg/errors"
)

func WithErr(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}
	return slog.String("err", fmt.Sprintf("%+v", errors.WithStack(err)))
}
