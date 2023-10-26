package logging

import (
	"fmt"
	"log/slog"

	"github.com/pkg/errors"
)

func WithStack(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}
	return slog.String("stack", fmt.Sprintf("%+v", errors.WithStack(err)))
}
