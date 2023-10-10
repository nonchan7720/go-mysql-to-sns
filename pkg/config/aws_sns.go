package config

import (
	"errors"
	"strings"

	"github.com/creasty/defaults"
)

type SNS struct {
	IEndpoint `yaml:",inline"`
	Topics    []Topic `yaml:"topics"`
}

type Topic struct {
	TopicArn               string       `yaml:"topicArn"`
	MessageGroupIdTemplate string       `yaml:"messageGroupIdTemplate"`
	TemplateType           TemplateType `yaml:"templateType" default:"Fast"`
	Transform              Transform    `yaml:"transform"`
}

func (t *Topic) GetMessageGroupId(mp map[string]interface{}) *string {
	if !t.IsFIFO() {
		return nil
	}
	value := templateExecute(t.TemplateType, t.MessageGroupIdTemplate, mp)
	return &value
}

func (t *Topic) IsFIFO() bool {
	return strings.HasSuffix(t.TopicArn, ".fifo")
}

func (t *Topic) Validation() error {
	if t.IsFIFO() && t.MessageGroupIdTemplate == "" {
		return errors.New("For FIFO topics, MessageGroupIdTemplate must be set.")
	}
	return nil
}

func (t *Topic) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(t); err != nil {
		return err
	}
	type plain Topic
	if err := unmarshal((*plain)(t)); err != nil {
		return err
	}
	return nil
}
