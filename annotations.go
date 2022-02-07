package enthasura

import (
	"github.com/minskylab/hasura-api/metadata"
)

const (
	// hasuraPermissionsAnnotationName     = "hasura-permissions"
	hasuraPermissionsRoleAnnotationName = "hasura-permissions-role"
	// hasuraNotInheritedAnnotationName    = "hasura-not-inherited"
)

type M map[string]interface{}

func Eq(val string) M {
	return M{
		"_eq": val,
	}
}

type A []string

// type PermissionsRoleAnnotation struct {
// 	Role             string            `json:"role"`
// 	InsertPermission *PermissionInsertAnnotation `json:"insert_permission,omitempty"`
// 	SelectPermission *PermissionSelectAnnotation `json:"select_permission,omitempty"`
// 	UpdatePermission *PermissionUpdateAnnotation `json:"update_permission,omitempty"`
// 	DeletePermission *PermissionDeleteAnnotation `json:"delete_permission,omitempty"`
// }

type PermissionsRoleAnnotation struct {
	Role             string            `json:"role"`
	InsertPermission *InsertPermission `json:"insert_permission,omitempty"`
	SelectPermission *SelectPermission `json:"select_permission,omitempty"`
	UpdatePermission *UpdatePermission `json:"update_permission,omitempty"`
	DeletePermission *DeletePermission `json:"delete_permission,omitempty"`
}

type NotInheritedPermissionsAnnotation struct{}

// func (PermissionsAnnotation) Name() string {
// 	return hasuraPermissionsAnnotationName
// }

func (PermissionsRoleAnnotation) Name() string {
	return hasuraPermissionsRoleAnnotationName
}

// func (NotInheritedPermissionsAnnotation) Name() string {
// 	return hasuraNotInheritedAnnotationName
// }

type InsertPermission struct {
	Check       M                          `json:"check"`
	Set         M                          `json:"set,omitempty"`
	Columns     metadata.PermissionColumns `json:"columns,omitempty"`
	BackendOnly bool                       `json:"backend_only,omitempty"`
}

type SelectPermission struct {
	Columns           metadata.PermissionColumns `json:"columns,omitempty"`
	ComputedFields    []string                   `json:"computed_fields,omitempty"`
	Filter            M                          `json:"filter"`
	Limit             int                        `json:"limit,omitempty"`
	AllowAggregations bool                       `json:"allow_aggregations,omitempty"`
}

type UpdatePermission struct {
	Columns metadata.PermissionColumns `json:"columns"`
	Filter  M                          `json:"filter"`
	Check   M                          `json:"check,omitempty"`
	Set     M                          `json:"set,omitempty"`
}

type DeletePermission struct {
	Filter M `json:"filter"`
}

type AllColumnsType string

func (s AllColumnsType) GetColumns() interface{} {
	return s
}

const AllColumns AllColumnsType = "*"

// func AllColumns() metadata.PermissionColumns {
// 	return allColumns
// }

func Columns(cols ...string) metadata.PGColumns {
	return cols
}
