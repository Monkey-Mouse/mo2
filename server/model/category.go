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
func (c *Category) UpdateParentId(id primitive.ObjectID) {
	c.ParentID = id
}
func (c *Category) UpdateName(name string) {
	c.Name = name
}
func (c *Category) Init() {
	c.ID = primitive.NewObjectID()
}
func (c *Category) InitWithName(name string) {
	c.Init()
	c.UpdateName(name)
}
func (c *Category) InitWithNameAndParent(name string, parentId primitive.ObjectID) {
	c.Init()
	c.UpdateName(name)
	c.UpdateParentId(parentId)

}
