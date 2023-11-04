package config

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestAWSConfig(t *testing.T) {
	data := `accessKey: dummy
secretKey: dummy
`
	var config AWS
	err := yaml.Unmarshal([]byte(data), &config)
	require.NoError(t, err)
	require.Equal(t, config.AccessKey, "dummy")
	require.Equal(t, config.SecretKet, "dummy")
}

func TestAWSValidation(t *testing.T) {
	aws := AWS{
		SNS: &SNS{},
		SQS: &SQS{},
	}
	err := Validate(&aws)
	require.Equal(t, "sns: (topics: cannot be blank.); sqs: (queues: cannot be blank.).", err.Error())
}
