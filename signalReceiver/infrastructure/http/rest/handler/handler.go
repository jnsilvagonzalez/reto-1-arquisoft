package handler

import (
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"signalReceiver/container"
	"signalReceiver/infrastructure/http/rest/receiveemergency"
	"signalReceiver/infrastructure/http/rest/receivesignal"
)

const (
	URLReceiveEmergency = "/receiver/v1/api/emergency"
	URLReceiveSignal    = "/receiver/v1/api/signal"
	PostOperation       = "POST"
)

type Endpoints struct {
	ReceiveEmergencyEndpoint func(dependencies *container.DependenciesContainer) endpoint.Endpoint
	ReceiveSignalEndpoint    func(dependencies *container.DependenciesContainer) endpoint.Endpoint
}

func MakeServerEndpoints() Endpoints {
	return Endpoints{
		ReceiveEmergencyEndpoint: receiveemergency.MakeReceiveEmergencyEndpoint,
		ReceiveSignalEndpoint:    receivesignal.MakeReceiveSignalEndpoint,
	}
}

func MakeHTTPHandler(dependencies *container.DependenciesContainer) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints()
	var serverOptions []kithttp.ServerOption

	r.Methods(PostOperation).Path(URLReceiveEmergency).Handler(kithttp.NewServer(
		e.ReceiveEmergencyEndpoint(dependencies),
		receiveemergency.DecodeRequest,
		receiveemergency.EncodeResponse,
		serverOptions...,
	))

	r.Methods(PostOperation).Path(URLReceiveSignal).Handler(kithttp.NewServer(
		e.ReceiveEmergencyEndpoint(dependencies),
		receivesignal.DecodeRequest,
		receivesignal.EncodeResponse,
		serverOptions...,
	))

	return r
}
