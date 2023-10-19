package client

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/ent"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
)

type Client struct {
	*ent.Client
	db *sql.DB
}

func NewDB(cfg *config.Config) (*Client, error) {
	dsn, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("rdb: failed to build connection string: %v\n", err)
	}
	db, err := otelsql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return newClient("mysql", db), nil
}

func newClient(dialect string, db *sql.DB) *Client {
	drv := entsql.OpenDB(dialect, db)
	SetDB(db)
	return &Client{
		Client: ent.NewClient(ent.Driver(drv), ent.Debug()),
		db:     db,
	}
}
