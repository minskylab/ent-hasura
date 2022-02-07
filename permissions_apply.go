package enthasura

import (
	"github.com/minskylab/hasura-api/metadata"
)

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
	return metadata.PgCreateDeletePermissionQuery(&metadata.PgCreateDeletePermissionArgs{
		Permission: metadata.GenericPermission(perm),
		Table: metadata.QualifiedTableName{
			Name:   tableName,
			Schema: schemaName,
		},
		Role:   role,
		Source: sourceName,
	})
}
