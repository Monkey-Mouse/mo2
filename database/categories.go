package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mo2/server/model"
)

var catCol = GetCollection("category")

func UpsertCategory(c *model.Category) {
	if c.ID.IsZero() {
		c.Init()
	}
	_, err := catCol.UpdateOne(context.TODO(), bson.D{{"_id", c.ID}},
		bson.M{"$set": bson.M{"parent_id": c.ParentID, "name": c.Name}},
		options.Update().SetUpsert(true))
	if err != nil {
		log.Fatal(err)
	}
}
func FindSubCategories(c model.Category) (cs []model.Category) {
	results, err := catCol.Find(context.TODO(), bson.M{"parent_id": c.ID})
	if err != nil {
		log.Fatal(err)
	}
	if err = results.All(context.TODO(), &cs); err != nil {
		log.Fatal(err)
	}
	return
}
