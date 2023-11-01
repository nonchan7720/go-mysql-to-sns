package config

import (
	"fmt"
	"time"

	"github.com/creasty/defaults"
)

type OutboxPollingConfig struct {
	TransformTable  `yaml:",inline"`
	ProducerName    string        `yaml:"producerName"`
	MaxRetryCount   int           `yaml:"retryCount" default:"10"`
	PollingInterval time.Duration `yaml:"pollingInterval" default:"5s"`
	RetryBackOff    time.Duration `yaml:"retryBackoff" default:"20s"`
}

func (c *OutboxPollingConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(c); err != nil {
		return err
	}
	type plain OutboxPollingConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if c.ProducerName == "" {
		return fmt.Errorf("producerName is required.")
	}
	return nil
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
	if err := config.Validation(); err != nil {
		return nil, err
	}
	return config, nil
}
