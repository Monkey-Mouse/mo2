package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Monkey-Mouse/go-abac/abac"
	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/Monkey-Mouse/mo2/services/accessControl"

	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/mo2utils"
	"github.com/Monkey-Mouse/mo2/mo2utils/mo2errors"
	"github.com/Monkey-Mouse/mo2/server/controller/badresponse"
	"github.com/Monkey-Mouse/mo2/server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

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
// @Success 204
// @Failure 400 {object} badresponse.ResponseError
// @Failure 401 {object} badresponse.ResponseError
// @Router /api/blogs/publish [post]
func (c *Controller) UpsertBlog(ctx *gin.Context) {
	// TODO 允许持有正确的token的用户保存文章
	isDraft := parseString2Bool(ctx.DefaultQuery("draft", "true"))
	var blog model.Blog
	if err := ctx.ShouldBindJSON(&blog); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("内容含非法字符，请检查"))
		return
	}
	if mErr := JudgeAuthorize(ctx, &blog); mErr.IsError() {
		if mErr.ErrorCode == mo2errors.Mo2NoLogin {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseError(mErr))
			return
		}
	} else {
		if mErr = database.UpsertBlog(&blog, isDraft); mErr.IsError() {
			ctx.AbortWithStatusJSON(http.StatusConflict, badresponse.SetResponseReason("访问冲突"))
		} else {
			ctx.Header("location", ctx.FullPath())
			ctx.JSON(http.StatusCreated, blog)
		}
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *Controller) SetDocType(ctx *gin.Context, u dto.LoginUserInfo) (status int, body interface{}, err error) {
	var blog model.Blog
	if err = ctx.ShouldBindJSON(&blog); err != nil {
		status = 400
		err = errors.New("内容含非法字符，请检查")
		return
	}
	if blog.IsYDoc {
		blog.YToken = primitive.NewObjectID()
	}
	_, err = database.DraftCol.UpdateOne(ctx,
		bson.D{{"_id", blog.ID}, {"author_id", u.ID}},
		bson.M{"$set": bson.M{
			"y_doc":    blog.YDoc,
			"y_token":  blog.YToken,
			"is_y_doc": blog.IsYDoc},
		},
	)
	if err != nil {
		status = http.StatusUnprocessableEntity
		err = errors.New("db error")
		return
	}
	return 200, gin.H{"token": blog.YToken}, nil
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
		if blog = database.FindBlogById(id, isDraft); blog.ID.IsZero() {
			ctx.AbortWithStatusJSON(http.StatusNotFound, badresponse.SetResponseReason("页面找不到了"))
			return
		}
	}
	if mErr := JudgeAuthorize(ctx, &blog); mErr.IsError() {
		if mErr.ErrorCode == mo2errors.Mo2NoLogin {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseError(mErr))
			return
		}
	} else {
		if mErr = database.DeleteBlogs(isDraft, blog.ID); mErr.IsError() {
			ctx.AbortWithStatusJSON(http.StatusConflict, badresponse.SetResponseReason("访问冲突"))
		} else {
			ctx.Status(http.StatusAccepted)
		}
		return
	}
	ctx.Status(http.StatusNoContent)
}

var (
	MErrNoLogin      = mo2errors.Init(mo2errors.Mo2NoLogin, "权限不足，请先登录")
	MErrUnauthorized = mo2errors.Init(mo2errors.Mo2Unauthorized, "没有权限修改文章")
)

// JudgeAuthorize only for user of same ID
// 只有本人id与blog的authorID一致可以继续
func JudgeAuthorize(ctx *gin.Context, blog *model.Blog) (mErr mo2errors.Mo2Errors) {
	userInfo, exist := mo2utils.GetUserInfo(ctx)
	if blog.AuthorID == primitive.NilObjectID {
		if exist {
			blog.AuthorID = userInfo.ID
		} else {

			return MErrNoLogin
		}
	} else {
		if blog.AuthorID != userInfo.ID {
			return MErrUnauthorized
		}
	}
	mErr.InitCode(mo2errors.Mo2NoError)
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
	var pass bool
	if id, err := primitive.ObjectIDFromHex(ctx.Param("id")); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
		return
	} else {
		blog = database.FindBlogById(id, isDraft)
		if info, exist := mo2utils.GetUserInfo(ctx); exist {
			pass, err = accessControl.Ctrl.CanOr(abac.IQueryInfo{
				Subject:  "account",
				Action:   abac.ActionUpdate,
				Resource: "blog",
				Context:  abac.DefaultContext{"allowOwn": accessControl.AllowOwn{ID: id, Filter: model.Filter{IsDraft: isDraft}, UserInfo: info, Resource: "blog"}},
			})
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusConflict, badresponse.SetResponseError(err))
				return
			}

		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseError(MErrUnauthorized))
			return
		}
	}
	if pass {
		if mErr := database.ProcessBlog(isDraft, &blog, ctx.Param(database.OperationKey)); mErr.IsError() {
			ctx.AbortWithStatusJSON(http.StatusConflict, badresponse.SetResponseError(mErr))
			return
		} else {
			if mErr = database.UpsertBlog(&blog, isDraft); mErr.IsError() {
				ctx.AbortWithStatusJSON(http.StatusConflict, badresponse.SetResponseError(mErr))
				return
			} else {
				ctx.Status(http.StatusAccepted)
				return
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
func (c *Controller) FindBlogsByID(ctx *gin.Context) {
	filter, mErr := parseFilter(ctx)
	if mErr.IsError() {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseError(mErr))
		return
	}
	field := ctx.Query("field")
	id, err := primitive.ObjectIDFromHex(ctx.Query("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
		return
	}
	if filter.IsDraft {
		if mErr = JudgeAuthorize(ctx, &model.Blog{AuthorID: id}); mErr.IsError() {
			if mErr.ErrorCode == mo2errors.Mo2NoLogin {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseError(mErr))
				return
			}
		}
	}
	if !filter.IsDraft || !mErr.IsError() {
		blogs := database.FindBlogsByValue(field, id, filter)
		ctx.JSON(http.StatusOK, blogs)
		return
	}
	ctx.Status(http.StatusNoContent)
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
	token, _ := primitive.ObjectIDFromHex(ctx.DefaultQuery("token", primitive.NilObjectID.Hex()))
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
	if token != primitive.NilObjectID && blog.YToken == token {
		ctx.JSON(http.StatusOK, blog)
		return
	}
	var mErr mo2errors.Mo2Errors
	if isDraft {
		if mErr = JudgeAuthorize(ctx, &blog); mErr.IsError() {
			if mErr.ErrorCode == mo2errors.Mo2NoLogin {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseError(mErr))
				return
			}
		}
	}
	uinfo, _ := mo2utils.GetUserInfo(ctx)
	if uinfo.ID != blog.AuthorID {
		blog.YToken = primitive.NilObjectID
	}
	if !isDraft || !mErr.IsError() {
		ctx.JSON(http.StatusOK, blog)
		return
	}
	ctx.Status(http.StatusNoContent)
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
	searchS := ctx.Query("search")
	if searchS != "" {
		hits := mo2utils.QueryBlog(searchS, filter.Page, filter.PageSize)
		var re = make([]map[string]interface{}, hits.Len())
		for i, v := range hits {
			re[i] = v.Fields
			for s, u := range v.Fragments {
				re[i][s] = u[0]
			}
		}

		ctx.JSON(http.StatusOK, re)
		return
	}
	blogs := database.FindBlogs(filter)
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
func (c *Controller) IndexAll(ctx *gin.Context) {
	mo2utils.IndexAccounts(database.FindAllAccounts())
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
	ids := ctx.QueryArray("ids")
	filter.Ids = make([]primitive.ObjectID, len(ids))
	for i, v := range ids {
		filter.Ids[i], _ = primitive.ObjectIDFromHex(v)
	}
	return
}
