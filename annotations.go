package hasura

import "entgo.io/ent"

const (
	hasuraPermissionsAnnotationName     = "hasura-permissions"
	hasuraPermissionsRoleAnnotationName = "hasura-permissions-role"
	hasuraNotInheritedAnnotationName    = "hasura-not-inherited"
)

type M map[string]interface{}

func Eq(val string) M {
	return M{
		"_eq": val,
	}
}

type A []string

type PermissionsAnnotation struct {
	InsertPermissions []*InsertPermission `json:"insert_permissions,omitempty"`
	SelectPermissions []*SelectPermission `json:"select_permissions,omitempty"`
	UpdatePermissions []*UpdatePermission `json:"update_permissions,omitempty"`
	DeletePermissions []*DeletePermission `json:"delete_permissions,omitempty"`
}

type PermissionsRoleAnnotation struct {
	Role             string            `json:"role"`
	InsertPermission *PermissionInsert `json:"insert_permission,omitempty"`
	SelectPermission *PermissionSelect `json:"select_permission,omitempty"`
	UpdatePermission *PermissionUpdate `json:"update_permission,omitempty"`
	DeletePermission *PermissionDelete `json:"delete_permission,omitempty"`
}

type NotInheritedPermissionsAnnotation struct{}

func (PermissionsAnnotation) Name() string {
	return hasuraPermissionsAnnotationName
}

func (PermissionsRoleAnnotation) Name() string {
	return hasuraPermissionsRoleAnnotationName
}

func (NotInheritedPermissionsAnnotation) Name() string {
	return hasuraNotInheritedAnnotationName
}

func AllFields(schema ent.Interface) []string {
	allFields := []string{}

	for _, f := range schema.Fields() {
		allFields = append(allFields, f.Descriptor().Name)
	}

	return allFields
}
