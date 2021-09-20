package hasura

import (
	"encoding/json"
	"io/ioutil"
)

func parseHasuraMetadata(filepath string) (*HasuraMetadata, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	meta := HasuraMetadata{}
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	return &meta, nil
}

func elementInArray(array []string, element string) bool {
	for _, el := range array {
		if el == element {
			return true
		}
	}

	return false
}
