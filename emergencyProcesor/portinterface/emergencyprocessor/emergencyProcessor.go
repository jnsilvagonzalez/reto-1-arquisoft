package emergencyprocessor

import (
	container "emergencyProcesor/conteiner"
	"emergencyProcesor/domain/model"
	"emergencyProcesor/usecases/emergencyprocessor"
)

type EmergencyProcessor interface {
	EmergencyProcessor(signal *model.ReqEmergencySignal) *model.ResEmergencySignal
}

type emergencyProcessor struct {
	dependencies *container.DependenciesContainer
}

func NewEmergencyProcessor(dependencies *container.DependenciesContainer) EmergencyProcessor {
	return &emergencyProcessor{
		dependencies: dependencies,
	}
}

func (e emergencyProcessor) EmergencyProcessor(signal *model.ReqEmergencySignal) *model.ResEmergencySignal {
	return e.dependencies.EmergencyProcessor.(emergencyprocessor.EmergencyProcessor).EmergencyProcessor(signal)

}
