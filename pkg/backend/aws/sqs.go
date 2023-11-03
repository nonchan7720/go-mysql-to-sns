package aws

import (
	"context"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/config"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/interfaces/aws"
)

func NewSQSClient(ctx context.Context, conf *config.AWS) (aws.SQSClient, error) {
	endpoint := NewEndpoint(WithSQSEndpoint(conf.SQS.Endpoint))
	awsConfig, err := NewConfig(ctx,
		endpoint.EndpointResolver(),
		awsConfig.WithCredentialsProvider(conf.WithStatic()),
	)
	if err != nil {
		return nil, err
	}
	return sqs.NewFromConfig(awsConfig), nil
}
