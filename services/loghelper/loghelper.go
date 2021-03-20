package loghelper

import (
	"context"
	"log"
	"os"

	"github.com/Monkey-Mouse/mo2log/service/logservice"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

func (l *LogClient) Init(portEnv string) {
	conn, err := grpc.Dial(":"+os.Getenv(portEnv), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	l.client = logservice.NewLogServiceClient(conn)
}

// LogClient as name
type LogClient struct {
	client logservice.LogServiceClient
}

type Log struct {
	Operator             primitive.ObjectID
	Operation            int32
	OperationTarget      primitive.ObjectID
	LogLevel             logservice.LogModel_Level
	ExtraMessage         string
	OperationTargetOwner primitive.ObjectID
}

func (l *LogClient) LogInfo(log Log) error {
	_, err := l.client.Log(context.TODO(), &logservice.LogModel{
		Operator:             log.Operator[:],
		Operation:            log.Operation,
		OperationTarget:      log.OperationTarget[:],
		OperationTargetOwner: log.OperationTargetOwner[:],
		LogLevel:             logservice.LogModel_INFO,
		ExtraMessage:         log.ExtraMessage,
	})
	return err
}
