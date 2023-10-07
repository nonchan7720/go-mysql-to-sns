package service

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/aws"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/interfaces"
)

type awsPublisher struct {
	client       aws.SNSClient
	conf         *config.AWS
	mpTableTopic map[string]config.Topic
}

var (
	_ interfaces.Publisher = (*awsPublisher)(nil)
)

func NewAWSPublisher(client aws.SNSClient, conf *config.AWS) interfaces.Publisher {
	return newAWSPublisher(client, conf)
}

func newAWSPublisher(client aws.SNSClient, conf *config.AWS) *awsPublisher {
	mpTableTopic := make(map[string]config.Topic, len(conf.SNS.Topics))
	for _, topic := range conf.SNS.Topics {
		mpTableTopic[topic.TableName] = topic
	}
	return &awsPublisher{
		client:       client,
		conf:         conf,
		mpTableTopic: mpTableTopic,
	}
}

func (p *awsPublisher) Publish(ctx context.Context, payload interfaces.Payload) error {
	slog.With(slog.String("Table", payload.Table)).InfoContext(ctx, "Publish.")
	topic, ok := p.mpTableTopic[payload.Table]
	if !ok {
		// 登録されていないテーブルは対象外
		return nil
	}
	for idx, row := range payload.Rows {
		v, err := payload.ToJson(idx)
		if err != nil {
			return err
		}
		input := &sns.PublishInput{
			Message:        &v,
			MessageGroupId: topic.GetMessageGroupId(row.MainRow(payload.Event)),
			TargetArn:      &topic.TopicArn,
		}
		if _, err := p.client.Publish(ctx, input); err != nil {
			return err
		}
	}
	return nil
}
