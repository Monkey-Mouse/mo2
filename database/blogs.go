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
func insertBlog(b *model.Blog) {
	b.Init()
	_, err := blogCol.InsertOne(context.TODO(), b)
	if err != nil {
		log.Fatal(err)
	}
}

// UpsertBlog upsert
func UpsertBlog(b *model.Blog) (success bool) {

	if b.ID == primitive.NilObjectID {
		insertBlog(b)
	} else {
		b.EntityInfo.Update()
		result, err := blogCol.UpdateOne(
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
		)
		if err != nil {
			log.Fatal(err)
		}
		if result.MatchedCount == 0 {
			log.Println("blog id do not match in database")
			success = false
			return success
		}
	}
	success = true
	return
}

// InsertDraft insert
func InsertDraft(d *model.Draft) {
	d.Init()
	draftCol.InsertOne(context.TODO(), d)
	return
}

// UpsertDraft upsert
func UpsertDraft(d *model.Draft) (success bool) {
	b := d.BlogObj
	if d.ID == primitive.NilObjectID {
		d.Init()
		// blog not exist
		if b.ID == primitive.NilObjectID {
			// insert new blog and update draft
			success = UpsertBlog(&b)
			d.BlogObj = b
			if !success {
				return
			}
		}
		InsertDraft(d)
	} else {
		d.EntityInfo.Update()
		result, err := draftCol.UpdateOne(
			context.TODO(),
			bson.D{{"_id", d.ID}},
			bson.D{{"$set", bson.M{
				"blog_obj":    d.BlogObj,
				"entity_info": d.EntityInfo,
			}}},
		)
		if err != nil {
			log.Fatal(err)
		}
		if result.MatchedCount == 0 {
			log.Println("blog id do not match in database")
			success = false
			return
		}
	}
	success = true
	return
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
