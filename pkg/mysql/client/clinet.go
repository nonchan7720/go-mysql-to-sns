package client

import (
	"context"
	"database/sql"

	"github.com/nonchan7720/go-storage-to-messenger/pkg/ent"
)

type Client struct {
	*ent.Client
	db *sql.DB
}

func (c *Client) PingContext(ctx context.Context) error {
	return c.db.PingContext(ctx)
}
