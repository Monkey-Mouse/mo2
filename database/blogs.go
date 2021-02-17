package database

import (
	"context"
	"log"
	"mo2/server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var blogCol *mongo.Collection = GetCollection("blog")

func ensureBlogIndex() {
	blogCol.Indexes().CreateMany(context.TODO(), append([]mongo.IndexModel{
		{Keys: bson.M{"ket_words": 1}},
	}, model.IndexModels...))
}

// AddBlog add
func AddBlog(b *model.Blog) {
	entity := model.InitEntity()
	b.EntityInfo = entity
	result, err := blogCol.InsertOne(context.TODO(), b)
	if err != nil {
		log.Fatal(err)
	}
	b.ID = result.InsertedID.(primitive.ObjectID)
}
