package hasura

import (
	"fmt"
	"strings"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
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
	endpoint := fmt.Sprintf("%s/v1/metadata", r.Config.Endpoint)

	res, err := r.Client.R().
		SetHeaders(map[string]string{
			"Content-Type":          "application/json",
			"X-Hasura-Role":         "admin",
			"X-Hasura-Admin-Secret": r.AdminSecret,
		}).
		SetBody(ActionBody{
			Type: pgTableCustomizationAction,
			Args: PGTableCustomizationArgs{
				Table:         table.Table.Name,
				Source:        source,
				Configuration: *table.Configuration,
			},
		}).
		Post(endpoint)
	if err != nil {
		logrus.Warn(errors.WithStack(err))
		logrus.Warn("response code: ", res.StatusCode())
		return nil
	}

	// logrus.Info("response code: ", res.StatusCode())

	return nil
}

type PGTableRenameArgs struct {
	Table   string `json:"table"`
	Name    string `json:"name"`
	Source  string `json:"source"`
	NewName string `json:"new_name"`
}

func (r *EphemeralRuntime) renamePGTableRelationshipsQuery(table TableDefinition, sourceName string, newName string) error {
	endpoint := fmt.Sprintf("%s/v1/metadata", r.Config.Endpoint)

	res, err := r.Client.R().
		SetHeaders(map[string]string{
			"Content-Type":          "application/json",
			"X-Hasura-Role":         "admin",
			"X-Hasura-Admin-Secret": r.AdminSecret,
		}).
		SetBody(ActionBody{
			Type: pgTableRenameRelationship,
			Args: PGTableRenameArgs{
				Table:  table.Table.Name,
				Source: sourceName,
				Name:   newName,
			},
		}).
		Post(endpoint)
	if err != nil {
		logrus.Warn(errors.WithStack(err))
		logrus.Warn("response code: ", res.StatusCode())
		return nil
	}

	// logrus.Info("response code: ", res.StatusCode())

	return nil
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

// type ForeignKeyConstraintOn struct {
// 	Table   string   `json:"table"`
// 	Columns []string `json:"columns"`
// }

func (r *EphemeralRuntime) genericHasuraMetadataQuery(body ActionBody) error {
	endpoint := fmt.Sprintf("%s/v1/metadata", r.Config.Endpoint)

	res, err := r.Client.R().
		SetHeaders(map[string]string{
			"Content-Type":          "application/json",
			"X-Hasura-Role":         "admin",
			"X-Hasura-Admin-Secret": r.AdminSecret,
		}).
		SetBody(body).
		Post(endpoint)
	if err != nil {
		logrus.Warn(errors.WithStack(err))
		logrus.Warn("response: ", res.StatusCode(), " ", res.String())
		return nil
	}

	logrus.Info("response: ", res.StatusCode(), " ", res.String())

	return nil
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

func (r *EphemeralRuntime) renamePGTableRelationships(table TableDefinition, sourceName string) error {
	for _, t := range table.ObjectRelationships {
		r.renamePGTableRelationshipsQuery(table, sourceName, t.Name)
	}

	for _, t := range table.ArrayRelationships {
		r.renamePGTableRelationshipsQuery(table, sourceName, t.Name)
	}

	return nil
}

func (r *EphemeralRuntime) TrackAllTables(schemaRoute, schemaName, sourceName string) error {
	graph, err := entc.LoadGraph(schemaRoute, &gen.Config{})
	if err != nil {
		return errors.WithStack(err)
	}

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

func (r *EphemeralRuntime) ApplyPGTableCustomizationForAllTables(schemaRoute, schemaName, sourceName string) error {
	r.TrackAllTables(schemaRoute, schemaName, sourceName)

	graph, err := entc.LoadGraph(schemaRoute, &gen.Config{})
	if err != nil {
		return errors.WithStack(err)
	}

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

	if err := r.ApplyPGPermissionsForAllTables(schemaRoute, schemaName, sourceName); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *EphemeralRuntime) ApplyPGPermissionsForAllTables(schemaRoute, schemaName, sourceName string) error {
	graph, err := entc.LoadGraph(schemaRoute, &gen.Config{})
	if err != nil {
		return errors.WithStack(err)
	}

	for _, node := range graph.Nodes {
		permAnn, isOk := node.Annotations[hasuraPermissionsRoleAnnotationName].(map[string]interface{})
		if !isOk {
			logrus.Warn("skipping node: ", node.Name, " as it does not have permissions annotation")
			continue
		}

		role, isOk := permAnn["role"].(string)
		if !isOk {
			logrus.Warn("skipping node: ", node.Name, " as it does not have permissions role name in annotation")
			continue
		}

		// fmt.Println(node.Name, role)

		insertPermission, isOk := permAnn["insert_permission"].(map[string]interface{})
		if isOk {
			logrus.Info("creating insert permission for table: ", node.Table(), " with role: ", role)
			if err := r.pgCreateInsertPermission(insertPermission, node.Table(), role, sourceName); err != nil {
				logrus.Warn(errors.WithStack(err))
				continue
			}
			for _, edge := range node.Edges {
				if !edge.IsInverse() && !edge.OwnFK() {
					tableName := edge.Rel.Table
					levelUp := strings.TrimSuffix(edge.Rel.Column(), "_id")
					insertPermission["columns"] = edge.Rel.Columns
					insertPermission["check"] = map[string]interface{}{
						levelUp: insertPermission["check"],
					}

					logrus.Info("creating [edge] insert permission for table: ", tableName, " with role: ", role)
					if err := r.pgCreateInsertPermission(insertPermission, tableName, role, sourceName); err != nil {
						logrus.Warn(errors.WithStack(err))
						continue
					}
				}
			}
		}

		selectPermission, isOk := permAnn["select_permission"].(map[string]interface{})
		if isOk {
			logrus.Info("creating select permission for table: ", node.Table(), " with role: ", role)
			if err := r.pgCreateSelectPermission(selectPermission, node.Table(), role, sourceName); err != nil {
				logrus.Warn(errors.WithStack(err))
				continue
			}

			for _, edge := range node.Edges {
				if !edge.IsInverse() && !edge.OwnFK() {
					tableName := edge.Rel.Table
					levelUp := strings.TrimSuffix(edge.Rel.Column(), "_id")
					selectPermission["columns"] = edge.Rel.Columns
					selectPermission["filter"] = map[string]interface{}{
						levelUp: selectPermission["filter"],
					}

					logrus.Info("creating [edge] select permission for table: ", tableName, " with role: ", role)
					if err := r.pgCreateSelectPermission(selectPermission, tableName, role, sourceName); err != nil {
						logrus.Warn(errors.WithStack(err))
						continue
					}
				}
			}
		}

		updatePermission, isOk := permAnn["update_permission"].(map[string]interface{})
		if isOk {
			logrus.Info("creating update permission for table: ", node.Table(), " with role: ", role)
			if err := r.pgCreateUpdatePermission(updatePermission, node.Table(), role, sourceName); err != nil {
				logrus.Warn(errors.WithStack(err))
				continue
			}

			for _, edge := range node.Edges {
				// logrus.Warn(edge.IsInverse())
				if !edge.IsInverse() && !edge.OwnFK() {
					// logrus.Warn("[USING]")
					tableName := edge.Rel.Table
					levelUp := strings.TrimSuffix(edge.Rel.Column(), "_id")
					updatePermission["columns"] = edge.Rel.Columns
					updatePermission["check"] = map[string]interface{}{
						levelUp: updatePermission["check"],
					}
					updatePermission["filter"] = map[string]interface{}{
						levelUp: updatePermission["filter"],
					}

					logrus.Info("creating [edge] update permission for table: ", tableName, " with role: ", role)
					if err := r.pgCreateUpdatePermission(updatePermission, tableName, role, sourceName); err != nil {
						logrus.Warn(errors.WithStack(err))
						continue
					}
				}
			}
		}

		deletePermission, isOk := permAnn["delete_permission"].(map[string]interface{})
		if isOk {
			logrus.Info("creating delete permission for table: ", node.Table(), " with role: ", role)
			if err := r.pgCreateDeletePermission(deletePermission, node.Table(), role, sourceName); err != nil {
				logrus.Warn(errors.WithStack(err))
				continue
			}

			for _, edge := range node.Edges {
				if !edge.IsInverse() && !edge.OwnFK() {
					tableName := edge.Rel.Table
					levelUp := strings.TrimSuffix(edge.Rel.Column(), "_id")
					deletePermission["filter"] = map[string]interface{}{
						levelUp: deletePermission["filter"],
					}

					logrus.Info("creating [edge] delete permission for table: ", tableName, " with role: ", role)
					if err := r.pgCreateDeletePermission(deletePermission, tableName, role, sourceName); err != nil {
						logrus.Warn(errors.WithStack(err))
						continue
					}
				}
			}
		}

	}

	return nil
}
