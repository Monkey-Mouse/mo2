package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Blog example
// 若修改字段，需注意依赖此model的使用地方
// important: dto.QueryBlog UpsertBlog()
type Blog struct {
	ID          primitive.ObjectID `json:"id,omitempty" example:"xxxxxxxxxxxxx==" bson:"_id,omitempty"`
	AuthorID    primitive.ObjectID `json:"authorId,omitempty" example:"xxxxxxxxxxxxx==" bson:"author_id"`
	Title       string             `json:"title,omitempty" example:"mouse ❤ monkey" bson:"title,omitempty"`
	Description string             `json:"description,omitempty" example:"mouse ❤ monkey" bson:"description,omitempty"`
	Content     string             `json:"content,omitempty" example:"xxxx\nxxxx" bson:"content,omitempty"`
	EntityInfo  Entity             `json:"entityInfo,omitempty" bson:"entity_info,omitempty"`
	Cover       string             `json:"cover,omitempty" example:"https://xxx/xxx" bson:"cover,omitempty"`
	KeyWords    []string           `json:"keyWords,omitempty" example:"xxx,xxx" bson:"key_words,omitempty"`
	Categories  []Category         `json:"categories,omitempty" bson:"categories,omitempty"`
}

func (b *Blog) Init() {
	b.ID = primitive.NewObjectID()
	b.EntityInfo = InitEntity()
}
func (b *Blog) Add2Category(category Category) {
	b.Categories = append(b.Categories, category)
}
func (b *Blog) Add2Categories(categories []Category) {
	b.Categories = append(b.Categories, categories...)
}
func (b *Blog) IsValid() (valid bool) {
	valid = true
	if b.ID.IsZero() {
		valid = false
	}
	return
}
