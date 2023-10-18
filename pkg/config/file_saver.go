package config

import (
	"errors"
	"os"
)

var (
	ErrUnsupportedBinlogFormat = errors.New("Unsupported type.")
)

type FileSaver struct {
	Name string `yaml:"name"`
}

func (s *FileSaver) Save(typ BinlogSaveFormatType, format BinlogSaveFormat) error {
	f, err := os.Create(s.Name)
	if err != nil {
		return err
	}
	defer f.Close()
	switch typ {
	case Position:
		return positionSaver{}.save(f, format)
	case GTID:
		return gtidSaver{}.save(f, format)
	default:
		return ErrUnsupportedBinlogFormat
	}
}

func (s *FileSaver) Load(typ BinlogSaveFormatType) (format *BinlogSaveFormat, err error) {
	f, err := os.Open(s.Name)
	if err != nil {
		return
	}
	defer f.Close()
	switch typ {
	case Position:
		return positionSaver{}.load(f)
	case GTID:
		return gtidSaver{}.load(f)
	default:
		return nil, ErrUnsupportedBinlogFormat
	}
}
