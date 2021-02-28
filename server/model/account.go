package model

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	//"fmt"
)

type Erole = string

const (
	GeneralAdmin Erole = "GeneralAdmin"
	OrdinaryUser Erole = "OrdinaryUser"
)
const (
	Token    = "token"
	Avatar   = "avatar"
	IsActive = "isActive"
	True     = "true"
	False    = "false"
)

// Account example
type Account struct {
	ID         primitive.ObjectID `json:"id,omitempty" example:"xxxxxxxxxxxxx==" bson:"_id,omitempty"`
	UserName   string             `json:"userName" example:"account name" bson:"username,omitempty"`
	Email      string             `json:"email" example:"email@qq.com" bson:"email,omitempty"`
	HashedPwd  string             `json:"hashedPassword" example:"$2a$10$rXMPcOyfgdU6y5n3pkYQAukc3avJE9CLsx1v0Kn99GKV1NpREvN2i" bson:"hashedpwd,omitempty"`
	EntityInfo Entity             `json:"entityInfo,omitempty" bson:"entity_info,omitempty" bson:"entity_info,omitempty"`
	Roles      []Erole            `json:"roles" bson:"roles"`
	Infos      map[string]string  `json:"infos" example:"'token': 'xxxxxxxx'(private data)" bson:"infos,omitempty"`
	Settings   map[string]string  `json:"settings" example:"'avatar': 'www.avatar.com/account_name','site':'www.limfx.com'(public data)" bson:"settings,omitempty"`
}

// AddAccount example
type AddAccount struct {
	UserName string `json:"userName" example:"account name"`
	Email    string `json:"email" example:"email@mo2.com"`
	Password string `json:"password" example:"p@ssword"`
}

// DeleteAccount example
type DeleteAccount struct {
	Email    string `json:"email" example:"email@mo2.com"`
	Password string `json:"password" example:"p@ssword"`
}

// AddAccountRole example
type AddAccountRole struct {
	ID       primitive.ObjectID `json:"id" example:"xxxxxxxxxxxxx==" `
	Roles    []Erole            `json:"roles"`
	SuperKey string             `json:"super_key" example:"special"`
}

// LoginAccount example
type LoginAccount struct {
	UserNameOrEmail string `json:"userNameOrEmail" example:"account name/email@mo2.com"`
	Password        string `json:"password" example:"p@ssword"`
}

// VerifyEmail
type VerifyEmail struct {
	Email string `json:"Email" example:"email@mo2.com"`
	Token string `json:"token" example:"p@ssword"`
}

// Validation example
func (a AddAccountRole) Validation() error {
	switch {
	case os.Getenv("MO2_SUPER_KEY") != a.SuperKey:
		return ErrPasswordInvalid
	case a.ID.IsZero():
		return ErrNameInvalid
	default:
		return nil
	}
}
func AddRoles(a *Account, roles ...Erole) {
	var new bool
	for _, role := range roles {
		new = true
		for _, existRole := range a.Roles {
			if role == existRole {
				new = false
				break
			}
		}
		if new {
			a.Roles = append(a.Roles, role)
		}
	}

}
func (a Account) IsValid() (valid bool) {
	valid = true
	if a.ID.IsZero() {
		valid = false
	}
	return
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
	ErrPasswordInvalid = errors.New("password is invalid")
)
