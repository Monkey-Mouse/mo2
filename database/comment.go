package database

import (
	"context"
	"log"
	"mo2/server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var commentCol *mongo.Collection = GetCollection("comment")

func init() {
	commentCol.Indexes().CreateMany(context.Background(), append([]mongo.IndexModel{
		{Keys: bson.M{"article": -1}},
		{Keys: bson.M{"author": -1}},
		{Keys: bson.M{"praise.up": -1}},
		{Keys: bson.M{"praise.weight": -1}},
		{Keys: bson.M{"subs.id": -1}},
	}, model.IndexModels...))
}

// Ge3tCommentNum 获取该文章下的一级评论数量
func GetCommentNum(articleID primitive.ObjectID) int64 {
	num, _ := commentCol.CountDocuments(context.TODO(),
		bson.M{"article": articleID})
	return num
}

// GetComments get comments
func GetComments(articleID primitive.ObjectID, page int64, pagesize int64) (cs []model.Comment) {
	cursor, err := commentCol.Find(
		context.TODO(),
		bson.M{"article": articleID},
		getPaginationOption(page, pagesize).SetSort(bson.M{"entity_info.update_time": -1}))
	if err != nil {
		log.Fatal(err)
	}
	cursor.All(context.TODO(), &cs)
	return
}

// UpsertComment 插入或更新评论
func UpsertComment(c *model.Comment) {
	if c.ID != primitive.NilObjectID {
		c.EntityInfo.Update()
		commentCol.UpdateOne(
			context.TODO(),
			bson.M{"_id": c.ID},
			bson.M{"$set": bson.M{"content": c.Content, "entity_info.updateTime": c.EntityInfo.UpdateTime}})
		return
	}
	c.Subs = []model.Subcomment{}
	c.EntityInfo.Create()
	re, err := commentCol.InsertOne(context.TODO(), c)
	if err != nil {
		log.Fatal(err)
	}
	c.ID = re.InsertedID.(primitive.ObjectID)
}

// UpsertSubComment as name
func UpsertSubComment(id primitive.ObjectID, c *model.Subcomment) {
	if c.ID != primitive.NilObjectID {
		c.EntityInfo.Update()
		commentCol.UpdateOne(
			context.TODO(),
			bson.D{{Key: "_id", Value: id}, {Key: "subs._id", Value: c.ID}},
			bson.M{"$set": bson.M{"subs.$.content": c.Content, "subs.$.entity_info.updateTime": c.EntityInfo.UpdateTime}})
		return
	}
	c.EntityInfo.Create()
	c.ID = primitive.NewObjectID()
	_, err := commentCol.UpdateOne(
		context.TODO(),
		bson.M{"_id": id},
		bson.M{"$push": bson.M{"subs": c}})
	if err != nil {
		log.Fatal(err)
	}
}
