package repositories

import "emergencyProcesor/domain/model"

type RulesRepository interface {
	ValidateRules(*model.ReqEmergencySignal) ([]model.ResponseRules, error)
}
