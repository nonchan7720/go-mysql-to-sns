package mysql

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/interfaces"
)

type OutboxPattern struct {
	*config.Outbox
	syncer *replication.BinlogSyncer
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
	var (
		file string
		pos  int
	)
	f, p, loadErr := outbox.Config.Saver.Load()
	if loadErr == nil {
		file = f
		pos = p
	}
	if file == "" && pos == 0 || loadErr != nil {
		file, pos, err = outbox.loadBinlog(conn)
		if err != nil {
			return
		}
	}
	serverId, err := outbox.findServerId(conn)
	if err != nil {
		return err
	}
	syncer, err := outbox.NewBinlogSyncer(serverId)
	if err != nil {
		return err
	}
	outbox.syncer = syncer

	streamer, err := syncer.StartSync(mysql.Position{Name: file, Pos: uint32(pos)})
	if err != nil {
		return err
	}

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

func (outbox *OutboxPattern) loadBinlog(conn *sql.DB) (file string, pos int, err error) {
	rows, err := conn.Query("show master status")

	if err != nil {
		return
	}

	defer rows.Close()
	columns, err := rows.Columns()

	if err != nil {
		return
	}

	colLen := len(columns)
	dest := make([]interface{}, colLen)
	dest[0] = &file
	dest[1] = &pos

	for i := 2; i < colLen; i++ {
		dest[i] = noopScanner{}
	}

	rows.Next()
	err = rows.Scan(dest...)

	if err != nil {
		return
	}
	return
}

func (outbox *OutboxPattern) findServerId(conn *sql.DB) (serverId int, err error) {
	rows, err := conn.Query(`SELECT @@server_id`)
	if err != nil {
		return
	}

	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return
	}

	colLen := len(columns)
	dest := make([]interface{}, colLen)
	dest[0] = &serverId

	for i := 1; i < colLen; i++ {
		dest[i] = noopScanner{}
	}

	rows.Next()
	err = rows.Scan(dest...)

	if err != nil {
		return
	}
	return
}

func (outbox *OutboxPattern) Close() {
	outbox.syncer.Close()
}

func (outbox *OutboxPattern) SavePosition() error {
	pos := outbox.syncer.GetNextPosition()
	return outbox.Config.Saver.Save(pos.Name, int(pos.Pos))
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