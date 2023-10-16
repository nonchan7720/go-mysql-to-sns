package mysql

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/interfaces"
)

type OutboxPattern struct {
	*config.Outbox
	streamer Streamer
}

func NewOutboxPattern(ctx context.Context, config *config.Outbox) (*OutboxPattern, error) {
	outbox := &OutboxPattern{
		Outbox: config,
	}
	return outbox, nil
}

func (outbox *OutboxPattern) Run(ctx context.Context, value chan interfaces.Outbox) (err error) {
	conn, err := outbox.Connect(ctx)
	if err != nil {
		return
	}
	defer conn.Close()
	info := NewTableInfo(conn)
	table, err := info.Get(outbox.Outbox.Schema, outbox.Outbox.TableName)
	if err != nil {
		return
	}
	streamer, err := NewStreamer(ctx, config.GTID, outbox.Config)
	if err != nil {
		return err
	}
	outbox.streamer = streamer

	for {
		ev, err := streamer.GetEvent(ctx)
		if err != nil {
			if err != context.Canceled {
				return nil
			}
			return err
		}
		var (
			values []*interfaces.Outbox
			evErr  error
		)
		switch ev.Header.EventType {
		case replication.WRITE_ROWS_EVENTv2:
			event := ev.Event.(*replication.RowsEvent)
			if strings.EqualFold(string(event.Table.Schema), outbox.Schema) && strings.EqualFold(string(event.Table.Table), outbox.TableName) {
				values, evErr = outbox.handleWriteRowsEvent(event, table.Columns)
			} else {
				continue
			}
		}
		if evErr != nil {
			return evErr
		}
		if len(values) > 0 {
			for _, v := range values {
				v := *v
				value <- v
			}
		}
	}
}

func (outbox *OutboxPattern) handleWriteRowsEvent(e *replication.RowsEvent, columns []Column) ([]*interfaces.Outbox, error) {
	values := make([]*interfaces.Outbox, 0, len(e.Rows))
	for _, row := range e.Rows {
		p, err := outboxTableToOutboxInterface(row, columns)
		if err != nil {
			return nil, err
		}
		values = append(values, p)
	}
	return values, nil
}

func (outbox *OutboxPattern) Close() {
	outbox.streamer.Close()
}

func (outbox *OutboxPattern) Save() error {
	return outbox.streamer.Save()
}

func outboxTableToOutboxInterface(row []interface{}, tableColumns []Column) (*interfaces.Outbox, error) {
	r := make(map[string]interface{}, len(row))
	for idx := 0; idx < len(row); idx++ {
		r[tableColumns[idx].Name] = row[idx]
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&r); err != nil {
		return nil, err
	}
	var outbox interfaces.Outbox
	if err := json.NewDecoder(&buf).Decode(&outbox); err != nil {
		return nil, err
	}
	return &outbox, nil
}
