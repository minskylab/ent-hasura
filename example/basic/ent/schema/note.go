package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	hasura "github.com/minskylab/ent-hasura"
)

// Note holds the schema definition for the Note entity.
type Note struct {
	ent.Schema
}

// Fields of the Note.
func (Note) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.String("title"),
		field.String("content"),
	}
}

// Edges of the Note.
func (Note) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("authors", User.Type).Ref("notes"),
	}
}

func (n Note) Annotations() []schema.Annotation {
	return []schema.Annotation{
		hasura.PermissionsRoleAnnotation{
			Role: "user",
			SelectPermission: &hasura.PermissionSelect{
				Filter:            hasura.M{"authors": hasura.M{"user": hasura.M{"id": hasura.Eq("X-Hasura-User-Id")}}},
				AllColumns:        true,
				AllowAggregations: true,
			},
			UpdatePermission: &hasura.PermissionUpdate{
				Check:           hasura.M{"authors": hasura.M{"user": hasura.M{"id": hasura.Eq("X-Hasura-User-Id")}}},
				Filter:          hasura.M{"authors": hasura.M{"user": hasura.M{"id": hasura.Eq("X-Hasura-User-Id")}}},
				AllColumns:      true,
				ExcludedColumns: []string{"content"},
			},
		},
	}
}
