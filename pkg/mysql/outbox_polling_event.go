package mysql

import (
	"time"

	"github.com/nonchan7720/go-storage-to-messenger/pkg/ent"
)

type outboxPollingEvent struct {
	ID            int64
	AggregateType string
	AggregateID   string
	Event         string
	Payload       string
	RetryAt       *time.Time
	RetryCount    int
}

func newOutboxEvent(v *ent.Outbox) *outboxPollingEvent {
	return &outboxPollingEvent{
		ID:            v.ID,
		AggregateType: v.AggregateType,
		AggregateID:   v.AggregateID,
		Event:         v.Event,
		Payload:       string(v.Payload),
		RetryAt:       v.RetryAt,
		RetryCount:    v.RetryCount,
	}
}

func (o *outboxPollingEvent) IncrementToRetryCount() {
	o.RetryCount++
}

func (o *outboxPollingEvent) CheckMaxRetryCount(retryCount int) bool {
	return o.RetryCount >= retryCount
}

func (o *outboxPollingEvent) CanNotRetry(value time.Time) bool {
	return o.RetryAt != nil && o.RetryAt.After(value)
}
