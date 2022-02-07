package enthasura

import (
	"strings"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/iancoleman/strcase"
	"github.com/minskylab/hasura-api/metadata"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (r *Runtime) PerformFullMetadataTransform(entSchemaPath string, sourceName, schemaName string) error {
	graph, err := entc.LoadGraph(entSchemaPath, &gen.Config{})
	if err != nil {
		return errors.WithStack(err)
	}

	logrus.Info("[1] Prelude, untracking tables or cleaning metadata")
	if err := r.PerformPrelude(graph, sourceName, schemaName, false); err != nil {
		return errors.WithMessage(err, "error at prelude")
	}

	logrus.Info("[2] Tracking all tables related to your Ent Schema")
	if err := r.TrackAllTables(graph, sourceName, schemaName); err != nil {
		return errors.WithMessage(err, "error at track all tables")
	}

	logrus.Info("[3] Customizing all tables with standard GraphQL casing")
	if err := r.CustomizeAllTables(graph, sourceName, schemaName); err != nil {
		return errors.WithMessage(err, "error at customize all tables")
	}

	logrus.Info("[4] Adding permission annotations to all tables")
	if err := r.PermissionsForAllTables(graph, sourceName, schemaName); err != nil {
		return errors.WithMessage(err, "error at permissions for all tables")
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

	return logAndResponseMetadataResponse(res, true)
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

	if len(trackBatch) > 0 {
		logrus.Infof("ready to TRACK %d tables", len(trackBatch))

		res, err := r.hasura.Metadata.Bulk(trackBatch)
		if err != nil {
			return errors.WithStack(err)
		}

		logAndResponseMetadataResponse(res, true)
	}

	return nil
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
					ForeignKeyConstraintOn: metadata.SameTable(rel.Using.ForeignKeyConstraintOn.(string)),
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

	if len(tableCustomizeBatch) > 0 {
		logrus.Infof("ready to set %d CUSTOMIZE TABLES", len(tableCustomizeBatch))

		res, err := r.hasura.Metadata.Bulk(tableCustomizeBatch)
		if err != nil {
			return errors.WithStack(err)
		}

		logAndResponseMetadataResponse(res, true)
	}

	if len(objectRelationBulk) > 0 {
		logrus.Infof("ready to set %d OBJECT RELATIONSHIPS", len(objectRelationBulk))

		res, err := r.hasura.Metadata.Bulk(objectRelationBulk)
		if err != nil {
			return errors.WithStack(err)
		}

		logAndResponseMetadataResponse(res, true)
	}

	if len(arrayRelationBulk) > 0 {
		logrus.Infof("ready to set %d ARRAY RELATIONSHIPS", len(arrayRelationBulk))

		res, err := r.hasura.Metadata.Bulk(arrayRelationBulk)
		if err != nil {
			return errors.WithStack(err)
		}

		logAndResponseMetadataResponse(res, true)
	}

	return nil
}

func (r *Runtime) PermissionsForAllTables(graph *gen.Graph, sourceName, schemaName string) error {
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
			// logrus.Debug("skipping node: ", node.Name, " as it does not have permissions annotation")
			continue
		}

		roleName, isOk := permAnn["role"].(string)
		if !isOk {
			logrus.Warn("skipping node: ", node.Name, " as it does not have permissions role name in annotation")
			continue
		}

		if insertPermission, isOk := permAnn["insert_permission"].(map[string]interface{}); isOk {
			// logrus.Info("creating insert permission for table: ", node.Table(), " with role: ", roleName)

			insertPermissionBulk = append(
				insertPermissionBulk,
				r.pgCreateInsertPermission(insertPermission, node.Table(), roleName, sourceName, schemaName),
			)

			insertPermissionBulk = append(
				insertPermissionBulk,
				r.createInsertPermissionForEdges(nodeTables, node, insertPermission, roleName, sourceName, schemaName)...,
			)
		}

		if selectPermission, isOk := permAnn["select_permission"].(map[string]interface{}); isOk {
			// logrus.Info("creating select permission for table: ", node.Table(), " with role: ", roleName)
			selectPermissionBulk = append(
				selectPermissionBulk,
				r.pgCreateSelectPermission(selectPermission, node.Table(), roleName, sourceName, schemaName),
			)

			selectPermissionBulk = append(
				selectPermissionBulk,
				r.createSelectPermissionForEdges(nodeTables, node, selectPermission, roleName, sourceName, schemaName)...,
			)
		}

		if updatePermission, isOk := permAnn["update_permission"].(map[string]interface{}); isOk {
			// logrus.Info("creating update permission for table: ", node.Table(), " with role: ", roleName)

			updatePermissionBulk = append(
				updatePermissionBulk,
				r.pgCreateUpdatePermission(updatePermission, node.Table(), roleName, sourceName, schemaName),
			)

			updatePermissionBulk = append(
				updatePermissionBulk,
				r.createUpdatePermissionForEdges(nodeTables, node, updatePermission, roleName, sourceName, schemaName)...,
			)

		}

		if deletePermission, isOk := permAnn["delete_permission"].(map[string]interface{}); isOk {
			// logrus.Info("creating delete permission for table: ", node.Table(), " with role: ", roleName)

			deletePermissionBulk = append(
				deletePermissionBulk,
				r.pgCreateDeletePermission(deletePermission, node.Table(), roleName, sourceName, schemaName),
			)

			deletePermissionBulk = append(
				deletePermissionBulk,
				r.createDeletePermissionForEdges(nodeTables, node, deletePermission, roleName, sourceName, schemaName)...,
			)
		}
	}

	if len(insertPermissionBulk) > 0 {
		logrus.Infof("ready to create %d INSERT permissions", len(insertPermissionBulk))

		res, err := r.hasura.Metadata.Bulk(insertPermissionBulk)
		if err != nil {
			return errors.WithStack(err)
		}

		logAndResponseMetadataResponse(res, true)
	}

	if len(selectPermissionBulk) > 0 {
		logrus.Infof("ready to create %d SELECT permissions", len(selectPermissionBulk))

		res, err := r.hasura.Metadata.Bulk(selectPermissionBulk)
		if err != nil {
			return errors.WithStack(err)
		}

		logAndResponseMetadataResponse(res, true)
	}

	if len(updatePermissionBulk) > 0 {
		logrus.Infof("ready to create %d UPDATE permissions", len(updatePermissionBulk))

		res, err := r.hasura.Metadata.Bulk(updatePermissionBulk)
		if err != nil {
			return errors.WithStack(err)
		}

		logAndResponseMetadataResponse(res, true)
	}

	if len(deletePermissionBulk) > 0 {
		logrus.Infof("ready to create %d DELETE permissions", len(deletePermissionBulk))

		res, err := r.hasura.Metadata.Bulk(deletePermissionBulk)
		if err != nil {
			return errors.WithStack(err)
		}

		logAndResponseMetadataResponse(res, true)
	}

	return nil
}

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

			// logrus.Info("creating [edge] insert permission for table: ", tableName, " with role: ", role)

			bulkEdgePermissions = append(bulkEdgePermissions, r.pgCreateInsertPermission(newPermission, tableName, role, sourceName, schemaName))
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

			// logrus.Info("creating [edge] insert permission for table: ", tableName, " with role: ", role)

			bulkEdgePermissions = append(bulkEdgePermissions, r.pgCreateSelectPermission(newPermission, tableName, role, sourceName, schemaName))
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

			// logrus.Info("creating [edge] insert permission for table: ", tableName, " with role: ", role)

			bulkEdgePermissions = append(bulkEdgePermissions, r.pgCreateUpdatePermission(newPermission, tableName, role, sourceName, schemaName))
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

			// logrus.Info("creating [edge] insert permission for table: ", tableName, " with role: ", role)

			bulkEdgePermissions = append(bulkEdgePermissions, r.pgCreateDeletePermission(newPermission, tableName, role, sourceName, schemaName))
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
