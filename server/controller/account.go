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

// AddAccountRole godoc
// @Summary Add role for an account
// @Description add by json account
// @Tags admin
// @Accept  json
// @Produce  json
// @Param account body model.AddAccountRole true "add new account info"
// @Success 200 {object} model.Account
// @Router /api/accounts/role [post]
func (c *Controller) AddAccountRole(ctx *gin.Context) {
	var addAccount model.AddAccountRole
	if err := ctx.ShouldBindJSON(&addAccount); err != nil {
		ctx.JSON(http.StatusUnauthorized, SetResponseError(err))
		return
	}
	if err := addAccount.Validation(); err != nil {
		ctx.JSON(http.StatusUnauthorized, SetResponseError(err))
		return
	}
	account, exist := database.FindAccount(addAccount.ID)
	if !exist {
		ctx.AbortWithStatusJSON(http.StatusNotFound, SetResponseReason("无此用户"))
		return
	}
	account.Roles = append(account.Roles, addAccount.Roles...)
	database.UpsertAccount(&account)
	ctx.JSON(http.StatusOK, account)
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
// @Summary Show account's info
// @Description get string by ID；若id为空，返回所有用户信息
// @Tags accounts
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path string false "Account ID"
// @Success 200 {object} []dto.UserInfo
// @Router /api/accounts/detail/{id} [get]
func (c *Controller) ShowAccount(ctx *gin.Context) {
	idStr := ctx.Param("id")
	var us []dto.UserInfo
	if idStr == "undefined" {
		us = database.FindAllAccountsInfo()
	} else {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("非法输入"))
			return
		}
		result, exist := database.FindAccountInfo(id)
		if !exist {
			ctx.AbortWithStatusJSON(http.StatusNotFound, SetResponseReason("无此用户"))
			return
		}
		us = append(us, result)
	}

	ctx.JSON(http.StatusOK, us)
}

// ListAccountsInfo godoc
// @Summary List accounts brief info
// @Description from a list of user ids [usage]:/api/accounts/listBrief?id=60223d4042d6febff9f276f0&id=60236866d2a68483adaccc38
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param userIDs path array true "user IDs list"
// @Success 200 {object} []dto.UserInfoBrief
// @Router /api/accounts/listBrief [get]
func (c *Controller) ListAccountsInfo(ctx *gin.Context) {
	userIDstrs, exist := ctx.GetQueryArray("id")
	var bs []dto.UserInfoBrief
	if !exist {
		bs = database.ListAllAccountsBrief()
	} else {
		bs = database.ListAccountsBrief(userIDstrs)
	}
	ctx.JSON(http.StatusOK, bs)
}
