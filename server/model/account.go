package model

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"fmt"
)

// Account example
type Account struct {
	ID        primitive.ObjectID `json:"id" example:"xxxxxxxxxxxxx==" bson:"_id,omitempty"`
	UserName  string             `json:"userName" example:"account name"`
	Email     string             `json:"email" example:"email@qq.com"`
	HashedPwd string             `json:"hashedPassword" example:"$2a$10$rXMPcOyfgdU6y5n3pkYQAukc3avJE9CLsx1v0Kn99GKV1NpREvN2i"`
}

// AddAccount example
type AddAccount struct {
	UserName string `json:"userName" example:"account name"`
	Email    string `json:"email" example:"email@qq.com"`
	Password string `json:"password" example:"p@ssword"`
}

// LoginAccount example
type LoginAccount struct {
	UserNameOrEmail string `json:"userNameOrEmail" example:"account name/email@qq.com"`
	Password        string `json:"password" example:"p@ssword"`
}

// Validation example
func (a AddAccount) Validation() error {
	switch {
	case len(a.UserName) == 0:
		return ErrNameInvalid
	case len(a.Email) == 0:
		return ErrEmailInvalid
	case len(a.Password) == 0:
		return ErrPasswordInvalid

	default:
		return nil
	}
}

//
//LoginAccount Validation
func (a LoginAccount) Validation() error {
	switch {
	case len(a.UserNameOrEmail) == 0:
		return ErrNameInvalid
	case len(a.Password) == 0:
		return ErrPasswordInvalid

	default:
		return nil
	}
}

//  example
var (
	ErrNameInvalid     = errors.New("name is empty")
	ErrEmailInvalid    = errors.New("email is empty")
	ErrPasswordInvalid = errors.New("password is empty")
)
