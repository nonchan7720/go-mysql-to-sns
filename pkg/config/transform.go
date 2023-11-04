package config

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Transform struct {
	Type   TransformType    `yaml:"type" default:"table"`
	Table  *TransformTable  `yaml:"table"`
	Column *TransformColumn `yaml:"column"`
	Outbox *TransformOutbox `yaml:"outbox"`
}

var (
	_ validation.Validatable = (*Transform)(nil)
)

func (t *Transform) IsTable() bool {
	return t.Type == TableType
}

func (t *Transform) IsOutbox() bool {
	return t.Type == OutboxPatternType
}

func (t Transform) Validate() error {
	fn := func(value string) bool {
		return value == "table" || value == "column" || value == "outbox"
	}

	return validation.ValidateStruct(&t,
		validation.Field(&t.Type,
			validation.NewStringRuleWithError(
				fn,
				validation.NewError("validation_is_transform_type", "must be a value with table, column or outbox"),
			),
		),
		validation.Field(&t.Table),
		validation.Field(&t.Column),
		validation.Field(&t.Outbox),
	)
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

func (t TransformTable) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Schema, validation.Required),
		validation.Field(&t.TableName, validation.Required),
	)
}

type TransformColumn struct {
	Table      TransformTable `yaml:",inline"`
	ColumnName string         `yaml:"columnName"`
	Value      string         `yaml:"value"`
}

func (t TransformColumn) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Table),
		validation.Field(&t.ColumnName, validation.Required),
		validation.Field(&t.Value, validation.Required),
	)
}

type TransformOutbox struct {
	AggregateType string `yaml:"aggregateType"`
}

func (t TransformOutbox) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.AggregateType, validation.Required),
	)
}
