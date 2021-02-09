package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
	//"github.com/swaggo/swag/example/celler/model"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"mo2/database"
	"mo2/examples"
	"mo2/server/model"
	"net/http"
	"strconv"
)

// @Summary 新增用户
// @Description 为新用户创建信息，加入数据库
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
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	if err := addAccount.Validation(); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	/*account := model.Account{
		Name: addAccount.Name,
	}
	lastID, err := account.Insert()*/
	account, err := database.AddAccount(addAccount)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	//account.ID = lastID
	ctx.JSON(http.StatusOK, account)
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
		httputil.NewError(ctx, http.StatusNotFound, err)
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
