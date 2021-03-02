package database

import (
	"context"
	"log"
	"mo2/dto"
	"mo2/server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var blogCol *mongo.Collection = GetCollection("blog")
var draftCol *mongo.Collection = GetCollection("draft")

func ensureBlogIndex() {
	blogCol.Indexes().CreateMany(context.TODO(), append([]mongo.IndexModel{
		{Keys: bson.M{"ket_words": 1}},
	}, model.IndexModels...))
}
func chooseCol(isDraft bool) (col *mongo.Collection) {
	col = draftCol
	if !isDraft {
		col = blogCol
	}
	return
}

// InsertBlog insert
func insertBlog(b *model.Blog, isDraft bool) (success bool) {
	b.Init()
	col := chooseCol(isDraft)
	success = true
	if _, err := col.InsertOne(context.TODO(), b); err != nil {
		log.Println(err)
		success = false
	}
	return
}

// upsertBlog
func upsertBlog(b *model.Blog, isDraft bool) (success bool) {
	col := chooseCol(isDraft)
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
			"categories":  b.CategoryIDs,
			"author_id":   b.AuthorID,
		}}},
		options.Update().SetUpsert(true),
	)
	success = true
	if err != nil {
		log.Println(err)
		success = false
	}
	if !isDraft {
		log.Println("发布时删除草稿" + b.ID.String())
		deleteBlog(*b, true)
	}
	if result.UpsertedCount != 0 {
		log.Println("新建文章" + b.ID.String())
	}
	return
}

// deleteBlog set flag of blog or draft to isDeleted
func deleteBlog(b model.Blog, isDraft bool) (success bool) {
	success = true
	res, err := chooseCol(isDraft).DeleteMany(context.TODO(), bson.M{"_id": b.ID})
	if err != nil {
		log.Fatal(err)
	}
	if res.DeletedCount == 0 {
		success = false
	}
	return
}

// UpsertBlog upsert blog or draft
func UpsertBlog(b *model.Blog, isDraft bool) (success bool) {

	if b.ID == primitive.NilObjectID {
		success = insertBlog(b, isDraft)
	} else {
		success = upsertBlog(b, isDraft)
	}
	return
}

//find blog by user
func FindBlogsByUser(u dto.LoginUserInfo, filter model.Filter) (b []model.Blog) {
	return FindBlogsByUserId(u.ID, filter)
}

//find blog by userId
func FindBlogsByUserId(id primitive.ObjectID, filter model.Filter) (b []model.Blog) {
	col := chooseCol(filter.IsDraft)
	opts := options.Find().SetSort(bson.D{{"entity_info", 1}})
	cursor, err := col.Find(context.TODO(), bson.M{"author_id": id, "entity_info.isdeleted": filter.IsDeleted}, opts)
	err = cursor.All(context.TODO(), &b)
	if err != nil {
		log.Fatal(err)
	}
	return
}

//find blog by id
func FindBlogById(id primitive.ObjectID, isDraft bool) (b model.Blog) {
	col := chooseCol(isDraft)
	err := col.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&b)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
	return
}

//find blog
func FindAllBlogs(filter model.Filter) (b []model.Blog) {
	col := chooseCol(filter.IsDraft)
	opts := options.Find().SetSort(bson.D{{"entity_info", 1}})
	cursor, err := col.Find(context.TODO(), bson.D{{"entity_info.isdeleted", filter.IsDeleted}}, opts)
	err = cursor.All(context.TODO(), &b)
	if err != nil {
		log.Fatal(err)
	}
	return
}
