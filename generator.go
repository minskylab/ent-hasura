package hasura

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/entc/gen"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
)

const (
	insertVerbName = "insert"
	updateVerbName = "update"
	deleteVerbName = "delete"
)

func basicDefinition(pluralize *pluralize.Client, tableName, nodeName, schemaName string) (*TableDefinition, error) {
	definition := &TableDefinition{}

	definition.Table.Name = tableName
	definition.Table.Schema = schemaName

	singularName := fixedName(pluralize.Singular(nodeName))
	pluralName := fixedName(pluralize.Plural(nodeName))

	if singularName == pluralName {
		logger.Warn("singular-plural equality found")
		logger.Warn(singularName, " table:", tableName, " node:", nodeName, " schema:", schemaName)
		pluralName = pluralName + "s"
	}

	definition.Configuration = &Configuration{
		CustomRootFields:  &CustomRootFields{},
		CustomColumnNames: map[string]string{},
	}

	definition.Configuration.CustomName = nodeName
	definition.Configuration.CustomRootFields = &CustomRootFields{
		Insert:          strcase.ToLowerCamel(insertVerbName + pluralName),
		InsertOne:       strcase.ToLowerCamel(insertVerbName + singularName),
		Select:          strcase.ToLowerCamel(pluralName),
		SelectByPk:      strcase.ToLowerCamel(singularName),
		SelectAggregate: strcase.ToLowerCamel(pluralName + "Aggregate"),
		Update:          strcase.ToLowerCamel(updateVerbName + pluralName),
		UpdateByPk:      strcase.ToLowerCamel(updateVerbName + singularName),
		Delete:          strcase.ToLowerCamel(deleteVerbName + pluralName),
		DeleteByPk:      strcase.ToLowerCamel(deleteVerbName + singularName),
	}

	return definition, nil
}

func hasuraTableMetadataFromNode(pluralize *pluralize.Client, node *gen.Type, schemaName string) (*TableDefinition, error) {
	definition, err := basicDefinition(pluralize, node.Table(), node.Name, schemaName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, field := range node.Fields {
		name := strcase.ToLowerCamel(field.Name)
		columnName := field.Column().Name

		definition.Configuration.CustomColumnNames[columnName] = name
	}

	for _, edge := range node.Edges {
		if edge.M2O() || edge.O2O() {
			name := edge.Rel.Column()
			realName := strcase.ToLowerCamel(edge.Name)

			var foreignKey interface{} = name

			fieldName := realName + "ID"

			if edge.M2O() {
				// fieldName = edge.Name
				definition.Configuration.CustomColumnNames[name] = fieldName
			}

			if edge.OwnFK() {
				fieldName = strcase.ToLowerCamel(realName) + "ID"
				definition.Configuration.CustomColumnNames[name] = fieldName
			}

			if !edge.OwnFK() {
				// fk, err := edge.ForeignKey()
				// if err != nil {
				// 	return nil, errors.WithStack(err)
				// }

				foreignKey = ForeignKeyConstraintOn{
					Column: name,
					Table: Table{
						Schema: schemaName,
						Name:   edge.Rel.Table,
					},
				}

				// definition.Configuration.CustomColumnNames[name] = fieldName
			}

			definition.ObjectRelationships = append(definition.ObjectRelationships, &ObjectRelationship{
				Name: realName,
				Using: Using{
					ForeignKeyConstraintOn: foreignKey,
				},
			})
		}

		if edge.M2M() || edge.O2M() {
			realName := strcase.ToLowerCamel(edge.Name)

			if edge.Ref == nil || len(edge.Ref.Rel.Columns) < 1 {
				continue
			}

			columnName := edge.Ref.Rel.Columns[0]

			if edge.M2M() {
				// columnName = edge.Ref.Rel.Columns[1]
				// if strings.EqualFold(node.Table(), "jobs") {
				// 	pp.Println(edge.Rel.Type.String())
				// 	pp.Println(edge.Constant())
				// 	pp.Println(edge.Rel)
				// 	// pp.Println(edge.Ref.Rel)
				// }

				// columnName = edge.Rel.Column()

				// pp.Println(, edge.Rel.Column())

				// for _, column := range edge.Rel.Columns {
				// if strings.HasPrefix(column, strings.ToLower(node.Name)) {
				// 	c := column
				// 	pp.Println("edge.Ref.Rel", c)
				// 	columnName = c
				// }
				// }

				columnName = strcase.ToSnake(node.Name) + "_id"
				// pp.Println(strcase.ToSnake(node.Name))

				// columnName = edge.Ref.Rel.Columns[1]

				// if edge.IsInverse() {
				// 	// logger.Info("abcd", edge)
				// 	columnName = edge.Rel.Columns[1]
				// }
			}

			// if edge.O2M() && !edge.OwnFK() {
			// if len(edge.Ref.Rel.Columns) > 1 && !edge.OwnFK() {
			// 	columnName = edge.Rel.Columns[1]
			// }
			// }

			tableName := edge.Rel.Table

			definition.ArrayRelationships = append(definition.ArrayRelationships, &ArrayRelationship{
				Name: realName,
				Using: UsingArray{
					ForeignKeyConstraintOn: ForeignKeyConstraintOn{
						Column: columnName,
						Table: Table{
							Schema: schemaName,
							Name:   tableName,
						},
					},
				},
			})
		}
	}

	// definition.SelectPermissions = []SelectPermission{}

	return definition, nil
}

func hasuraTableFromRelationalTable(pluralize *pluralize.Client, table *schema.Table, schemaName string) (*TableDefinition, error) {
	customName := fixedName(pluralize.Singular(table.Name))
	definition, err := basicDefinition(pluralize, table.Name, customName, schemaName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, field := range table.Columns {
		words := strings.Split(field.Name, "_")

		for i, word := range words {
			if strings.ToLower(word) == "id" {
				words[i] = "ID"
			}
		}

		newFieldName := strings.Join(words, "_")

		name := strcase.ToLowerCamel(newFieldName)
		nameWithoutID := strcase.ToLowerCamel(strings.ReplaceAll(newFieldName, "ID", ""))

		definition.Configuration.CustomColumnNames[field.Name] = name

		var foreignKey interface{} = field.Name

		// if
		// foreignKey = ForeignKeyConstraintOn{
		// 	Column: field.Name,
		// }

		definition.ObjectRelationships = append(definition.ObjectRelationships, &ObjectRelationship{
			Name: nameWithoutID,
			Using: Using{
				ForeignKeyConstraintOn: foreignKey,
			},
		})
	}

	// for _, f := range table.ForeignKeys {
	// 	logger.Info(f.Symbol)
	// }

	return definition, nil
}

func obtainHasuraTablesFromEntSchema(schema *gen.Graph, schemaName string) ([]*TableDefinition, error) {
	pluralize := pluralize.NewClient()
	// strcase.ConfigureAcronym("ID", "Id")

	tables := []*TableDefinition{}

	mappedNodes := []string{}

	for _, node := range schema.Nodes {
		definition, err := hasuraTableMetadataFromNode(pluralize, node, schemaName)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		mappedNodes = append(mappedNodes, node.Table())
		tables = append(tables, definition)
	}

	schemaTables, err := schema.Tables()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, table := range schemaTables {
		if elementInArray(mappedNodes, table.Name) {
			continue
		}

		// logger.Infof("table not mapped: %s", table.Name)
		definition, err := hasuraTableFromRelationalTable(pluralize, table, schemaName)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		tables = append(tables, definition)
	}

	// logger.Info("done")

	return tables, nil
}

func hasuraMetadataFromEntSchema(schema *gen.Graph, schemaName string) (*Metadata, error) {
	metadata := &Metadata{}

	tables, err := obtainHasuraTablesFromEntSchema(schema, schemaName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	metadata.Sources = append(metadata.Sources, &Source{
		Tables: tables,
	})

	return metadata, nil
}

func generateFile(metadata HasuraMetadata, outputFile string) error {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return errors.WithStack(err)
	}

	return ioutil.WriteFile(outputFile, []byte(data), 0644)
}

func generateRawMetadata(graph *gen.Graph, schemaName, outputFile string) error {
	hMetadata := HasuraMetadata{}

	metadata, err := hasuraMetadataFromEntSchema(graph, schemaName)
	if err != nil {
		return errors.WithStack(err)
	}

	hMetadata.Metadata = metadata
	return generateFile(hMetadata, outputFile)
}
