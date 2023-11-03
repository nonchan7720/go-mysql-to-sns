package client

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/config"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/ent"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
)

func NewDB(ctx context.Context, cfg *config.Config) (*Client, error) {
	dsn, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("rdb: failed to build connection string: %v\n", err)
	}
	db, err := otelsql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(timeoutCtx); err != nil {
		return nil, err
	}
	return newClient(ctx, "mysql", db), nil
}

func newClient(ctx context.Context, dialect string, db *sql.DB) *Client {
	drv := entsql.OpenDB(dialect, db)
	SetDB(db)
	opts := []ent.Option{
		ent.Driver(drv),
	}
	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		opts = append(opts, ent.Debug())
	}
	return &Client{
		Client: ent.NewClient(opts...),
		db:     db,
	}
}
