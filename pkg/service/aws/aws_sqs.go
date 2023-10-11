package aws

import (
	"context"
	"strings"

	originAWS "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/interfaces"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/interfaces/aws"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/utils"
)

type awsSQS struct {
	client aws.SQSClient
	conf   *config.AWS
}

var (
	_ interfaces.BackendPublisher = (*awsSQS)(nil)
)

func NewAWSSQS(ctx context.Context, client aws.SQSClient, conf *config.AWS) (interfaces.BackendPublisher, error) {
	return newAWSSQS(ctx, client, conf)
}

func newAWSSQS(ctx context.Context, client aws.SQSClient, conf *config.AWS) (*awsSQS, error) {
	return &awsSQS{
		client: client,
		conf:   conf,
	}, nil
}

func (p *awsSQS) IsTarget(ctx context.Context, payload interfaces.SendPayload) bool {
	_, ok := p.findQueue(payload)
	return ok
}

func (p *awsSQS) findQueue(payload interfaces.SendPayload) (config.Queue, bool) {
	// 最初に見つかったtopicを使用する
	for _, queue := range p.conf.SQS.Queues {
		if queue.Transform.IsTable() {
			if queue.Transform.Table.IsEnabled(payload.Schema, payload.Table) {
				return queue, true
			}
		} else {
			if !queue.Transform.Column.Table.IsEnabled(payload.Schema, payload.Table) {
				continue
			}
			row := utils.Mapper(payload.Row.MainRow(payload.Event))
			v, ok := row.Find(queue.Transform.Column.ColumnName)
			if !ok {
				continue
			}
			value, ok := v.(string)
			if !ok {
				continue
			}
			if strings.EqualFold(queue.Transform.Column.Value, value) {
				return queue, true
			}
		}
	}
	return config.Queue{}, false
}

func (p *awsSQS) PublishBinlog(ctx context.Context, event interfaces.Event, payload interfaces.SendPayload) (string, error) {
	queue, _ := p.findQueue(payload)
	v, err := payload.ToJson()
	if err != nil {
		return "", err
	}
	input := &sqs.SendMessageInput{
		MessageBody:    &v,
		MessageGroupId: queue.GetMessageGroupId(payload.Row.MainRow(payload.Event)),
		QueueUrl:       &queue.QueueUrl,
	}
	if output, err := p.client.SendMessage(ctx, input); err != nil {
		return "", err
	} else {
		return *output.MessageId, nil
	}
}

func (p *awsSQS) PublishOutbox(ctx context.Context, outbox interfaces.Outbox) (string, error) {
	queue, err := p.conf.SQS.FindOutboxQueue(outbox.AggregateType)
	if err != nil {
		return "", err
	}
	input := &sqs.SendMessageInput{
		MessageBody:    originAWS.String(outbox.Payload),
		MessageGroupId: originAWS.String(outbox.AggregateId),
		QueueUrl:       originAWS.String(queue.QueueUrl),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"Event": {
				DataType:    originAWS.String("String"),
				StringValue: originAWS.String(outbox.EventType),
			},
		},
	}
	if output, err := p.client.SendMessage(ctx, input); err != nil {
		return "", err
	} else {
		return *output.MessageId, nil
	}
}
