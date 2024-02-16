package signalreceive

import (
	"signalReceiver/container"
	"signalReceiver/infrastructure/http/model"
	"signalReceiver/usecases/signalreceive"
)

type SignalReceive interface {
	SignalReceive(*model.ReqPostReceiveSignal) *model.ResPostReceiveSignal
}

type signalReceive struct {
	dependencies *container.DependenciesContainer
}

func NewSignalReceive(dependencies *container.DependenciesContainer) SignalReceive {
	return &signalReceive{
		dependencies: dependencies,
	}
}

func (cs *signalReceive) SignalReceive(dtoReq *model.ReqPostReceiveSignal) *model.ResPostReceiveSignal {
	return cs.dependencies.SignalReceive.(signalreceive.SignalReceive).SignalReceive(dtoReq)
}
