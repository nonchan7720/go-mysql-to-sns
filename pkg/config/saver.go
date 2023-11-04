package config

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Saver interface {
	Save(file string, position int) error
	Load() (file string, position int, err error)
}

type BinlogSaver struct {
	File *FileSaver `yaml:"file"`
}

func (s *BinlogSaver) Save(file string, position int) error {
	if s.File != nil {
		return s.File.Save(file, position)
	}
	return nil
}

var (
	ErrNotSelected = errors.New("Not selected.")
)

func (s *BinlogSaver) Load() (string, int, error) {
	if s.File != nil {
		return s.File.Load()
	}
	return "", 0, ErrNotSelected
}

func (s BinlogSaver) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.File),
	)
}
