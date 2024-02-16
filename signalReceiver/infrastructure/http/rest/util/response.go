package util

import (
	"signalReceiver/errcatalogs"
	"signalReceiver/infrastructure/http/model"
	"time"
)

func BuildErrorResponse(err error, headersIntoStruct *model.Headers) (interface{}, error) {

	switch t := err.(type) {
	case *errcatalogs.APIBadRequestError:
		return model.ErrorCodeModel{
			HTTPCode:        t.GetStatusCode(),
			HTTPMessage:     t.GetHTTPMessage(),
			MoreInformation: t.GetMoreInformation(),
		}, nil
	case *errcatalogs.APIInternalServerError:

		return model.ErrorCodeModel{
			HTTPCode:        t.GetHTTPCode(),
			HTTPMessage:     t.GetHTTPMessage(),
			MoreInformation: t.GetMoreInformation(),
		}, nil

	case *errcatalogs.APIBusinessError:

		status := model.StatusModel{
			StatusCode:       t.GetStatusCode(),
			ServerStatusCode: t.GetServerStatusCode(),
			StatusDesc:       t.GetHTTPMessage(),
			EndDt:            time.Now(),
		}

		return model.MsgRsHdrModel{
			RqUID:  headersIntoStruct.RquID,
			Status: &status,
		}, nil
	default:
		return model.ErrorCodeModel{
			HTTPCode:        500,
			HTTPMessage:     "Error General",
			MoreInformation: "Error General",
		}, nil
	}
}
