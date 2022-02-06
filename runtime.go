package hasura

import (
	hasura_api "github.com/minskylab/hasura-api"
	"github.com/pkg/errors"
)

type Runtime struct {
	hasura *hasura_api.HasuraClient
}

func NewRuntime(options ...hasura_api.HasuraClientOption) (*Runtime, error) {
	client, err := hasura_api.NewHasuraClient(options...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Runtime{
		hasura: client,
	}, nil
}

// type EphemeralRuntime struct {
// 	Client      *resty.Client
// 	Config      *HasuraConfig
// 	AdminSecret string

// 	operatedTables map[string]map[string]time.Time
// }

// type EphemeralRuntimeOptions struct {
// 	configFilepaths []string
// 	envFilepaths    []string
// }

// type EphemeralRuntimeOption func(*EphemeralRuntimeOptions)

// func WithConfigFilepath(filepath ...string) EphemeralRuntimeOption {
// 	return func(options *EphemeralRuntimeOptions) {
// 		options.configFilepaths = filepath
// 	}
// }

// func WithEnvFilepath(filepath ...string) EphemeralRuntimeOption {
// 	return func(options *EphemeralRuntimeOptions) {
// 		options.envFilepaths = filepath
// 	}
// }

// func NewEphemeralRuntime(options ...EphemeralRuntimeOption) (*EphemeralRuntime, error) {
// 	optionsStruct := new(EphemeralRuntimeOptions)

// 	for _, opt := range options {
// 		opt(optionsStruct)
// 	}

// 	if err := godotenv.Load(optionsStruct.envFilepaths...); err != nil {
// 		logrus.Warn("Error loading .env file", err)
// 	}

// 	conf, err := ConfigFromFile(optionsStruct.configFilepaths...)
// 	if err != nil {
// 		return nil, errors.WithStack(err)
// 	}

// 	logrus.Info(conf.Endpoint)
// 	return &EphemeralRuntime{
// 		Client:      resty.New(),
// 		AdminSecret: config.Getenv("HASURA_GRAPHQL_ADMIN_SECRET", ""),
// 		Config:      conf,
// 	}, nil
// }
