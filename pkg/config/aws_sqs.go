package config

import (
	"errors"
	"strings"
	"sync"

	"github.com/creasty/defaults"
)

var (
	ErrNotFoundAggregateTypeQueue = errors.New("Not found aggregate type queue.")
)

type SQS struct {
	IEndpoint `yaml:",inline"`
	Queues    []Queue `yaml:"queues"`

	mpAggregateTypeTopic   map[string]Queue
	lockAggregateTypeTopic sync.Mutex
	onceMapTopic           sync.Once
}

func (sqs *SQS) FindOutboxQueue(aggregateType string) (Queue, error) {
	sqs.lockAggregateTypeTopic.Lock()
	defer sqs.lockAggregateTypeTopic.Unlock()
	sqs.onceMapTopic.Do(func() {
		for _, queue := range sqs.Queues {
			if queue.Transform.IsOutbox() {
				sqs.mpAggregateTypeTopic[queue.Transform.Outbox.AggregateType] = queue
			}
		}
	})

	if v, ok := sqs.mpAggregateTypeTopic[aggregateType]; ok {
		return v, nil
	}
	return Queue{}, ErrNotFoundAggregateTypeTopic
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
