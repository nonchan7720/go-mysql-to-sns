package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type FileSaver struct {
	Name string `yaml:"name"`
}

type fileValue struct {
	File     string `yaml:"file"`
	Position int    `yaml:"position"`
}

func (s *FileSaver) Save(file string, position int) error {
	f, err := os.Create(s.Name)
	if err != nil {
		return err
	}
	defer f.Close()
	value := fileValue{
		File:     file,
		Position: position,
	}
	return yaml.NewEncoder(f).Encode(&value)
}

func (s *FileSaver) Load() (file string, position int, err error) {
	var value fileValue
	f, err := os.Open(s.Name)
	if err != nil {
		return
	}
	defer f.Close()
	err = yaml.NewDecoder(f).Decode(&value)
	return value.File, value.Position, err
}
