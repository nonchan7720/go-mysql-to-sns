package config

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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
	err := validation.Validate(&aws)
	require.Equal(t, "SNS: (Topics: cannot be blank.); SQS: (Queues: cannot be blank.).", err.Error())
}
