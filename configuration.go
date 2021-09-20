package hasura

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/pkg/errors"
)

func enhanceHasuraTable(source *Source, table *TableDefinition) {
	for i, iTable := range source.Tables {
		if iTable.Table.Schema != table.Table.Schema || iTable.Table.Name != table.Table.Name {
			continue
		}

		if source.Tables[i].Configuration == nil {
			source.Tables[i].Configuration = &Configuration{}
		}

		if source.Tables[i].Configuration.CustomName == "" {
			source.Tables[i].Configuration.CustomName = table.Configuration.CustomName
		}

		if source.Tables[i].Configuration.CustomRootFields == nil {
			source.Tables[i].Configuration.CustomRootFields = &CustomRootFields{}
		}

		if source.Tables[i].Configuration.CustomRootFields.Insert == "" {
			source.Tables[i].Configuration.CustomRootFields.Insert = table.Configuration.CustomRootFields.Insert
		}

		if source.Tables[i].Configuration.CustomRootFields.InsertOne == "" {
			source.Tables[i].Configuration.CustomRootFields.InsertOne = table.Configuration.CustomRootFields.InsertOne
		}

		if source.Tables[i].Configuration.CustomRootFields.Update == "" {
			source.Tables[i].Configuration.CustomRootFields.Update = table.Configuration.CustomRootFields.Update
		}

		if source.Tables[i].Configuration.CustomRootFields.UpdateByPk == "" {
			source.Tables[i].Configuration.CustomRootFields.UpdateByPk = table.Configuration.CustomRootFields.UpdateByPk
		}

		if source.Tables[i].Configuration.CustomRootFields.Select == "" {
			source.Tables[i].Configuration.CustomRootFields.Select = table.Configuration.CustomRootFields.Select
		}

		if source.Tables[i].Configuration.CustomRootFields.Delete == "" {
			source.Tables[i].Configuration.CustomRootFields.Delete = table.Configuration.CustomRootFields.Delete
		}

		if source.Tables[i].Configuration.CustomRootFields.DeleteByPk == "" {
			source.Tables[i].Configuration.CustomRootFields.DeleteByPk = table.Configuration.CustomRootFields.DeleteByPk
		}

		if source.Tables[i].Configuration.CustomRootFields.SelectAggregate == "" {
			source.Tables[i].Configuration.CustomRootFields.SelectAggregate = table.Configuration.CustomRootFields.SelectAggregate
		}

		if source.Tables[i].Configuration.CustomRootFields.SelectByPk == "" {
			source.Tables[i].Configuration.CustomRootFields.SelectByPk = table.Configuration.CustomRootFields.SelectByPk
		}

		source.Tables[i].Configuration.CustomColumnNames = table.Configuration.CustomColumnNames

		source.Tables[i].ObjectRelationships = table.ObjectRelationships
		source.Tables[i].ArrayRelationships = table.ArrayRelationships
	}
}

func enhancedHasuraConfigurationAndRelationships(initial *HasuraMetadata, schema *gen.Graph, sourceName, schemaName string, overrideTables bool) error {
	initial.ResourceVersion += 1

	tables, err := obtainHasuraTablesFromEntSchema(schema, schemaName)
	if err != nil {
		return err
	}

	for _, source := range initial.Metadata.Sources {
		if source.Name == sourceName {
			if overrideTables {
				source.Tables = tables
				break
			}

			for _, table := range tables {
				enhanceHasuraTable(source, table)
			}
		}
	}

	return nil
}

func GenerateHasuraConfigurationAndRelationships(schemaRoute string, outputFile, inputFile, source, schemaName string, overrideTables bool, defaultRole string) error {
	graph, err := entc.LoadGraph(schemaRoute, &gen.Config{})
	if err != nil {
		return errors.WithStack(err)
	}

	if inputFile == "" { // If input file is not specified, use the default
		return generateRawMetadata(graph, schemaName, outputFile)
	}

	initialMetadata, err := parseHasuraMetadata(inputFile)
	if err != nil {
		return errors.WithStack(err)
	}

	if defaultRole != "" {
		err := enhancedHasuraPermissions(initialMetadata, graph, source, defaultRole, schemaName)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	if err := enhancedHasuraConfigurationAndRelationships(initialMetadata, graph, source, schemaName, overrideTables); err != nil {
		return errors.WithStack(err)
	}

	return generateFile(*initialMetadata, outputFile)
}
