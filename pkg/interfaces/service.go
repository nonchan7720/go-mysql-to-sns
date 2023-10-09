package interfaces

import "context"

type Publisher interface {
	Publish(ctx context.Context, payload Payload) error
}

type BackendPublisher interface {
	Publish(ctx context.Context, event Event, payload SendPayload) (string, error)
	IsTarget(ctx context.Context, payload SendPayload) bool
}
