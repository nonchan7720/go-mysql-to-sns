package config

import (
	"errors"
	"strings"
	"sync"
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
		sqs.mpAggregateTypeTopic = make(map[string]Queue)
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

func (sqs *SQS) FindOutboxQueueUrl(aggregateType string) (string, error) {
	q, err := sqs.FindOutboxQueue(aggregateType)
	if err != nil {
		return "", err
	}
	return q.QueueUrl, nil
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
