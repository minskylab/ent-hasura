package hasura

// import (
// 	"github.com/minskylab/hasura-api/metadata"
// 	"github.com/sirupsen/logrus"
// )

// func (r *Runtime) pgTrackTable(tableName, source string) error {
// 	res, err := r.hasura.MetadataClient.PgTrackTable(&metadata.PgTrackTableArgs{
// 		Table:  metadata.TableName(tableName),
// 		Source: sourceName(source),
// 	})
// 	if err != nil {
// 		logrus.Warn(err)
// 		return nil
// 	}

// 	logAndResponseMetadataResponse(res)

// 	return nil
// }

// func (r *Runtime) pgBulkTrackTables(tables []*TableDefinition, sourceName metadata.SourceName) error {
// 	toTrackTables := make([]metadata.MetadataQuery, 0)

// 	for _, def := range tables {
// 		toTrackTables = append(toTrackTables,
// 			metadata.PgTrackTableQuery(&metadata.PgTrackTableArgs{
// 				Table:  metadata.TableName(def.Table.Name),
// 				Source: &sourceName,
// 			}),
// 		)
// 	}

// 	res, err := r.hasura.MetadataClient.Bulk(toTrackTables)
// 	if err != nil {
// 		logrus.Warn(err)
// 		return nil
// 	}

// 	logAndResponseMetadataResponse(res)

// 	return nil
// }
