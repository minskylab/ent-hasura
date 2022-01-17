package hasura

type HasuraMetadata struct {
	ResourceVersion int       `json:"resource_version"`
	Metadata        *Metadata `json:"metadata"`
}

type Table struct {
	Schema string `json:"schema"`
	Name   string `json:"name"`
}

type CustomRootFields struct {
	Insert          string `json:"insert,omitempty"`
	SelectAggregate string `json:"select_aggregate,omitempty"`
	InsertOne       string `json:"insert_one,omitempty"`
	SelectByPk      string `json:"select_by_pk,omitempty"`
	Select          string `json:"select,omitempty"`
	Delete          string `json:"delete,omitempty"`
	Update          string `json:"update,omitempty"`
	DeleteByPk      string `json:"delete_by_pk,omitempty"`
	UpdateByPk      string `json:"update_by_pk,omitempty"`
}

type Configuration struct {
	CustomRootFields  *CustomRootFields `json:"custom_root_fields,omitempty"`
	CustomName        string            `json:"custom_name,omitempty"`
	CustomColumnNames map[string]string `json:"custom_column_names,omitempty"`
}

type Using struct {
	// ForeignKeyConstraintOn string `json:"foreign_key_constraint_on"`
	ForeignKeyConstraintOn interface{} `json:"foreign_key_constraint_on,omitempty"`
}

type ObjectRelationship struct {
	Name  string `json:"name"`
	Using Using  `json:"using"`
}

type ID struct {
	Eq string `json:"_eq"`
}

type PermissionInsert struct {
	Check           map[string]interface{} `json:"check"`
	Set             map[string]interface{} `json:"set,omitempty"`
	Columns         []string               `json:"columns,omitempty"`
	AllColumns      bool                   `json:"all_columns,omitempty"`
	ExcludedColumns []string               `json:"excluded_columns,omitempty"`
	BackendOnly     bool                   `json:"backend_only,omitempty"`
}

type InsertPermission struct {
	Role       string           `json:"role,omitempty"`
	Permission PermissionInsert `json:"permission,omitempty"`
}

type PermissionSelect struct {
	Columns           []string               `json:"columns,omitempty"`
	AllColumns        bool                   `json:"all_columns,omitempty"`
	ExcludedColumns   []string               `json:"excluded_columns,omitempty"`
	Filter            map[string]interface{} `json:"filter,omitempty"`
	AllowAggregations bool                   `json:"allow_aggregations,omitempty"`
}

type SelectPermission struct {
	Role       string           `json:"role,omitempty"`
	Permission PermissionSelect `json:"permission,omitempty"`
}

type PermissionUpdate struct {
	Columns         []string               `json:"columns,omitempty"`
	AllColumns      bool                   `json:"all_columns,omitempty"`
	ExcludedColumns []string               `json:"excluded_columns,omitempty"`
	Filter          map[string]interface{} `json:"filter,omitempty"`
	Check           map[string]interface{} `json:"check,omitempty"`
}

type UpdatePermission struct {
	Role       string           `json:"role,omitempty"`
	Permission PermissionUpdate `json:"permission,omitempty"`
}

type PermissionDelete struct {
	Filter map[string]interface{} `json:"filter,omitempty"`
}

type DeletePermission struct {
	Role       string           `json:"role,omitempty"`
	Permission PermissionDelete `json:"permission,omitempty"`
}

type ForeignKeyConstraintOn struct {
	Column          string   `json:"column,omitempty"`
	AllColumns      bool     `json:"all_columns,omitempty"`
	ExcludedColumns []string `json:"excluded_columns,omitempty"`
	Table           Table    `json:"table,omitempty"`
}

type UsingArray struct {
	ForeignKeyConstraintOn ForeignKeyConstraintOn `json:"foreign_key_constraint_on"`
}

type ArrayRelationship struct {
	Name  string     `json:"name"`
	Using UsingArray `json:"using"`
}

type TableDefinition struct {
	Table               Table                 `json:"table"`
	Configuration       *Configuration        `json:"configuration,omitempty"`
	ObjectRelationships []*ObjectRelationship `json:"object_relationships,omitempty"`
	InsertPermissions   []*InsertPermission   `json:"insert_permissions,omitempty"`
	SelectPermissions   []*SelectPermission   `json:"select_permissions,omitempty"`
	UpdatePermissions   []*UpdatePermission   `json:"update_permissions,omitempty"`
	DeletePermissions   []*DeletePermission   `json:"delete_permissions,omitempty"`
	ArrayRelationships  []*ArrayRelationship  `json:"array_relationships,omitempty"`
}

type DatabaseURL struct {
	FromEnv string `json:"from_env"`
}

type PoolSettings struct {
	ConnectionLifetime int `json:"connection_lifetime"`
	Retries            int `json:"retries"`
	IdleTimeout        int `json:"idle_timeout"`
	MaxConnections     int `json:"max_connections"`
}

type ConnectionInfo struct {
	UsePreparedStatements bool         `json:"use_prepared_statements"`
	DatabaseURL           DatabaseURL  `json:"database_url"`
	IsolationLevel        string       `json:"isolation_level"`
	PoolSettings          PoolSettings `json:"pool_settings"`
}

type ConfigurationSource struct {
	ConnectionInfo ConnectionInfo `json:"connection_info"`
}

type Source struct {
	Name          string              `json:"name"`
	Kind          string              `json:"kind"`
	Tables        []*TableDefinition  `json:"tables"`
	Configuration ConfigurationSource `json:"configuration"`
}

type Definition struct {
	URL                  string `json:"url"`
	TimeoutSeconds       int    `json:"timeout_seconds"`
	ForwardClientHeaders bool   `json:"forward_client_headers"`
}

type RemoteSchemas struct {
	Name       string     `json:"name"`
	Definition Definition `json:"definition"`
}

type Metadata struct {
	Version       int              `json:"version"`
	Sources       []*Source        `json:"sources"`
	RemoteSchemas []*RemoteSchemas `json:"remote_schemas"`
}
