package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Entity example
type Entity struct {
	CreateTime time.Time `json:"createTime" example:"2020-10-1" bson:"create_time,omitempty"`
	UpdateTime time.Time `json:"updateTime" example:"2020-10-1" bson:"update_time,omitempty"`
}

// IndexModels Index Models to index entity
var IndexModels = []mongo.IndexModel{
	{Keys: bson.M{"entity_info.create_time": 1}},
	{Keys: bson.M{"entity_info.update_time": 1}},
}

// InitEntity init new entity
func InitEntity() Entity {
	t := time.Now()
	return Entity{t, t}
}
