package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"mo2/dto"
	"mo2/server/model"

	"mo2/mo2utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var accCol = GetCollection("accounts")

//already check the validation in controller
//if add a newAccount success, return account info
func AddAccount(newAccount model.AddAccount) (account model.Account, err error) {
	collection := GetCollection("accounts")
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
	}
	_, err = collection.Indexes().CreateMany(context.TODO(), indexModel)
	if err != nil {
		log.Fatal(err)
	}
	//var account model.Account
	account.Email = newAccount.Email
	account.UserName = newAccount.UserName
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newAccount.Password), bcrypt.DefaultCost)
	account.EntityInfo = model.InitEntity()
	account.Roles = append(account.Roles, model.OrdinaryUser) // default role: OrdinaryUser
	account.Infos = make(map[string]string)
	account.Infos["avatar"] = ""        // default pic
	account.Infos["isActive"] = "false" // default pic
	if err != nil {
		log.Fatal(err)
		return
	}
	account.HashedPwd = string(hashedPwd)
	if err != nil {
		log.Fatal(err)
		return
	}
	insertResult, err := collection.InsertOne(context.TODO(), account)
	account.ID = insertResult.InsertedID.(primitive.ObjectID)
	return
}

// UpsertAccount
func UpsertAccount(a *model.Account) (success bool) {
	a.EntityInfo.Update()
	result, err := accCol.UpdateOne(context.TODO(), bson.M{"_id": a.ID}, bson.M{
		"$set": bson.M{
			"username":    a.UserName,
			"email":       a.Email,
			"hashedpwd":   a.HashedPwd,
			"entity_info": a.EntityInfo,
			"roles":       a.Roles,
			"infos":       a.Infos,
		},
	}, options.Update().SetUpsert(true))
	success = true
	if err != nil {
		log.Println(err)
		success = false
	}
	if result.MatchedCount == 0 {
		log.Println("blog id do not match in database")
		success = false
	}
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
		Roles:      nil,
		Infos:      nil,
	}
	return
}

// GenerateEmailToken generate token for email of an account
func GenerateEmailToken(addAccount model.AddAccount) (account *model.Account, err error) {
	//use email to verify
	collection := GetCollection("accounts")
	// verify email
	if err = collection.FindOne(context.TODO(), bson.D{{"email", addAccount.Email}}).Decode(&account); err != nil {
		return
	}
	if account.Infos == nil {
		account.Infos = make(map[string]string)
	}
	account.Infos["token"] = mo2utils.GenerateJwtToken(account.Email)
	account.Infos["isActive"] = "false"
	return
}

//verify email of an account
func VerifyEmail(info model.VerifyEmail) (account model.Account, err error) {
	//use email to verify
	collection := GetCollection("accounts")
	email := info.Email

	// verify email
	if err = collection.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&account); err != nil {
		return
	}
	if account.Infos["token"] == info.Token {
		account.Infos["isActive"] = "true"
		delete(account.Infos, "token")
		UpsertAccount(&account)
	} else {
		err = errors.New("token不符")
	}
	return
}

//verify an account
func VerifyAccount(info model.LoginAccount) (account model.Account, err error) {

	//first use username, then use email to verify
	//var
	collection := GetCollection("accounts")
	userNameOrEmail := info.UserNameOrEmail
	err = collection.FindOne(context.TODO(), bson.D{{"username", userNameOrEmail}}).Decode(&account)
	if account.Infos["isActive"] == "false" {
		err = errors.New("邮箱暂未激活")
		return
	}
	if err != nil {
		//then verify email
		if err == mongo.ErrNoDocuments {
			log.Println(err)
			err = collection.FindOne(context.TODO(), bson.D{{"email", userNameOrEmail}}).Decode(&account)
			if err != nil {
				// no chance
				return
			}
		} else {
			log.Fatal(err)
			return
		}

	}

	password := info.Password
	hashedPassword := account.HashedPwd
	//judge hash with hashed password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println(err)
		return
	}
	return
}

// FindAccount find
func FindAccount(id primitive.ObjectID) (a model.Account, exist bool) {
	exist = false
	if err := accCol.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&a); err != nil {
		if err != mongo.ErrNoDocuments {
			log.Fatal(err)
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
		us = append(us, dto.Account2UserInfo(account))
	}
	return
}

// FindAllAccounts find
func FindAllAccounts() (as []model.Account) {
	results, err := accCol.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	if err = results.All(context.TODO(), &as); err != nil {
		log.Fatal(err)
	}
	return
}

// FindAccountInfo find
func FindAccountInfo(id primitive.ObjectID) (u dto.UserInfo, exist bool) {
	a, exist := FindAccount(id)
	if exist {
		u = dto.Account2UserInfo(a)
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
