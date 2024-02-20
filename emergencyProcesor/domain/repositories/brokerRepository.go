package repositories

import "emergencyProcesor/domain/model"

type BrokerRepository interface {
	PublishMessage(signal *[]model.ResponseRules) (model.ResEmergencySignal, error)
}
