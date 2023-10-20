package config

import (
	"time"

	"github.com/creasty/defaults"
)

type OutboxConfig struct {
	TransformTable  `yaml:",inline"`
	MaxRetryCount   int           `yaml:"retryCount" default:"10"`
	PollingInterval time.Duration `yaml:"pollingInterval" default:"5s"`
	RetryBackOff    time.Duration `yaml:"retryBackoff" default:"20s"`
}

func (c *OutboxConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(c); err != nil {
		return err
	}
	type plain OutboxConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	return nil
}

type OutboxPolling struct {
	Config       `yaml:",inline"`
	OutboxConfig OutboxConfig `yaml:"outbox"`
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
