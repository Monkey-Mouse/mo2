package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Blog example
// 若修改字段，需注意依赖此model的使用地方
// important: dto.QueryBlog UpsertBlog()
type Blog struct {
	ID          primitive.ObjectID   `json:"id,omitempty" example:"xxxxxxxxxxxxx==" bson:"_id,omitempty"`
	AuthorID    primitive.ObjectID   `json:"authorId,omitempty" example:"xxxxxxxxxxxxx==" bson:"author_id"`
	Title       string               `json:"title,omitempty" example:"mouse ❤ monkey" bson:"title,omitempty"`
	Description string               `json:"description,omitempty" example:"mouse ❤ monkey" bson:"description,omitempty"`
	Content     string               `json:"content,omitempty" example:"xxxx\nxxxx" bson:"content,omitempty"`
	EntityInfo  Entity               `json:"entityInfo,omitempty" bson:"entity_info,omitempty"`
	Cover       string               `json:"cover,omitempty" example:"https://xxx/xxx" bson:"cover,omitempty"`
	KeyWords    []string             `json:"keyWords,omitempty" example:"xxx,xxx" bson:"key_words,omitempty"`
	CategoryIDs []primitive.ObjectID `json:"categories,omitempty" bson:"categories,omitempty"`
	YDoc        string               `json:"y_doc,omitempty" bson:"y_doc,omitempty"`       // 用于collaboration
	IsYDoc      bool                 `json:"is_y_doc,omitempty" bson:"is_y_doc,omitempty"` // 用于collaboration
	YToken      primitive.ObjectID   `json:"y_token,omitempty" bson:"y_token,omitempty"`   //用于collaboration
}

type Filter struct {
	IsDraft   bool `json:"is_draft" example:"false"`
	IsDeleted bool `json:"is_deleted" example:"false"`
	Page      int  `json:"page"`
	PageSize  int  `json:"page_size"`
	Ids       []primitive.ObjectID
}

func (b *Blog) Init() {
	b.ID = primitive.NewObjectID()
	b.EntityInfo.Create()
}
func (b *Blog) Add2Category(categoryID primitive.ObjectID) {
	b.CategoryIDs = append(b.CategoryIDs, categoryID)
}
func (b *Blog) Add2Categories(categoryIDs []primitive.ObjectID) {
	b.CategoryIDs = append(b.CategoryIDs, categoryIDs...)
}
func (b *Blog) IsValid() (valid bool) {
	valid = true
	if b.ID.IsZero() {
		valid = false
	}
	return
}
