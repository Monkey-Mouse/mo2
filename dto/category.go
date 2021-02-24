package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mo2/server/model"
)

type AddBlogs2Categories struct {
	BlogIDs     []primitive.ObjectID `json:"blog_ids,omitempty"`
	CategoryIDs []primitive.ObjectID `json:"category_ids,omitempty"`
}

type AddCategory2User struct {
	UserID     primitive.ObjectID `json:"user_id,omitempty"`
	CategoryID string             `json:"category_id,omitempty" example:"xxxxxxx"`
}

type AddCategory2Category struct {
	ParentCategoryID primitive.ObjectID `json:"parent_id,omitempty" example:"xxxxxxx"`
	Category         model.Category     `json:"category_id,omitempty"`
}
