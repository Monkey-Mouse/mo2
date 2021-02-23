package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID       primitive.ObjectID `json:"id,omitempty" example:"xxxxxxxxxxxxxx==" bson:"_id,omitempty"`
	ParentID primitive.ObjectID `json:"parent_id,omitempty" example:"xxxxxxxxxxxxxx==" bson:"parent_id,omitempty"`
	Name     string             `json:"name,omitempty" example:"records" bson:"name,omitempty"`
}
type CategoryUser struct {
	ID     primitive.ObjectID `json:"id,omitempty" example:"xxxxxxxxxxxxxx==" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"user_id,omitempty" example:"xxxxxxxxxxxxxx==" bson:"user_id,omitempty"`
}

func (c *Category) UpdateParent(parent Category) {
	c.ParentID = parent.ID
}
func (c *Category) UpdateName(name string) {
	c.Name = name
}
func (c *Category) Init() {
	c.ID = primitive.NewObjectID()
}
