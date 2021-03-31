package database

import (
	"context"
	"github.com/Monkey-Mouse/mo2/mo2utils/mo2errors"
	"github.com/Monkey-Mouse/mo2/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const accessManagerStr = "accessManager"

var AccessManagerCol = GetCollection(accessManagerStr)

// UpsertManager upsert Manager
func UpsertManager(manager model.AccessManager) mo2errors.Mo2Errors {
	if _, err := AccessManagerCol.UpdateOne(context.TODO(), bson.M{"_id": manager.ID}, bson.M{"$set": bson.M{
		"entity_info": manager.EntityInfo,
		"role_map":    manager.RoleMap,
	}}, options.Update().SetUpsert(true)); err != nil {
		return mo2errors.InitError(err)
	}
	return mo2errors.InitNoError("upsert success")
}

// FindManager find Manager
func FindManager(id primitive.ObjectID) (manager model.AccessManager, mErr mo2errors.Mo2Errors) {
	if err := AccessManagerCol.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&manager); err != nil {
		if err == mongo.ErrNoDocuments {
			mErr.InitCode(mo2errors.Mo2NotFound)
			return
		}
		mErr.InitError(err)
		return
	}
	mErr.InitNoError("find one")
	return
}

// DeleteManagerByID  delete Manager by id
func DeleteManagerByID(id primitive.ObjectID) mo2errors.Mo2Errors {
	if _, err := AccessManagerCol.DeleteOne(context.TODO(), bson.M{"_id": id}); err != nil {
		return mo2errors.InitError(err)
	}
	return mo2errors.InitNoError("delete success")
}
