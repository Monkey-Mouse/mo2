package database

import (
	"context"
	"github.com/Monkey-Mouse/mo2/mo2utils/mo2errors"
	"github.com/Monkey-Mouse/mo2/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const groupStr = "group"

var GroupCol = GetCollection(groupStr)

// UpsertGroup upsert group
func UpsertGroup(group model.Group) mo2errors.Mo2Errors {
	if _, err := GroupCol.UpdateOne(context.TODO(), group, options.Update().SetUpsert(true)); err != nil {
		return mo2errors.InitError(err)
	}
	return mo2errors.InitNoError("upsert success")
}

// DeleteGroupByID  delete group by id
func DeleteGroupByID(id primitive.ObjectID) mo2errors.Mo2Errors {
	if _, err := GroupCol.DeleteOne(context.TODO(), bson.M{"_id": id}); err != nil {
		return mo2errors.InitError(err)
	}
	return mo2errors.InitNoError("delete success")
}
