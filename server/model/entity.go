package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Praiseable 可被点赞的
type Praiseable struct {
	Up       uint64 `json:"up" bson:"up"`
	Down     uint64 `json:"down" bson:"down"`
	Weighted uint64 `json:"weighted" bson:"weighted"`
}

// Entity example
type Entity struct {
	CreateTime time.Time `json:"createTime" example:"2020-10-1" bson:"create_time,omitempty"`
	UpdateTime time.Time `json:"updateTime" example:"2020-10-1" bson:"update_time,omitempty"`
	IsDeleted  bool      `json:"is_deleted,omitempty" example:"true" bson:"is_deleted"`
}

// IndexModels Index Models to index entity
var IndexModels = []mongo.IndexModel{
	{Keys: bson.M{"entity_info.create_time": -1}},
	{Keys: bson.M{"entity_info.update_time": -1}},
	{Keys: bson.M{"entity_info.is_deleted": 1}},
}

// InitEntity init new entity
func InitEntity() Entity {
	t := time.Now()
	return Entity{t, t, false}
}

// Create create entity
func (e *Entity) Create() {
	e.CreateTime = time.Now()
	e.UpdateTime = e.CreateTime
	e.IsDeleted = false
}

// Set set entity with exist create time
func (e *Entity) Set(createTime time.Time) {
	e.CreateTime = createTime
	e.UpdateTime = e.CreateTime
	e.IsDeleted = false
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
