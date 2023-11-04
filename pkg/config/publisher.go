package config

import "errors"

var (
	ErrNotFoundProducer = errors.New("Not found producer")
)

type Publisher struct {
	AWS *AWS `yaml:"aws"`
}

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
