package hasura

const (
	pgTrackTable = "pg_track_table"
)

type PGTackTableArg struct {
	Table  string `json:"table"`
	Source string `json:"source"`
}

func (r *EphemeralRuntime) pgTrackTable(tableName, source string) error {
	return r.genericHasuraMetadataQuery(ActionBody{
		Type: pgTrackTable,
		Args: PGTackTableArg{
			Table:  tableName,
			Source: source,
		},
	})
}
