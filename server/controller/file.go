package controller

import (
	"fmt"
	"io/ioutil"
	_ "mime/multipart" //in godoc comment
	"net/http"
	"time"

	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/mo2utils"
	"github.com/Monkey-Mouse/mo2/server/controller/badresponse"
	"github.com/Monkey-Mouse/mo2/services/importService"
	"github.com/gin-gonic/gin"
)

// Upload godoc
// @Summary simple test
// @Description say something
// @Accept multipart/form-data
// @Produce  json
// @Param form body string true "file"
// @Success 200 {object} model.Blog
// @Failure 400 {object} badresponse.ResponseError
// @Failure 401 {object} badresponse.ResponseError
// @Failure 404 {object} badresponse.ResponseError
// @Router /api/file [post]
func (c *Controller) Upload(ctx *gin.Context) {
	userInfo, exist := mo2utils.GetUserInfo(ctx)
	if !exist {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseReason("请先登录"))
	}
	// Multipart form
	file, _ := ctx.FormFile("upload[]")

	t1 := time.Now()
	src, err := file.Open()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseError(err))
	}
	defer src.Close()
	data, err := ioutil.ReadAll(src)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseError(err))
	}
	blog := importService.Transform(data)
	fmt.Println(time.Since(t1))
	blog.AuthorID = userInfo.ID
	isDraft := true
	database.UpsertBlog(&blog, isDraft)
	ctx.JSON(http.StatusOK, blog)
}
