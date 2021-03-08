package dto

import (
	"mo2/server/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QueryBlog struct {
	ID          primitive.ObjectID   `json:"id,omitempty" example:"xxxxxxxxxxxxx==" `
	AuthorID    primitive.ObjectID   `json:"authorId,omitempty" example:"xxxxxxxxxxxxx==" `
	Title       string               `json:"title,omitempty" example:"mouse ❤ monkey" `
	Description string               `json:"description,omitempty" example:"mouse ❤ monkey" `
	EntityInfo  model.Entity         `json:"entityInfo,omitempty"`
	Cover       string               `json:"cover,omitempty" example:"https://xxx/xxx" `
	KeyWords    []string             `json:"keyWords,omitempty" example:"xxx,xxx"`
	CategoryIDs []primitive.ObjectID `json:"categories,omitempty" bson:"categories,omitempty"`
}
type QueryBlogs struct {
	blogs []QueryBlog `json:"blogs,omitempty"`
}

func MapBlog2QueryBlog(b model.Blog) (q QueryBlog) {
	q.ID = b.ID
	q.AuthorID = b.AuthorID
	q.Title = b.Title
	q.Description = b.Description
	q.EntityInfo = b.EntityInfo
	q.Cover = b.Cover
	q.KeyWords = b.KeyWords
	q.CategoryIDs = b.CategoryIDs
	return
}
