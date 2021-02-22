package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mo2/dto"
	"mo2/server/model"
)

var blogCol *mongo.Collection = GetCollection("blog")
var draftCol *mongo.Collection = GetCollection("draft")

func ensureBlogIndex() {
	blogCol.Indexes().CreateMany(context.TODO(), append([]mongo.IndexModel{
		{Keys: bson.M{"ket_words": 1}},
	}, model.IndexModels...))
}

// InsertBlog insert
func insertBlog(b *model.Blog, isDraft bool) {
	b.Init()
	if isDraft {
		if _, err := draftCol.InsertOne(context.TODO(), b); err != nil {
			log.Fatal(err)
		}

	} else {
		if _, err := blogCol.InsertOne(context.TODO(), b); err != nil {
			log.Fatal(err)
		}
	}
}

// upsertBlog
func upsertBlog(b *model.Blog, isDraft bool) {
	col := draftCol
	if !isDraft {
		col = blogCol
	}
	b.EntityInfo.Update()
	result, err := col.UpdateOne(
		context.TODO(),
		bson.D{{"_id", b.ID}},
		bson.D{{"$set", bson.M{
			"entity_info": b.EntityInfo,
			"title":       b.Title,
			"description": b.Description,
			"content":     b.Content,
			"cover":       b.Cover,
			"key_words":   b.KeyWords,
		}}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Fatal(err)
	}
	if result.MatchedCount == 0 {
		log.Println("blog id do not match in database")
	}
}

// UpsertBlog upsert blog or draft
func UpsertBlog(b *model.Blog, isDraft bool) {

	if b.ID == primitive.NilObjectID {
		insertBlog(b, isDraft)
	} else {
		upsertBlog(b, isDraft)
	}
}

//find blog by user
func FindBlogsByUser(u dto.LoginUserInfo) (b []model.Blog) {
	opts := options.Find().SetSort(bson.D{{"entity_info", 1}})
	cursor, err := blogCol.Find(context.TODO(), bson.D{{"author_id", u.ID}}, opts)
	err = cursor.All(context.TODO(), &b)
	if err != nil {
		log.Fatal(err)
	}
	return
}

//find blog by id
func FindBlogById(id primitive.ObjectID) (b model.Blog) {
	err := blogCol.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&b)
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

//find draft
func FindAllDrafts() (d []model.Draft) {
	opts := options.Find().SetSort(bson.D{{"entity_info", 1}})
	cursor, err := draftCol.Find(context.TODO(), bson.D{{}}, opts)
	err = cursor.All(context.TODO(), &d)
	if err != nil {
		log.Fatal(err)
	}
	return
}

//find draft by id
func FindDraftById(id primitive.ObjectID) (d model.Draft) {
	err := draftCol.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&d)
	if err != nil {
		log.Fatal(err)
	}
	return
}
