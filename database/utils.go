package database

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var globalClient *mongo.Client

func GetClient() *mongo.Client {
	if globalClient == nil {
		connectMongoDB()
		CreateAccountIndex()
	}

	return globalClient
}
func GetCollection(colName string) *mongo.Collection {
	client := GetClient()
	return client.Database(dbName).Collection(colName)
}

func sanitizePagesize(pagesize int64) int64 {
	if pagesize > 100 {
		pagesize = 100
	}
	return pagesize
}

func getPaginationOption(page int64, pagesize int64) *options.FindOptions {
	return options.Find().SetSkip(page * pagesize).SetLimit(pagesize)
}

func If(condition bool, trueVal interface{}, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	} else {
		return falseVal
	}
}
