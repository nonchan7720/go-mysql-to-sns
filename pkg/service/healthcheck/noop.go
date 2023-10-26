package healthcheck

import "context"

type noop struct {
	err error
}

var (
	_ IPing = (*noop)(nil)
)

func NewNoop() IPing {
	return &noop{}
}

func (n *noop) PingContext(_ context.Context) error {
	return n.err
}

func newNoopWithErr(err error) IPing {
	return &noop{
		err: err,
	}
}
