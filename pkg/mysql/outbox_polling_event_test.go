package mysql

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOutboxPollingEvent(t *testing.T) {
	v, _ := time.Parse("2006-01-02", "2023-10-16")
	before, _ := time.Parse("2006-01-02", "2023-10-15")
	after := time.Now()
	event := outboxPollingEvent{
		RetryAt: &v,
	}
	result := event.CanNotRetry(before)
	assert.True(t, result)
	result = event.CanNotRetry(after)
	assert.False(t, result)
}

func TestOutboxPollingRetryCount(t *testing.T) {
	event := outboxPollingEvent{
		RetryCount: 0,
	}
	assert.False(t, event.CheckMaxRetryCount(1))
	event.IncrementToRetryCount()
	assert.True(t, event.CheckMaxRetryCount(1))
}
