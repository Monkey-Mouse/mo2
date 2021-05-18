package loghelper

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/Monkey-Mouse/mo2log/helpers"
	"github.com/Monkey-Mouse/mo2log/logmodel"
	"github.com/Monkey-Mouse/mo2log/service/logservice"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

const (
	COMMENT   = 1
	LIKE_BLOG = 2
)

// Init log client
func (l *LogClient) Init(targetEnv string) {
	conn, err := grpc.Dial(os.Getenv(targetEnv), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	l.Client = logservice.BuildLogClient(conn)
}

// LogClient as name
type LogClient struct {
	Client logservice.LogClient
}

// ProtoToLog log proto to log model
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

// Log log struct
type Log struct {
	Operator             primitive.ObjectID
	Operation            int32
	OperationTarget      primitive.ObjectID
	LogLevel             logservice.LogModel_Level
	ExtraMessage         string
	OperationTargetOwner primitive.ObjectID
}

// LogInfo log at info level
func (l *LogClient) LogInfo(log Log) error {
	log.LogLevel = logservice.LogModel_INFO
	err := l.LogMsg(context.TODO(), log)
	return err
}

// LogMsg log a message
func (l *LogClient) LogMsg(ctx context.Context, log Log) error {
	_, err := l.Client.Log(ctx, &logservice.LogModel{
		Operator:             log.Operator[:],
		Operation:            log.Operation,
		OperationTarget:      log.OperationTarget[:],
		OperationTargetOwner: log.OperationTargetOwner[:],
		LogLevel:             log.LogLevel,
		ExtraMessage:         log.ExtraMessage,
	})
	return err
}

type errLogger struct {
	c *LogClient
}

// BuildWriter Create an io.Writer write logs to log service
func BuildWriter(target string) io.Writer {
	l := &LogClient{}
	l.Init(target)
	return &errLogger{l}
}

// Write write log
func (l *errLogger) Write(p []byte) (n int, err error) {
	var nilID [12]byte
	err = l.c.LogMsg(context.TODO(), Log{
		Operator:             nilID,
		OperationTarget:      nilID,
		OperationTargetOwner: nilID,
		Operation:            -1,
		ExtraMessage:         string(p),
		LogLevel:             logservice.LogModel_FATAL,
	})
	return
}
