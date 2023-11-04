package config

import (
	"errors"
	"strings"
	"sync"
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
