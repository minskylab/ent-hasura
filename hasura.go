package hasura

type HasuraMetadataConfig struct {
	SchemaPath string
	SchemaName string
	Source     string

	OutputMetadataFile string

	MetadataInput  string
	OverrideTables bool
	DefaultRole    string
}

var DefaultHasuraMetadataConfig HasuraMetadataConfig = HasuraMetadataConfig{
	SchemaPath: "./ent/schema",
	SchemaName: "public",
	Source:     "default",

	OutputMetadataFile: "hasura/metadata.json",

	MetadataInput:  "",
	OverrideTables: false,
	DefaultRole:    "",
}

func CreateDefaultMetadataFromSchema(config *HasuraMetadataConfig) error {
	return GenerateHasuraConfigurationAndRelationships(
		config.SchemaPath,
		config.OutputMetadataFile,
		config.MetadataInput,
		config.Source,
		config.SchemaName,
		config.OverrideTables,
		config.DefaultRole,
	)
}
