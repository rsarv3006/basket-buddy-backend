package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// AppUser holds the schema definition for the AppUser entity.
type AppUser struct {
	ent.Schema
}

// Fields of the AppUser.
func (AppUser) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("username").
			MaxLen(255).
			NotEmpty(),
		field.String("email").
			MaxLen(255).
			NotEmpty(),
		field.String("role").
			MaxLen(255).
			NotEmpty(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Bool("is_active").
			Default(false),
	}
}

// Edges of the AppUser.
func (AppUser) Edges() []ent.Edge {
	return nil
}
