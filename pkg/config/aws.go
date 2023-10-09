package config

import (
	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type IEndpoint struct {
	Endpoint string `yaml:"endpoint"`
}

type AWS struct {
	AccessKey string `yaml:"accessKey"`
	SecretKet string `yaml:"secretKey"`
	SNS       *SNS   `yaml:"sns"`
	SQS       *SQS   `yaml:"sqs"`
}

func (conf *AWS) WithStatic() awsv2.CredentialsProvider {
	if conf.AccessKey == "" && conf.SecretKet == "" {
		return nil
	}
	return credentials.NewStaticCredentialsProvider(conf.AccessKey, conf.SecretKet, "")
}

func (conf *AWS) Validation() error {
	if conf.IsSNS() {
		for _, topic := range conf.SNS.Topics {
			if err := topic.Validation(); err != nil {
				return err
			}
		}
	}
	if conf.IsSQS() {
		for _, queue := range conf.SQS.Queues {
			if err := queue.Validation(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (conf *AWS) IsSNS() bool {
	return conf.SNS != nil
}

func (conf *AWS) IsSQS() bool {
	return conf.SQS != nil
}
