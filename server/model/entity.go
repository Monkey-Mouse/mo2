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

// Create create entity
func (e *Entity) Create() {
	e.CreateTime = time.Now()
	e.UpdateTime = e.CreateTime
}

// Update update entity
func (e *Entity) Update() {
	if IsTimeValid(e.CreateTime) {
		e.UpdateTime = time.Now()
	} else {
		e.Create()
	}
}

// Judge whether time exists
func IsTimeValid(time2 time.Time) (valid bool) {
	valid = false
	if time2.After(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)) {
		valid = true
	}
	return
}
