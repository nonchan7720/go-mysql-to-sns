package config

import (
	"errors"
	"strings"
	"sync"

	"github.com/creasty/defaults"
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

func (sns *SNS) FindOutboxTopic(aggregateType string) (Topic, error) {
	sns.lockAggregateTypeTopic.Lock()
	defer sns.lockAggregateTypeTopic.Unlock()
	sns.onceMapTopic.Do(func() {
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
