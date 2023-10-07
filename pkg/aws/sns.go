package aws

import (
	"context"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
)

func NewSNSClient(ctx context.Context, conf *config.AWS) (SNSClient, error) {
	endpoint := NewEndpoint(WithSNSEndpoint(conf.SNS.Endpoint))
	awsConfig, err := NewConfig(ctx,
		endpoint.EndpointResolver(),
		awsConfig.WithCredentialsProvider(conf.WithStatic()),
	)
	if err != nil {
		return nil, err
	}
	return sns.NewFromConfig(awsConfig), nil
}
