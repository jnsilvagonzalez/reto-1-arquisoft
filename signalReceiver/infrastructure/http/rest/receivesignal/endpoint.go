package receivesignal

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"log"
	"net/http"
	"signalReceiver/container"
	"signalReceiver/errcatalogs"
	"signalReceiver/infrastructure/http/model"
	"signalReceiver/infrastructure/http/rest/util"
	"signalReceiver/portinterface/signalreceive"
)

func MakeReceiveSignalEndpoint(dependencies *container.DependenciesContainer) endpoint.Endpoint {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var body model.ReqPostReceiveSignal

		req := request.(*http.Request)

		log.Printf("starting impl ReceiveEmergency endpoint req: %v", req)

		headersIntoStruct, err := util.ExtractHeadersFromRequest(req)

		if err != nil {
			return util.BuildErrorResponse(err, &headersIntoStruct)
		}

		err = json.NewDecoder(req.Body).Decode(&body)

		if err != nil {
			err = errcatalogs.MakeBadRequestResponseError(err.Error())
			return util.BuildErrorResponse(err, &headersIntoStruct)
		}

		signalResponse := dependencies.ReceiveSignalPort.(signalreceive.SignalReceive).SignalReceive(&body)

		if signalResponse.Err != nil {
			return util.BuildErrorResponse(signalResponse.Err, &headersIntoStruct)
		}
		return *signalResponse, nil
	}
}
