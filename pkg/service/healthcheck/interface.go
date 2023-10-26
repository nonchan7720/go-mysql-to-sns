package healthcheck

import "context"

type IPing interface {
	PingContext(ctx context.Context) error
}
