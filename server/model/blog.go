package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Blog example
type Blog struct {
	ID          primitive.ObjectID `json:"id,omitempty" example:"xxxxxxxxxxxxx==" bson:"_id,omitempty"`
	AuthorID    primitive.ObjectID `json:"authorId,omitempty" example:"xxxxxxxxxxxxx==" bson:"author_id"`
	Title       string             `json:"title,omitempty" example:"mouse ❤ monkey" bson:"title,omitempty"`
	Description string             `json:"description,omitempty" example:"mouse ❤ monkey" bson:"description,omitempty"`
	Content     string             `json:"content,omitempty" example:"xxxx\nxxxx" bson:"content,omitempty"`
	EntityInfo  Entity             `json:"entityInfo,omitempty" bson:"entity_info,omitempty"`
	Cover       string             `json:"cover,omitempty" example:"https://xxx/xxx" bson:"cover,omitempty"`
	KeyWords    []string           `json:"keyWords,omitempty" example:"xxx,xxx" bson:"key_words,omitempty"`
	IsPublished bool               `json:"is_published,omitempty" example:"false" bson:"is_published,omitempty"`
}

// Draft
type Draft struct {
	ID         primitive.ObjectID `json:"id,omitempty" example:"xxxxxxxxxxxxx==" bson:"_id,omitempty"`
	EntityInfo Entity             `json:"entityInfo,omitempty" bson:"entity_info,omitempty"`
	BlogObj    Blog               `json:"blog_obj" bson:"blog_obj"`
}

func (b *Blog) Init() {
	b.ID = primitive.NewObjectID()
	b.EntityInfo = InitEntity()
}
func (d *Draft) Init() {
	d.ID = primitive.NewObjectID()
	d.EntityInfo = InitEntity()
}
