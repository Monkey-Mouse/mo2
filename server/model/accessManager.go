package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const accessManagerStr = "AccessManager"

type AccessManager struct {
	EntityInfo Entity                          `json:"entityInfo,omitempty" bson:"entity_info,omitempty"`
	RoleMap    map[string][]primitive.ObjectID `json:"role_map,omitempty" example:"'admin':xxxxx 'write':xxxxx" bson:"role_map,omitempty"`
}
