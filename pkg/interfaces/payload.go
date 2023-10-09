package interfaces

import (
	"encoding/json"
	"strings"
)

type Event int

const (
	Create Event = iota
	Update
	Delete
)

var (
	mpEvent = map[Event]string{
		Create: "CREATE",
		Update: "UPDATE",
		Delete: "DELETE",
	}
)

func (e Event) String() string {
	return mpEvent[e]
}

func (e Event) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

type Payload struct {
	Event  Event        `json:"event"`
	Schema string       `json:"schema"`
	Table  string       `json:"table"`
	Rows   []PayloadRow `json:"rows"`
}

type PayloadRow struct {
	OldRow Row `json:"old_row"`
	NewRow Row `json:"new_row"`
}

func (r *PayloadRow) MainRow(event Event) Row {
	var mp map[string]interface{}
	switch event {
	case Create:
		mp = map[string]interface{}{}
		for key, value := range r.NewRow {
			mp[strings.ToLower(key)] = value
		}
	case Update:
		mp = map[string]interface{}{}
		for key, value := range r.OldRow {
			mp[strings.ToLower(key)] = value
		}
	case Delete:
		mp = map[string]interface{}{}
		for key, value := range r.OldRow {
			mp[strings.ToLower(key)] = value
		}
	}
	if mp != nil {
		mp["payload_event"] = event.String()
	}
	return mp
}

func NewPayloadRow(oldRow Row, newRow Row) PayloadRow {
	if oldRow == nil {
		oldRow = Row{}
	}
	if newRow == nil {
		newRow = Row{}
	}
	return PayloadRow{
		OldRow: oldRow,
		NewRow: newRow,
	}
}

type Row map[string]interface{}

type Column struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func (p *Payload) SendPayload(rowIdx int) SendPayload {
	return SendPayload{
		Event:  p.Event,
		Schema: p.Schema,
		Table:  p.Table,
		Row:    p.Rows[rowIdx],
	}
}

type SendPayload struct {
	Event  Event      `json:"event"`
	Schema string     `json:"schema"`
	Table  string     `json:"table"`
	Row    PayloadRow `json:"row"`
}

func (p *SendPayload) ToJson() (string, error) {
	if buf, err := json.Marshal(p); err != nil {
		return "", err
	} else {
		return string(buf), nil
	}
}

func (p *SendPayload) GetSchema() string {
	return strings.ToLower(p.Schema)
}

func (p *SendPayload) GetTable() string {
	return strings.ToLower(p.Table)
}
