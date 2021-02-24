package dto

import (
	"mo2/server/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Contains test slice contain string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}
func (user LoginUserInfo) IsUserInRole(role model.Erole) bool {
	return Contains(user.Roles, role)
}

//todo implement error
type LoginUserInfo struct {
	ID    primitive.ObjectID `json:"id" example:"xxxxxxxxxxxxx==" bson:"_id,omitempty"`
	Name  string             `json:"name" example:"account name"`
	Email string             `json:"email" example:"email@qq.com"`
	Roles []model.Erole      `json:"roles" example:"ordinaryUser"  bson:"roles"`
}

func Account2SuccessLogin(a model.Account) (s LoginUserInfo) {
	s.ID = a.ID
	s.Name = a.UserName
	s.Roles = a.Roles
	s.Email = a.Email

	return s
}

type UserInfoBrief struct {
	ID   primitive.ObjectID `json:"id" example:"xxxxxxxxxxxxx==" bson:"_id,omitempty"`
	Name string             `json:"name" example:"account name"`
}

func MapAccount2InfoBrief(a model.Account) (b UserInfoBrief) {
	b.ID = a.ID
	b.Name = a.UserName
	return b
}
