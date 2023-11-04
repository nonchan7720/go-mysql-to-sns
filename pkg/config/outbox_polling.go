package config

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type OutboxPollingConfig struct {
	TransformTable  `yaml:",inline"`
	ProducerName    string        `yaml:"producerName"`
	MaxRetryCount   int           `yaml:"retryCount" default:"10"`
	PollingInterval time.Duration `yaml:"pollingInterval" default:"5s"`
	RetryBackOff    time.Duration `yaml:"retryBackoff" default:"20s"`
}

var (
	_ validation.Validatable = (*OutboxPollingConfig)(nil)
)

func (o OutboxPollingConfig) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(o.ProducerName, validation.Required),
	)
}

type OutboxPolling struct {
	Config       `yaml:",inline"`
	OutboxConfig OutboxPollingConfig `yaml:"outbox"`
}

func LoadOutboxPollingConfig(filePath string) (*OutboxPolling, error) {
	config, err := loadConfig[OutboxPolling](filePath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (conf OutboxPolling) Validate() error {
	return validation.ValidateStruct(&conf,
		validation.Field(&conf.Config),
		validation.Field(&conf.OutboxConfig),
	)
}
