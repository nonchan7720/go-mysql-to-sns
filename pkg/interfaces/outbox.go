package interfaces

type Outbox struct {
	AggregateId string `json:"aggregate_id"`
	EventType   string `json:"event_type"`
	Payload     string `json:"payload"`
}
