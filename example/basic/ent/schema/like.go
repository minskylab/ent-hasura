package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	hasura "github.com/minskylab/ent-hasura"
	"github.com/sirupsen/logrus"
)

// Like holds the schema definition for the Like entity.
type Like struct {
	ent.Schema
}

func (l Like) Annotations() []schema.Annotation {
	logrus.Info("edges")
	logrus.Info(l.Edges()[0].Descriptor().Field)

	return []schema.Annotation{
		hasura.PermissionsRoleAnnotation{
			Role: "user",
			SelectPermission: &hasura.SelectPermission{
				Columns:        hasura.AllColumns,
				ComputedFields: []string{},
				Filter:         hasura.M{"creator": hasura.M{"id": hasura.Eq("X-Hasura-User-Id")}},
			},
			UpdatePermission: &hasura.UpdatePermission{
				Check:   hasura.M{"creator": hasura.M{"id": hasura.Eq("X-Hasura-User-Id")}},
				Filter:  hasura.M{"creator": hasura.M{"id": hasura.Eq("X-Hasura-User-Id")}},
				Columns: hasura.AllColumns,
			},
		},
	}
}

// Fields of the Like.
func (Like) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Like.
func (Like) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("creator", User.Type).Unique().Required().Ref("likes"),
	}
}
