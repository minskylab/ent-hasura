package hasura

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// applyHasuraMetadata execute the replace_metadata method to our hasura instance and metadata.
func applyHasuraMetadata(hasuraHost string, metadataFilepath string) error {
	endpoint := fmt.Sprintf("%s/v1/api", hasuraHost)
	if !strings.HasPrefix(endpoint, "http") {
		endpoint = "http://" + endpoint
	}

	metadataData, err := os.ReadFile(metadataFilepath)
	if err != nil {
		return errors.WithStack(err)
	}

	buff := bytes.NewBufferString(fmt.Sprintf(`{
		"type" : "replace_metadata",
		"args": %s
	}`, string(metadataData)))

	req, err := http.NewRequest(http.MethodPost, endpoint, buff)
	if err != nil {
		return errors.WithStack(err)
	}

	head := req.Header.Clone()

	head.Add("Content-Type", "application/json")
	head.Add("X-Hasura-Role", "admin")

	req.Header = head

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}

	if res.StatusCode != 200 {
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return errors.WithStack(err)
		}

		return errors.New(fmt.Sprintf("error response, code: %d, res: %s", res.StatusCode, string(resBody)))
	}

	return nil
}
