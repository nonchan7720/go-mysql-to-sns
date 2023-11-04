package config

import (
	"errors"
	"strings"
	"sync"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
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

var (
	_ validation.Validatable = (*SQS)(nil)
)

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

func (sqs *SQS) Validate() error {
	return validation.ValidateStruct(sqs,
		validation.Field(&sqs.Queues, validation.Required),
	)
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

func (t Queue) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.QueueName, validation.Required.When(t.QueueUrl == "")),
		validation.Field(&t.QueueUrl, validation.Required.When(t.QueueName == "")),
		validation.Field(&t.QueueUrl, validation.When(t.QueueUrl != "", is.URL)),
		validation.Field(&t.MessageGroupIdTemplate,
			validation.Required.When(t.IsFIFO() && !t.Transform.IsOutbox()),
		),
		validation.Field(&t.Transform),
	)
}
