package model

import (
	"errors"
	//"fmt"

	uuid "github.com/gofrs/uuid"
)

// Account example
type Account struct {
	ID        int       `json:"id" example:"1" format:"int64"`
	UserName  string    `json:"userName" example:"account name"`
	Email     string    `json:"email" example:"email@qq.com"`
	HashedPwd string    `json:"hashedPassword" example:"sdfsfx"`
	UUID      uuid.UUID `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
}

// AddAccount example
type AddAccount struct {
	UserName string `json:"userName" example:"account name"`
	Email    string `json:"email" example:"email@qq.com"`
	Password string `json:"password" example:"p@ssword"`
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

//  example
var (
	ErrNameInvalid     = errors.New("name is empty")
	ErrEmailInvalid    = errors.New("email is empty")
	ErrPasswordInvalid = errors.New("password is empty")
)
