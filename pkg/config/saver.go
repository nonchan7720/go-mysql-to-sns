package config

import (
	"errors"
)

type BinlogSaveFormat struct {
	File     string
	Position int
	GTID     []byte
}

type Saver interface {
	Save(format BinlogSaveFormat) error
	Load() (format *BinlogSaveFormat, err error)
}

type BinlogSaveFormatType string

const (
	GTID     = BinlogSaveFormatType("gtid")
	Position = BinlogSaveFormatType("position")
)

type BinlogSaver struct {
	File *FileSaver `yaml:"file"`
}

func (s *BinlogSaver) Save(typ BinlogSaveFormatType, format BinlogSaveFormat) error {
	if s.File != nil {
		return s.File.Save(typ, format)
	}
	return nil
}

var (
	ErrNotSelected = errors.New("Not selected.")
)

func (s *BinlogSaver) Load(typ BinlogSaveFormatType) (*BinlogSaveFormat, error) {
	if s.File != nil {
		return s.File.Load(typ)
	}
	return nil, ErrNotSelected
}
