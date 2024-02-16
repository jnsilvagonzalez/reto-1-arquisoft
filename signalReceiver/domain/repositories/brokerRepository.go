package repositories

import "signalReceiver/infrastructure/http/model"

type BrokerRepository interface {
	PublishMessage(*model.ReqPostReceiveSignal, bool) (model.ResPostReceiveSignal, error)
}
