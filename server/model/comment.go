package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Comment a comment
type Comment struct {
	ID         primitive.ObjectID `json:"id,omitempty" example:"xxxxxxxxxxxxxx==" bson:"_id,omitempty"`
	Content    string             `json:"content,omitempty" example:"a comment" bson:"content,omitempty"`
	EntityInfo Entity             `json:"entity_info,omitempty" bson:"entity_info,omitempty"`
	Praise     Praiseable         `json:"praise,omitempty" bson:"praise,omitempty"`
	Aurhor     primitive.ObjectID `json:"aurhor,omitempty" bson:"aurhor,omitempty"`
	Article    primitive.ObjectID `json:"article,omitempty" bson:"article,omitempty"`
	Subs       []Subcomment       `json:"subs" bson:"subs"`
}

// Subcomment level 2 comment
type Subcomment struct {
	ID         primitive.ObjectID `json:"id,omitempty" example:"xxxxxxxxxxxxxx==" bson:"_id,omitempty"`
	Content    string             `json:"content,omitempty" example:"a comment" bson:"content,omitempty"`
	Aurhor     primitive.ObjectID `json:"aurhor,omitempty" bson:"aurhor,omitempty"`
	EntityInfo Entity             `json:"entity_info,omitempty" bson:"entity_info,omitempty"`
	Praise     Praiseable         `json:"praise,omitempty" bson:"praise,omitempty"`
}
