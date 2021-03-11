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

func getPaginationOption(page int64, pagesize int64) *options.FindOptions {
	return options.Find().SetSkip(page * pagesize).SetLimit(pagesize)
}

//TODO implement
//get max ID in collection
func GetMaxID(collection string) int {
	return 1

}
