package hasura

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/entc/gen"
	"github.com/pkg/errors"
)

func minimalInsertPermission(role string, columns []string) *InsertPermission {
	return &InsertPermission{
		Role: role,
		Permission: PermissionInsert{
			Check:       map[string]interface{}{},
			Set:         map[string]interface{}{},
			Columns:     columns,
			BackendOnly: false,
		},
	}
}

func minimalSelectPermission(role string, columns []string) *SelectPermission {
	return &SelectPermission{
		Role: role,
		Permission: PermissionSelect{
			Filter:            map[string]interface{}{},
			Columns:           columns,
			AllowAggregations: true,
		},
	}
}

func minimalUpdatePermission(role string, columns []string) *UpdatePermission {
	return &UpdatePermission{
		Role: role,
		Permission: PermissionUpdate{
			Check:   map[string]interface{}{},
			Filter:  map[string]interface{}{},
			Columns: columns,
		},
	}
}

func minimalDeletePermission(role string, columns []string) *DeletePermission {
	return &DeletePermission{
		Role: role,
		Permission: PermissionDelete{
			Filter: map[string]interface{}{},
		},
	}
}

func enhanceHasuraTableWithPermissions(source *Source, table *TableDefinition) {
	for i, iTable := range source.Tables {
		if iTable.Table.Schema == table.Table.Schema && iTable.Table.Name == table.Table.Name {
			source.Tables[i].InsertPermissions = table.InsertPermissions
			source.Tables[i].UpdatePermissions = table.UpdatePermissions
			source.Tables[i].SelectPermissions = table.SelectPermissions
			source.Tables[i].DeletePermissions = table.DeletePermissions

			break
		}
	}
}

func hasuraPermissionFromRelationalTable(table *schema.Table, roleName, schemaName string) (*TableDefinition, error) {
	columns := []string{}

	for _, column := range table.Columns {
		columns = append(columns, column.Name)
	}

	definition := &TableDefinition{Table: Table{Schema: schemaName, Name: table.Name}}

	definition.InsertPermissions = append(definition.InsertPermissions, minimalInsertPermission(roleName, columns))
	definition.UpdatePermissions = append(definition.UpdatePermissions, minimalUpdatePermission(roleName, columns))
	definition.SelectPermissions = append(definition.SelectPermissions, minimalSelectPermission(roleName, columns))
	definition.DeletePermissions = append(definition.DeletePermissions, minimalDeletePermission(roleName, columns))

	return definition, nil
}

func obtainPermissionsTableFromEntSchema(schema *gen.Graph, roleName, schemaName string) ([]*TableDefinition, error) {
	tableDefinitions := []*TableDefinition{}

	tables, err := schema.Tables()
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain tables from schema")
	}

	for _, table := range tables {
		definition, err := hasuraPermissionFromRelationalTable(table, roleName, schemaName)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		tableDefinitions = append(tableDefinitions, definition)
	}

	return tableDefinitions, nil
}

func enhancedHasuraPermissions(initial *HasuraMetadata, schema *gen.Graph, sourceName, roleName, schemaName string) error {
	initial.ResourceVersion += 1

	tables, err := obtainPermissionsTableFromEntSchema(schema, roleName, schemaName)
	if err != nil {
		return err
	}

	for _, source := range initial.Metadata.Sources {
		if source.Name == sourceName {
			for _, table := range tables {
				enhanceHasuraTableWithPermissions(source, table)
			}
		}
	}

	return nil
}
