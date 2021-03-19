package loghelper

import (
	"context"
	"mo2/database"
	"os"
	"testing"

	"github.com/Monkey-Mouse/mo2log/logmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestLogInfo(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skipf("Skip for ci")
		return
	}
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	id3 := primitive.NewObjectID()
	type args struct {
		log Log
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "text log info", args: args{Log{
			Operation:            1,
			Operator:             id1,
			OperationTarget:      id2,
			OperationTargetOwner: id3,
			ExtraMessage:         "Hello",
		}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LogInfo(tt.args.log); (err != nil) != tt.wantErr {
				t.Errorf("LogInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	col := database.GetCollection(os.Getenv("LOG_COL"))
	m := &logmodel.LogModel{}
	err := col.FindOneAndDelete(context.TODO(), bson.M{"operator_id": id1}).Decode(m)
	if err != nil {
		t.Errorf("cannot find document in db, err: %v", err)
	}
	if m.OperationTargetID != id2 || m.OperationTargetOwnerID != id3 {
		t.Errorf("Data doesn't match!")
	}
}
