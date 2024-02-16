package receiveemergency

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
	"signalReceiver/portinterface/emergencyreceive"
)

func MakeReceiveEmergencyEndpoint(dependencies *container.DependenciesContainer) endpoint.Endpoint {

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

		emergencyResponse := dependencies.ReceiveEmergencyPort.(emergencyreceive.EmergencyReceive).EmergencyReceive(&body)

		if emergencyResponse.Err != nil {
			return util.BuildErrorResponse(emergencyResponse.Err, &headersIntoStruct)
		}
		return *emergencyResponse, nil
	}
}
