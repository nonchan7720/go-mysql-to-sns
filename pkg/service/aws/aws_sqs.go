package service

import (
	"context"
	"log/slog"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/interfaces"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/interfaces/aws"
)

type awsSQS struct {
	client       aws.SQSClient
	conf         *config.AWS
	mpTableQueue map[string]config.Queue
}

var (
	_ interfaces.Publisher = (*awsPublisher)(nil)
)

func NewAWSSQS(ctx context.Context, client aws.SQSClient, conf *config.AWS) (interfaces.Publisher, error) {
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

func (p *awsSQS) Publish(ctx context.Context, payload interfaces.Payload) error {
	slog.With(slog.String("Table", payload.Table)).InfoContext(ctx, "Receive payload.")
	queue, ok := p.mpTableQueue[strings.ToLower(payload.Table)]
	if !ok {
		// 登録されていないテーブルは対象外
		return nil
	}
	slog.With(slog.String("Table", payload.Table)).InfoContext(ctx, "Publish.")
	for idx, row := range payload.Rows {
		v, err := payload.ToJson(idx)
		if err != nil {
			return err
		}
		input := &sqs.SendMessageInput{
			MessageBody:    &v,
			MessageGroupId: queue.GetMessageGroupId(row.MainRow(payload.Event)),
			QueueUrl:       &queue.QueueUrl,
		}
		if _, err := p.client.SendMessage(ctx, input); err != nil {
			return err
		}
	}
	return nil
}
