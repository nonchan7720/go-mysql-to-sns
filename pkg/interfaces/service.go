package interfaces

import "context"

type Publisher interface {
	Publish(ctx context.Context, payload Payload) error
}
