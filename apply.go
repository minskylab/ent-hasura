package hasura

import (
	"strings"

	"entgo.io/ent/entc/gen"
	"github.com/iancoleman/strcase"
	"github.com/minskylab/hasura-api/metadata"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (r *Runtime) tableNameFromDefinition(table metadata.TableDefinition) (string, error) {
	tableName := ""
	switch tName := table.Table.(type) {
	case metadata.TableName:
		tableName = string(tName)
	case metadata.QualifiedTableName:
		tableName = string(tName.Name)
	default:
		return "", errors.Errorf("unexpected type for table name: %T", tName)
	}

	return tableName, nil
}

func (r *Runtime) clearMetadata() error {
	_, err := r.hasura.Metadata.ClearMetadata(&metadata.ClearMetadataArgs{})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Runtime) PerformPrelude(graph *gen.Graph, sourceName, schemaName string, clearMetadata bool) error {
	allTables, err := graph.Tables()
	if err != nil {
		return errors.WithStack(err)
	}

	if clearMetadata {
		if err := r.clearMetadata(); err != nil {
			return errors.WithStack(err)
		}

		return nil
	}

	untrackBatch := []metadata.MetadataQuery{}
	for _, table := range allTables {
		untrackBatch = append(untrackBatch, metadata.PgUntrackTableQuery(&metadata.PgUntrackTableArgs{
			Table: metadata.QualifiedTableName{
				Name:   table.Name,
				Schema: schemaName,
			},
			Cascade: true,
			Source:  sourceName,
		}))
	}

	res, err := r.hasura.Metadata.Bulk(untrackBatch)
	if err != nil {
		return errors.WithStack(err)
	}

	return logAndResponseMetadataResponse(res)
}

func (r *Runtime) TrackAllTables(graph *gen.Graph, sourceName, schemaName string) error {
	allTables, err := graph.Tables()
	if err != nil {
		return errors.WithStack(err)
	}

	trackBatch := []metadata.MetadataQuery{}
	for _, table := range allTables {
		trackBatch = append(trackBatch, metadata.PgTrackTableQuery(&metadata.PgTrackTableArgs{
			Table: metadata.QualifiedTableName{
				Name:   table.Name,
				Schema: schemaName,
			},
			Source: sourceName,
		}))
	}

	res, err := r.hasura.Metadata.Bulk(trackBatch)
	if err != nil {
		return errors.WithStack(err)
	}

	return logAndResponseMetadataResponse(res)
}

func (r *Runtime) CustomizeAllTables(graph *gen.Graph, sourceName, schemaName string) error {
	tables, err := obtainHasuraTablesFromEntSchema(graph, schemaName)
	if err != nil {
		return errors.WithStack(err)
	}

	tableCustomizeBatch := []metadata.MetadataQuery{}
	objectRelationBulk := []metadata.MetadataQuery{}
	arrayRelationBulk := []metadata.MetadataQuery{}

	for _, def := range tables {
		tableName, err := r.tableNameFromDefinition(*def)
		if err != nil {
			return errors.WithStack(err)
		}

		tableCustomizeBatch = append(tableCustomizeBatch, metadata.PgSetTableCustomizationQuery(&metadata.PgSetTableCustomizationArgs{
			Table: metadata.QualifiedTableName{
				Name:   tableName,
				Schema: schemaName,
			},
			Source:        sourceName,
			Configuration: def.Configuration,
		}))

		for _, rel := range def.ObjectRelationships {
			objectRelationBulk = append(objectRelationBulk, metadata.PgCreateObjectRelationshipQuery(&metadata.PgCreateObjectRelationshipArgs{
				Table: metadata.QualifiedTableName{
					Name:   tableName,
					Schema: schemaName,
				},
				Name:   rel.Name,
				Source: sourceName,
				Using: metadata.ObjRelUsing{
					ForeignKeyConstraintOn: rel.Using.ForeignKeyConstraintOn.(metadata.IObjRelUsingChoice),
				},
			}))
		}

		for _, rel := range def.ArrayRelationships {
			arrayRelationBulk = append(arrayRelationBulk, metadata.PgCreateArrayRelationshipQuery(&metadata.PgCreateArrayRelationshipArgs{
				Table: metadata.QualifiedTableName{
					Name:   tableName,
					Schema: schemaName,
				},
				Name:   rel.Name,
				Source: sourceName,
				Using: metadata.ArrRelUsing{
					ForeignKeyConstraintOn: &metadata.ArrRelUsingFKeyOn{
						Table:  rel.Using.ForeignKeyConstraintOn.Table,
						Column: rel.Using.ForeignKeyConstraintOn.Column,
					},
				},
			}))
		}
	}

	res, err := r.hasura.Metadata.Bulk(tableCustomizeBatch)
	if err != nil {
		return errors.WithStack(err)
	}

	logAndResponseMetadataResponse(res)

	res, err = r.hasura.Metadata.Bulk(objectRelationBulk)
	if err != nil {
		return errors.WithStack(err)
	}

	logAndResponseMetadataResponse(res)

	res, err = r.hasura.Metadata.Bulk(arrayRelationBulk)
	if err != nil {
		return errors.WithStack(err)
	}

	logAndResponseMetadataResponse(res)

	return nil
}

////////////////////////////////

////////////////////////////////

func (r *Runtime) ApplyPGPermissionsForAllTables(graph *gen.Graph, schemaName, sourceName string) error {
	insertPermissionBulk := []metadata.MetadataQuery{}
	selectPermissionBulk := []metadata.MetadataQuery{}
	updatePermissionBulk := []metadata.MetadataQuery{}
	deletePermissionBulk := []metadata.MetadataQuery{}

	nodeTables := []string{} // "permissions"

	for _, n := range graph.Nodes {
		nodeTables = append(nodeTables, n.Table())
	}

	for _, node := range graph.Nodes {
		permAnn, isOk := node.Annotations[hasuraPermissionsRoleAnnotationName].(map[string]interface{})
		if !isOk {
			logrus.Debug("skipping node: ", node.Name, " as it does not have permissions annotation")
			continue
		}

		roleName, isOk := permAnn["role"].(string)
		if !isOk {
			logrus.Debug("skipping node: ", node.Name, " as it does not have permissions role name in annotation")
			continue
		}

		if insertPermission, isOk := permAnn["insert_permission"].(map[string]interface{}); isOk {
			logrus.Info("creating insert permission for table: ", node.Table(), " with role: ", roleName)

			insertPermissionBulk = append(
				insertPermissionBulk,
				r.pgCreateInsertPermission(insertPermission, node.Table(), roleName, schemaName, sourceName),
			)

			insertPermissionBulk = append(
				insertPermissionBulk,
				r.createInsertPermissionForEdges(nodeTables, node, permAnn, roleName, sourceName, schemaName)...,
			)
		}

		if selectPermission, isOk := permAnn["select_permission"].(map[string]interface{}); isOk {
			logrus.Info("creating select permission for table: ", node.Table(), " with role: ", roleName)

			selectPermissionBulk = append(
				selectPermissionBulk,
				r.pgCreateSelectPermission(selectPermission, node.Table(), roleName, schemaName, sourceName),
			)

			selectPermissionBulk = append(
				selectPermissionBulk,
				r.createSelectPermissionForEdges(nodeTables, node, permAnn, roleName, sourceName, schemaName)...,
			)
		}

		if updatePermission, isOk := permAnn["update_permission"].(map[string]interface{}); isOk {
			logrus.Info("creating update permission for table: ", node.Table(), " with role: ", roleName)

			updatePermissionBulk = append(
				updatePermissionBulk,
				r.pgCreateUpdatePermission(updatePermission, node.Table(), roleName, schemaName, sourceName),
			)

			updatePermissionBulk = append(
				updatePermissionBulk,
				r.createUpdatePermissionForEdges(nodeTables, node, permAnn, roleName, sourceName, schemaName)...,
			)

		}

		if deletePermission, isOk := permAnn["delete_permission"].(map[string]interface{}); isOk {
			logrus.Info("creating delete permission for table: ", node.Table(), " with role: ", roleName)

			deletePermissionBulk = append(
				deletePermissionBulk,
				r.pgCreateDeletePermission(deletePermission, node.Table(), roleName, schemaName, sourceName),
			)

			deletePermissionBulk = append(
				deletePermissionBulk,
				r.createDeletePermissionForEdges(nodeTables, node, permAnn, roleName, sourceName, schemaName)...,
			)
		}
	}

	return nil
}

// func (r *Runtime) pgCreateAllXPermissionforNode(op HasuraOperation, nodeTables []string, perm map[string]interface{}, node *gen.Type, role string, sourceName string, allColumns []string) error {
// 	if err := r.pgCreateXPermission(op, perm, node.Table(), role, sourceName, allColumns...); err != nil {
// 		return errors.WithStack(err)
// 	}

// 	return r.createXPermissionForEdges(op, nodeTables, node, perm, role, sourceName)
// }

func (r *Runtime) isNodeTable(nodeTables []string, tableName string) bool {
	for _, nodeTable := range nodeTables {
		if nodeTable == tableName {
			return true
		}
	}

	return false
}

func (r *Runtime) createInsertPermissionForEdges(nodeTables []string, node *gen.Type, permission map[string]interface{}, role string, sourceName, schemaName string) []metadata.MetadataQuery {
	bulkEdgePermissions := []metadata.MetadataQuery{}

	for _, edge := range node.Edges {
		if !edge.IsInverse() && !edge.OwnFK() {
			tableName := edge.Rel.Table
			if r.isNodeTable(nodeTables, tableName) {
				continue
			}

			tableName, newPermission := r.tableAndPermissionsFromEdge(edge, nodeTables, permission)

			logrus.Info("creating [edge] insert permission for table: ", tableName, " with role: ", role)

			bulkEdgePermissions = append(bulkEdgePermissions, r.pgCreateInsertPermission(newPermission, tableName, role, schemaName, sourceName))
		}
	}

	return bulkEdgePermissions
}

func (r *Runtime) createSelectPermissionForEdges(nodeTables []string, node *gen.Type, permission map[string]interface{}, role string, sourceName, schemaName string) []metadata.MetadataQuery {
	bulkEdgePermissions := []metadata.MetadataQuery{}

	for _, edge := range node.Edges {
		if !edge.IsInverse() && !edge.OwnFK() {
			tableName := edge.Rel.Table
			if r.isNodeTable(nodeTables, tableName) {
				continue
			}

			tableName, newPermission := r.tableAndPermissionsFromEdge(edge, nodeTables, permission)

			logrus.Info("creating [edge] insert permission for table: ", tableName, " with role: ", role)

			bulkEdgePermissions = append(bulkEdgePermissions, r.pgCreateSelectPermission(newPermission, tableName, role, schemaName, sourceName))
		}
	}

	return bulkEdgePermissions
}

func (r *Runtime) createUpdatePermissionForEdges(nodeTables []string, node *gen.Type, permission map[string]interface{}, role string, sourceName, schemaName string) []metadata.MetadataQuery {
	bulkEdgePermissions := []metadata.MetadataQuery{}

	for _, edge := range node.Edges {
		if !edge.IsInverse() && !edge.OwnFK() {
			tableName := edge.Rel.Table
			if r.isNodeTable(nodeTables, tableName) {
				continue
			}

			tableName, newPermission := r.tableAndPermissionsFromEdge(edge, nodeTables, permission)

			logrus.Info("creating [edge] insert permission for table: ", tableName, " with role: ", role)

			bulkEdgePermissions = append(bulkEdgePermissions, r.pgCreateUpdatePermission(newPermission, tableName, role, schemaName, sourceName))
		}
	}

	return bulkEdgePermissions
}

func (r *Runtime) createDeletePermissionForEdges(nodeTables []string, node *gen.Type, permission map[string]interface{}, role string, sourceName, schemaName string) []metadata.MetadataQuery {
	bulkEdgePermissions := []metadata.MetadataQuery{}

	for _, edge := range node.Edges {
		if !edge.IsInverse() && !edge.OwnFK() {
			tableName := edge.Rel.Table
			if r.isNodeTable(nodeTables, tableName) {
				continue
			}

			tableName, newPermission := r.tableAndPermissionsFromEdge(edge, nodeTables, permission)

			logrus.Info("creating [edge] insert permission for table: ", tableName, " with role: ", role)

			bulkEdgePermissions = append(bulkEdgePermissions, r.pgCreateDeletePermission(newPermission, tableName, role, schemaName, sourceName))
		}
	}

	return bulkEdgePermissions
}

func (r *Runtime) tableAndPermissionsFromEdge(edge *gen.Edge, nodeTables []string, permission map[string]interface{}) (string, map[string]interface{}) {
	tableName := edge.Rel.Table

	levelUp := strcase.ToLowerCamel(strings.TrimSuffix(edge.Rel.Column(), "_id"))
	newPermission := make(map[string]interface{})

	for k, v := range permission {
		newPermission[k] = v
	}

	newPermission["columns"] = edge.Rel.Columns

	if newPermission["check"] != nil || levelUp == "" {
		newPermission["check"] = map[string]interface{}{
			levelUp: newPermission["check"],
		}
	}

	if newPermission["filter"] != nil || levelUp == "" {
		newPermission["filter"] = map[string]interface{}{
			levelUp: newPermission["filter"],
		}
	}

	return tableName, newPermission
}

func (r *Runtime) allColumnsOfNode(node *gen.Type) []string {
	columns := []string{node.ID.Name}

	for _, f := range node.Fields {
		columns = append(columns, f.Column().Name)
	}

	for _, e := range node.Edges {
		if e.OwnFK() {
			columns = append(columns, e.Rel.Columns...)
		}
	}

	return columns
}
