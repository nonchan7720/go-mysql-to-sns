package service

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/golang/mock/gomock"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/interfaces"
	mockAws "github.com/nonchan7720/go-mysql-to-sns/pkg/mock/aws"
	"github.com/stretchr/testify/require"
)

func TestAWSSQS(t *testing.T) {
	var tables = []struct {
		name    string
		payload interfaces.Payload
		queue   config.Queue
		fn      func(client *mockAws.MockSQSClient, require *require.Assertions)
	}{
		{
			name: "Create Row",
			payload: interfaces.Payload{
				Event:  interfaces.Create,
				Schema: "public",
				Table:  "test",
				Rows: []interfaces.PayloadRow{
					interfaces.NewPayloadRow(nil, map[string]interface{}{"key": "value"}),
				},
			},
			queue: config.Queue{
				TableName: "test",
				QueueUrl:  "http://localhost:4566/000000000000/test-sqs",
			},
			fn: func(client *mockAws.MockSQSClient, require *require.Assertions) {
				output := &sqs.SendMessageOutput{}
				client.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, input *sqs.SendMessageInput, optFns ...func(*sqs.Options)) {
					require.NotNil(input)
					require.NotNil(input.MessageBody)
					require.NotNil(input.QueueUrl)
					require.Equal(*input.QueueUrl, "http://localhost:4566/000000000000/test-sqs")
					require.Equal(*input.MessageBody, `{"event":"CREATE","schema":"public","table":"test","row":{"old_row":{},"new_row":{"key":"value"}}}`)
				}).Return(output, nil).Times(1)
			},
		},
		{
			name: "Create Rows",
			payload: interfaces.Payload{
				Event:  interfaces.Create,
				Schema: "public",
				Table:  "test",
				Rows: []interfaces.PayloadRow{
					interfaces.NewPayloadRow(nil, map[string]interface{}{"key": "value"}),
					interfaces.NewPayloadRow(nil, map[string]interface{}{"key": "value"}),
				},
			},
			queue: config.Queue{
				TableName: "test",
				QueueUrl:  "http://localhost:4566/000000000000/test-sqs",
			},
			fn: func(client *mockAws.MockSQSClient, require *require.Assertions) {
				output := &sqs.SendMessageOutput{}
				client.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, input *sqs.SendMessageInput, optFns ...func(*sqs.Options)) {
					require.NotNil(input)
					require.NotNil(input.MessageBody)
					require.NotNil(input.QueueUrl)
					require.Equal(*input.QueueUrl, "http://localhost:4566/000000000000/test-sqs")
					require.Equal(*input.MessageBody, `{"event":"CREATE","schema":"public","table":"test","row":{"old_row":{},"new_row":{"key":"value"}}}`)
				}).Return(output, nil).Times(2)
			},
		},
		{
			name: "Update Row",
			payload: interfaces.Payload{
				Event:  interfaces.Update,
				Schema: "public",
				Table:  "test",
				Rows: []interfaces.PayloadRow{
					interfaces.NewPayloadRow(map[string]interface{}{"key": "oldValue"}, map[string]interface{}{"key": "value"}),
				},
			},
			queue: config.Queue{
				TableName: "test",
				QueueUrl:  "http://localhost:4566/000000000000/test-sqs",
			},
			fn: func(client *mockAws.MockSQSClient, require *require.Assertions) {
				output := &sqs.SendMessageOutput{}
				client.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, input *sqs.SendMessageInput, optFns ...func(*sqs.Options)) {
					require.NotNil(input)
					require.NotNil(input.MessageBody)
					require.NotNil(input.QueueUrl)
					require.Equal(*input.QueueUrl, "http://localhost:4566/000000000000/test-sqs")
					require.Equal(*input.MessageBody, `{"event":"UPDATE","schema":"public","table":"test","row":{"old_row":{"key":"oldValue"},"new_row":{"key":"value"}}}`)
				}).Return(output, nil).Times(1)
			},
		},
		{
			name: "Delete Row",
			payload: interfaces.Payload{
				Event:  interfaces.Delete,
				Schema: "public",
				Table:  "test",
				Rows: []interfaces.PayloadRow{
					interfaces.NewPayloadRow(map[string]interface{}{"key": "oldValue"}, nil),
				},
			},
			queue: config.Queue{
				TableName: "test",
				QueueUrl:  "http://localhost:4566/000000000000/test-sqs",
			},
			fn: func(client *mockAws.MockSQSClient, require *require.Assertions) {
				output := &sqs.SendMessageOutput{}
				client.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, input *sqs.SendMessageInput, optFns ...func(*sqs.Options)) {
					require.NotNil(input)
					require.NotNil(input.MessageBody)
					require.NotNil(input.QueueUrl)
					require.Equal(*input.QueueUrl, "http://localhost:4566/000000000000/test-sqs")
					require.Equal(*input.MessageBody, `{"event":"DELETE","schema":"public","table":"test","row":{"old_row":{"key":"oldValue"},"new_row":{}}}`)
				}).Return(output, nil).Times(1)
			},
		},
	}
	ctx := context.Background()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	for _, tbl := range tables {
		tbl := tbl
		t.Run(tbl.name, func(t *testing.T) {
			require := require.New(t)
			client := mockAws.NewMockSQSClient(mockCtl)
			tbl.fn(client, require)
			conf := config.AWS{
				SQS: &config.SQS{
					Queues: []config.Queue{
						tbl.queue,
					},
				},
			}
			p, err := newAWSSQS(ctx, client, &conf)
			require.NoError(err)
			err = p.Publish(ctx, tbl.payload)
			require.NoError(err)
		})
	}
}
