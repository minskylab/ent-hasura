package hasura

import (
	"github.com/go-resty/resty/v2"
	"github.com/gookit/config/v2"
	"github.com/pkg/errors"
)

type EphemeralRuntime struct {
	Client      *resty.Client
	Config      *HasuraConfig
	AdminSecret string
}

func NewEphemeralRuntime(filepath ...string) (*EphemeralRuntime, error) {
	conf, err := ConfigFromFile(filepath...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &EphemeralRuntime{
		Client:      resty.New(),
		AdminSecret: config.Getenv("HASURA_GRAPHQL_ADMIN_SECRET", ""),
		Config:      conf,
	}, nil
}
