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
	mp := map[string]interface{}{}
	switch event {
	case Create:
		for key, value := range r.NewRow {
			mp[strings.ToLower(key)] = value
		}
	case Update:
		for key, value := range r.OldRow {
			mp[strings.ToLower(key)] = value
		}
	case Delete:
		for key, value := range r.OldRow {
			mp[strings.ToLower(key)] = value
		}
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

func (p *Payload) ToJson(rowIdx int) (string, error) {
	type payload struct {
		Event  Event      `json:"event"`
		Schema string     `json:"schema"`
		Table  string     `json:"table"`
		Row    PayloadRow `json:"row"`
	}
	value := payload{
		Event:  p.Event,
		Schema: p.Schema,
		Table:  p.Table,
		Row:    p.Rows[rowIdx],
	}
	if buf, err := json.Marshal(&value); err != nil {
		return "", err
	} else {
		return string(buf), nil
	}
}
