package enthasura

import (
	"strings"

	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/entc/gen"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/minskylab/hasura-api/metadata"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
)

const (
	insertVerbName    = "insert"
	updateVerbName    = "update"
	deleteVerbName    = "delete"
	aggregateVerbName = "aggregate"
)

func basicDefinition(pluralize *pluralize.Client, tableName, nodeName, schemaName string) (*metadata.TableDefinition, error) {
	definition := &metadata.TableDefinition{}

	definition.Table = metadata.QualifiedTableName{
		Name:   tableName,
		Schema: schemaName,
	}

	singularName := strcase.ToCamel(pluralize.Singular(nodeName))
	pluralName := strcase.ToCamel(pluralize.Plural(nodeName))

	aggregationSuffix := strcase.ToCamel(aggregateVerbName)

	if singularName == pluralName {
		logger.Warn("singular-plural equality found")
		logger.Warn(singularName, " table:", tableName, " node:", nodeName, " schema:", schemaName)
		pluralName = pluralName + "s"
	}

	definition.Configuration = &metadata.TableConfiguration{
		CustomRootFields:  &metadata.CustomRootFields{},
		CustomColumnNames: map[string]string{},
	}

	definition.Configuration.CustomName = nodeName
	definition.Configuration.CustomRootFields = &metadata.CustomRootFields{
		Insert:          strcase.ToLowerCamel(insertVerbName + pluralName),
		InsertOne:       strcase.ToLowerCamel(insertVerbName + singularName),
		Select:          strcase.ToLowerCamel(pluralName),
		SelectByPk:      strcase.ToLowerCamel(singularName),
		SelectAggregate: strcase.ToLowerCamel(pluralName + aggregationSuffix),
		Update:          strcase.ToLowerCamel(updateVerbName + pluralName),
		UpdateByPk:      strcase.ToLowerCamel(updateVerbName + singularName),
		Delete:          strcase.ToLowerCamel(deleteVerbName + pluralName),
		DeleteByPk:      strcase.ToLowerCamel(deleteVerbName + singularName),
	}

	return definition, nil
}

func hasuraTableMetadataFromNode(pluralize *pluralize.Client, node *gen.Type, schemaName string) (*metadata.TableDefinition, error) {
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
				definition.Configuration.CustomColumnNames[name] = fieldName
			}

			if edge.OwnFK() {
				fieldName = strcase.ToLowerCamel(realName) + "ID"
				definition.Configuration.CustomColumnNames[name] = fieldName
			}

			if !edge.OwnFK() {
				foreignKey = metadata.RemoteTable{
					Column: name,
					Table: metadata.QualifiedTableName{
						Schema: schemaName,
						Name:   edge.Rel.Table,
					},
				}
			}

			definition.ObjectRelationships = append(definition.ObjectRelationships, &metadata.ObjectRelationship{
				Name: realName,
				Using: metadata.Using{
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
				columnName = strcase.ToSnake(node.Name) + "_id"
			}

			tableName := edge.Rel.Table

			definition.ArrayRelationships = append(definition.ArrayRelationships, &metadata.ArrayRelationship{
				Name: realName,
				Using: metadata.UsingArray{
					ForeignKeyConstraintOn: metadata.RemoteTable{
						Column: columnName,
						Table: metadata.QualifiedTableName{
							Schema: schemaName,
							Name:   tableName,
						},
					},
				},
			})
		}
	}

	return definition, nil
}

func hasuraTableFromRelationalTable(pluralize *pluralize.Client, table *schema.Table, schemaName string) (*metadata.TableDefinition, error) {
	customName := strcase.ToCamel(pluralize.Singular(table.Name))

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

		definition.ObjectRelationships = append(definition.ObjectRelationships, &metadata.ObjectRelationship{
			Name: nameWithoutID,
			Using: metadata.Using{
				ForeignKeyConstraintOn: foreignKey,
			},
		})
	}

	return definition, nil
}

func obtainHasuraTablesFromEntSchema(schema *gen.Graph, schemaName string) ([]*metadata.TableDefinition, error) {
	pluralize := pluralize.NewClient()

	tables := []*metadata.TableDefinition{}
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

		definition, err := hasuraTableFromRelationalTable(pluralize, table, schemaName)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		tables = append(tables, definition)
	}

	return tables, nil
}

// func hasuraMetadataFromEntSchema(schema *gen.Graph, schemaName string) (*Metadata, error) {
// 	metadata := &Metadata{}

// 	tables, err := obtainHasuraTablesFromEntSchema(schema, schemaName)
// 	if err != nil {
// 		return nil, errors.WithStack(err)
// 	}

// 	metadata.Sources = append(metadata.Sources, &Source{
// 		Tables: tables,
// 	})

// 	return metadata, nil
// }

// func generateFile(metadata HasuraMetadata, outputFile string) error {
// 	data, err := json.MarshalIndent(metadata, "", "  ")
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}

// 	if err := os.MkdirAll(filepath.Dir(outputFile), os.ModePerm); err != nil {
// 		return errors.WithStack(err)
// 	}

// 	return ioutil.WriteFile(outputFile, []byte(data), 0644)
// }

// func generateRawMetadata(graph *gen.Graph, schemaName, outputFile string) error {
// 	hMetadata := HasuraMetadata{}

// 	metadata, err := hasuraMetadataFromEntSchema(graph, schemaName)
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}

// 	hMetadata.Metadata = metadata
// 	return generateFile(hMetadata, outputFile)
// }
