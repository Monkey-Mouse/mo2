package loghelper

import (
	"context"
	"log"
	"os"

	"github.com/Monkey-Mouse/mo2log/service/logservice"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

var client logservice.LogServiceClient

func init() {
	conn, err := grpc.Dial(":"+os.Getenv("LOG_PORT"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	client = logservice.NewLogServiceClient(conn)
}

type Log struct {
	Operator             primitive.ObjectID
	Operation            int32
	OperationTarget      primitive.ObjectID
	LogLevel             logservice.LogModel_Level
	ExtraMessage         string
	OperationTargetOwner primitive.ObjectID
}

func LogInfo(log Log) error {
	_, err := client.Log(context.TODO(), &logservice.LogModel{
		Operator:             log.Operator[:],
		Operation:            log.Operation,
		OperationTarget:      log.OperationTarget[:],
		OperationTargetOwner: log.OperationTargetOwner[:],
		LogLevel:             logservice.LogModel_INFO,
		ExtraMessage:         log.ExtraMessage,
	})
	return err
}
