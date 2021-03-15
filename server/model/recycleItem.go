package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// RecycleItem 回收站中的对象信息，记录加入回收站时间和预计被删除时间，以及处理函数
type RecycleItem struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ItemID     primitive.ObjectID `json:"item_id,omitempty" bson:"item_id,omitempty"`
	CreateTime time.Time          `json:"create_time,omitempty" example:"2020-10-1" bson:"create_time,omitempty"`
	DeleteTime time.Time          `json:"delete_time,omitempty" example:"2020-10-1" bson:"delete_time,omitempty"`
	Handler    string             `json:"handler,omitempty" example:"blog" bson:"handler,omitempty"`
}
