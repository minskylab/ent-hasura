package permissions

import "entgo.io/ent/schema"

type InsertPermissionAnnotation struct{}

func (ann *InsertPermissionAnnotation) Name() string {
	return ""
}

func annotation() schema.Annotation {
	return &InsertPermissionAnnotation{}
}
