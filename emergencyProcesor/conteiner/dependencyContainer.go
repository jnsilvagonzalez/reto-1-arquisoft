package container

type DependenciesContainer struct {
	EmergencyProcessorPort interface{}

	EmergencyProcessor interface{}

	BrokerRepository interface{}
	RulesRepository  interface{}
	SNSClient        interface{}
	DynamoClient     interface{}
}
