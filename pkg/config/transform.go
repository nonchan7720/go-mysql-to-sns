package config

import (
	"strings"

	"github.com/creasty/defaults"
)

type Transform struct {
	Type   TransformType    `yaml:"type" default:"table"`
	Table  *TransformTable  `yaml:"table"`
	Column *TransformColumn `yaml:"column"`
}

func (t *Transform) IsTable() bool {
	return t.Type == TableType
}

func (t *Transform) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(t); err != nil {
		return err
	}
	type plain Transform
	if err := unmarshal((*plain)(t)); err != nil {
		return err
	}
	return nil
}

type TransformType string

const (
	TableType  TransformType = TransformType("table")
	ColumnType TransformType = TransformType("column")
)

type TransformTable struct {
	Schema    string `yaml:"schema"`
	TableName string `yaml:"tableName"`
}

func (t *TransformTable) IsEnabled(schema, tableName string) bool {
	return strings.EqualFold(t.Schema, schema) && strings.EqualFold(t.TableName, tableName)
}

type TransformColumn struct {
	Table      TransformTable `yaml:",inline"`
	ColumnName string         `yaml:"columnName"`
	Value      string         `yaml:"value"`
}
