package dynamodb

import (
	"context"
	container "emergencyProcesor/conteiner"
	"emergencyProcesor/domain/model"
	"emergencyProcesor/domain/repositories"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"os"
)

type DynamoRepositoryImpl struct {
	clientDynamo dynamodbiface.DynamoDBAPI
}

type Message struct {
	Default string `json:"default"`
}

func NewDynamoRepository(dependencies *container.DependenciesContainer) repositories.RulesRepository {
	return &DynamoRepositoryImpl{
		clientDynamo: dependencies.DynamoClient.(dynamodbiface.DynamoDBAPI),
	}
}

func (d DynamoRepositoryImpl) ValidateRules(signal *model.ReqEmergencySignal) ([]model.ResponseRules, error) {

	table := "Rules_table"
	searchVehicle := signal.IdVehicle
	params := &dynamodb.QueryInput{
		TableName: aws.String(table),
		KeyConditions: map[string]*dynamodb.Condition{
			"VehiculoID": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(searchVehicle),
					},
				},
			},
		},

		ProjectionExpression: aws.String("Actions"),
	}

	resp, err := d.clientDynamo.QueryWithContext(context.Background(), params)
	if err != nil {
		fmt.Println("Query API call failed:", err)
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("Query API call failed:", err)
		os.Exit(1)
	}

	var rules []model.ResponseRules
	err = dynamodbattribute.UnmarshalListOfMaps(resp.Items, &rules)
	if err != nil {
		fmt.Println("Failed to unmarshal Query result items:", err)
		os.Exit(1)
	}

	return rules, nil

}
