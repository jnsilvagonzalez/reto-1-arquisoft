package signalreceive

import (
	"signalReceiver/container"
	"signalReceiver/domain/repositories"
	"signalReceiver/infrastructure/http/model"
)

type SignalReceive interface {
	SignalReceive(emergency *model.ReqPostReceiveSignal) *model.ResPostReceiveSignal
}

type signalReceive struct {
	dependencies *container.DependenciesContainer
}

func NewSignalReceive(dependencies *container.DependenciesContainer) SignalReceive {
	return &signalReceive{
		dependencies: dependencies,
	}
}

func (s signalReceive) SignalReceive(reqDto *model.ReqPostReceiveSignal) *model.ResPostReceiveSignal {
	response, err := (s.dependencies.BrokerRepository.(repositories.BrokerRepository)).
		PublishMessage(reqDto, true)

	if err != nil {
		return &model.ResPostReceiveSignal{
			Err: err,
		}
	}
	return &response
}
