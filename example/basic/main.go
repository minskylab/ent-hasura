package main

import (
	hasura "github.com/minskylab/ent-hasura"
	"github.com/pkg/errors"
)

func main() {
	defaultConfig := hasura.DefaultHasuraMetadataConfig

	if err := hasura.CreateDefaultMetadataFromSchema(&defaultConfig); err != nil {
		panic(errors.WithStack(err))
	}
}
