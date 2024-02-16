package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	snsAWS "github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"os/signal"
	"signalReceiver/container"
	"signalReceiver/infrastructure/broker/sns"
	"signalReceiver/infrastructure/http/rest/handler"
	emergencyreceiveport "signalReceiver/portinterface/emergencyreceive"
	signalreceiveport "signalReceiver/portinterface/signalreceive"
	emergencyreceive "signalReceiver/usecases/emergencyreceive"
	"signalReceiver/usecases/signalreceive"
	"syscall"
)

func main() {

	var (
		httpAddr = flag.String("http.addr", ":8126", "HTTP listen address")
	)

	dependenciesContainer := &container.DependenciesContainer{}
	dependenciesContainer.SNSClient = newSNSSession()
	dependenciesContainer.BrokerRepository = sns.NewSNSRepository(dependenciesContainer)
	dependenciesContainer.EmergencyReceive = emergencyreceive.NewEmergencyReceive(dependenciesContainer)
	dependenciesContainer.ReceiveEmergencyPort = emergencyreceiveport.NewEmergencyReceive(dependenciesContainer)
	dependenciesContainer.SignalReceive = signalreceive.NewSignalReceive(dependenciesContainer)
	dependenciesContainer.ReceiveSignalPort = signalreceiveport.NewSignalReceive(dependenciesContainer)
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	var h = handler.MakeHTTPHandler(dependenciesContainer)

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		_ = logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	_ = logger.Log("exit", <-errs)
}

func newSNSSession() snsiface.SNSAPI {
	config := &aws.Config{
		Region: aws.String("us-east-1"),
	}
	sess := session.Must(session.NewSession(config))
	snsSvc := snsAWS.New(sess)
	xray.AWS(snsSvc.Client)
	return snsSvc
}
