package service

import (
	"context"
	"log/slog"

	"github.com/nonchan7720/go-mysql-to-sns/pkg/interfaces"
)

type service struct {
	publisher interfaces.BackendPublisher
}

func New(publisher interfaces.BackendPublisher) interfaces.Publisher {
	return &service{
		publisher: publisher,
	}
}

func (svc *service) Publish(ctx context.Context, payload interfaces.Payload) error {
	slog.With(slog.String("Table", payload.Table)).InfoContext(ctx, "Receive payload.")
	for idx := range payload.Rows {
		p := payload.SendPayload(idx)
		if !svc.publisher.IsTarget(ctx, p) {
			break
		}
		slog.With(
			slog.String("Table", payload.Table),
			slog.String("Event", payload.Event.String()),
			slog.Int("Row", idx+1),
		).InfoContext(ctx, "Publish.")
		msgId, err := svc.publisher.Publish(ctx, payload.Event, p)
		if err != nil {
			return err
		}
		slog.Info("Published.", "MessageId", msgId)
	}
	return nil
}
