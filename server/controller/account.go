package controller

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	dto "mo2/dto"
	"mo2/mo2utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	//"github.com/swaggo/swag/example/celler/model"
	"log"
	"mo2/database"
	"mo2/server/model"
	"net/http"
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
// @Success 200 {object} dto.LoginUserInfo
// @Router /api/logs [get]
func (c *Controller) Log(ctx *gin.Context) {

	jwtToken, err := ctx.Cookie("jwtToken")
	var s dto.LoginUserInfo
	if err != nil {
		//allocate an anonymous account
		account := database.CreateAnonymousAccount()
		s = dto.Account2SuccessLogin(account)
		jwtToken = mo2utils.GenerateJwtCode(s)
		//login success: to record the state
		ctx.SetCookie("jwtToken", jwtToken, cookieExpiredTime, "/", ctx.Request.Host, false, true)
	} else {
		//parse jwtToken and get user info
		s, err = mo2utils.ParseJwt(jwtToken)
		if err != nil {
			log.Println(err)
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
// @Param account body model.AddAccount true "add new account info"
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
// @Success 200 {object} dto.LoginUserInfo
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
	jwtToken := mo2utils.GenerateJwtCode(s)
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
// @Summary Show a account's info
// @Description get string by ID
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path string true "Account ID"
// @Success 200 {object} model.Account
// @Router /api/accounts/detail/{id} [get]
func (c *Controller) ShowAccount(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("非法输入"))
		return
	}
	result := database.FindAccount(id)
	if !result.IsValid() {
		ctx.AbortWithStatusJSON(http.StatusNotFound, SetResponseReason("无此用户"))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// ListAccountsInfo godoc
// @Summary List accounts brief info
// @Description from a list of user ids
// @Accept  json
// @Produce  json
// @Param userIDs body []primitive.ObjectID false "user IDs list"
// @Success 200 {object} []dto.UserInfoBrief
// @Router /api/accounts/listBrief [post]
func (c *Controller) ListAccountsInfo(ctx *gin.Context) {
	var userIDs []primitive.ObjectID
	if err := ctx.ShouldBindJSON(&userIDs); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("非法输入"))
		return
	}
	bs := database.FindAccounts(userIDs)
	ctx.JSON(http.StatusOK, bs)
}
