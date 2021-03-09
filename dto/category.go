package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
