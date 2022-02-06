package hasura

import (
	"github.com/minskylab/hasura-api/metadata"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func elementInArray(array []string, element string) bool {
	for _, el := range array {
		if el == element {
			return true
		}
	}

	return false
}

func logAndResponseMetadataResponse(res metadata.MetadataResponse) error {
	switch msg := res.GetResponse().(type) {
	case metadata.SuccessResponse:
		logrus.Debug(msg)
		return nil
	case metadata.BadRequestResponse:
		logrus.Error(msg)
		return errors.New(msg.Error)
	case metadata.InternalServerErrorResponse:
		logrus.Error(msg)
		return errors.New(msg.Error)
	case metadata.UnauthorizedResponse:
		logrus.Error(msg)
		return errors.New(msg.Error)
	default:
		logrus.Warn(msg)
	}

	return nil
}
