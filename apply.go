package hasura

import (
	"fmt"
)

const pgTableCustomizationAction = "pg_set_table_customization"

type PGTableCustomizationArgs struct {
	Table         string        `json:"table"`
	Source        string        `json:"source"`
	Configuration Configuration `json:"configuration"`
}

func (r *EphemeralRuntime) setPGTableCustomization(hasuraHost string, table TableDefinition, source ...string) error {
	endpoint := fmt.Sprintf("%s/v1/metadata", hasuraHost)
	// if !strings.HasPrefix(endpoint, "http://") {
	// 	endpoint = "http://" + endpoint
	// }

	selectedSource := "default"
	if len(source) > 0 {
		selectedSource = source[0]
	}

	r.Client.R().
		SetHeaders(map[string]string{
			"Content-Type":          "application/json",
			"X-Hasura-Role":         "admin",
			"X-Hasura-Admin-Secret": r.AdminSecret,
		}).
		SetBody(ActionBody{
			Type: pgTableCustomizationAction,
			Args: PGTableCustomizationArgs{
				Table:         table.Table.Name,
				Source:        selectedSource,
				Configuration: *table.Configuration,
			},
		}).
		Post(endpoint)

	return nil
}

// func applyHasuraMetadata(hasuraHost string, metadataFilepath string) error {
// 	endpoint := fmt.Sprintf("%s/v1/api", hasuraHost)
// 	if !strings.HasPrefix(endpoint, "http") {
// 		endpoint = "http://" + endpoint
// 	}

// 	metadataData, err := os.ReadFile(metadataFilepath)
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}

// 	buff := bytes.NewBufferString(fmt.Sprintf(`{
// 		"type" : "replace_metadata",
// 		"args": %s
// 	}`, string(metadataData)))

// 	req, err := http.NewRequest(http.MethodPost, endpoint, buff)
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}

// 	head := req.Header.Clone()

// 	head.Add("Content-Type", "application/json")
// 	head.Add("X-Hasura-Role", "admin")

// 	req.Header = head

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}

// 	if res.StatusCode != 200 {
// 		resBody, err := io.ReadAll(res.Body)
// 		if err != nil {
// 			return errors.WithStack(err)
// 		}

// 		return errors.New(fmt.Sprintf("error response, code: %d, res: %s", res.StatusCode, string(resBody)))
// 	}

// 	return nil
// }
