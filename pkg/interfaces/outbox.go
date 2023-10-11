package interfaces

type Outbox struct {
	ID            int64  `json:"id"`
	AggregateType string `json:"aggregate_type"`
	AggregateId   string `json:"aggregate_id"`
	EventType     string `json:"event_type"`
	Payload       string `json:"payload"`
}
