package createschedule

import (
	"signalReceiver/container"
	"signalReceiver/domain/repositories"
	"signalReceiver/infrastructure/http/model"
)

type EmergencyReceive interface {
	EmergencyReceive(emergency *model.ReqPostReceiveSignal) *model.ResPostReceiveSignal
}

type emergencyReceive struct {
	dependencies *container.DependenciesContainer
}

func NewEmergencyReceive(dependencies *container.DependenciesContainer) EmergencyReceive {
	return &emergencyReceive{
		dependencies: dependencies,
	}
}

func (er *emergencyReceive) EmergencyReceive(reqDto *model.ReqPostReceiveSignal) *model.ResPostReceiveSignal {

	response, err := (er.dependencies.BrokerRepository.(repositories.BrokerRepository)).
		PublishMessage(reqDto, true)

	if err != nil {
		return &model.ResPostReceiveSignal{
			Err: err,
		}
	}
	return &response
}
