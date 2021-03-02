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
func (qs *QueryBlogs) Init(bs []model.Blog) {
	for _, b := range bs {
		qs.blogs = append(qs.blogs, MapBlog2QueryBlog(b))
		//qs.blogs[i].BlogMap(b)
	}
}
func (qs *QueryBlogs) Query(page int, pageSize int) (rs QueryBlogs, exist bool) {
	// check whether index out of scope
	numWith0 := len(qs.blogs) - 1
	begin := page * pageSize
	end := (page + 1) * pageSize
	if begin > numWith0 {
		exist = false
		return
	}
	exist = true
	// if the last page, return all of this page
	if end > numWith0 {
		rs.blogs = qs.blogs[begin:]
	} else {
		rs.blogs = qs.blogs[begin:end]
	}
	return
}
func (qs *QueryBlogs) GetBlogs() (bs []QueryBlog) {
	return qs.blogs
}
