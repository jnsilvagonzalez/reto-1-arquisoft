package emergencyreceive

import (
	"signalReceiver/container"
	"signalReceiver/infrastructure/http/model"
	emergencyreceive "signalReceiver/usecases/emergencyreceive"
)

type EmergencyReceive interface {
	EmergencyReceive(*model.ReqPostReceiveSignal) *model.ResPostReceiveSignal
}

type emergencyReceive struct {
	dependencies *container.DependenciesContainer
}

func NewEmergencyReceive(dependencies *container.DependenciesContainer) EmergencyReceive {
	return &emergencyReceive{
		dependencies: dependencies,
	}
}

func (cs *emergencyReceive) EmergencyReceive(dtoReq *model.ReqPostReceiveSignal) *model.ResPostReceiveSignal {
	return cs.dependencies.EmergencyReceive.(emergencyreceive.EmergencyReceive).EmergencyReceive(dtoReq)
}
