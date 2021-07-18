package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"

	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/Monkey-Mouse/mo2/mo2utils"
	"github.com/Monkey-Mouse/mo2/mo2utils/mo2errors"
	"github.com/Monkey-Mouse/mo2/server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var accCol = GetCollection("accounts")

func CreateAccountIndex() (err error) {
	//ensure index
	indexModel := []mongo.IndexModel{
		{
			Keys:    bson.D{{"username", 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{"email", 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{"settings.github_id", 1}},
		},
	}
	_, err = GetCollection("accounts").Indexes().CreateMany(context.TODO(), indexModel)
	return
}

// FindAccountByEmail get account by email
func FindAccountByEmail(email string) (account model.Account, e mo2errors.Mo2Errors) {
	if err := accCol.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&account); err != nil {
		if err == mongo.ErrNoDocuments {
			e.Init(mo2errors.Mo2NotFound, err.Error())
		} else {
			e.Init(mo2errors.Mo2Error, err.Error())
		}
	}
	e.InitCode(mo2errors.Mo2NoError)
	return
}

// EnsureEmailUnique return true if unique, else return false
// 	if not unique, check whether the user exists
// 								if active, already exists, return false
//								if not active, delete document, return ture
// 	if unique, return true
func EnsureEmailUnique(email string) (unique bool, e mo2errors.Mo2Errors) {
	user, e := FindAccountByEmail(email)
	if e.ErrorCode != mo2errors.Mo2NotFound {
		unique = false
		if user.Infos[model.IsActive] == model.True {
			return
		} else {
			if _, e = DeleteAccountByEmail(user.Email); e.IsError() {
				if e.ErrorCode != mo2errors.Mo2NotFound {
					return
				}
			}
		}
	}
	unique = true
	e.InitCode(mo2errors.Mo2NoError)
	return
}

//already check the validation in controller
//if add a newAccount success, return account info
func InitAccount(newAccount model.AddAccount, token string) (account model.Account, err error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newAccount.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
		return
	}
	//var account model.Account
	account = model.Account{
		UserName:   newAccount.UserName,
		Email:      newAccount.Email,
		HashedPwd:  string(hashedPwd),
		EntityInfo: model.InitEntity(),
		Infos:      map[string]string{model.Token: token, model.IsActive: model.False},
		Settings:   map[string]string{model.Avatar: ""},
	}
	model.AddRoles(&account, model.OrdinaryUser)
	insertResult, err := accCol.InsertOne(context.TODO(), account)
	if err != nil {
		merr := err.(mongo.WriteException).WriteErrors[0]
		if merr.Code == 11000 {
			acc, _ := FindAccountByName(account.UserName)
			if acc.Infos[model.IsActive] == model.True {
				err = mo2errors.New(mo2errors.Mo2Conflict, "Name已被注册！")
				return
			}
			DeleteAccount(acc.ID)
			insertResult, err = accCol.InsertOne(context.TODO(), account)
		} else {
			return
		}
	} else {
		if insertResult.InsertedID != nil {
			id, ext := insertResult.InsertedID.(primitive.ObjectID)
			if ext {
				account.ID = id
			}
		}
	}
	return
}

// FindAccountByName as name
func FindAccountByName(name string) (a model.Account, exist bool) {
	exist = false
	if err := accCol.FindOne(context.TODO(), bson.D{{"username", name}}).Decode(&a); err != nil {
		if err != mongo.ErrNoDocuments {
			panic(err)
		}
	}
	if a.IsValid() {
		exist = true
	}
	return
}

// UpsertAccount
func UpsertAccount(a *model.Account) (merr mo2errors.Mo2Errors) {
	merr = UpsertAccountWithF(a, bson.M{"_id": a.ID})
	return
}
func UpsertAccountWithF(a *model.Account, filter interface{}) (merr mo2errors.Mo2Errors) {
	a.EntityInfo.Update()
	result, err := accCol.UpdateOne(context.TODO(), filter, bson.M{
		"$set": bson.M{
			"username":    a.UserName,
			"email":       a.Email,
			"hashedpwd":   a.HashedPwd,
			"entity_info": a.EntityInfo,
			"roles":       a.Roles,
			"infos":       a.Infos,
			"settings":    a.Settings,
		},
	}, options.Update().SetUpsert(true))
	if err != nil {
		merr.Init(mo2errors.Mo2Error, err.Error())
		return
	}
	if result.UpsertedID != nil {
		id, ext := result.UpsertedID.(primitive.ObjectID)
		if ext {
			a.ID = id
		}
	}
	mo2utils.IndexAccount(a)
	if result.ModifiedCount != 0 {
		merr.Init(mo2errors.Mo2NoError, "更新完成")
	} else {
		merr.Init(mo2errors.Mo2NoError, "没有任何更改")
	}
	return
}

func DeleteAccount(id primitive.ObjectID) (a model.Account, e mo2errors.Mo2Errors) {
	if err := accCol.FindOneAndDelete(context.TODO(), bson.D{{"_id", id}}).Decode(&a); err != nil {
		if err == mongo.ErrNoDocuments {
			e.Init(mo2errors.Mo2NotFound, err.Error())
		} else {
			e.Init(mo2errors.Mo2Error, err.Error())
		}
	}
	return
}

func DeleteAccountByEmail(email string) (a model.Account, e mo2errors.Mo2Errors) {
	if err := accCol.FindOneAndDelete(context.TODO(), bson.D{{"email", email}}).Decode(&a); err != nil {
		if err == mongo.ErrNoDocuments {
			e.Init(mo2errors.Mo2NotFound, err.Error())
		} else {
			e.Init(mo2errors.Mo2Error, err.Error())
		}
	}
	mo2utils.DeleteAccountIndex(a.ID.Hex())
	e.InitCode(mo2errors.Mo2NoError)
	return
}

//create an anonymous account
func CreateAnonymousAccount() (a model.Account) {
	a = model.Account{
		ID:         primitive.NewObjectID(),
		UserName:   "visitor",
		Email:      fmt.Sprint(rand.Int()) + "@mo2.com",
		HashedPwd:  "#2a$10$rXMPcOyfgdU6y5n3pkYQAukc3avJE9CLsx1v0Kn99GKV1NpREvN2i",
		EntityInfo: model.InitEntity(),
		Roles:      []string{model.Anonymous},
		Infos:      nil,
	}
	return
}

// GenerateEmailToken generate token for email of an account
func GenerateEmailToken(addAccount model.AddAccount) (account *model.Account, err error) {
	//use email to verify
	if err = accCol.FindOne(context.TODO(), bson.D{{"email", addAccount.Email}}).Decode(&account); err != nil {
		return
	}
	if account.Infos == nil {
		account.Infos = make(map[string]string)
	}
	account.Infos[model.Token] = mo2utils.GenerateVerifyJwtToken(account.Email)
	account.Infos[model.IsActive] = model.False
	return
}

//verify email of an account
func VerifyEmail(info model.VerifyEmail) (account model.Account, err error) {
	//use email to verify
	email := info.Email
	// verify email
	if err = accCol.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&account); err != nil {
		return
	}
	if account.Infos[model.Token] == info.Token {
		account.Infos[model.IsActive] = model.True
		delete(account.Infos, model.Token)
		UpsertAccount(&account)
	} else {
		err = errors.New("token不符")
	}
	return
}

//verify an account
func VerifyAccount(info model.LoginAccount) (account model.Account, err error) {

	//first use username, then use email to verify
	collection := GetCollection("accounts")
	userNameOrEmail := info.UserNameOrEmail
	err = collection.FindOne(context.TODO(), bson.D{{"username", userNameOrEmail}}).Decode(&account)
	if err != nil {
		//then verify email
		if err == mongo.ErrNoDocuments {
			log.Println(err)
			err = collection.FindOne(context.TODO(), bson.D{{"email", userNameOrEmail}}).Decode(&account)
			if err != nil {
				err = errors.New("用户名或email不存在")
				// no chance
				return
			}
		} else {
			err = errors.New("未知错误：数据异常")
			// panic(err)
			return
		}

	}
	if account.Infos[model.IsActive] == model.False {
		err = errors.New("邮箱暂未激活！")
		return
	}
	password := info.Password
	hashedPassword := account.HashedPwd
	//judge hash with hashed password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		err = errors.New("密码错误！")
		return
	}
	return
}

// FindAccount find
func FindAccount(id primitive.ObjectID) (a model.Account, exist bool) {
	exist = false
	if err := accCol.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&a); err != nil {
		if err != mongo.ErrNoDocuments {
			panic(err)
		}
	}
	if a.IsValid() {
		exist = true
	}
	return
}

// FindAllAccountsInfo find
func FindAllAccountsInfo() (us []dto.UserInfo) {
	as := FindAllAccounts()
	for _, account := range as {
		us = append(us, dto.Account2UserPublicInfo(account))
	}
	return
}

// FindAllAccounts find
func FindAllAccounts() (as []model.Account) {
	results, err := accCol.Find(context.TODO(), bson.D{{}})
	if err != nil {
		panic(err)
	}
	if err = results.All(context.TODO(), &as); err != nil {
		panic(err)
	}
	return
}

// FindAccountInfo find
func FindAccountInfo(id primitive.ObjectID) (u dto.UserInfo, exist bool) {
	a, exist := FindAccount(id)
	if exist {
		u = dto.Account2UserPublicInfo(a)
	}
	return
}

// ListAccountsBrief find from a list of id
func ListAccountsBrief(idStrs []string) (bs []dto.UserInfoBrief) {
	ids := make([]primitive.ObjectID, len(idStrs))
	i := 0
	for _, idStr := range idStrs {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return
		}
		ids[i] = id
		i++
	}
	cursor, err := accCol.Find(context.TODO(),
		bson.D{
			{"_id",
				bson.D{
					{"$in", ids},
				}}})
	if err != nil {
		fmt.Println(err)
	}
	cursor.All(context.TODO(), &bs)
	return
}

// ListAllAccountsBrief find from a list of id
func ListAllAccountsBrief() (bs []dto.UserInfoBrief) {
	as := FindAllAccounts()
	for _, account := range as {
		bs = append(bs, dto.MapAccount2InfoBrief(account))
	}
	return
}
