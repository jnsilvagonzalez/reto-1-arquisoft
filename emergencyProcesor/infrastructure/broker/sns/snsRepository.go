package sns

import (
	"context"
	container "emergencyProcesor/conteiner"
	"emergencyProcesor/domain/model"
	"emergencyProcesor/domain/repositories"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

type SNSRepositoryImpl struct {
	clientSNS snsiface.SNSAPI
}

type Message struct {
	Default string `json:"default"`
}

func NewSNSRepository(dependencies *container.DependenciesContainer) repositories.BrokerRepository {
	return &SNSRepositoryImpl{
		clientSNS: dependencies.SNSClient.(snsiface.SNSAPI),
	}
}

func (s SNSRepositoryImpl) PublishMessage(signal *[]model.ResponseRules) (model.ResEmergencySignal, error) {
	signalStr, err := json.Marshal(signal)

	if err != nil {
		return model.ResEmergencySignal{
			Err: err,
		}, nil
	}

	message := Message{
		Default: string(signalStr),
	}

	messageBytes, _ := json.Marshal(message)
	messageStr := string(messageBytes)

	_, err = s.clientSNS.PublishWithContext(context.Background(), &sns.PublishInput{
		TopicArn:         aws.String("arn:aws:sns:us-east-1:992382691015:actions"),
		Message:          aws.String(messageStr),
		MessageStructure: aws.String("json"),
	})

	if err != nil {
		return model.ResEmergencySignal{
			Err: err,
		}, nil
	}
	return model.ResEmergencySignal{
		RqUID:   "",
		Err:     nil,
		Message: "PROCESSED",
	}, nil
}
