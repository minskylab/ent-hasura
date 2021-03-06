// Code generated by entc, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/minskylab/ent-hasura/example/basic/ent/like"
	"github.com/minskylab/ent-hasura/example/basic/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	likeFields := schema.Like{}.Fields()
	_ = likeFields
	// likeDescCreatedAt is the schema descriptor for created_at field.
	likeDescCreatedAt := likeFields[1].Descriptor()
	// like.DefaultCreatedAt holds the default value on creation for the created_at field.
	like.DefaultCreatedAt = likeDescCreatedAt.Default.(func() time.Time)
}
