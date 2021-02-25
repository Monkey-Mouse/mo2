package controller

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mo2/database"
	"mo2/dto"
	"mo2/mo2utils"
	"mo2/server/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UpsertBlog godoc
// @Summary Publish Blog
// @Description add by json
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param draft query bool false "bool true" true
// @Param account body model.Blog true "Add blog"
// @Success 201 {object} model.Blog
// @Failure 204 {object} model.Blog
// @Failure 400 {object} ResponseError
// @Failure 401 {object} ResponseError
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
	userInfo, exist := mo2utils.GetUserInfo(ctx)
	if b.AuthorID == primitive.NilObjectID {
		if exist {
			b.AuthorID = userInfo.ID
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, SetResponseReason("权限不足，请先登录"))
			return
		}
	} else {
		if b.AuthorID != userInfo.ID {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, SetResponseReason("没有权限修改文章"))
			return
		}
	}
	if success := database.UpsertBlog(&b, isDraft); success {
		ctx.Header("location", ctx.FullPath())
		ctx.JSON(http.StatusCreated, b)
	} else {
		ctx.JSON(http.StatusNoContent, b)
	}

}

// DeleteBlog godoc
// @Summary delete Blog
// @Description delete by path
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param draft path bool true "bool true" true
// @Param id path string false "string xxxxxxxx" "xxxxxxx"
// @Success 200 {object} model.Blog
// @Router /api/blogs/{id} [delete]
func (c *Controller) DeleteBlog(ctx *gin.Context) {
	isDraftStr := ctx.Param("draft")
	isDraft := true
	if isDraftStr == "false" {
		isDraft = false
	}
	idStr := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("非法输入"))
		return
	}
	blog := database.FindBlogById(id, isDraft)
	if blog.ID.IsZero() {
		ctx.AbortWithStatusJSON(http.StatusNotFound, SetResponseReason("页面找不到了"))
		return
	}
	if blog.AuthorID == primitive.NilObjectID {
		userInfo, exist := mo2utils.GetUserInfo(ctx)
		if exist {
			blog.AuthorID = userInfo.ID
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, SetResponseReason("权限不足"))
			return
		}
	}
	blog.EntityInfo.IsDeleted = true
	database.UpsertBlog(&blog, isDraft)
	ctx.JSON(http.StatusOK, blog)
}

// RestoreBlog godoc
// @Summary restore Blog
// @Description restore by path
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param draft path bool true "bool true" true
// @Param id path string false "string xxxxxxxx" "xxxxxxx"
// @Success 200 {object} model.Blog
// @Router /api/blogs/{id} [put]
func (c *Controller) RestoreBlog(ctx *gin.Context) {
	isDraftStr := ctx.Param("draft")
	isDraft := true
	if isDraftStr == "false" {
		isDraft = false
	}
	idStr := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("非法输入"))
		return
	}
	blog := database.FindBlogById(id, isDraft)
	if blog.ID.IsZero() {
		ctx.AbortWithStatusJSON(http.StatusNotFound, SetResponseReason("页面找不到了"))
		return
	}
	if blog.AuthorID == primitive.NilObjectID {
		userInfo, exist := mo2utils.GetUserInfo(ctx)
		if exist {
			blog.AuthorID = userInfo.ID
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, SetResponseReason("权限不足"))
			return
		}
	}
	blog.EntityInfo.IsDeleted = false
	database.UpsertBlog(&blog, isDraft)
	ctx.JSON(http.StatusOK, blog)
}

// FindBlogsByUser godoc
// @Summary find Blog
// @Description
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param draft query bool false "bool true" true
// @Param page query int false "int 0" 0
// @Param pageSize query int false "int 5" 5
// @Success 200 {object} []dto.QueryBlogs
// @Router /api/blogs/find/own [get]
func (c *Controller) FindBlogsByUser(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "0")
	pageSizeStr := ctx.DefaultQuery("pageSize", "5")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("query with page"))
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("query with pageSize"))
	}
	isDraftStr := ctx.DefaultQuery("draft", "true")
	isDraft := parseString2Bool(isDraftStr)
	// get user info due to cookie information
	info, ext := mo2utils.GetUserInfo(ctx)
	if !ext {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, SetResponseReason("权限不足，请先登录"))
		return
	}

	blogs := database.FindBlogsByUser(info, isDraft)
	qBlogs := dto.QueryBlogs{}
	qBlogs.Init(blogs)
	results, _ := qBlogs.Query(page, pageSize)
	ctx.JSON(http.StatusOK, results.GetBlogs())
}

// FindBlogsByUserId godoc
// @Summary find Blog by userid
// @Description
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param draft query bool false "bool true" true
// @Param id query string false "string xxxxxxxx" "xxxxxxx"
// @Param page query int false "int 0" 0
// @Param pageSize query int false "int 5" 5
// @Success 200 {object} []dto.QueryBlogs
// @Router /api/blogs/find/userId [get]
func (c *Controller) FindBlogsByUserId(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "0")
	pageSizeStr := ctx.DefaultQuery("pageSize", "5")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("query with page"))
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("query with pageSize"))
	}
	isDraftStr := ctx.DefaultQuery("draft", "true")
	isDraft := parseString2Bool(isDraftStr)
	idStr := ctx.Query("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("非法输入"))
		return
	}
	// todo check if the user has right
	// get user info due to cookie information
	blogs := database.FindBlogsByUserId(id, isDraft)
	qBlogs := dto.QueryBlogs{}
	qBlogs.Init(blogs)
	results, _ := qBlogs.Query(page, pageSize)
	ctx.JSON(http.StatusOK, results.GetBlogs())
}

// FindBlogById godoc
// @Summary find Blog by id
// @Description
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param draft query bool false "bool true" true
// @Param id query string false "string xxxxxxxx" "xxxxxxx"
// @Success 200 {object} model.Blog
// @Router /api/blogs/find/id [get]
func (c *Controller) FindBlogById(ctx *gin.Context) {
	isDraftStr := ctx.DefaultQuery("draft", "true")
	isDraft := parseString2Bool(isDraftStr)
	idStr := ctx.Query("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("非法输入"))
		return
	}
	// todo check if the user has right
	// get user info due to cookie information
	blog := database.FindBlogById(id, isDraft)
	if blog.ID.IsZero() {
		ctx.AbortWithStatusJSON(http.StatusNotFound, SetResponseReason("页面找不到了"))
		return
	}
	ctx.JSON(http.StatusOK, blog)
}

// QueryBlogs godoc
// @Summary find all Blogs
// @Description find
// @Tags blogs
// @Produce  json
// @Param draft query bool false "bool true" true
// @Param page query int false "int 0" 0
// @Param pageSize query int false "int 5" 5
// @Success 200 {object} []dto.QueryBlog
// @Router /api/blogs/query [get]
func (c *Controller) QueryBlogs(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "0")
	pageSizeStr := ctx.DefaultQuery("pageSize", "5")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("query with page"))
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("query with pageSize"))
	}

	isDraftStr := ctx.DefaultQuery("draft", "true")
	isDraft := parseString2Bool(isDraftStr)

	blogs := database.FindAllBlogs(isDraft)
	qBlogs := dto.QueryBlogs{}
	qBlogs.Init(blogs)
	results, _ := qBlogs.Query(page, pageSize)
	ctx.JSON(http.StatusOK, results.GetBlogs())
}

func parseString2Bool(s string) (b bool) {

	if s == "false" {
		b = false
	} else {
		b = true
	}
	return
}
