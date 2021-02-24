package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"mo2/dto"
	"mo2/server/model"
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
	account.Infos["avatar"] = "" // default pic
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

//create an anonymous account
func CreateAnonymousAccount() (a model.Account) {
	a = model.Account{
		ID:         primitive.NewObjectID(),
		UserName:   "visitor",
		Email:      string(rand.Int()) + "@mo2.com",
		HashedPwd:  "#2a$10$rXMPcOyfgdU6y5n3pkYQAukc3avJE9CLsx1v0Kn99GKV1NpREvN2i",
		EntityInfo: model.InitEntity(),
		Roles:      nil,
		Infos:      nil,
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

// FindAccountInfo find
func FindAccountInfo(id primitive.ObjectID) (u dto.UserInfo, exist bool) {
	a, exist := FindAccount(id)
	if exist {
		u = dto.Account2UserInfo(a)
	}
	return
}

// FindAccounts find from a list of id
func FindAccounts(ids []primitive.ObjectID) (bs []dto.UserInfoBrief) {
	for _, id := range ids {
		a, exist := FindAccount(id)
		if exist {
			bs = append(bs, dto.MapAccount2InfoBrief(a))
		}
	}
	return
}
