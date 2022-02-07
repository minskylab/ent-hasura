package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	hasura "github.com/minskylab/ent-hasura"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.String("email").Unique(),
		field.String("name"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("notes", Note.Type),
		edge.To("likes", Like.Type),
	}
}

func (u User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		hasura.PermissionsRoleAnnotation{
			Role: "user",
			SelectPermission: &hasura.SelectPermission{
				Columns:           hasura.AllColumns,
				Filter:            hasura.M{"id": hasura.Eq("X-Hasura-User-Id")},
				AllowAggregations: true,
			},
			UpdatePermission: &hasura.UpdatePermission{
				Columns: hasura.AllColumns,
				Check:   hasura.M{"id": hasura.Eq("X-Hasura-User-Id")},
				Filter:  hasura.M{"id": hasura.Eq("X-Hasura-User-Id")},
			},
		},
	}
}
