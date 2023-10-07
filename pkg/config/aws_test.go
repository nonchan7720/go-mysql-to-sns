package config

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestAWSConfig(t *testing.T) {
	data := `accessKey: dummy
secretKey: dummy
sns:
  endpoint: http://localstack:4566
  topics:
    - tableName: tbl1
      topicArn: arn:aws:sns:ap-northeast-1:000000000000:test-topic
`
	var config AWS
	err := yaml.Unmarshal([]byte(data), &config)
	require.NoError(t, err)
	require.Equal(t, config.AccessKey, "dummy")
	require.Equal(t, config.SecretKet, "dummy")
	require.Equal(t, config.SNS.Endpoint, "http://localstack:4566")
	require.Len(t, config.SNS.Topics, 1)
	require.Equal(t, config.SNS.Topics[0].TableName, "tbl1")
	require.Equal(t, config.SNS.Topics[0].TopicArn, "arn:aws:sns:ap-northeast-1:000000000000:test-topic")
}
