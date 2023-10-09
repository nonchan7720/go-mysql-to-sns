package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/interfaces"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/interfaces/aws"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/service"
)

type awsSNS struct {
	client       aws.SNSClient
	conf         *config.AWS
	mpTableTopic map[string]config.Topic
}

var (
	_ interfaces.BackendPublisher = (*awsSNS)(nil)
)

func NewAWSSNS(ctx context.Context, client aws.SNSClient, conf *config.AWS) interfaces.BackendPublisher {
	return newAWSSNS(ctx, client, conf)
}

func newAWSSNS(ctx context.Context, client aws.SNSClient, conf *config.AWS) *awsSNS {
	mpTableTopic := make(map[string]config.Topic, len(conf.SNS.Topics))
	for _, topic := range conf.SNS.Topics {
		mpTableTopic[topic.TableName] = topic
	}
	return &awsSNS{
		client:       client,
		conf:         conf,
		mpTableTopic: mpTableTopic,
	}
}

func (p *awsSNS) IsTarget(ctx context.Context, payload interfaces.SendPayload) bool {
	return service.IsTarget(p.mpTableTopic, payload)
}

func (p *awsSNS) Publish(ctx context.Context, event interfaces.Event, payload interfaces.SendPayload) (string, error) {
	topic := service.FindTarget(p.mpTableTopic, payload)
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
