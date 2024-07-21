package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Share holds the schema definition for the Share entity.
type Share struct {
	ent.Schema
}

// Fields of the Share.
func (Share) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Time("created_at").
			Default(time.Now),
		field.Time("expiration").
			Default(time.Now().Add(15 * time.Minute)),
		field.String("share_code").
			Unique(),
		field.JSON("data", []map[string]interface{}(nil)),
		field.UUID("creator_id", uuid.UUID{}),
	}
}

// Edges of the Share.
func (Share) Edges() []ent.Edge {
	return nil
}
