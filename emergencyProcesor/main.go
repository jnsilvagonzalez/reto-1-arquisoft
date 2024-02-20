package main

import (
	"context"
	container "emergencyProcesor/conteiner"
	"emergencyProcesor/domain/model"
	sns2 "emergencyProcesor/infrastructure/broker/sns"
	dynamodb2 "emergencyProcesor/infrastructure/database/dynamodb"
	emergencyProcesorPort "emergencyProcesor/portinterface/emergencyprocessor"
	"emergencyProcesor/usecases/emergencyprocessor"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	snsAWS "github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-xray-sdk-go/xray"
	"os"
)

var (
	sess      = session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))
	dynamoSvc = dynamodb.New(sess)
	snsSvc    = snsAWS.New(sess)
)

func init() {
	xray.AWS(dynamoSvc.Client)
	xray.AWS(snsSvc.Client)
}

func HandleRequest(ctx context.Context, snsEvent events.SNSEvent) {

	for _, record := range snsEvent.Records {

		snsRecord := record.SNS
		fmt.Println("received message from sns", snsRecord.MessageID, "with body", snsRecord.Message)

		snsMessage := model.SnsMessage{}
		err := json.Unmarshal([]byte(snsRecord.Message), &snsMessage)

		if err != nil {
			fmt.Println("Failed to process message :", err)
			os.Exit(1)
		}

		dependenciesContainer := &container.DependenciesContainer{}
		dependenciesContainer.SNSClient = snsSvc
		dependenciesContainer.DynamoClient = dynamoSvc
		dependenciesContainer.BrokerRepository = sns2.NewSNSRepository(dependenciesContainer)
		dependenciesContainer.RulesRepository = dynamodb2.NewDynamoRepository(dependenciesContainer)
		dependenciesContainer.EmergencyProcessor = emergencyprocessor.NewEmergencyProcessor(dependenciesContainer)
		dependenciesContainer.EmergencyProcessorPort = emergencyProcesorPort.NewEmergencyProcessor(dependenciesContainer)

		reqEmergencySiganl := model.ReqEmergencySignal{
			RqUID:     snsMessage.RqUID,
			IdVehicle: snsMessage.IdVehicle,
			Speed:     snsMessage.Speed,
			Address:   snsMessage.Address,
			Latitude:  snsMessage.Latitude,
			Longitude: snsMessage.Longitude,
		}

		emergencyResponse := dependenciesContainer.EmergencyProcessorPort.(emergencyProcesorPort.EmergencyProcessor).EmergencyProcessor(&reqEmergencySiganl)

		if emergencyResponse.Err != nil {
			fmt.Println("Failed to send message to topic of actions:", emergencyResponse.Err)
			os.Exit(1)
		}

		fmt.Println("Message procesed:", emergencyResponse.Err)

	}
}

func main() {
	lambda.Start(HandleRequest)
}
