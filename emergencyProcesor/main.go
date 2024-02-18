package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sns"
	snsAWS "github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-xray-sdk-go/xray"
	"log"
	"os"
)

var table string

type Rule struct {
	VehiculoID string
	ReglaId    string
	Actions    []Action
}

type Action struct {
	ActionId  string `json:"actionId"`
	Parameter string `json:"parameter"`
	Value     string `json:"value"`
}

type Message struct {
	Default string `json:"default"`
}

type SnsMessage struct {
	RqUID     string `json:"RqUID"`
	IdVehicle string `json:"IdVehicle"`
	Speed     int32  `json:"StatusCode"`
	Address   string `json:"Address"`
	Latitude  string `json:"Severity"`
	Longitude string `json:"Longitude,omitempty"`
}

func HandleRequest(ctx context.Context, snsEvent events.SNSEvent) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		fmt.Println("Error al iniciar sesión en AWS:", err)
		os.Exit(1)
	}

	svc := dynamodb.New(sess)
	xray.AWS(svc.Client)

	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		fmt.Println("received message from sns", snsRecord.MessageID, "with body", snsRecord.Message)
		fmt.Println("storing message info to dynamodb table", table)

		snsMessage := SnsMessage{}

		err := json.Unmarshal([]byte(snsRecord.Message), &snsMessage)

		fmt.Println("ID-VEHICULO", snsMessage.IdVehicle)

		if err != nil {
			log.Println(err)
			return
		}

		table = "Rules_table"
		searchAuthor := snsMessage.IdVehicle
		params := &dynamodb.QueryInput{
			TableName: aws.String(table),
			KeyConditions: map[string]*dynamodb.Condition{
				"VehiculoID": {
					ComparisonOperator: aws.String("EQ"),
					AttributeValueList: []*dynamodb.AttributeValue{
						{
							S: aws.String(searchAuthor),
						},
					},
				},
			},

			ProjectionExpression: aws.String("Actions"),
		}

		resp, err := svc.QueryWithContext(ctx, params)
		if err != nil {
			fmt.Println("Query API call failed:", err)
			os.Exit(1)
		}

		itemsStr, err := json.Marshal(resp.Items)

		if err != nil {
			fmt.Println("Query API call failed:", err)
			os.Exit(1)
		}

		fmt.Printf("ITEMS ->: %s", string(itemsStr))

		var rules []Rule
		err = dynamodbattribute.UnmarshalListOfMaps(resp.Items, &rules)
		if err != nil {
			fmt.Println("Failed to unmarshal Query result items:", err)
			os.Exit(1)
		}

		// Imprime los resultados
		for _, rule := range rules {
			fmt.Printf("VehiculoID: %s, ReglaId: %s, Actions: %+v\n", rule.VehiculoID, rule.ReglaId, rule.Actions)
		}

		snsSvc := snsAWS.New(sess)
		xray.AWS(snsSvc.Client)

		message := Message{
			Default: fmt.Sprintf("%+v", rules),
		}

		messageBytes, _ := json.Marshal(message)
		messageStr := string(messageBytes)

		fmt.Printf("MESSAGE FOR ACTIONS  ->: %s", messageStr)

		_, err = snsSvc.PublishWithContext(ctx, &sns.PublishInput{
			TopicArn:         aws.String("arn:aws:sns:us-east-1:992382691015:actions"),
			Message:          aws.String(messageStr),
			MessageStructure: aws.String("json"),
		})

		if err != nil {
			fmt.Println("Failed to send message to topic of actions:", err)
			os.Exit(1)
		}

	}
}

func main() {
	lambda.Start(HandleRequest)
}
