package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
)

type JSON map[string]interface{}

// Outbox holds the schema definition for the Outbox entity.
type Outbox struct {
	ent.Schema
}

// Fields of the Outbox.
func (Outbox) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("aggregate_type").SchemaType(
			map[string]string{
				dialect.MySQL: "varchar(255)",
			},
		),
		field.String("aggregate_id").SchemaType(map[string]string{
			dialect.MySQL: "varchar(128)",
		}),
		field.String("event").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)",
		}),
		field.JSON("payload", &JSON{}),
		field.Time("retry_at").Nillable(),
		field.Int("retry_count").Nillable(),
	}
}

// Edges of the Outbox.
func (Outbox) Edges() []ent.Edge {
	return nil
}
