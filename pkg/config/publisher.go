package config

type Publisher struct {
	AWS *AWS `yaml:"aws"`
}

func (p *Publisher) IsAWS() bool {
	return p.AWS != nil
}

func (p *Publisher) Validation() error {
	if p.AWS != nil {
		return p.AWS.Validation()
	}
	return nil
}
