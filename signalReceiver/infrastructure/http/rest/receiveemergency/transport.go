package receiveemergency

import (
	"context"
	"encoding/json"
	"net/http"
	"signalReceiver/infrastructure/http/model"
)

func DecodeRequest(_ context.Context, receiveEmergencyRequest *http.Request) (interface{}, error) {
	return receiveEmergencyRequest, nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, receiveEmergencyResponse interface{}) error {
	switch res := receiveEmergencyResponse.(type) {
	case model.ErrorCodeModel:
		w.WriteHeader(res.HTTPCode)
	case model.MsgRsHdrModel:
		w.WriteHeader(res.Status.ServerStatusCode)
	case model.ResPostReceiveSignal:
		w.WriteHeader(http.StatusCreated)
		return json.NewEncoder(w).Encode(receiveEmergencyResponse)
	}
	if receiveEmergencyResponse != nil {
		return json.NewEncoder(w).Encode(receiveEmergencyResponse)
	}
	return nil
}
