package main

import hasura "github.com/minskylab/ent-hasura"

const defaultSchemaPath = "./ent/schema"

func main() {
	input := ""

	output := "hasura/metadata.json"
	source := "default"

	overrideTables := false
	schemaName := "public"

	defaultRole := ""

	hasura.GenerateHasuraConfigurationAndRelationships(
		defaultSchemaPath,
		output,
		input,
		source,
		schemaName,
		overrideTables,
		defaultRole,
	)
}
