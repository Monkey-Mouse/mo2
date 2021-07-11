package database

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbName = "mo2"

func connectMongoDB() {
	// 设置客户端连接配置
	conn := os.Getenv("MO2_MONGO_URL")
	if len(conn) == 0 {
		conn = "mongodb://localhost"
	}
	clientOptions := options.Client().ApplyURI(conn)

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	globalClient = client
	if err != nil {
		panic(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB!")

}

func disconnectMongoDB() {
	err := globalClient.Disconnect(context.TODO())

	if err != nil {
		panic(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
