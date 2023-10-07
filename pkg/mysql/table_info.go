package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

type TableInfo struct {
	Conn  *sql.DB
	Cache *cache.Cache
}

func NewTableInfo(conn *sql.DB) *TableInfo {
	cache := cache.New(5*time.Minute, 10*time.Minute)
	return &TableInfo{Conn: conn, Cache: cache}
}

type Column struct {
	Name string
}

func (info *TableInfo) Get(schema string, name string) (*Table, error) {
	key := schema + "." + name
	t, found := info.Cache.Get(key)

	if found {
		return t.(*Table), nil
	}

	columns, err := info.getColumns(schema, name)
	if err != nil {
		return nil, err
	}
	if len(columns) == 0 {
		return nil, nil
	}
	tbl := &Table{
		Schema:      schema,
		Name:        name,
		Columns:     columns,
		ColumnCount: len(columns),
	}
	info.Cache.Set(key, tbl, cache.DefaultExpiration)
	return tbl, nil
}

func (info *TableInfo) getColumns(schema string, name string) ([]Column, error) {
	rows, err := info.Conn.Query(`
		SELECT
			COLUMN_NAME
		FROM
			information_schema.COLUMNS
		WHERE
			TABLE_SCHEMA = ?
			AND TABLE_NAME = ?
		ORDER BY
			ORDINAL_POSITION
`, schema, name)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	columns := []Column{}

	for rows.Next() {
		var colName string
		err = rows.Scan(&colName)

		if err != nil {
			return nil, err
		}

		columns = append(columns, Column{Name: colName})
	}

	if len(columns) == 0 {
		return nil, nil
	}

	return columns, nil
}

type Table struct {
	Schema      string
	Name        string
	Columns     []Column
	ColumnCount int
}

func (tbl *Table) Get() string {
	return fmt.Sprintf("%s.%s", tbl.Schema, tbl.Name)
}
