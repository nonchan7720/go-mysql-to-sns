package config

import (
	"errors"
	"strings"
	"sync"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrNotFoundAggregateTypeTopic = errors.New("Not found aggregate type topic.")
)

type SNS struct {
	IEndpoint `yaml:",inline"`
	Topics    []Topic `yaml:"topics"`

	mpAggregateTypeTopic   map[string]Topic
	lockAggregateTypeTopic sync.Mutex
	onceMapTopic           sync.Once
}

var (
	_ validation.Validatable = (*SNS)(nil)
)

func (sns *SNS) FindOutboxTopic(aggregateType string) (Topic, error) {
	sns.lockAggregateTypeTopic.Lock()
	defer sns.lockAggregateTypeTopic.Unlock()
	sns.onceMapTopic.Do(func() {
		sns.mpAggregateTypeTopic = make(map[string]Topic)
		for _, topic := range sns.Topics {
			if topic.Transform.IsOutbox() {
				sns.mpAggregateTypeTopic[topic.Transform.Outbox.AggregateType] = topic
			}
		}
	})

	if v, ok := sns.mpAggregateTypeTopic[aggregateType]; ok {
		return v, nil
	}
	return Topic{}, ErrNotFoundAggregateTypeTopic
}

func (sns *SNS) FindOutboxTopicArn(aggregateType string) (string, error) {
	t, err := sns.FindOutboxTopic(aggregateType)
	if err != nil {
		return "", err
	}
	return t.TopicArn, nil
}

func (sns *SNS) Validate() error {
	return validation.ValidateStruct(sns,
		validation.Field(&sns.Topics, validation.Required),
	)
}

type Topic struct {
	TopicArn               string       `yaml:"topicArn"`
	MessageGroupIdTemplate string       `yaml:"messageGroupIdTemplate"`
	TemplateType           TemplateType `yaml:"templateType" default:"Fast"`
	Transform              Transform    `yaml:"transform"`
}

var (
	_ validation.Validatable = (*Topic)(nil)
)

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

func (t Topic) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.TopicArn, validation.Required),
		validation.Field(&t.MessageGroupIdTemplate,
			validation.Required.When(t.IsFIFO() && !t.Transform.IsOutbox()),
		),
		validation.Field(&t.Transform),
	)
}
