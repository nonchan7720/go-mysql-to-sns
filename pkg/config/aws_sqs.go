package config

import (
	"errors"
	"strings"

	"github.com/creasty/defaults"
)

type SQS struct {
	IEndpoint `yaml:",inline"`
	Queues    []Queue `yaml:"queues"`
}

type Queue struct {
	QueueName              string       `yaml:"queueName"`
	QueueUrl               string       `yaml:"queueUrl"`
	MessageGroupIdTemplate string       `yaml:"messageGroupIdTemplate"`
	TemplateType           TemplateType `yaml:"templateType" default:"Fast"`
	Transform              Transform    `yaml:"transform"`
}

func (t *Queue) GetMessageGroupId(mp map[string]interface{}) *string {
	if !t.IsFIFO() {
		return nil
	}
	value := templateExecute(t.TemplateType, t.MessageGroupIdTemplate, mp)
	return &value
}

func (t *Queue) IsFIFO() bool {
	return strings.HasSuffix(t.QueueUrl, ".fifo")
}

func (t *Queue) Validation() error {
	if t.IsFIFO() && t.MessageGroupIdTemplate == "" {
		return errors.New("For FIFO topics, MessageGroupIdTemplate must be set.")
	}
	return nil
}

func (t *Queue) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(t); err != nil {
		return err
	}
	type plain Queue
	if err := unmarshal((*plain)(t)); err != nil {
		return err
	}
	return nil
}
