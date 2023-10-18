package config

type Outbox struct {
	*Config        `yaml:",inline"`
	TransformTable `yaml:"outbox"`
}

func LoadOutboxConfig(filePath string) (*Outbox, error) {
	config, err := loadConfig[Outbox](filePath)
	if err != nil {
		return nil, err
	}
	if err := config.Validation(); err != nil {
		return nil, err
	}
	return config, nil
}
