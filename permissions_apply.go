package hasura

const (
	pgCreateInsertPermission = "pg_create_insert_permission"
	pgCreateSelectPermission = "pg_create_select_permission"
	pgCreateUpdatePermission = "pg_create_update_permission"
	pgCreateDeletePermission = "pg_create_delete_permission"
)

type PGCreateInsertPermissionArgs struct {
	Table      string                 `json:"table"`
	Source     string                 `json:"source"`
	Role       string                 `json:"role"`
	Permission map[string]interface{} `json:"permission"`
}

type PGCreateSelectPermissionArgs struct {
	Table      string                 `json:"table"`
	Source     string                 `json:"source"`
	Role       string                 `json:"role"`
	Permission map[string]interface{} `json:"permission"`
}

type PGCreateUpdatePermissionArgs struct {
	Table      string                 `json:"table"`
	Source     string                 `json:"source"`
	Role       string                 `json:"role"`
	Permission map[string]interface{} `json:"permission"`
}

type PGCreateDeletePermissionArgs struct {
	Table      string                 `json:"table"`
	Source     string                 `json:"source"`
	Role       string                 `json:"role"`
	Permission map[string]interface{} `json:"permission"`
}

func (r *EphemeralRuntime) pgCreateInsertPermission(perm map[string]interface{}, tableName, role, source string) error {
	return r.genericHasuraMetadataQuery(ActionBody{
		Type: pgCreateInsertPermission,
		Args: PGCreateInsertPermissionArgs{
			Table:      tableName,
			Source:     source,
			Role:       role,
			Permission: perm,
		},
	})
}

func (r *EphemeralRuntime) pgCreateSelectPermission(perm map[string]interface{}, tableName, role, source string) error {
	return r.genericHasuraMetadataQuery(ActionBody{
		Type: pgCreateSelectPermission,
		Args: PGCreateSelectPermissionArgs{
			Table:      tableName,
			Source:     source,
			Role:       role,
			Permission: perm,
		},
	})
}

func (r *EphemeralRuntime) pgCreateUpdatePermission(perm map[string]interface{}, tableName, role, source string) error {
	return r.genericHasuraMetadataQuery(ActionBody{
		Type: pgCreateUpdatePermission,
		Args: PGCreateUpdatePermissionArgs{
			Table:      tableName,
			Source:     source,
			Role:       role,
			Permission: perm,
		},
	})
}

func (r *EphemeralRuntime) pgCreateDeletePermission(perm map[string]interface{}, tableName, role, source string) error {
	return r.genericHasuraMetadataQuery(ActionBody{
		Type: pgCreateDeletePermission,
		Args: PGCreateDeletePermissionArgs{
			Table:      tableName,
			Source:     source,
			Role:       role,
			Permission: perm,
		},
	})
}
