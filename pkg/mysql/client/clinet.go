package client

import (
	"context"
	"database/sql"

	"github.com/nonchan7720/go-mysql-to-sns/pkg/ent"
)

type Client struct {
	*ent.Client
	db *sql.DB
}

func (c *Client) PingContext(ctx context.Context) error {
	return c.db.PingContext(ctx)
}
