package util

import (
	"net/http"
	"signalReceiver/infrastructure/http/model"
)

const (
	RqUID = "X-RqUID"
)

func ExtractHeadersFromRequest(req *http.Request) (model.Headers, error) {
	headersIntoStruct, err := ExtractHeadersIntoStruct(req)
	if err != nil {
		return model.Headers{}, err
	}
	return headersIntoStruct, err
}

func ExtractHeadersIntoStruct(request *http.Request) (model.Headers, error) {

	headers := model.Headers{RquID: request.Header.Get(RqUID)}

	return headers, nil
}
