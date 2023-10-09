package config

import (
	"errors"
	"strings"
)

type SNS struct {
	IEndpoint `yaml:",inline"`
	Topics    []Topic `yaml:"topics"`
}

type Topic struct {
	TableName              string `yaml:"tableName"`
	TopicArn               string `yaml:"topicArn"`
	MessageGroupIdTemplate string `yaml:"messageGroupIdTemplate"`
}

func (t *Topic) GetMessageGroupId(mp map[string]interface{}) *string {
	if !t.IsFIFO() {
		return nil
	}
	value := templateExecute(t.MessageGroupIdTemplate, mp)
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
