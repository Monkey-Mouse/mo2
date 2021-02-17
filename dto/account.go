package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mo2/server/model"
)

//todo implement error
type SuccessLogin struct {
	ID    primitive.ObjectID `json:"id" example:"xxxxxxxxxxxxx==" bson:"_id,omitempty"`
	Name  string             `json:"name" example:"account name"`
	Email string             `json:"email" example:"email@qq.com"`
	Roles []model.Erole      `json:"roles" example:"ordinaryUser"  bson:"roles"`
}

func Account2SuccessLogin(a model.Account) (s SuccessLogin) {
	s.ID = a.ID
	s.Name = a.UserName
	s.Roles = a.Roles
	s.Email = a.Email

	return s
}
