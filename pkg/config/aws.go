package config

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type IEndpoint struct {
	Endpoint string `yaml:"endpoint"`
}

type SNS struct {
	IEndpoint `yaml:",inline"`
	Topics    []Topic `yaml:"topics"`
}

type AWS struct {
	AccessKey string `yaml:"accessKey"`
	SecretKet string `yaml:"secretKey"`
	SNS       SNS    `yaml:"sns"`
}

func (conf *AWS) WithStatic() aws.CredentialsProvider {
	if conf.AccessKey == "" && conf.SecretKet == "" {
		return nil
	}
	return credentials.NewStaticCredentialsProvider(conf.AccessKey, conf.SecretKet, "")
}

func (conf *AWS) Validation() error {
	for _, topic := range conf.SNS.Topics {
		if err := topic.Validation(); err != nil {
			return err
		}
	}
	return nil
}

type Topic struct {
	TableName              string `yaml:"tableName"`
	TopicArn               string `yaml:"topicArn"`
	MessageGroupIdTemplate string `yaml:"messageGroupIdTemplate"`
}

func (t *Topic) GetMessageGroupId(mp map[string]interface{}) *string {
	if !t.IsFIFO() {
		return nil
	}
	value := templateExecute(t.MessageGroupIdTemplate, mp)
	return &value
}

func (t *Topic) IsFIFO() bool {
	return strings.HasSuffix(t.TopicArn, ".fifo")
}

func (t *Topic) Validation() error {
	if t.IsFIFO() && t.MessageGroupIdTemplate == "" {
		return errors.New("For FIFO topics, MessageGroupIdTemplate must be set.")
	}
	return nil
}
