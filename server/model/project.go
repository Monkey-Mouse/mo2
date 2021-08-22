package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Project struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Name        string               `bson:"name"`
	Tags        []string             `bson:"tags"`
	OwnerID     primitive.ObjectID   `bson:"owner_id"`
	ManagerIDs  []primitive.ObjectID `bson:"manager_i_ds"`
	MemberIDs   []primitive.ObjectID `bson:"member_i_ds"`
	BlogIDs     []primitive.ObjectID `bson:"blog_i_ds"`
	Description string               `bson:"description"`
	Avatar      string               `bson:"avatar"`
	EntityInfo  Entity               `bson:"entity_info"`
}
