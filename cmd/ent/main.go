package main

import (
	"log"
	"os"

	hasura "github.com/minskylab/ent-hasura"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "ent-hasura",
		Usage: "make an explosive entrance",
		Commands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "generate a default metadata file",
				Flags: []cli.Flag{
					stringFlag("schema", "s", "./ent/schema"),
					stringFlag("name", "n", "public"),
					stringFlag("source", "c", "default"),
					stringFlag("output", "o", "hasura/metadata.json"),
					stringFlag("input", "i", ""),
					stringFlag("role", "r", ""),
					boolFlag("override", "ov", false),
				},
				Action: generateCommand,
			},
			{
				Name:  "apply",
				Usage: "apply metadata generate from ent to a Hasura GraphQL Engine",
				Flags: []cli.Flag{
					stringFlag("schema", "s", "./ent/schema"),
					stringFlag("name", "n", "public"),
					stringFlag("source", "c", "default"),
					stringFlag("envfile", "e", ".env"),
					stringFlag("configfile", "f", ""),
				},
				Action: applyCommand,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(errors.WithStack(err))
	}
}

func stringFlag(name, alias, defaultValue string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    name,
		Value:   defaultValue,
		Aliases: []string{alias},
	}
}

func boolFlag(name, alias string, defaultValue bool) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:    name,
		Value:   defaultValue,
		Aliases: []string{alias},
	}
}

func generateCommand(c *cli.Context) error {
	defaultConfig := hasura.DefaultHasuraMetadataConfig

	schema := c.String("schema")
	name := c.String("name")
	source := c.String("source")
	output := c.String("output")
	input := c.String("input")
	role := c.String("role")
	override := c.Bool("override")

	if schemaOverride := c.Args().First(); schemaOverride != "" {
		schema = schemaOverride
	}

	defaultConfig.SchemaPath = schema
	defaultConfig.SchemaName = name
	defaultConfig.Source = source
	defaultConfig.OutputMetadataFile = output
	defaultConfig.MetadataInput = input
	defaultConfig.DefaultRole = role
	defaultConfig.OverrideTables = override

	if err := hasura.CreateDefaultMetadataFromSchema(&defaultConfig); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func applyCommand(c *cli.Context) error {
	envFile := c.String("envfile")
	logrus.Info(envFile)
	run, err := hasura.NewEphemeralRuntime(
		hasura.WithEnvFilepath(envFile),
	)
	if err != nil {
		return errors.WithStack(err)
	}

	schema := c.String("schema")
	name := c.String("name")
	source := c.String("source")

	if schemaOverride := c.Args().First(); schemaOverride != "" {
		schema = schemaOverride
	}

	run.ApplyPGTableCustomizationForAllTables(schema, name, source)

	return nil
}
