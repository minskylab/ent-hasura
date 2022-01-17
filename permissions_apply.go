package hasura

import (
	"fmt"
)

type HasuraOperation string

const (
	pgCreateInsertPermission HasuraOperation = "pg_create_insert_permission"
	pgCreateSelectPermission HasuraOperation = "pg_create_select_permission"
	pgCreateUpdatePermission HasuraOperation = "pg_create_update_permission"
	pgCreateDeletePermission HasuraOperation = "pg_create_delete_permission"
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

func (r *EphemeralRuntime) pgCreateXPermission(
	operation HasuraOperation,
	perm map[string]interface{},
	tableName,
	role,
	source string,
	completeColumns ...string,
) error {
	if operation == pgCreateDeletePermission {
		return r.pgCreateDeletePermission(perm, tableName, role, source)
	}

	selectedColumns := []string{}

	if allColumnsFlag, isOk := perm["all_columns"].(bool); isOk && allColumnsFlag && len(completeColumns) > 0 {
		selectedColumns = completeColumns
	}

	excludedColumns, _ := perm["excluded_columns"].([]interface{})
	if len(excludedColumns) > 0 {
		// TODO: improve this shit
		for i, column := range selectedColumns {
			for _, excludedColumn := range excludedColumns {
				if column == excludedColumn.(string) {
					selectedColumns = append(selectedColumns[:i], selectedColumns[i+1:]...)
					break
				}
			}
		}
	}

	perm["columns"] = selectedColumns

	switch operation {
	case pgCreateInsertPermission:
		return r.pgCreateInsertPermission(perm, tableName, role, source)
	case pgCreateSelectPermission:
		return r.pgCreateSelectPermission(perm, tableName, role, source)
	case pgCreateUpdatePermission:
		return r.pgCreateUpdatePermission(perm, tableName, role, source)
	case pgCreateDeletePermission:
		return r.pgCreateDeletePermission(perm, tableName, role, source)
	default:
		return fmt.Errorf("unsupported operation: %s", operation)
	}
}

func (r *EphemeralRuntime) pgCreateInsertPermission(perm map[string]interface{}, tableName, role, source string) error {
	return r.genericHasuraMetadataQuery(ActionBody{
		Type: string(pgCreateInsertPermission),
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
		Type: string(pgCreateSelectPermission),
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
		Type: string(pgCreateUpdatePermission),
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
		Type: string(pgCreateDeletePermission),
		Args: PGCreateDeletePermissionArgs{
			Table:      tableName,
			Source:     source,
			Role:       role,
			Permission: perm,
		},
	})
}
