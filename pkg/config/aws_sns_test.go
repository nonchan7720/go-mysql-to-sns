package config

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestAWSSNSConfig(t *testing.T) {
	data := `endpoint: http://localstack:4566
topics:
  - tableName: tbl1
    topicArn: arn:aws:sns:ap-northeast-1:000000000000:test-topic
    messageGroupIdTemplate: id-{{id}}
`
	var config SNS
	err := yaml.Unmarshal([]byte(data), &config)
	require.NoError(t, err)
	require.Equal(t, config.Endpoint, "http://localstack:4566")
	require.Len(t, config.Topics, 1)
	require.Equal(t, config.Topics[0].TableName, "tbl1")
	require.Equal(t, config.Topics[0].TopicArn, "arn:aws:sns:ap-northeast-1:000000000000:test-topic")
	require.Equal(t, config.Topics[0].MessageGroupIdTemplate, "id-{{id}}")
	require.Equal(t, config.Topics[0].TemplateType, FastTemplate)
}

func TestAWSSNSConfigWithGoTemplate(t *testing.T) {
	data := `endpoint: http://localstack:4566
topics:
  - tableName: tbl1
    topicArn: arn:aws:sns:ap-northeast-1:000000000000:test-topic
    messageGroupIdTemplate: id-{{id}}
    templateType: Go
`
	var config SNS
	err := yaml.Unmarshal([]byte(data), &config)
	require.NoError(t, err)
	require.Equal(t, config.Endpoint, "http://localstack:4566")
	require.Len(t, config.Topics, 1)
	require.Equal(t, config.Topics[0].TableName, "tbl1")
	require.Equal(t, config.Topics[0].TopicArn, "arn:aws:sns:ap-northeast-1:000000000000:test-topic")
	require.Equal(t, config.Topics[0].MessageGroupIdTemplate, "id-{{id}}")
	require.Equal(t, config.Topics[0].TemplateType, GoTemplate)
}
