package controller

import (
	"mo2/database"
	"mo2/server/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
)

// PublishBlog godoc
// @Summary Publish Blog
// @Description add by json
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param account body model.Blog true "Add blog"
// @Success 200 {object} model.Blog
// @Router /api/blogs [post]
func (c *Controller) PublishBlog(ctx *gin.Context) {
	b := model.Blog{}
	if err := ctx.ShouldBindJSON(&b); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	// to be implemented: set blog's author due to cookie infomation

	database.AddBlog(&b)
	ctx.JSON(http.StatusOK, &b)
}
