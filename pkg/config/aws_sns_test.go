package config

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAWSSNSConfig(t *testing.T) {
	data := `endpoint: http://localstack:4566
topics:
  - topicArn: arn:aws:sns:ap-northeast-1:000000000000:test-topic
    messageGroupIdTemplate: id-{{id}}
    transform:
      table:
        schema: public
        tableName: tbl1
      column:
        schema: public
        tableName: tbl1
        columnName: column1
        value: test
`
	var config SNS
	err := loadYaml(bytes.NewBufferString(data), &config)
	require.NoError(t, err)
	require.Equal(t, config.Endpoint, "http://localstack:4566")
	require.Len(t, config.Topics, 1)
	require.Equal(t, config.Topics[0].Transform.Table.Schema, "public")
	require.Equal(t, config.Topics[0].Transform.Table.TableName, "tbl1")
	require.Equal(t, config.Topics[0].Transform.Column.Table.Schema, "public")
	require.Equal(t, config.Topics[0].Transform.Column.Table.TableName, "tbl1")
	require.Equal(t, config.Topics[0].Transform.Column.ColumnName, "column1")
	require.Equal(t, config.Topics[0].Transform.Column.Value, "test")
	require.Equal(t, config.Topics[0].TopicArn, "arn:aws:sns:ap-northeast-1:000000000000:test-topic")
	require.Equal(t, config.Topics[0].MessageGroupIdTemplate, "id-{{id}}")
	require.Equal(t, config.Topics[0].TemplateType, FastTemplate)
}

func TestAWSSNSConfigWithGoTemplate(t *testing.T) {
	data := `endpoint: http://localstack:4566
topics:
  - topicArn: arn:aws:sns:ap-northeast-1:000000000000:test-topic
    messageGroupIdTemplate: id-{{id}}
    templateType: Go
    transform:
      table:
        schema: public
        tableName: tbl1
      column:
        schema: public
        tableName: tbl1
        columnName: column1
        value: test
`
	var config SNS
	err := loadYaml(bytes.NewBufferString(data), &config)
	require.NoError(t, err)
	require.Equal(t, config.Endpoint, "http://localstack:4566")
	require.Len(t, config.Topics, 1)
	require.Equal(t, config.Topics[0].Transform.Table.Schema, "public")
	require.Equal(t, config.Topics[0].Transform.Table.TableName, "tbl1")
	require.Equal(t, config.Topics[0].Transform.Column.Table.Schema, "public")
	require.Equal(t, config.Topics[0].Transform.Column.Table.TableName, "tbl1")
	require.Equal(t, config.Topics[0].Transform.Column.ColumnName, "column1")
	require.Equal(t, config.Topics[0].Transform.Column.Value, "test")
	require.Equal(t, config.Topics[0].TopicArn, "arn:aws:sns:ap-northeast-1:000000000000:test-topic")
	require.Equal(t, config.Topics[0].MessageGroupIdTemplate, "id-{{id}}")
	require.Equal(t, config.Topics[0].TemplateType, GoTemplate)
}
