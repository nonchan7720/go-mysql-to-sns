package mysql

import (
	"context"
	"database/sql"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/config"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/interfaces"
	"go.uber.org/multierr"
)

type Binlog struct {
	*config.Config
	pos  mysql.Position
	conn *sql.DB
}

func NewBinlog(ctx context.Context, config *config.Config) (*Binlog, error) {
	binlog := &Binlog{
		Config: config,
	}
	conn, err := binlog.Connect(ctx)
	if err != nil {
		return nil, err
	}
	binlog.conn = conn
	return binlog, nil
}

func (binlog *Binlog) Run(ctx context.Context, value chan interfaces.Payload) (err error) {
	var (
		file string
		pos  int
	)
	if f, p, loadErr := binlog.Config.Saver.Load(); loadErr == nil {
		file = f
		pos = p
	}
	if file == "" && pos == 0 {
		file, pos, err = binlog.loadBinlog(binlog.conn)
		if err != nil {
			return
		}
	}
	info := NewTableInfo(binlog.conn)
	serverId, err := binlog.findServerId(binlog.conn)
	if err != nil {
		return err
	}
	syncer, err := binlog.NewBinlogSyncer(serverId)
	if err != nil {
		return err
	}
	defer func() {
		pos := syncer.GetNextPosition()
		if saveErr := binlog.Config.Saver.Save(pos.Name, int(pos.Pos)); saveErr != nil {
			err = multierr.Append(err, saveErr)
		}
		syncer.Close()
	}()
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

func (binlog *Binlog) loadBinlog(conn *sql.DB) (file string, pos int, err error) {
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

func (binlog *Binlog) findServerId(conn *sql.DB) (serverId int, err error) {
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

func (binlog *Binlog) SavePosition(fn func(file string, position int)) {
	fn(binlog.pos.Name, int(binlog.pos.Pos))
}

func (binlog *Binlog) Close() {
	_ = binlog.conn.Close()
}

func (binlog *Binlog) PingContext(ctx context.Context) error {
	return binlog.conn.PingContext(ctx)
}

func binlogRowToPayloadRow(row []interface{}, tableColumns []Column) interfaces.Row {
	payloadRow := make(map[string]interface{}, len(row))
	for idx := 0; idx < len(row); idx++ {
		payloadRow[tableColumns[idx].Name] = row[idx]
	}
	return payloadRow
}
