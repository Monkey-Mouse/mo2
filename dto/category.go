package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mo2/server/model"
)

type AddBlogs2Categories struct {
	BlogIDs     []primitive.ObjectID `json:"blog_ids,omitempty"`
	CategoryIDs []primitive.ObjectID `json:"category_ids,omitempty"`
}

// RelateEntity2Entity 将单实体关联到单实体dto
type RelateEntity2Entity struct {
	RelatedID  primitive.ObjectID `json:"related_id,omitempty"`
	RelateToID primitive.ObjectID `json:"relateTo_id,omitempty"`
}

// RelateEntitySet2EntitySet 关联两个实体集dto
type RelateEntitySet2EntitySet struct {
	RelatedIDs  []primitive.ObjectID `json:"related_ids,omitempty"`
	RelateToIDs []primitive.ObjectID `json:"relateTo_ids,omitempty"`
}

// RelateEntity2EntitySet 关联单实体到多实体集dto
type RelateEntity2EntitySet struct {
	RelatedID   primitive.ObjectID   `json:"related_id,omitempty"`
	RelateToIDs []primitive.ObjectID `json:"relateTo_ids,omitempty"`
}

// RelateEntitySet2Entity 关联实体集到单实体dto
type RelateEntitySet2Entity struct {
	RelatedIDs []primitive.ObjectID `json:"related_ids,omitempty"`
	RelateToID primitive.ObjectID   `json:"relateTo_id,omitempty"`
}
type AddCategory2User struct {
	UserID     primitive.ObjectID `json:"user_id,omitempty"`
	CategoryID string             `json:"category_id,omitempty" example:"xxxxxxx"`
}

type AddCategory2Category struct {
	ParentCategoryID primitive.ObjectID `json:"parent_id,omitempty" example:"xxxxxxx"`
	Category         model.Category     `json:"category_id,omitempty"`
}
