package config

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrNotFoundProducer = errors.New("Not found producer")
)

type Publisher struct {
	AWS *AWS `yaml:"aws"`
}

var (
	_ validation.Validatable = (*Publisher)(nil)
)

func (p *Publisher) IsAWS() bool {
	return p.AWS != nil
}

func (p *Publisher) FindProducer(value string) (string, error) {
	if p.IsAWS() {
		if p.AWS.IsSNS() {
			return p.AWS.SNS.FindOutboxTopicArn(value)
		}
		if p.AWS.IsSQS() {
			return p.AWS.SQS.FindOutboxQueueUrl(value)
		}
	}
	return "", ErrNotFoundProducer
}

func (p Publisher) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.AWS),
	)
}
