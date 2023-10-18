package config

import (
	"io"
)

type gtidSaver struct{}

func (p gtidSaver) save(f io.Writer, format BinlogSaveFormat) error {
	_, err := f.Write(format.GTID)
	return err
}

func (p gtidSaver) load(f io.Reader) (format *BinlogSaveFormat, err error) {
	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	format = &BinlogSaveFormat{
		GTID: buf,
	}
	return format, err
}
