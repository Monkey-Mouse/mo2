package demo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Trainer struct {
	Name    string
	Age     int
	Address string
}

var globalClient *mongo.Client

func GetClient() *mongo.Client {
	if globalClient == nil {
		connectMongoDB()
	}

	return globalClient
}

func connectMongoDB() {
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	globalClient = client
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

}

func disconnectMongoDB() {
	err := globalClient.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
func find() {
	collection := globalClient.Database("test").Collection("trainers")

	// Pass these options to the Find method
	findOptions := options.Find()
	//findOptions.SetLimit(2)

	// Here's an array in which you can store the decoded documents
	var results []*Trainer
	//"age","10"
	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Trainer
		err := cur.Decode(&elem)
		fmt.Println(elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	for _, result := range results {
		fmt.Println(result)

	}

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
}

func test() {
	collection := globalClient.Database("test").Collection("trainers")
	collection.DeleteMany(context.TODO(), Trainer{})

	ash := Trainer{"Ash", 10, "Pallet Town"}
	collection.InsertOne(context.TODO(), ash)

	filter := bson.D{{"name", "Ash"}}

	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	/*//ash := trainer{"Ash", 10, "Pallet Town"}
	misty := trainer{"Mist", 0, "Cerulean City"}
	brock := trainer{"Brock", 15, "Pewter City"}

	trainers := []interface{}{misty, brock}

	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

	_,err :=collection.InsertOne(context.TODO(),misty)
	if err!=nil{
		log.Fatal(err)
	}

	filter:=bson.D{{"name","Ash"}}

	var result trainer
	//err:=collection.Find(context.TODO(),filter)//.Decode(&result)


	//opts := options.Find().SetSort(bson.D{{"age", 1}})
	cursor, err := collection.Find(context.TODO(), filter)//, opts)
	if err != nil {
		log.Fatal(err)
	}  // get a list of all returned documents and print them out
	// see the mongo.Cursor documentation for more examples of using cursors
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); 	err != nil {     log.Fatal(err) }

	//
	for _, result := range results {
		fmt.Println(result)
	}


	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("this one: ",result)*/

}
func insert() {
	collection := globalClient.Database("test").Collection("trainers")
	var ash string
	insertResult, err := collection.InsertOne(context.TODO(), ash)
	//insertResult.
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertResult.InsertedID)
}
