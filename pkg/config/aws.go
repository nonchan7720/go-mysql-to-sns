package config

import (
	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type IEndpoint struct {
	Endpoint string `yaml:"endpoint"`
}

func (e IEndpoint) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Endpoint, validation.When(e.Endpoint != "", is.URL)),
	)
}

type AWS struct {
	AccessKey string `yaml:"accessKey"`
	SecretKet string `yaml:"secretKey"`
	SNS       *SNS   `yaml:"sns"`
	SQS       *SQS   `yaml:"sqs"`
}

var (
	_ validation.Validatable = (*AWS)(nil)
)

func (conf *AWS) WithStatic() awsv2.CredentialsProvider {
	if conf.AccessKey == "" && conf.SecretKet == "" {
		return nil
	}
	return credentials.NewStaticCredentialsProvider(conf.AccessKey, conf.SecretKet, "")
}

func (conf *AWS) IsSNS() bool {
	return conf.SNS != nil
}

func (conf *AWS) IsSQS() bool {
	return conf.SQS != nil
}

func (conf AWS) Validate() error {
	return validation.ValidateStruct(&conf,
		validation.Field(&conf.SNS),
		validation.Field(&conf.SQS),
	)
}
