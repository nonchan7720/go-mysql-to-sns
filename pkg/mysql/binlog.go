package mysql

import (
	"context"

	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/interfaces"
)

type Binlog struct {
	*config.Config
	streamer Streamer
}

func NewBinlog(ctx context.Context, config *config.Config) (*Binlog, error) {
	binlog := &Binlog{
		Config: config,
	}
	return binlog, nil
}

func (binlog *Binlog) Run(ctx context.Context, value chan interfaces.Payload) (err error) {
	conn, err := binlog.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	info := NewTableInfo(conn)
	streamer, err := NewStreamer(ctx, config.Position, binlog.Config)
	if err != nil {
		return err
	}
	binlog.streamer = streamer

	for {
		ev, err := streamer.GetEvent(ctx)
		if err != nil {
			if err != context.Canceled {
				return nil
			}
			return err
		}
		var (
			v     *interfaces.Payload
			evErr error
		)
		switch ev.Header.EventType {
		case replication.WRITE_ROWS_EVENTv2:
			v, evErr = binlog.handleWriteRowsEvent(ev.Event.(*replication.RowsEvent), info)
		case replication.UPDATE_ROWS_EVENTv2:
			v, evErr = binlog.handleUpdateRowsEvent(ev.Event.(*replication.RowsEvent), info)
		case replication.DELETE_ROWS_EVENTv2:
			v, evErr = binlog.handleDeleteRowsEvent(ev.Event.(*replication.RowsEvent), info)
		}
		if evErr != nil {
			return evErr
		}
		if v != nil {
			value <- *v
		}
	}
}

func (binlog *Binlog) handleWriteRowsEvent(e *replication.RowsEvent, info *TableInfo) (*interfaces.Payload, error) {
	rowNum := len(e.Rows)
	table, err := info.Get(string(e.Table.Schema), string(e.Table.Table))
	if err != nil {
		return nil, err
	}
	payload := interfaces.Payload{
		Event:  interfaces.Create,
		Schema: table.Schema,
		Table:  table.Name,
		Rows:   make([]interfaces.PayloadRow, rowNum),
	}
	for i := 0; i < rowNum; i++ {
		row := e.Rows[i]
		payload.Rows[i] = interfaces.NewPayloadRow(nil, binlogRowToPayloadRow(row, table.Columns))
	}
	return &payload, nil
}

func (binlog *Binlog) handleUpdateRowsEvent(e *replication.RowsEvent, info *TableInfo) (*interfaces.Payload, error) {
	rowNum := len(e.Rows)
	table, err := info.Get(string(e.Table.Schema), string(e.Table.Table))
	if err != nil {
		return nil, err
	}
	payload := interfaces.Payload{
		Event:  interfaces.Update,
		Schema: table.Schema,
		Table:  table.Name,
		Rows:   make([]interfaces.PayloadRow, 0, rowNum),
	}
	for i := 0; i < rowNum; i += 2 {
		oldRow := e.Rows[i]
		newRow := e.Rows[i+1]
		payload.Rows = append(payload.Rows,
			interfaces.NewPayloadRow(
				binlogRowToPayloadRow(oldRow, table.Columns),
				binlogRowToPayloadRow(newRow, table.Columns),
			),
		)
	}
	return &payload, nil
}

func (binlog *Binlog) handleDeleteRowsEvent(e *replication.RowsEvent, info *TableInfo) (*interfaces.Payload, error) {
	rowNum := len(e.Rows)
	table, err := info.Get(string(e.Table.Schema), string(e.Table.Table))
	if err != nil {
		return nil, err
	}
	payload := interfaces.Payload{
		Event:  interfaces.Delete,
		Schema: table.Schema,
		Table:  table.Name,
		Rows:   make([]interfaces.PayloadRow, rowNum),
	}
	for i := 0; i < rowNum; i++ {
		row := e.Rows[i]
		payload.Rows[i] = interfaces.NewPayloadRow(binlogRowToPayloadRow(row, table.Columns), nil)
	}
	return &payload, nil
}

func (binlog *Binlog) Save() error {
	return binlog.streamer.Save()
}

func binlogRowToPayloadRow(row []interface{}, tableColumns []Column) interfaces.Row {
	payloadRow := make(map[string]interface{}, len(row))
	for idx := 0; idx < len(row); idx++ {
		payloadRow[tableColumns[idx].Name] = row[idx]
	}
	return payloadRow
}
