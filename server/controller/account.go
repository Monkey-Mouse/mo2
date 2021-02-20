package controller

import (
	"context"
	"encoding/json"
	"fmt"
	dto "mo2/dto"
	"mo2/mo2utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	//"github.com/swaggo/swag/example/celler/model"
	"log"
	"mo2/database"
	demo "mo2/examples"
	"mo2/server/model"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

const cookieExpiredTime int = 300000

// @Summary simple test
// @Description say something
// @Produce  json
// @Success 200 {string} json
// @Router /sayHello [get]
func SayHello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello! Welcome to Mo2!",
	})

}

// @Summary 新增用户
// @Description 为新用户创建信息，加入数据库
// @Produce  json
// @Success 200 {string} json
// @Router /api/accounts/addUser [post]
func (c *Controller) AddMo2User(ctx *gin.Context) {
	message := ctx.PostForm("message")
	nick := ctx.DefaultPostForm("nick", "anonymous")
	ctx.JSON(http.StatusOK, gin.H{
		"action":  "posted",
		"message": message,
		"nick":    nick,
	})

}

// Log godoc
// @Summary get user info
// @Description get by check cookies
// @Tags logs
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.SuccessLogin
// @Router /api/logs [get]
func (c *Controller) Log(ctx *gin.Context) {

	jwtToken, err := ctx.Cookie("jwtToken")
	var s dto.SuccessLogin
	if err != nil {
		//allocate an anonymous account
		account := database.CreateAnonymousAccount()
		s = dto.Account2SuccessLogin(account)
		jwtToken = mo2utils.GenerateJwtCode(account.UserName, s)
		//login success: to record the state
		ctx.SetCookie("jwtToken", jwtToken, cookieExpiredTime, "/", ctx.Request.Host, false, true)
	} else {
		//parse jwtToken and get user info
		userInfo, err := mo2utils.ParseJwt(jwtToken)
		if err != nil {
			log.Println(err)
		}
		err = json.Unmarshal(userInfo.Infos, &s)
		if err != nil {
			log.Fatal(err)
		}
	}
	ctx.JSON(http.StatusOK, s)
}

// AddAccount godoc
// @Summary Add an account
// @Description add by json account
// @Tags accounts
// @Accept  json
// @Produce  json
//// @Param account body model.AddAccount true "Add account"
// @Success 200 {object} model.Account
// @Router /api/accounts [post]
func (c *Controller) AddAccount(ctx *gin.Context) {
	var addAccount model.AddAccount
	if err := ctx.ShouldBindJSON(&addAccount); err != nil {
		ctx.JSON(http.StatusUnauthorized, SetResponseError(err))
		return
	}
	if err := addAccount.Validation(); err != nil {
		ctx.JSON(http.StatusUnauthorized, SetResponseError(err))
		return
	}
	nano := time.Now().Nanosecond()
	account, err := database.AddAccount(addAccount)

	if err != nil {
		if strings.Contains(err.Error(), "username") {
			fmt.Println(time.Now().Nanosecond() - nano)
			ctx.JSON(http.StatusUnauthorized, SetResponseReason("用户名已被使用"))
		} else {
			ctx.JSON(http.StatusUnauthorized, SetResponseReason("email已被使用"))
		}

		return
	}
	ctx.JSON(http.StatusOK, account)
}

// LoginAccount godoc
// @Summary login an account
// @Description login by json model.LoginAccount and set cookies
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param account body model.LoginAccount true "login account"
// @Success 200 {object} dto.SuccessLogin
// @Router /api/accounts/login [post]
func (c *Controller) LoginAccount(ctx *gin.Context) {
	var loginAccount model.LoginAccount
	if err := ctx.ShouldBindJSON(&loginAccount); err != nil {
		ctx.JSON(http.StatusBadRequest, SetResponseError(err))
		return
	}
	if err := loginAccount.Validation(); err != nil {
		ctx.JSON(http.StatusBadRequest, SetResponseError(err))
		return
	}
	account, err := database.VerifyAccount(loginAccount)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, SetResponseReason("用户名或密码错误"))
		return
	}
	var s = dto.Account2SuccessLogin(account)
	jwtToken := mo2utils.GenerateJwtCode(account.UserName, s)
	//login success: to record the state
	ctx.SetCookie("jwtToken", jwtToken, cookieExpiredTime, "/", ctx.Request.Host, false, true)
	ctx.JSON(http.StatusOK, s)
}

// LogoutAccount godoc
// @Summary logout
// @Description logout and delete cookies
// @Tags accounts
// @Produce  json
// @Success 200
// @Router /api/accounts/logout [get]
func (c *Controller) LogoutAccount(ctx *gin.Context) {

	ctx.SetCookie("jwtToken", "true", -1, "/", ctx.Request.Host, false, true)
	ctx.JSON(http.StatusOK, gin.H{"message": "logout success"})
}

// ShowAccount godoc
// @Summary Show a account
// @Description get string by ID
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path int true "Account ID"
// @Success 200 {string} json
// @Header 200 {string} Token "qwerty"
// @Router /api/accounts/{id} [get]
func (c *Controller) ShowAccount(ctx *gin.Context) {

	demo.Find()
	id := ctx.Param("id")

	aid, err := strconv.Atoi(id)

	fmt.Println(aid)
	if err != nil {
		log.Fatal(err)
		return
	}
	cli := demo.GetClient()
	col := cli.Database("test").Collection("trainers")

	filter := bson.D{{"name", "Ash"}}

	result := col.FindOne(context.TODO(), filter)
	//account, err := model.AccountOne(aid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, SetResponseError(err))
		return
	}
	fmt.Println(result)
	fmt.Println(result.Decode(demo.Trainer{}))

	ctx.JSON(http.StatusOK, result.Decode(demo.Trainer{}))

}

/*// ListAccounts godoc
// @Summary List accounts
// @Description get accounts
// @Accept  json
// @Produce  json
// @Param q query string false "name search by q"
// @Success 200 {string} json
// @Header 200 {string} Token "qwerty"
// @Router /accounts [get]
func (c *Controller) ListAccounts(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "list",
	})
}*/
