package schema

import (
	"basket-buddy-backend/model"
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
		field.Time("expiration"),
		field.String("share_code").
			Unique(),
		field.JSON("data", []map[string]interface{}(nil)),
		field.UUID("creator_id", uuid.UUID{}),
		field.String("status").
			Default(model.ShareStatusCreated),
	}
}

// Edges of the Share.
func (Share) Edges() []ent.Edge {
	return nil
}
