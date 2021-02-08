package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"mo2/examples"
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
// @Router /addUser [get]
func AddMo2User(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "great",
	})

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

// @Router /accounts/{id} [get]
func (c *Controller) ShowAccount(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)

	fmt.Println(aid)
	if err != nil {
		log.Fatal(err)
		return
	}
	cli := demo.GetClient()
	col := cli.Database("test").Collection("trainers")
	result := col.FindOne(context.TODO(), bson.D{{"Name", "Ash"}})
	//account, err := model.AccountOne(aid)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}
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
