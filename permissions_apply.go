package enthasura

import (
	"github.com/minskylab/hasura-api/metadata"
)

// func (r *Runtime) pgCreatePermission(
// 	operation HasuraOperation,
// 	perm map[string]interface{},
// 	tableName,
// 	role,
// 	sourceName,
// 	schemaName string,
// 	completeColumns ...string,
// ) (metadata.MetadataQuery, error) {
// 	if operation == pgCreateDeletePermission {
// 		return r.pgCreateDeletePermission(perm, tableName, role, sourceName, schemaName), nil
// 	}

// 	selectedColumns := []string{}

// 	if cols, ok := perm["columns"].([]string); ok {
// 		selectedColumns = append(selectedColumns, cols...)
// 	}

// 	if allColumnsFlag, isOk := perm["all_columns"].(bool); isOk && allColumnsFlag && len(completeColumns) > 0 {
// 		selectedColumns = completeColumns
// 	}

// 	excludedColumns, _ := perm["excluded_columns"].([]interface{})
// 	if len(excludedColumns) > 0 {
// 		// TODO: improve this shit
// 		for i, column := range selectedColumns {
// 			for _, excludedColumn := range excludedColumns {
// 				if column == excludedColumn.(string) {
// 					selectedColumns = append(selectedColumns[:i], selectedColumns[i+1:]...)
// 					break
// 				}
// 			}
// 		}
// 	}

// 	perm["columns"] = selectedColumns

// 	switch operation {
// 	case pgCreateInsertPermission:
// 		return r.pgCreateInsertPermission(perm, tableName, role, sourceName, schemaName), nil
// 	case pgCreateSelectPermission:
// 		return r.pgCreateSelectPermission(perm, tableName, role, sourceName, schemaName), nil
// 	case pgCreateUpdatePermission:
// 		return r.pgCreateUpdatePermission(perm, tableName, role, sourceName, schemaName), nil
// 	case pgCreateDeletePermission:
// 		return r.pgCreateDeletePermission(perm, tableName, role, sourceName, schemaName), nil
// 	}

// 	return metadata.MetadataQuery{}, fmt.Errorf("unsupported operation: %s", operation)
// }

func (r *Runtime) pgCreateInsertPermission(perm map[string]interface{}, tableName, roleName, sourceName, schemaName string) metadata.MetadataQuery {
	return metadata.PgCreateInsertPermissionQuery(&metadata.PgCreateInsertPermissionArgs{
		Permission: metadata.GenericPermission(perm),
		Table: metadata.QualifiedTableName{
			Name:   tableName,
			Schema: schemaName,
		},
		Role:   roleName,
		Source: sourceName,
	})
}

func (r *Runtime) pgCreateSelectPermission(perm map[string]interface{}, tableName, role, sourceName, schemaName string) metadata.MetadataQuery {
	return metadata.PgCreateSelectPermissionQuery(&metadata.PgCreateSelectPermissionArgs{
		Permission: metadata.GenericPermission(perm),
		Table: metadata.QualifiedTableName{
			Name:   tableName,
			Schema: schemaName,
		},
		Role:   role,
		Source: sourceName,
	})
}

func (r *Runtime) pgCreateUpdatePermission(perm map[string]interface{}, tableName, role, sourceName, schemaName string) metadata.MetadataQuery {
	return metadata.PgCreateUpdatePermissionQuery(&metadata.PgCreateUpdatePermissionArgs{
		Permission: metadata.GenericPermission(perm),
		Table: metadata.QualifiedTableName{
			Name:   tableName,
			Schema: schemaName,
		},
		Role:   role,
		Source: sourceName,
	})
}

func (r *Runtime) pgCreateDeletePermission(perm map[string]interface{}, tableName, role, sourceName, schemaName string) metadata.MetadataQuery {
	return metadata.PgCreateUpdatePermissionQuery(&metadata.PgCreateUpdatePermissionArgs{
		Permission: metadata.GenericPermission(perm),
		Table: metadata.QualifiedTableName{
			Name:   tableName,
			Schema: schemaName,
		},
		Role:   role,
		Source: sourceName,
	})
}
