package emergencyprocessor

import (
	container "emergencyProcesor/conteiner"
	"emergencyProcesor/domain/model"
	"emergencyProcesor/domain/repositories"
)

type EmergencyProcessor interface {
	EmergencyProcessor(emergency *model.ReqEmergencySignal) *model.ResEmergencySignal
}

type emergencyProcessor struct {
	dependencies *container.DependenciesContainer
}

func NewEmergencyProcessor(dependencies *container.DependenciesContainer) EmergencyProcessor {
	return &emergencyProcessor{
		dependencies: dependencies,
	}
}

func (e emergencyProcessor) EmergencyProcessor(emergency *model.ReqEmergencySignal) *model.ResEmergencySignal {

	responseRules, err := (e.dependencies.RulesRepository.(repositories.RulesRepository)).
		ValidateRules(emergency)

	response, err := (e.dependencies.BrokerRepository.(repositories.BrokerRepository)).
		PublishMessage(&responseRules)

	if err != nil {
		return &model.ResEmergencySignal{
			Err: err,
		}
	}
	return &response
}
