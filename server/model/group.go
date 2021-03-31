package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Group struct {
	ID      primitive.ObjectID `json:"id" bson:"id"`
	OwnerID primitive.ObjectID `json:"owner_id" bson:"owner_id"`
}
