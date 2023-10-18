package config

import (
	"io"

	"gopkg.in/yaml.v3"
)

type positionSaver struct{}

type fileValue struct {
	File     string `yaml:"file"`
	Position int    `yaml:"position"`
}

func (p positionSaver) save(f io.Writer, format BinlogSaveFormat) error {
	value := fileValue{
		File:     format.File,
		Position: format.Position,
	}
	return yaml.NewEncoder(f).Encode(&value)
}

func (p positionSaver) load(f io.Reader) (format *BinlogSaveFormat, err error) {
	var value fileValue
	err = yaml.NewDecoder(f).Decode(&value)
	if err != nil {
		return nil, err
	}
	format = &BinlogSaveFormat{
		File:     value.File,
		Position: value.Position,
	}
	return format, err
}
