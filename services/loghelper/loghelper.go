package loghelper

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Monkey-Mouse/mo2log/helpers"
	"github.com/Monkey-Mouse/mo2log/logmodel"
	"github.com/Monkey-Mouse/mo2log/service/logservice"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

const COMMENT = 1

func (l *LogClient) Init(portEnv string) {
	conn, err := grpc.Dial(os.Getenv(portEnv), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	l.Client = logservice.NewLogServiceClient(conn)
}

// LogClient as name
type LogClient struct {
	Client logservice.LogServiceClient
}

func ProtoToLog(log *logservice.LogModel) *logmodel.LogModel {
	return &logmodel.LogModel{
		Operation:              log.Operation,
		OperatorID:             helpers.BytesToMongoID(log.Operator),
		OperationTargetID:      helpers.BytesToMongoID(log.OperationTarget),
		OperationTargetOwnerID: helpers.BytesToMongoID(log.OperationTargetOwner),
		LogLevel:               int32(log.LogLevel),
		ExtraMessage:           log.ExtraMessage,
		CreateTime:             time.Unix(0, log.CreateTime),
		UpdateTime:             time.Unix(0, log.UpdateTime),
		Processed:              log.Processed,
	}
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
	_, err := l.Client.Log(context.TODO(), &logservice.LogModel{
		Operator:             log.Operator[:],
		Operation:            log.Operation,
		OperationTarget:      log.OperationTarget[:],
		OperationTargetOwner: log.OperationTargetOwner[:],
		LogLevel:             logservice.LogModel_INFO,
		ExtraMessage:         log.ExtraMessage,
	})
	return err
}
