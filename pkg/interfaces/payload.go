package interfaces

import "encoding/json"

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

func (p *Payload) ToJson() (string, error) {
	if buf, err := json.Marshal(p); err != nil {
		return "", err
	} else {
		return string(buf), nil
	}
}
