package database

import "go.mongodb.org/mongo-driver/mongo"

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

//TODO implement
//get max ID in collection
func GetMaxID(collection string) int {
	return 1

}
