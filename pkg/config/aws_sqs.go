package config

import (
	"errors"
	"strings"
)

type SQS struct {
	IEndpoint `yaml:",inline"`
	Queues    []Queue `yaml:"queues"`
}

type Queue struct {
	TableName              string `yaml:"tableName"`
	QueueName              string `yaml:"queueName"`
	QueueUrl               string `yaml:"queueUrl"`
	MessageGroupIdTemplate string `yaml:"messageGroupIdTemplate"`
}

func (t *Queue) GetMessageGroupId(mp map[string]interface{}) *string {
	if !t.IsFIFO() {
		return nil
	}
	value := templateExecute(t.MessageGroupIdTemplate, mp)
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
