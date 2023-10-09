package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/interfaces"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/interfaces/aws"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/service"
)

type awsSQS struct {
	client       aws.SQSClient
	conf         *config.AWS
	mpTableQueue map[string]config.Queue
}

var (
	_ interfaces.BackendPublisher = (*awsSQS)(nil)
)

func NewAWSSQS(ctx context.Context, client aws.SQSClient, conf *config.AWS) (interfaces.BackendPublisher, error) {
	return newAWSSQS(ctx, client, conf)
}

func newAWSSQS(ctx context.Context, client aws.SQSClient, conf *config.AWS) (*awsSQS, error) {
	mpTableQueue := make(map[string]config.Queue, len(conf.SQS.Queues))
	for idx := range conf.SQS.Queues {
		queue := conf.SQS.Queues[idx]
		if queue.QueueUrl == "" {
			resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: &queue.QueueName})
			if err != nil {
				return nil, err
			}
			queue.QueueUrl = *resp.QueueUrl
			conf.SQS.Queues[idx] = queue
		}
		mpTableQueue[strings.ToLower(queue.TableName)] = queue
	}
	return &awsSQS{
		client:       client,
		conf:         conf,
		mpTableQueue: mpTableQueue,
	}, nil
}

func (p *awsSQS) IsTarget(ctx context.Context, payload interfaces.SendPayload) bool {
	return service.IsTarget(p.mpTableQueue, payload)
}

func (p *awsSQS) Publish(ctx context.Context, event interfaces.Event, payload interfaces.SendPayload) (string, error) {
	queue := service.FindTarget(p.mpTableQueue, payload)
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
