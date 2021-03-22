package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbName = "github.com/Monkey-Mouse/mo2"

func connectMongoDB() {
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI(os.Getenv("MO2_MONGO_URL"))

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
