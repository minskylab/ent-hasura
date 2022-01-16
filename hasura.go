package hasura

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

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
