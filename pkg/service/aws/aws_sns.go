package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/interfaces"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/interfaces/aws"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/utils"
)

type awsSNS struct {
	client aws.SNSClient
	conf   *config.AWS
}

var (
	_ interfaces.BackendPublisher = (*awsSNS)(nil)
)

func NewAWSSNS(ctx context.Context, client aws.SNSClient, conf *config.AWS) interfaces.BackendPublisher {
	return newAWSSNS(ctx, client, conf)
}

func newAWSSNS(ctx context.Context, client aws.SNSClient, conf *config.AWS) *awsSNS {
	return &awsSNS{
		client: client,
		conf:   conf,
	}
}

func (p *awsSNS) IsTarget(ctx context.Context, payload interfaces.SendPayload) bool {
	_, ok := p.findTopic(payload)
	return ok
}

func (p *awsSNS) findTopic(payload interfaces.SendPayload) (config.Topic, bool) {
	// 最初に見つかったtopicを使用する
	for _, topic := range p.conf.SNS.Topics {
		if topic.Transform.IsTable() {
			if topic.Transform.Table.IsEnabled(payload.Schema, payload.Table) {
				return topic, true
			}
		} else {
			if !topic.Transform.Column.Table.IsEnabled(payload.Schema, payload.Table) {
				continue
			}
			row := utils.Mapper(payload.Row.MainRow(payload.Event))
			v, ok := row.Find(topic.Transform.Column.ColumnName)
			if !ok {
				continue
			}
			value, ok := v.(string)
			if !ok {
				continue
			}
			if strings.EqualFold(topic.Transform.Column.Value, value) {
				return topic, true
			}
		}
	}
	return config.Topic{}, false
}

func (p *awsSNS) Publish(ctx context.Context, event interfaces.Event, payload interfaces.SendPayload) (string, error) {
	topic, _ := p.findTopic(payload)
	v, err := payload.ToJson()
	if err != nil {
		return "", err
	}
	input := &sns.PublishInput{
		Message:        &v,
		MessageGroupId: topic.GetMessageGroupId(payload.Row.MainRow(payload.Event)),
		TargetArn:      &topic.TopicArn,
	}
	if output, err := p.client.Publish(ctx, input); err != nil {
		return "", err
	} else {
		return *output.MessageId, nil
	}
}
