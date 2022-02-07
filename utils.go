package enthasura

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

func logAndResponseMetadataResponse(res metadata.MetadataResponse, soft bool) error {
	if res == nil {
		if soft {
			return nil
		}
		return errors.New("metadata response is nil")
	}

	switch response := res.(type) {
	case metadata.RestyResponse:
		logrus.WithFields(logrus.Fields{
			"status": response.StatusCode(),
			"body":   string(response.Body()),
			"url":    response.Request.URL,
		}).Debug("metadata response")
	}

	// switch msg := res.(type) {
	// case metadata.ObjectResponse, metadata.ArrayResponse:
	// 	logrus.Debug(msg)
	// 	return nil
	// case metadata.MetadataResponses:
	// 	logrus.Debug(msg)
	// case metadata.BadRequestResponse:
	// 	logrus.Error(msg)
	// 	if soft {
	// 		return nil
	// 	}
	// 	return errors.New(msg.Error)
	// case metadata.InternalServerErrorResponse:
	// 	logrus.Error(msg)
	// 	if soft {
	// 		return nil
	// 	}
	// 	return errors.New(msg.Error)
	// case metadata.UnauthorizedResponse:
	// 	logrus.Error(msg)
	// 	if soft {
	// 		return nil
	// 	}
	// 	return errors.New(msg.Error)
	// default:
	// 	logrus.Warn(msg)
	// }

	return nil
}
