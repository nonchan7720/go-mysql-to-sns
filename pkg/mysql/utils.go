package mysql

import "database/sql"

func findServerId(conn *sql.DB) (serverId int, err error) {
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

func loadBinlogPosition(conn *sql.DB) (file string, pos int, err error) {
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

func loadBinlogGTID(conn *sql.DB) (gtid string, err error) {
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
	dest[4] = &gtid

	for i := 0; i < colLen-1; i++ {
		dest[i] = noopScanner{}
	}

	rows.Next()
	err = rows.Scan(dest...)

	if err != nil {
		return
	}
	return
}

func loadGlobalBinlogGTID(conn *sql.DB) (gtid string, err error) {
	rows, err := conn.Query("select @@GLOBAL.gtid_executed")

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
	dest[0] = &gtid

	for i := 0; i < colLen-1; i++ {
		dest[i] = noopScanner{}
	}

	rows.Next()
	err = rows.Scan(dest...)

	if err != nil {
		return
	}
	return
}
