package mysql

import (
	"testing"
	"time"

	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestNextRetryDuration(t *testing.T) {
	var tables = []struct {
		name       string
		backoff    time.Duration
		retryCount int
		expect     time.Duration
	}{
		{
			name:       "Case1",
			backoff:    10 * time.Second,
			retryCount: 1,
			expect:     20 * time.Second,
		},
		{
			name:       "Case2",
			backoff:    10 * time.Second,
			retryCount: 2,
			expect:     40 * time.Second,
		},
		{
			name:       "Case3",
			backoff:    10 * time.Second,
			retryCount: 3,
			expect:     80 * time.Second,
		},
		{
			name:       "Case4",
			backoff:    3 * time.Second,
			retryCount: 1,
			expect:     6 * time.Second,
		},
		{
			name:       "Case5",
			backoff:    3 * time.Second,
			retryCount: 2,
			expect:     12 * time.Second,
		},
		{
			name:       "Case6",
			backoff:    3 * time.Second,
			retryCount: 3,
			expect:     24 * time.Second,
		},
	}
	for idx := range tables {
		tt := tables[idx]
		t.Run(tt.name, func(t *testing.T) {
			result := getNextRetryDuration(tt.backoff, tt.retryCount)
			assert.Equal(t, result, tt.expect)
		})
	}
}

func testMockEvents() []*outboxPollingEvent {
	return []*outboxPollingEvent{
		{
			ID:            1,
			AggregateType: "test-a",
			AggregateID:   "",
			Event:         "Create",
			Payload:       "{}",
		},
		{
			ID:            2,
			AggregateType: "test-a",
			AggregateID:   "",
			Event:         "Create",
			Payload:       "{}",
		},
		{
			ID:            3,
			AggregateType: "test-a.fifo",
			AggregateID:   "aaa-bbb",
			Event:         "Create",
			Payload:       "{}",
		},
		{
			ID:            4,
			AggregateType: "test-a",
			AggregateID:   "",
			Event:         "Create",
			Payload:       "{}",
		},
		{
			ID:            6,
			AggregateType: "test-a.fifo",
			AggregateID:   "aaa-bbb",
			Event:         "Delete",
			Payload:       "{}",
		},
		{
			ID:            5,
			AggregateType: "test-a.fifo",
			AggregateID:   "aaa-bbb",
			Event:         "Update",
			Payload:       "{}",
		},
	}
}

func testEventToMapProducerEvent() map[string][]*outboxPollingEvent {
	mpProducer := map[string]string{
		"test-a":      "arn:aws:sns:ap-northeast-1:000000000000:test-a",
		"test-a.fifo": "arn:aws:sns:ap-northeast-1:000000000000:test-a.fifo",
	}
	events := testMockEvents()
	vv := eventToMapProducerEvent(func(s string) (string, error) {
		v, ok := mpProducer[s]
		if ok {
			return v, nil
		}
		return "", config.ErrNotFoundProducer
	}, events)
	return vv
}

func TestEventToMapProducerEvent(t *testing.T) {
	mpEvents := testEventToMapProducerEvent()
	assert.Len(t, mpEvents, 2)
	assert.Len(t, mpEvents["arn:aws:sns:ap-northeast-1:000000000000:test-a"], 3)
	assert.Len(t, mpEvents["arn:aws:sns:ap-northeast-1:000000000000:test-a.fifo"], 3)
}

func TestEventToGroupingAggregateIdAndSort(t *testing.T) {
	mpEvents := testEventToMapProducerEvent()
	events := mpEvents["arn:aws:sns:ap-northeast-1:000000000000:test-a.fifo"]
	e := eventToGroupingAggregateId(events)
	assert.Len(t, e, 1)
	assert.Len(t, e["aaa-bbb"], 3)
	assert.Equal(t, e["aaa-bbb"][0].ID, int64(3))
	assert.Equal(t, e["aaa-bbb"][1].ID, int64(6))
	assert.Equal(t, e["aaa-bbb"][2].ID, int64(5))

	eventsSort(e["aaa-bbb"])
	assert.Equal(t, e["aaa-bbb"][0].ID, int64(3))
	assert.Equal(t, e["aaa-bbb"][1].ID, int64(5))
	assert.Equal(t, e["aaa-bbb"][2].ID, int64(6))
}
