package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Directory 归档
type Directory struct {
	ID       primitive.ObjectID   `json:"id,omitempty" example:"xxxxxxxxxxxxxx==" bson:"_id,omitempty"`
	ParentID primitive.ObjectID   `json:"parent_id,omitempty" example:"xxxxxxxxxxxxxx==" bson:"parent_id,omitempty"`
	Name     string               `json:"name,omitempty" example:"records" bson:"name,omitempty"`
	OwnerIDs []primitive.ObjectID `json:"owner_ids,omitempty"  bson:"owner_ids,omitempty"`
}

// UpdateParent 以父目录更新当前目录
func (c *Directory) UpdateParent(parent Directory) {
	c.ParentID = parent.ID
}

func (c *Directory) UpdateParentId(id primitive.ObjectID) {
	c.ParentID = id
}
func (c *Directory) UpdateName(name string) {
	c.Name = name
}
func (c *Directory) Init() {
	c.ID = primitive.NewObjectID()
}
func (c *Directory) InitWithName(name string) {
	c.Init()
	c.UpdateName(name)
}
func (c *Directory) InitWithNameAndParent(name string, parentId primitive.ObjectID) {
	c.Init()
	c.UpdateName(name)
	c.UpdateParentId(parentId)

}
func (c *Directory) IsValid() (valid bool) {
	valid = true
	if c.ID.IsZero() {
		valid = false
	}
	return
}
