// Code generated by MockGen. DO NOT EDIT.
// Source: aws.go

// Package aws is a generated GoMock package.
package aws

import (
	context "context"
	reflect "reflect"

	sns "github.com/aws/aws-sdk-go-v2/service/sns"
	sqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	gomock "github.com/golang/mock/gomock"
)

// MockSNSClient is a mock of SNSClient interface.
type MockSNSClient struct {
	ctrl     *gomock.Controller
	recorder *MockSNSClientMockRecorder
}

// MockSNSClientMockRecorder is the mock recorder for MockSNSClient.
type MockSNSClientMockRecorder struct {
	mock *MockSNSClient
}

// NewMockSNSClient creates a new mock instance.
func NewMockSNSClient(ctrl *gomock.Controller) *MockSNSClient {
	mock := &MockSNSClient{ctrl: ctrl}
	mock.recorder = &MockSNSClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSNSClient) EXPECT() *MockSNSClientMockRecorder {
	return m.recorder
}

// Publish mocks base method.
func (m *MockSNSClient) Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Publish", varargs...)
	ret0, _ := ret[0].(*sns.PublishOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Publish indicates an expected call of Publish.
func (mr *MockSNSClientMockRecorder) Publish(ctx, params interface{}, optFns ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockSNSClient)(nil).Publish), varargs...)
}

// MockSQSClient is a mock of SQSClient interface.
type MockSQSClient struct {
	ctrl     *gomock.Controller
	recorder *MockSQSClientMockRecorder
}

// MockSQSClientMockRecorder is the mock recorder for MockSQSClient.
type MockSQSClientMockRecorder struct {
	mock *MockSQSClient
}

// NewMockSQSClient creates a new mock instance.
func NewMockSQSClient(ctrl *gomock.Controller) *MockSQSClient {
	mock := &MockSQSClient{ctrl: ctrl}
	mock.recorder = &MockSQSClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSQSClient) EXPECT() *MockSQSClientMockRecorder {
	return m.recorder
}

// GetQueueUrl mocks base method.
func (m *MockSQSClient) GetQueueUrl(ctx context.Context, params *sqs.GetQueueUrlInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetQueueUrl", varargs...)
	ret0, _ := ret[0].(*sqs.GetQueueUrlOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQueueUrl indicates an expected call of GetQueueUrl.
func (mr *MockSQSClientMockRecorder) GetQueueUrl(ctx, params interface{}, optFns ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQueueUrl", reflect.TypeOf((*MockSQSClient)(nil).GetQueueUrl), varargs...)
}

// SendMessage mocks base method.
func (m *MockSQSClient) SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SendMessage", varargs...)
	ret0, _ := ret[0].(*sqs.SendMessageOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockSQSClientMockRecorder) SendMessage(ctx, params interface{}, optFns ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockSQSClient)(nil).SendMessage), varargs...)
}
