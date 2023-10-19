package schema

import (
	"encoding/json"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type JSON map[string]interface{}

func (mp *JSON) ToJson() (string, error) {
	buf, err := json.Marshal(mp)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// Outbox holds the schema definition for the Outbox entity.
type Outbox struct {
	ent.Schema
}

// Fields of the Outbox.
func (Outbox) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("aggregate_type").
			SchemaType(
				map[string]string{
					dialect.MySQL: "varchar(255)",
				},
			),
		field.String("aggregate_id").
			SchemaType(map[string]string{
				dialect.MySQL: "varchar(128)",
			}),
		field.String("event").
			SchemaType(map[string]string{
				dialect.MySQL: "varchar(255)",
			}),
		field.Bytes("payload").
			SchemaType(map[string]string{
				dialect.MySQL: "json",
			}),
		field.Time("retry_at").
			Nillable(),
		field.Int("retry_count").
			Nillable(),
	}
}

// Edges of the Outbox.
func (Outbox) Edges() []ent.Edge {
	return nil
}

func (Outbox) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "outbox"},
	}
}
