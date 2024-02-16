package container

type DependenciesContainer struct {
	ReceiveEmergencyPort interface{}
	ReceiveSignalPort    interface{}

	EmergencyReceive interface{}
	SignalReceive    interface{}

	BrokerRepository interface{}
	SNSClient        interface{}
}
