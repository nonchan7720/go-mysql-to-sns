package config

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Outbox struct {
	Config          `yaml:",inline"`
	*TransformTable `yaml:"outbox"`
}

func LoadOutboxConfig(filePath string) (*Outbox, error) {
	config, err := loadConfig[Outbox](filePath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (conf Outbox) Validate() error {
	return validation.ValidateStruct(&conf,
		validation.Field(&conf.Config),
		validation.Field(&conf.TransformTable, validation.NotNil),
	)
}
