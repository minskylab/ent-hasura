package hasura

import (
	"fmt"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	pgTableCustomizationAction = "pg_set_table_customization"
	pgTableRenameRelationship  = "pg_rename_relationship"
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
	}

	logrus.Info("response code: ", res.StatusCode())

	return nil
}

type PGTableRenameArgs struct {
	Table   string `json:"table"`
	Name    string `json:"name"`
	Source  string `json:"source"`
	NewName string `json:"new_name"`
}

func (r *EphemeralRuntime) renamePGTableRelationshipsQuery(table TableDefinition, sourceName string, newName string) {
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
	}
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

func (r *EphemeralRuntime) ApplyPGTableCustomizationForAllTables(schemaRoute, schemaName, sourceName string) error {
	graph, err := entc.LoadGraph(schemaRoute, &gen.Config{})
	if err != nil {
		return errors.WithStack(err)
	}

	tables, err := obtainHasuraTablesFromEntSchema(graph, schemaName)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, table := range tables {
		logrus.Info("pushing table customization for table: ", table.Table.Name)
		if err := r.setPGTableCustomization(*table, sourceName); err != nil {
			return errors.WithStack(err)
		}

		logrus.Info("pushing table relationships for table: ", table.Table.Name)
		if err := r.renamePGTableRelationships(*table, sourceName); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
