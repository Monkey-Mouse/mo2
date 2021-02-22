package controller

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mo2/database"
	"mo2/mo2utils"
	"mo2/server/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PublishBlog godoc
// @Summary Publish Blog
// @Description add by json
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param account body model.Blog true "Add blog"
// @Success 200 {object} model.Blog
// @Router /api/blogs/publish [post]
func (c *Controller) PublishBlog(ctx *gin.Context) {
	b := model.Blog{}
	if err := ctx.ShouldBindJSON(&b); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("内容含非法字符，请检查"))
		return
	}
	// set blog's author due to cookie information
	info, ext := mo2utils.GetUserInfo(ctx)
	if !ext {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, SetResponseReason("权限不足，请先登录"))
		return
	}
	b.AuthorID = info.ID
	success := database.UpsertBlog(&b)
	if success {
		ctx.JSON(http.StatusOK, &b)
	} else {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("网络繁忙，请稍后重试"))
	}
}

// SaveDraft godoc
// @Summary save Draft
// @Description add by json
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param account body model.Draft true "Save draft"
// @Success 200 {object} model.Draft
// @Router /api/blogs/saveDraft [post]
func (c *Controller) SaveDraft(ctx *gin.Context) {
	d := model.Draft{}
	if err := ctx.ShouldBindJSON(&d); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("内容含非法字符，请检查"))
		return
	}
	success := database.UpsertDraft(&d)
	if success {
		if d.BlogObj.AuthorID == primitive.NilObjectID {
			// set blog's author due to cookie information
			info, ext := mo2utils.GetUserInfo(ctx)
			if !ext {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, SetResponseReason("权限不足，请先登录"))
				return
			}
			d.BlogObj.AuthorID = info.ID
		}
		ctx.JSON(http.StatusOK, &d)
	} else {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("网络繁忙，请稍后重试"))
	}
}

// FindBlogsByUser godoc
// @Summary find Blog
// @Description add by json
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

	blogs := database.FindBlogs(info)
	ctx.JSON(http.StatusOK, blogs)
}

// FindAllBlogs godoc
// @Summary find all Blogs
// @Description find
// @Tags blogs
// @Produce  json
// @Success 200 {object} []model.Blog
// @Router /api/blogs/find/all [get]
func (c *Controller) FindAllBlogs(ctx *gin.Context) {
	// get user info due to cookie information
	//info, ext := mo2utils.GetUserInfo(ctx)
	//if !ext {
	//	ctx.AbortWithStatusJSON(http.StatusUnauthorized, SetResponseReason("权限不足，请先登录"))
	//	return
	//}

	blogs := database.FindAllBlogs()
	ctx.JSON(http.StatusOK, blogs)
}
