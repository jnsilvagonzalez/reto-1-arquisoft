package sns

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-xray-sdk-go/xray"
	"signalReceiver/container"
	"signalReceiver/domain/repositories"
	"signalReceiver/infrastructure/http/model"
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

func (s *SNSRepositoryImpl) PublishMessage(signal *model.ReqPostReceiveSignal, emergency bool) (model.ResPostReceiveSignal, error) {

	message := Message{
		Default: fmt.Sprintf("%+v", signal),
	}
	messageBytes, _ := json.Marshal(message)
	messageStr := string(messageBytes)

	ctx, seg := xray.BeginSegment(context.Background(), "PublishMessage")
	defer seg.Close(nil)

	var topic string

	if emergency {
		topic = "arn:aws:sns:us-east-1:992382691015:emergency"
	} else {
		topic = "arn:aws:sns:us-east-1:992382691015:signals"
	}

	err := xray.Capture(ctx, "Publish.Emergency", func(ctx context.Context) error {
		_, err := s.clientSNS.PublishWithContext(ctx, &sns.PublishInput{
			TopicArn:         aws.String(topic),
			Message:          aws.String(messageStr),
			MessageStructure: aws.String("json"),
		})
		return err // Devuelve solo el error.
	})

	if err != nil {
		return model.ResPostReceiveSignal{
			Err: err,
		}, nil
	}

	return model.ResPostReceiveSignal{
		RqUID:   signal.RqUID,
		Err:     nil,
		Message: "RECEIVED",
	}, nil

}
