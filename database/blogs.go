package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mo2/dto"
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
func AddBlog(b *model.Blog) (success bool, err error) {
	entity := model.InitEntity()
	b.EntityInfo = entity
	if b.ID == primitive.NilObjectID {
		result, err := blogCol.InsertOne(context.TODO(), bson.M{
			"_id":         primitive.NewObjectID(),
			"author_id":   b.AuthorID,
			"title":       b.Title,
			"description": b.Description,
			"content":     b.Content,
			"cover":       b.Cover,
			"key_words":   b.KeyWords,
		})
		if err != nil {
			log.Fatal(err)
		}
		b.ID = result.InsertedID.(primitive.ObjectID)
	} else {
		result, err := blogCol.UpdateOne(
			context.TODO(),
			bson.D{{"_id", b.ID}},
			bson.D{{"$set", bson.M{
				"title":       b.Title,
				"description": b.Description,
				"content":     b.Content,
				"cover":       b.Cover,
				"key_words":   b.KeyWords,
			}}},
			//options.Update().SetUpsert(true),
		)
		if err != nil {
			log.Fatal(err)
		}
		if result.MatchedCount == 0 {
			log.Println("blog id do not match in database")
			success = false
			return success, err
		}
		b.ID = result.UpsertedID.(primitive.ObjectID)
	}
	success = true
	return
}

//find blog
func FindBlogs(u dto.LoginUserInfo) (b []model.Blog) {
	opts := options.Find().SetSort(bson.D{{"entity_info", 1}})
	cursor, err := blogCol.Find(context.TODO(), bson.D{{"author_id", u.ID}}, opts)
	err = cursor.All(context.TODO(), &b)
	if err != nil {
		log.Fatal(err)
	}
	return
}

//find blog
func FindAllBlogs() (b []model.Blog) {
	opts := options.Find().SetSort(bson.D{{"entity_info", 1}})
	cursor, err := blogCol.Find(context.TODO(), bson.D{{}}, opts)
	err = cursor.All(context.TODO(), &b)
	if err != nil {
		log.Fatal(err)
	}
	return
}
