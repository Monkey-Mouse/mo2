package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/Monkey-Mouse/mo2/mo2utils"
	"github.com/Monkey-Mouse/mo2/server/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var githubID = os.Getenv("GITHUB_ID")
var githubSecret = os.Getenv("GITHUB_SECRET")

type GithubUser struct {
	Login      string `json:"login,omitempty"`
	ID         int    `json:"id,omitempty"`
	AvatarUrl  string `json:"avatar_url,omitempty"`
	Location   string `json:"location,omitempty"`
	Bio        string `json:"bio,omitempty"`
	GithubHome string `json:"html_url,omitempty"`
}

func (c *Controller) GithubOauth(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.Redirect(307, "/oautherr")
		return
	}
	s := fmt.Sprintf(`
	{
		"client_id":"%s",
		"client_secret":"%s",
		"code":"%s"
	}
	`, githubID, githubSecret, code)
	req1, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(s))
	req1.Header.Add("Accept", "application/json")
	req1.Header.Add("Content-Type", "application/json")
	re, err := http.DefaultClient.Do(req1)
	if err != nil {
		ctx.Redirect(307, "/oautherr")
		fmt.Println(err)
		return
	}
	defer re.Body.Close()
	data, err := ioutil.ReadAll(re.Body)
	if err != nil {
		ctx.Redirect(307, "/oautherr")
		fmt.Println(err, s)
		return
	}
	token := struct {
		AccessToken string `json:"access_token"`
	}{}
	err = json.Unmarshal(data, &token)
	if err != nil {
		ctx.Redirect(307, "/oautherr")
		fmt.Println(err, string(data), "67")
		return
	}
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", token.AccessToken))
	req.Header.Add("Accept", "application/json")
	re1, err := http.DefaultClient.Do(req)
	if err != nil {
		ctx.Redirect(307, "/oautherr")
		fmt.Println(err, "76")
		return
	}
	defer re1.Body.Close()
	udata, err := ioutil.ReadAll(re1.Body)
	if err != nil {
		ctx.Redirect(307, "/oautherr")
		fmt.Println(err, "83")
		return
	}
	guser := GithubUser{}
	fmt.Println(string(udata))
	err = json.Unmarshal(udata, &guser)
	if err != nil || guser.ID == 0 {
		ctx.Redirect(307, "/oautherr")
		fmt.Println(err, string(udata), "90")
		return
	}
	fmt.Println(guser)
	account := model.Account{
		Roles:      []string{model.OrdinaryUser},
		UserName:   guser.Login,
		Email:      "@" + primitive.NewObjectID().Hex(),
		HashedPwd:  "some pass",
		EntityInfo: model.InitEntity(),
		Infos:      map[string]string{model.IsActive: model.True},
		Settings: map[string]string{
			model.Avatar: guser.AvatarUrl,
			"bio":        guser.Bio,
			"location":   guser.Location,
			"github":     guser.GithubHome,
			"github_id":  fmt.Sprint(guser.ID),
		},
	}
	database.UpsertAccountWithF(&account, bson.M{"settings.github_id": fmt.Sprint(guser.ID)})
	if account.ID == primitive.NilObjectID {
		account, _ = database.FindAccountByName(account.UserName)
	}
	var su = dto.Account2SuccessLogin(account)
	jwtToken := mo2utils.GenerateJwtCode(su)
	//login success: to record the state
	ctx.SetCookie("jwtToken", jwtToken, cookieExpiredTime, "/", ctx.Request.Host, false, true)
	ctx.Redirect(307, "/account")
}
