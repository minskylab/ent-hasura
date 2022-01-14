package hasura

import (
	"github.com/go-resty/resty/v2"
	"github.com/gookit/config/v2"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type EphemeralRuntime struct {
	Client      *resty.Client
	Config      *HasuraConfig
	AdminSecret string
}

type EphemeralRuntimeOptions struct {
	configFilepaths []string
	envFilepaths    []string
}

type EphemeralRuntimeOption func(*EphemeralRuntimeOptions)

func WithConfigFilepath(filepath ...string) EphemeralRuntimeOption {
	return func(options *EphemeralRuntimeOptions) {
		options.configFilepaths = filepath
	}
}

func WithEnvFilepath(filepath ...string) EphemeralRuntimeOption {
	return func(options *EphemeralRuntimeOptions) {
		options.envFilepaths = filepath
	}
}

func NewEphemeralRuntime(options ...EphemeralRuntimeOption) (*EphemeralRuntime, error) {
	optionsStruct := new(EphemeralRuntimeOptions)

	for _, opt := range options {
		opt(optionsStruct)
	}

	if err := godotenv.Load(optionsStruct.envFilepaths...); err != nil {
		logrus.Warn("Error loading .env file", err)
	}

	conf, err := ConfigFromFile(optionsStruct.configFilepaths...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	logrus.Info(conf.Endpoint)
	return &EphemeralRuntime{
		Client:      resty.New(),
		AdminSecret: config.Getenv("HASURA_GRAPHQL_ADMIN_SECRET", ""),
		Config:      conf,
	}, nil
}
