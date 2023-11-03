package service

import (
	"context"
	"log/slog"

	"github.com/nonchan7720/go-storage-to-messenger/pkg/interfaces"
)

type service struct {
	publisher interfaces.BackendPublisher
}

func New(publisher interfaces.BackendPublisher) interfaces.Publisher {
	return &service{
		publisher: publisher,
	}
}

func (svc *service) PublishBinlog(ctx context.Context, payload interfaces.Payload) error {
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
		msgId, err := svc.publisher.PublishBinlog(ctx, payload.Event, p)
		if err != nil {
			return err
		}
		slog.Info("Published.", "MessageId", msgId)
	}
	return nil
}

func (svc *service) PublishOutbox(ctx context.Context, producer string, outbox interfaces.Outbox) error {
	slog.With(
		slog.String("Producer", producer),
		slog.String("AggregateId", outbox.AggregateId),
	).InfoContext(ctx, "Receive outbox.")
	msgId, err := svc.publisher.PublishOutbox(ctx, producer, outbox)
	if err != nil {
		return err
	}
	slog.Info("Published outbox.", "MessageId", msgId)
	return nil
}
