//go:generate go run -mod=mod github.com/golang/mock/mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSClient interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}
