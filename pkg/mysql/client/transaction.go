package client

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/ent"
	"go.uber.org/multierr"
)

func RunInTransaction(ctx context.Context, db *ent.Client, execFunc func(ctx context.Context, tx *ent.Tx) error) (err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			slog.Error(fmt.Sprintf("recover from panic: %+v", pErr))
		}
	}()
	newTx, err := db.BeginTx(ctx, &entsql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			if txErr := newTx.Rollback(); txErr != nil {
				err = multierr.Append(err, txErr)
			}
		} else {
			if txErr := newTx.Commit(); txErr != nil {
				err = txErr
			}
		}
	}()

	return execFunc(ctx, newTx)
}

func MockRunInTransaction(err error) func(ctx context.Context, db *ent.Client, execFunc func(ctx context.Context, tx *ent.Tx) error) (err error) {
	return func(ctx context.Context, db *ent.Client, execFunc func(ctx context.Context, tx *ent.Tx) error) (err error) {
		return err
	}
}
