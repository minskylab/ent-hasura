package hasura

import (
	"strings"
	"time"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	pgTableCustomizationAction  = "pg_set_table_customization"
	pgTableRenameRelationship   = "pg_rename_relationship"
	pgTableCreateObjectRelation = "pg_create_object_relationship"
	pgTableCreateArrayRelation  = "pg_create_array_relationship"
)

type PGTableCustomizationArgs struct {
	Table         string        `json:"table"`
	Source        string        `json:"source"`
	Configuration Configuration `json:"configuration"`
}

func (r *EphemeralRuntime) setPGTableCustomization(table TableDefinition, source string) error {
	return r.genericHasuraMetadataQuery(ActionBody{
		Type: pgTableCustomizationAction,
		Args: PGTableCustomizationArgs{
			Table:         table.Table.Name,
			Source:        source,
			Configuration: *table.Configuration,
		},
	})
}

type PGCreateRelationship struct {
	Table  string      `json:"table"`
	Name   string      `json:"name"`
	Source string      `json:"source"`
	Using  interface{} `json:"using"`
}

type PGCreateObjectUsing struct {
	ForeignKeyConstraintOn []string `json:"foreign_key_constraint_on"`
}

type PGCreateArrayUsing struct {
	ForeignKeyConstraintOn ForeignKeyConstraintOn `json:"foreign_key_constraint_on"`
}

func (r *EphemeralRuntime) createPGObjectRelationships(table TableDefinition, rel *ObjectRelationship, sourceName string) error {
	return r.genericHasuraMetadataQuery(ActionBody{
		Type: pgTableCreateObjectRelation,
		Args: PGCreateRelationship{
			Table:  table.Table.Name,
			Source: sourceName,
			Name:   rel.Name,
			Using:  rel.Using,
		},
	})
}

func (r *EphemeralRuntime) createPGArrayRelationships(table TableDefinition, rel *ArrayRelationship, sourceName string) error {
	return r.genericHasuraMetadataQuery(ActionBody{
		Type: pgTableCreateArrayRelation,
		Args: PGCreateRelationship{
			Table:  table.Table.Name,
			Source: sourceName,
			Name:   rel.Name,
			Using:  rel.Using,
		},
	})
}

func (r *EphemeralRuntime) createPGTableRelationships(table TableDefinition, sourceName string) error {
	for _, rel := range table.ObjectRelationships {
		r.createPGObjectRelationships(table, rel, sourceName)
	}

	for _, rel := range table.ArrayRelationships {
		r.createPGArrayRelationships(table, rel, sourceName)
	}

	return nil
}

func (r *EphemeralRuntime) ApplyPGFullProcessForAllTables(schemaRoute, schemaName, sourceName string) error {
	graph, err := entc.LoadGraph(schemaRoute, &gen.Config{})
	if err != nil {
		return errors.WithStack(err)
	}

	if err := r.resetMetadata(); err != nil {
		return errors.WithStack(err)
	}

	if err := r.TrackAllTables(graph, schemaName, sourceName); err != nil {
		return errors.WithStack(err)
	}

	if err := r.ApplyPGTableCustomizationForAllTables(graph, schemaName, sourceName); err != nil {
		return errors.WithStack(err)
	}

	if err := r.ApplyPGPermissionsForAllTables(graph, schemaName, sourceName); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *EphemeralRuntime) TrackAllTables(graph *gen.Graph, schemaName, sourceName string) error {
	tables, err := obtainHasuraTablesFromEntSchema(graph, schemaName)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, def := range tables {
		logrus.Info("tracking table: ", def.Table.Name)
		if err := r.pgTrackTable(def.Table.Name, sourceName); err != nil {
			logrus.Warn(errors.WithStack(err))
		}
	}

	return nil
}

func (r *EphemeralRuntime) ApplyPGTableCustomizationForAllTables(graph *gen.Graph, schemaName, sourceName string) error {
	tables, err := obtainHasuraTablesFromEntSchema(graph, schemaName)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, def := range tables {
		logrus.Info("pushing set table customization for table: ", def.Table.Name)
		if err := r.setPGTableCustomization(*def, sourceName); err != nil {
			return errors.WithStack(err)
		}

		logrus.Info("pushing create table relationships for table: ", def.Table.Name)
		if err := r.createPGTableRelationships(*def, sourceName); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (r *EphemeralRuntime) ApplyPGPermissionsForAllTables(graph *gen.Graph, schemaName, sourceName string) error {
	if r.operatedTables == nil {
		r.operatedTables = make(map[string]map[string]time.Time)
	}

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

		role, isOk := permAnn["role"].(string)
		if !isOk {
			logrus.Debug("skipping node: ", node.Name, " as it does not have permissions role name in annotation")
			continue
		}

		allColumns := r.autocompleteWithAllColumns(node)

		if insertPermission, isOk := permAnn["insert_permission"].(map[string]interface{}); isOk {
			logrus.Info("creating insert permission for table: ", node.Table(), " with role: ", role)
			if err := r.pgCreateAllXPermissionforNode(pgCreateInsertPermission, nodeTables, insertPermission, node, role, sourceName, allColumns); err != nil {
				logrus.Warn(errors.WithStack(err))
			}
		}

		if selectPermission, isOk := permAnn["select_permission"].(map[string]interface{}); isOk {
			logrus.Info("creating select permission for table: ", node.Table(), " with role: ", role)
			if err := r.pgCreateAllXPermissionforNode(pgCreateSelectPermission, nodeTables, selectPermission, node, role, sourceName, allColumns); err != nil {
				logrus.Warn(errors.WithStack(err))
			}
		}

		if updatePermission, isOk := permAnn["update_permission"].(map[string]interface{}); isOk {
			logrus.Info("creating update permission for table: ", node.Table(), " with role: ", role)
			if err := r.pgCreateAllXPermissionforNode(pgCreateUpdatePermission, nodeTables, updatePermission, node, role, sourceName, allColumns); err != nil {
				logrus.Warn(errors.WithStack(err))
			}
		}

		if deletePermission, isOk := permAnn["delete_permission"].(map[string]interface{}); isOk {
			logrus.Info("creating delete permission for table: ", node.Table(), " with role: ", role)
			if err := r.pgCreateAllXPermissionforNode(pgCreateDeletePermission, nodeTables, deletePermission, node, role, sourceName, allColumns); err != nil {
				logrus.Warn(errors.WithStack(err))
			}
		}
	}

	return nil
}

func (r *EphemeralRuntime) pgCreateAllXPermissionforNode(op HasuraOperation, nodeTables []string, perm map[string]interface{}, node *gen.Type, role string, sourceName string, allColumns []string) error {
	if err := r.pgCreateXPermission(op, perm, node.Table(), role, sourceName, allColumns...); err != nil {
		return errors.WithStack(err)
	}

	return r.createXPermissionForEdges(op, nodeTables, node, perm, role, sourceName)
}

func (r *EphemeralRuntime) isNodeTable(nodeTables []string, tableName string) bool {
	for _, nodeTable := range nodeTables {
		if nodeTable == tableName {
			return true
		}
	}

	return false
}

func (r *EphemeralRuntime) createXPermissionForEdges(op HasuraOperation, nodeTables []string, node *gen.Type, permission map[string]interface{}, role string, sourceName string) error {
	for _, edge := range node.Edges {
		if !edge.IsInverse() && !edge.OwnFK() {
			tableName := edge.Rel.Table

			if r.isNodeTable(nodeTables, tableName) {
				continue
			}

			levelUp := strcase.ToLowerCamel(strings.TrimSuffix(edge.Rel.Column(), "_id"))
			newPermission := make(map[string]interface{})

			for k, v := range permission {
				newPermission[k] = v
			}

			newPermission["columns"] = edge.Rel.Columns

			if newPermission["check"] != nil || levelUp == "" {
				// newPermission["check"].(map[string]interface{})
				newPermission["check"] = map[string]interface{}{
					levelUp: newPermission["check"],
				}
			}

			if newPermission["filter"] != nil || levelUp == "" {
				newPermission["filter"] = map[string]interface{}{
					levelUp: newPermission["filter"],
				}
			}

			logrus.Info("creating [edge] insert permission for table: ", tableName, " with role: ", role)
			if err := r.pgCreateXPermission(op, newPermission, tableName, role, sourceName); err != nil {
				logrus.Warn(errors.WithStack(err))
				continue
			}
		}
	}

	return nil
}

func (r *EphemeralRuntime) autocompleteWithAllColumns(node *gen.Type) []string {
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
