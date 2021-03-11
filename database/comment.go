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

// GetComments get comments
func GetComments(articleID primitive.ObjectID, page int64, pagesize int64) (cs []model.Comment) {
	cursor, err := commentCol.Find(
		context.TODO(),
		bson.M{"article": articleID},
		getPaginationOption(page, pagesize))
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
