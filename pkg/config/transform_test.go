package config

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestTransform(t *testing.T) {
	data := `type: column
table:
  schema: public
  tableName: test
column:
  schema: public
  tableName: test
  columnName: column1
  value: test
`
	var transform Transform
	err := yaml.Unmarshal([]byte(data), &transform)
	require.NoError(t, err)
	require.Equal(t, transform.Type, ColumnType)
	require.Equal(t, transform.Table.Schema, "public")
	require.Equal(t, transform.Table.TableName, "test")
	require.Equal(t, transform.Column.Table.Schema, "public")
	require.Equal(t, transform.Column.Table.TableName, "test")
	require.Equal(t, transform.Column.ColumnName, "column1")
	require.Equal(t, transform.Column.Value, "test")
	require.True(t, transform.Table.IsEnabled("Public", "Test"))
	require.True(t, transform.Column.Table.IsEnabled("Public", "Test"))
	require.False(t, transform.Table.IsEnabled("fPublic", "Test"))
	require.False(t, transform.Column.Table.IsEnabled("fPublic", "Test"))
}

func TestTransformValidation(t *testing.T) {
	transform := Transform{
		Table:  &TransformTable{},
		Column: &TransformColumn{},
		Outbox: &TransformOutbox{},
	}
	err := Validate(&transform)
	require.Equal(t, "column: (Table: (schema: cannot be blank; tableName: cannot be blank.); columnName: cannot be blank; value: cannot be blank.); outbox: (aggregateType: cannot be blank.); table: (schema: cannot be blank; tableName: cannot be blank.).", err.Error())
}
