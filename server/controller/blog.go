package controller

import (
	"mo2/database"
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
// @Success 204
// @Failure 400 {object} badresponse.ResponseError
// @Failure 401 {object} badresponse.ResponseError
// @Router /api/blogs/publish [post]
func (c *Controller) UpsertBlog(ctx *gin.Context) {
	isDraft := parseString2Bool(ctx.DefaultQuery("draft", "true"))
	var b model.Blog
	if err := ctx.ShouldBindJSON(&b); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("内容含非法字符，请检查"))
		return
	}
	JudgeAuthorize(ctx, &b)
	if passAuthValue, passAuthExist := ctx.Get(passAuthKey); passAuthExist {
		if passAuthValue.(bool) {
			if mErr := database.UpsertBlog(&b, isDraft); mErr.IsError() {
				ctx.AbortWithStatusJSON(http.StatusConflict, badresponse.SetResponseReason("访问冲突"))
			} else {
				ctx.Header("location", ctx.FullPath())
				ctx.JSON(http.StatusCreated, b)
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
// @Summary 彻底删除blog
// @Description delete by id path(draft/blog)
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param draft query bool true "bool true" true
// @Param id path string false "string xxxxxxxx" "xxxxxxx"
// @Success 202
// @Success 204
// @Failure 304
// @Failure 400 {object} badresponse.ResponseError
// @Failure 401 {object} badresponse.ResponseError
// @Failure 404 {object} badresponse.ResponseError
// @Router /api/blogs/{id} [delete]
func (c *Controller) DeleteBlog(ctx *gin.Context) {
	isDraft := parseString2Bool(ctx.DefaultQuery("draft", "true"))
	var blog model.Blog
	if id, err := primitive.ObjectIDFromHex(ctx.Param("id")); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
		return
	} else {
		blog = database.FindBlogById(id, isDraft)
		if blog.ID.IsZero() {
			ctx.AbortWithStatusJSON(http.StatusNotFound, badresponse.SetResponseReason("页面找不到了"))
			return
		}
	}
	JudgeAuthorize(ctx, &blog)
	if passAuth, passAuthExist := ctx.Get(passAuthKey); passAuthExist {
		if passAuth.(bool) {
			blog.EntityInfo.IsDeleted = true
			mo2utils.DeleteBlogIndex(blog.ID.Hex())
			if mErr := database.DeleteBlogs(isDraft, blog.ID); mErr.IsError() {
				ctx.AbortWithStatusJSON(http.StatusConflict, badresponse.SetResponseReason("访问冲突"))
			} else {
				ctx.Status(http.StatusAccepted)
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

// ProcessBlog godoc
// @Summary restore Blog
// @Description restore by path
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param draft query bool true "bool true" true
// @Param operation path string false "recycle/restore" "different type operation"
// @Param id path string false "Blog id" "objectID"
// @Success 202
// @Failure 400 {object} badresponse.ResponseError
// @Failure 401 {object} badresponse.ResponseError
// @Failure 404 {object} badresponse.ResponseError
// @Router /api/blogs/{operation}/{id} [put]
func (c *Controller) ProcessBlog(ctx *gin.Context) {
	isDraft := parseString2Bool(ctx.DefaultQuery("draft", "true"))
	var blog model.Blog
	if id, err := primitive.ObjectIDFromHex(ctx.Param("id")); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
		return
	} else {
		blog = database.FindBlogById(id, isDraft)
		if blog.ID.IsZero() {
			ctx.AbortWithStatusJSON(http.StatusNotFound, badresponse.SetResponseReason("页面找不到了"))
			return
		}
	}
	JudgeAuthorize(ctx, &blog)
	if passAuth, passAuthExist := ctx.Get(passAuthKey); passAuthExist {
		if passAuth.(bool) {
			if mErr := database.ProcessBlog(isDraft, &blog, ctx.Param(database.OperationKey)); mErr.IsError() {
				ctx.AbortWithStatusJSON(http.StatusConflict, badresponse.SetResponseError(mErr))
			} else {
				if mErr = database.UpsertBlog(&blog, isDraft); mErr.IsError() {
					ctx.AbortWithStatusJSON(http.StatusConflict, badresponse.SetResponseError(mErr))
				} else {
					ctx.Status(http.StatusAccepted)
				}
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
// @Success 200 {object} []model.Blog
// @Success 204 {object} []model.Blog
// @Failure 400 {object} badresponse.ResponseError
// @Failure 401 {object} badresponse.ResponseError
// @Failure 404 {object} badresponse.ResponseError
// @Router /api/blogs/find/own [get]
func (c *Controller) FindBlogsByUser(ctx *gin.Context) {
	filter, mErr := parseFilter(ctx)
	if mErr.IsError() {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseError(mErr))
		return
	}
	// get user info due to cookie information
	info, ext := mo2utils.GetUserInfo(ctx)
	if !ext {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseReason("权限不足，请先登录"))
		return
	}

	blogs := database.FindBlogsByUser(info, filter)
	ctx.JSON(http.StatusOK, blogs)
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
// @Success 200 {object} []model.Blog
// @Success 204
// @Failure 400 {object} badresponse.ResponseError
// @Failure 401 {object} badresponse.ResponseError
// @Failure 404 {object} badresponse.ResponseError
// @Router /api/blogs/find/userId [get]
func (c *Controller) FindBlogsByUserId(ctx *gin.Context) {
	filter, mErr := parseFilter(ctx)
	if mErr.IsError() {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseError(mErr))
		return
	}
	idStr := ctx.Query("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
		return
	}
	var passAuth bool
	if filter.IsDraft {
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
	if !filter.IsDraft || passAuth {
		blogs := database.FindBlogsByUserId(id, filter)
		ctx.JSON(http.StatusOK, blogs)
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
// @Failure 400 {object} badresponse.ResponseError
// @Failure 401 {object} badresponse.ResponseError
// @Failure 404 {object} badresponse.ResponseError
// @Router /api/blogs/find/id [get]
func (c *Controller) FindBlogById(ctx *gin.Context) {
	isDraft := parseString2Bool(ctx.DefaultQuery("draft", "true"))
	id, err := primitive.ObjectIDFromHex(ctx.Query("id"))
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
// @Param search query string false "aaaa" "aaa"
// @Success 200 {object} []model.Blog
// @Success 204 {object} []model.Blog
// @Failure 400 {object} badresponse.ResponseError
// @Failure 401 {object} badresponse.ResponseError
// @Failure 404 {object} badresponse.ResponseError
// @Router /api/blogs/query [get]
func (c *Controller) QueryBlogs(ctx *gin.Context) {
	filter, mErr := parseFilter(ctx)
	if mErr.IsError() {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseError(mErr))
		return
	}
	if (filter.IsDraft || filter.IsDeleted) && !mo2utils.IsInRole(ctx, model.GeneralAdmin) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseReason("权限不足，请联系管理员"))
		return
	}
	search := ctx.Query("search")
	var ids []primitive.ObjectID
	if search != "" {
		hits := mo2utils.QueryBlog(search)
		len := filter.PageSize
		if hits.Len() < filter.PageSize {
			len = hits.Len()
		}
		ids = make([]primitive.ObjectID, len)
		for i, v := range hits {
			if i == len {
				break
			}
			ids[i], _ = primitive.ObjectIDFromHex(v.ID)
		}
	}
	blogs := database.FindBlogs(model.Filter{
		IsDraft:   filter.IsDraft,
		IsDeleted: filter.IsDeleted,
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		Ids:       ids,
	})
	ctx.JSON(http.StatusOK, blogs)
}

// FindBlogsByType godoc
// @Summary find blogs by given type
// @Description  根据type返回不同结果：[category] 所有category包含的blog
// @Tags relation
// @Accept  json
// @Produce  json
// @Param type path string true "find by category"
// @Param ID path string true "ID"
// @Success 200 {object} []model.Blog
// @Failure 400 {object} badresponse.ResponseError
// @Failure 404 {object} badresponse.ResponseError
// @Router /api/relation/blogs/{type}/{ID} [get]
func (c *Controller) FindBlogsByType(ctx *gin.Context) {
	idStr := ctx.Param("ID")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
		return
	}
	var blogs []model.Blog
	var mErr mo2errors.Mo2Errors
	switch ctx.Param(typeKey) {
	case typeCategory:
		blogs, mErr = database.FindBlogsByCategoryId(id, model.Filter{IsDraft: false, IsDeleted: false})
	}
	if mErr.IsError() {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseError(mErr))
		return
	}
	ctx.JSON(http.StatusOK, blogs)
}

// IndexAllBlogs godoc
// @Summary index all blogs
// @Description  none
// @Tags admin
// @Success 200
// @Router /api/admin/indexblogs [post]
func (c *Controller) IndexAllBlogs(ctx *gin.Context) {
	mo2utils.IndexBlogs(database.FindBlogs(model.Filter{IsDeleted: false, IsDraft: false, Page: 0, PageSize: 1000}))
	ctx.Status(200)
}

func parseString2Bool(s string) (b bool) {
	b = true
	if s == "false" {
		b = false
	}
	return
}

func parseFilter(ctx *gin.Context) (filter model.Filter, mErr mo2errors.Mo2Errors) {
	if page, err := strconv.Atoi(ctx.DefaultQuery("page", "0")); err != nil {
		mErr.Init(mo2errors.Mo2Error, "query with page")
		return
	} else {
		filter.Page = page
	}
	if pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "5")); err != nil {
		mErr.Init(mo2errors.Mo2Error, "query with pageSize")
		return
	} else {
		filter.PageSize = pageSize
	}
	filter.IsDraft = parseString2Bool(ctx.DefaultQuery("draft", "true"))
	filter.IsDeleted = parseString2Bool(ctx.DefaultQuery("deleted", "false"))
	return
}
