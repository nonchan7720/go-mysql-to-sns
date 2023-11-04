package config

import (
	"time"
)

type OutboxPollingConfig struct {
	TransformTable  `yaml:",inline"`
	ProducerName    string        `yaml:"producerName" validate:"required"`
	MaxRetryCount   int           `yaml:"retryCount" default:"10"`
	PollingInterval time.Duration `yaml:"pollingInterval" default:"5s"`
	RetryBackOff    time.Duration `yaml:"retryBackoff" default:"20s"`
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
