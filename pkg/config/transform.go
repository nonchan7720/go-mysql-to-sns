package config

import (
	"strings"
)

type Transform struct {
	Type   TransformType    `yaml:"type" default:"table"`
	Table  *TransformTable  `yaml:"table"`
	Column *TransformColumn `yaml:"column"`
	Outbox *TransformOutbox `yaml:"outbox"`
}

func (t *Transform) IsTable() bool {
	return t.Type == TableType
}

func (t *Transform) IsOutbox() bool {
	return t.Type == OutboxPatternType
}

type TransformType string

const (
	TableType         TransformType = TransformType("table")
	ColumnType        TransformType = TransformType("column")
	OutboxPatternType TransformType = TransformType("outbox")
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

type TransformOutbox struct {
	AggregateType string `yaml:"aggregateType"`
}
