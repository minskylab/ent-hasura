package hasura

const resetMetadataOperation = "clear_metadata"

type ResetMetadataArgs struct{}

func (r *EphemeralRuntime) resetMetadata() error {
	return r.genericHasuraMetadataQuery(ActionBody{
		Type: string(resetMetadataOperation),
		Args: ResetMetadataArgs{},
	})
}
