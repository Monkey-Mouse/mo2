package database

import (
	"context"
	"fmt"
	"log"
	"mo2/dto"
	"mo2/mo2utils"
	"mo2/server/model"

	"mo2/mo2utils/mo2errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var blogCol *mongo.Collection = GetCollection("blog")
var draftCol *mongo.Collection = GetCollection("draft")

func init() {
	createBlogIndexes(blogCol)
	createBlogIndexes(draftCol)
}

func createBlogIndexes(col *mongo.Collection) {
	col.Indexes().CreateMany(context.TODO(), append([]mongo.IndexModel{
		{Keys: bson.M{"ket_words": 1}},
		{Keys: bson.M{"author_id": 1}},
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
func insertBlog(b *model.Blog, isDraft bool) (mErr mo2errors.Mo2Errors) {
	col := chooseCol(isDraft)
	if res, err := col.InsertOne(context.TODO(), b); err != nil {
		log.Println(err)
		mErr.InitError(err)
	} else {
		mErr.InitNoError("insert %v", res.InsertedID)
	}
	return
}

// upsertBlog
func upsertBlog(b *model.Blog, isDraft bool) (mErr mo2errors.Mo2Errors) {
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
	if err != nil {
		log.Println(err)
		mErr.InitError(err)
		return
	}
	if !isDraft {
		log.Println("发布时删除草稿" + b.ID.String())
		mErr = deleteBlogs(true, b.ID)
	}
	if result.UpsertedCount != 0 {
		mErr.InitNoError("新建文章" + b.ID.String())
	}
	return
}

// deleteBlogs set flag of blog or draft to isDeleted
func deleteBlogs(isDraft bool, blogIDs ...primitive.ObjectID) (mErr mo2errors.Mo2Errors) {
	if res, err := chooseCol(isDraft).DeleteMany(context.TODO(), bson.M{"_id": bson.M{"$in": blogIDs}}); err != nil {
		log.Println(err)
		mErr.InitError(err)
	} else {
		tip := fmt.Sprintf("delete %v %v blogs\n", res.DeletedCount, isDraft)
		log.Printf(tip)
		mErr.InitNoError(tip)
	}
	return
}

// UpsertBlog upsert blog or draft
func UpsertBlog(b *model.Blog, isDraft bool) (mErr mo2errors.Mo2Errors) {
	if b.ID == primitive.NilObjectID {
		b.Init()
		mErr = insertBlog(b, isDraft)
	} else {
		mErr = upsertBlog(b, isDraft)
	}
	if !isDraft {
		mo2utils.IndexBlog(b)
	}
	return
}

//find blog by user
func FindBlogsByUser(u dto.LoginUserInfo, filter model.Filter) (b []model.Blog) {
	return FindBlogsByUserId(u.ID, filter)
}
func getBlogListQueryOption() *options.FindOptions {
	return options.Find().SetSort(bson.D{{"entity_info.update_time", -1}}).SetProjection(bson.D{{"content", 0}})
}

//find blog by userId
func FindBlogsByUserId(id primitive.ObjectID, filter model.Filter) (b []model.Blog) {
	col := chooseCol(filter.IsDraft)
	opts := getBlogListQueryOption().SetSkip(int64(filter.Page * filter.PageSize)).SetLimit(int64(filter.PageSize))
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
func FindBlogs(filter model.Filter) (b []model.Blog) {
	col := chooseCol(filter.IsDraft)
	opts := getBlogListQueryOption().SetSkip(int64(filter.Page * filter.PageSize)).SetLimit(int64(filter.PageSize))
	f := bson.D{{"entity_info.isdeleted", filter.IsDeleted}}
	if filter.Ids != nil && len(filter.Ids) > 0 {
		f = bson.D{
			{"entity_info.isdeleted", filter.IsDeleted},
			{"_id", bson.M{"$in": filter.Ids}},
		}
	}
	cursor, err := col.Find(context.TODO(), f, opts)
	err = cursor.All(context.TODO(), &b)
	if err != nil {
		log.Fatal(err)
	}
	return
}
