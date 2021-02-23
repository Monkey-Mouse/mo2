package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddBlogs2Categories struct {
	BlogIDs     []primitive.ObjectID `json:"blog_ids,omitempty"`
	CategoryIDs []primitive.ObjectID `json:"category_ids,omitempty"`
}
