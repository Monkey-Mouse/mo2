package controller

import (
	"mo2/database"
	"mo2/dto"
	"mo2/mo2utils"
	"mo2/mo2utils/mo2errors"
	"mo2/server/controller/badresponse"
	"mo2/server/model"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

const passAuthKey = "passAuth"
const reasonKey = "reason"

// UpsertBlog godoc
// @Summary Publish Blog
// @Description add by json
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param draft query bool false "bool true" true
// @Param account body model.Blog true "Add blog"
// @Success 201 {object} model.Blog
// @Success 204 {object}
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("内容含非法字符，请检查"))
		return
	}
	JudgeAuthorize(ctx, &b)
	if passAuthValue, passAuthExist := ctx.Get(passAuthKey); passAuthExist {
		if passAuthValue.(bool) {
			if success := database.UpsertBlog(&b, isDraft); success {
				ctx.Header("location", ctx.FullPath())
				ctx.JSON(http.StatusCreated, b)
			} else {
				ctx.AbortWithStatusJSON(http.StatusConflict, badresponse.SetResponseReason("访问冲突"))
			}
			return
		} else {
			if err, ok := ctx.Get(reasonKey); ok {
				merr, ok := err.(mo2errors.Mo2Errors)
				if ok && merr.ErrorCode == mo2errors.Mo2NoLogin {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseReason(merr.ErrorTip))
					return
				}
			}
		}

	}
	ctx.Status(http.StatusNoContent)
}

// DeleteBlog godoc
// @Summary delete Blog
// @Description delete by path
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param draft query bool true "bool true" true
// @Param id path string false "string xxxxxxxx" "xxxxxxx"
// @Success 202
// @Success 204
// @Failure 304
// @Failure 400 {object} ResponseError
// @Failure 401 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Router /api/blogs/{id} [delete]
func (c *Controller) DeleteBlog(ctx *gin.Context) {
	isDraftStr := ctx.DefaultQuery("draft", "true")
	isDraft := true
	if isDraftStr == "false" {
		isDraft = false
	}
	idStr := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
		return
	}
	blog := database.FindBlogById(id, isDraft)
	if blog.ID.IsZero() {
		ctx.AbortWithStatusJSON(http.StatusNotFound, badresponse.SetResponseReason("页面找不到了"))
		return
	}

	JudgeAuthorize(ctx, &blog)
	if passAuth, passAuthExist := ctx.Get(passAuthKey); passAuthExist {
		if passAuth.(bool) {
			blog.EntityInfo.IsDeleted = true
			if success := database.UpsertBlog(&blog, isDraft); success {
				ctx.Status(http.StatusAccepted)
			} else {
				ctx.AbortWithStatusJSON(http.StatusConflict, badresponse.SetResponseReason("访问冲突"))
			}
			return
		} else {
			if err, ok := ctx.Get(reasonKey); ok {
				if merr := err.(mo2errors.Mo2Errors); merr.ErrorCode == mo2errors.Mo2NoLogin {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseReason(merr.ErrorTip))
					return
				}
			}
		}
	}
	ctx.Status(http.StatusNoContent)
}

// JudgeAuthorize only for user of same ID
// 只有本人id与blog的authorID一致可以继续
func JudgeAuthorize(ctx *gin.Context, blog *model.Blog) {
	userInfo, exist := mo2utils.GetUserInfo(ctx)
	if blog.AuthorID == primitive.NilObjectID {
		if exist {
			blog.AuthorID = userInfo.ID
		} else {
			ctx.Set(passAuthKey, false)
			ctx.Set(reasonKey, *mo2errors.New(mo2errors.Mo2NoLogin, "权限不足，请先登录"))
			return
		}
	} else {
		if blog.AuthorID != userInfo.ID {
			ctx.Set(passAuthKey, false)
			ctx.Set(reasonKey, *mo2errors.New(mo2errors.Mo2Unauthorized, "没有权限修改文章"))
			return
		}
	}
	ctx.Set(passAuthKey, true)
	return
}

// RestoreBlog godoc
// @Summary restore Blog
// @Description restore by path
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param draft query bool true "bool true" true
// @Param id path string false "string xxxxxxxx" "xxxxxxx"
// @Success 200 {object} model.Blog
// @Failure 400 {object} ResponseError
// @Failure 401 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Router /api/blogs/{id} [put]
func (c *Controller) RestoreBlog(ctx *gin.Context) {
	isDraftStr := ctx.DefaultQuery("draft", "true")
	isDraft := true
	if isDraftStr == "false" {
		isDraft = false
	}
	idStr := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
		return
	}
	blog := database.FindBlogById(id, isDraft)
	if blog.ID.IsZero() {
		ctx.AbortWithStatusJSON(http.StatusNotFound, badresponse.SetResponseReason("页面找不到了"))
		return
	}

	JudgeAuthorize(ctx, &blog)
	if passAuth, passAuthExist := ctx.Get(passAuthKey); passAuthExist {
		if passAuth.(bool) {
			blog.EntityInfo.IsDeleted = false
			if success := database.UpsertBlog(&blog, isDraft); success {
				ctx.Status(http.StatusAccepted)
			} else {
				ctx.AbortWithStatusJSON(http.StatusConflict, badresponse.SetResponseReason("访问冲突"))
			}
			return
		} else {
			if err, ok := ctx.Get(reasonKey); ok {
				if merr := err.(mo2errors.Mo2Errors); merr.ErrorCode == mo2errors.Mo2NoLogin {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseReason(merr.ErrorTip))
					return
				}
			}
		}
	}
	ctx.Status(http.StatusNoContent)
}

// FindBlogsByUser godoc
// @Summary find Blog
// @Description
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param draft query bool false "bool true" true
// @Param deleted query bool false "bool default false" false
// @Param page query int false "int 0" 0
// @Param pageSize query int false "int 5" 5
// @Success 200 {object} []dto.QueryBlogs
// @Success 204 {object} []dto.QueryBlogs
// @Failure 400 {object} ResponseError
// @Failure 401 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Router /api/blogs/find/own [get]
func (c *Controller) FindBlogsByUser(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "0")
	pageSizeStr := ctx.DefaultQuery("pageSize", "5")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("query with page"))
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("query with pageSize"))
	}
	isDraftStr := ctx.DefaultQuery("draft", "true")
	isDraft := parseString2Bool(isDraftStr)
	isDeletedStr := ctx.DefaultQuery("deleted", "false")
	isDeleted := parseString2Bool(isDeletedStr)
	// get user info due to cookie information
	info, ext := mo2utils.GetUserInfo(ctx)
	if !ext {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseReason("权限不足，请先登录"))
		return
	}

	blogs := database.FindBlogsByUser(info, model.Filter{
		IsDraft:   isDraft,
		IsDeleted: isDeleted,
	})
	qBlogs := dto.QueryBlogs{}
	qBlogs.Init(blogs)
	if results, exist := qBlogs.Query(page, pageSize); exist {
		ctx.JSON(http.StatusOK, results.GetBlogs())
	} else {
		ctx.JSON(http.StatusNoContent, results.GetBlogs())
	}
}

// FindBlogsByUserId godoc
// @Summary find Blog by userid
// @Description
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param draft query bool false "bool true" true
// @Param deleted query bool false "bool default false" false
// @Param id query string false "string xxxxxxxx" "xxxxxxx"
// @Param page query int false "int 0" 0
// @Param pageSize query int false "int 5" 5
// @Success 200 {object} []dto.QueryBlogs
// @Success 204 {object}
// @Failure 400 {object} ResponseError
// @Failure 401 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Router /api/blogs/find/userId [get]
func (c *Controller) FindBlogsByUserId(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "0")
	pageSizeStr := ctx.DefaultQuery("pageSize", "5")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("query with page"))
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("query with pageSize"))
	}
	isDraftStr := ctx.DefaultQuery("draft", "true")
	isDraft := parseString2Bool(isDraftStr)
	isDeletedStr := ctx.DefaultQuery("deleted", "false")
	isDeleted := parseString2Bool(isDeletedStr)
	idStr := ctx.Query("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
		return
	}
	var passAuth bool
	if isDraft {
		JudgeAuthorize(ctx, &model.Blog{AuthorID: id})
		if passAuthValue, passAuthExist := ctx.Get(passAuthKey); passAuthExist {
			passAuth, _ = passAuthValue.(bool)
			if !passAuth {
				if err, ok := ctx.Get(reasonKey); ok {
					if merr := err.(mo2errors.Mo2Errors); merr.ErrorCode == mo2errors.Mo2NoLogin {
						ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseReason(merr.ErrorTip))
						return
					}
				}
			}
		}
		ctx.Status(http.StatusNoContent)
	}
	if !isDraft || passAuth {
		blogs := database.FindBlogsByUserId(id, model.Filter{
			IsDraft:   isDraft,
			IsDeleted: isDeleted,
		})
		qBlogs := dto.QueryBlogs{}
		qBlogs.Init(blogs)
		if results, exist := qBlogs.Query(page, pageSize); exist {
			ctx.JSON(http.StatusOK, results.GetBlogs())
		} else {
			ctx.Status(http.StatusNoContent)
		}
	}
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
// @Failure 400 {object} ResponseError
// @Failure 401 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Router /api/blogs/find/id [get]
func (c *Controller) FindBlogById(ctx *gin.Context) {
	isDraftStr := ctx.DefaultQuery("draft", "true")
	isDraft := parseString2Bool(isDraftStr)
	idStr := ctx.Query("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
		return
	}
	blog := database.FindBlogById(id, isDraft)
	if blog.ID.IsZero() {
		ctx.AbortWithStatusJSON(http.StatusNotFound, badresponse.SetResponseReason("页面找不到了"))
		return
	}
	var passAuth bool
	if isDraft {
		JudgeAuthorize(ctx, &blog)
		if passAuthValue, passAuthExist := ctx.Get(passAuthKey); passAuthExist {
			passAuth, _ = passAuthValue.(bool)
			if !passAuth {
				if err, ok := ctx.Get(reasonKey); ok {
					if merr := err.(mo2errors.Mo2Errors); merr.ErrorCode == mo2errors.Mo2NoLogin {
						ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseReason(merr.ErrorTip))
						return
					}
				}
			}
		}
		ctx.Status(http.StatusNoContent)
	}
	if !isDraft || passAuth {
		ctx.JSON(http.StatusOK, blog)
	}
}

// QueryBlogs godoc
// @Summary find all Blogs
// @Description find
// @Tags blogs
// @Produce  json
// @Param draft query bool false "bool default false" false
// @Param deleted query bool false "bool default false" false
// @Param page query int false "int 0" 0
// @Param pageSize query int false "int 5" 5
// @Success 200 {object} []dto.QueryBlog
// @Success 204 {object} []dto.QueryBlog
// @Failure 400 {object} ResponseError
// @Failure 401 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Router /api/blogs/query [get]
func (c *Controller) QueryBlogs(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "0")
	pageSizeStr := ctx.DefaultQuery("pageSize", "5")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("query with page"))
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("query with pageSize"))
	}

	isDraftStr := ctx.DefaultQuery("draft", "false")
	isDraft := parseString2Bool(isDraftStr)
	isDeletedStr := ctx.DefaultQuery("deleted", "false")
	isDeleted := parseString2Bool(isDeletedStr)
	if (isDraft || isDeleted) && !mo2utils.IsInRole(ctx, model.GeneralAdmin) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseReason("权限不足，请联系管理员"))
		return
	}
	blogs := database.FindAllBlogs(model.Filter{
		IsDraft:   isDraft,
		IsDeleted: isDeleted,
	})
	qBlogs := dto.QueryBlogs{}
	qBlogs.Init(blogs)
	if results, exist := qBlogs.Query(page, pageSize); exist {
		ctx.JSON(http.StatusOK, results.GetBlogs())
	} else {
		ctx.JSON(http.StatusNoContent, results.GetBlogs())
	}
}

func parseString2Bool(s string) (b bool) {
	b = true
	if s == "false" {
		b = false
	}
	return
}
