package controller

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mo2/database"
	"mo2/mo2utils"
	"mo2/server/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpsertBlog godoc
// @Summary Publish Blog
// @Description add by json
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param draft query bool false "bool false" false
// @Param account body model.Blog true "Add blog"
// @Success 200 {object} model.Blog
// @Router /api/blogs/publish [post]
func (c *Controller) UpsertBlog(ctx *gin.Context) {
	isDraftStr := ctx.DefaultQuery("draft", "true")
	isDraft := true
	if isDraftStr == "false" {
		isDraft = false
	}
	var b model.Blog
	if err := ctx.ShouldBindJSON(&b); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("内容含非法字符，请检查"))
		return
	}
	if b.AuthorID == primitive.NilObjectID {
		userInfo, exist := mo2utils.GetUserInfo(ctx)
		if exist {
			b.AuthorID = userInfo.ID
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, SetResponseReason("权限不足，请先登录"))
			return
		}
	}
	database.UpsertBlog(&b, isDraft)
	ctx.JSON(http.StatusOK, b)
}

// FindBlogsByUser godoc
// @Summary find Blog
// @Description
// @Tags blogs
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.Blog
// @Router /api/blogs/find/byUser [get]
func (c *Controller) FindBlogsByUser(ctx *gin.Context) {
	// get user info due to cookie information
	info, ext := mo2utils.GetUserInfo(ctx)
	if !ext {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, SetResponseReason("权限不足，请先登录"))
		return
	}

	blogs := database.FindBlogsByUser(info)
	ctx.JSON(http.StatusOK, blogs)
}

// FindBlogsByUserId godoc
// @Summary find Blog by userid
// @Description
// @Tags blogs
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.Blog
// @Router /api/blogs/find/byUser [get]
func (c *Controller) FindBlogsByUserId(ctx *gin.Context) {
	// get user info due to cookie information
	info, ext := mo2utils.GetUserInfo(ctx)
	if !ext {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, SetResponseReason("权限不足，请先登录"))
		return
	}

	blogs := database.FindBlogsByUser(info)
	ctx.JSON(http.StatusOK, blogs)
}

// FindBlogById godoc
// @Summary find Blog by id
// @Description
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param account body primitive.ObjectID true "Save draft"
// @Success 200 {object} model.Blog
// @Router /api/blogs/find/id [post]
func (c *Controller) FindBlogById(ctx *gin.Context) {
	id := primitive.ObjectID{}
	if err := ctx.ShouldBindJSON(&id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("非法输入"))
		return
	}
	// todo check if the user has right
	// get user info due to cookie information

	blog := database.FindBlogById(id)
	ctx.JSON(http.StatusOK, blog)
}

// FindAllBlogs godoc
// @Summary find all Blogs
// @Description find
// @Tags blogs
// @Produce  json
// @Success 200 {object} []model.Blog
// @Router /api/blogs/find/all [get]
func (c *Controller) FindAllBlogs(ctx *gin.Context) {
	blogs := database.FindAllBlogs()
	ctx.JSON(http.StatusOK, blogs)
}
