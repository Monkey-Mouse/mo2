package controller

import (
	"mo2/database"
	"mo2/dto"
	"mo2/server/controller/badresponse"
	"mo2/server/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpsertCategory godoc
// @Summary upsert category
// @Description add by json
// @Tags category
// @Accept  json
// @Produce  json
// @Param account body model.Category true "Add category"
// @Success 200 {object} model.Category
// @Router /api/blogs/addCategory [post]
func (c *Controller) UpsertCategory(ctx *gin.Context) {
	var cat model.Category
	if err := ctx.ShouldBindJSON(&cat); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
		return
	}
	database.UpsertCategory(&cat)
	ctx.JSON(http.StatusOK, cat)
}

// FindAllCategories godoc
// @Summary find categories
// @Description 若id为空，返回所有categories；若id不为空，返回该id的category
// @Tags category
// @Produce  json
// @Param id query string false "string ObjectID" ""
// @Success 200 {object} []model.Category
// @Router /api/blogs/findAllCategories [get]
func (c *Controller) FindAllCategories(ctx *gin.Context) {
	idStr := ctx.Query("id")
	var cats []model.Category
	if len(idStr) == 0 {
		cats = database.FindAllCategories()
	} else {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
		}
		var ids []primitive.ObjectID
		ids = append(ids, id)
		cats = database.FindCategories(ids)
	}
	ctx.JSON(http.StatusOK, cats)

}

// AddBlogs2Categories godoc
// @Summary add blogs to chosen categories
// @Description blogs 与 categories皆为id列表，方便批量操作
// @Tags category
// @Accept  json
// @Produce  json
// @Param id body dto.AddBlogs2Categories true "dto.AddBlogs2Categories"
// @Success 200 {object} []dto.QueryBlog
// @Router /api/blogs/addBlogs2Categories [post]
func (c *Controller) AddBlogs2Categories(ctx *gin.Context) {
	var ab2c dto.AddBlogs2Categories
	if err := ctx.ShouldBindJSON(&ab2c); err != nil {
		ctx.JSON(http.StatusBadRequest, badresponse.SetResponseReason("非法参数"))
	}
	results := database.AddBlogs2Categories(ab2c)
	ctx.JSON(http.StatusOK, results)

}

// FindCategoryByUserId godoc
// @Summary find category by user id
// @Description  return (main category)个人的主存档 于前端不可见，用于后端存储
// @Tags category
// @Produce  json
// @Param userId query string false "string ObjectID" ""
// @Success 200 {object} model.Category
// @Router /api/blogs/findCategoryByUserId [get]
func (c *Controller) FindCategoryByUserId(ctx *gin.Context) {
	idStr := ctx.Query("userId")
	var cat model.Category
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
		return
	}
	cat = database.FindCategoryByUserId(id)
	ctx.JSON(http.StatusOK, cat)
}

// AddCategory2User godoc
// @Summary add category to user
// @Description user 与 category 皆为id
// @Tags category
// @Accept  json
// @Produce  json
// @Param id body dto.AddCategory2User true "category id and user id"
// @Success 200 {object} dto.AddCategory2User
// @Router /api/blogs/addCategory2User [post]
func (c *Controller) AddCategory2User(ctx *gin.Context) {
	var c2u dto.AddCategory2User
	if err := ctx.ShouldBindJSON(&c2u); err != nil {
		ctx.JSON(http.StatusBadRequest, badresponse.SetResponseReason("非法参数"))
	}
	// todo 唯一性
	database.AddCategoryIdStr2User(c2u.CategoryID, c2u.UserID)
	ctx.JSON(http.StatusOK, c2u)

}

// FindCategoriesByUserId godoc
// @Summary find categories by user id
// @Description  return (main category)个人的主存档 于前端不可见，用于后端存储
// @Tags category
// @Produce  json
// @Param userId query string false "string ObjectID" ""
// @Success 200 {object} map[string][]model.Category
// @Router /api/blogs/findCategoriesByUserId [get]
func (c *Controller) FindCategoriesByUserId(ctx *gin.Context) {
	idStr := ctx.Query("userId")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
		return
	}
	m := database.FindCategoriesByUserId(id)
	if m == nil {
		ctx.JSON(http.StatusNotFound, badresponse.SetResponseReason("没有任何归档"))
		return
	}
	ctx.JSON(http.StatusOK, m)
}

// AddCategory2Category godoc
// @Summary add category to parent category
// @Description category为model.Category(若id存在，直接存放；否则新建) parent category 为id
// @Tags category
// @Accept  json
// @Produce  json
// @Param id body dto.AddCategory2Category true "category info and parent id"
// @Success 200 {object} model.Category
// @Router /api/blogs/addCategory2Category [post]
func (c *Controller) AddCategory2Category(ctx *gin.Context) {
	var c2c dto.AddCategory2Category
	if err := ctx.ShouldBindJSON(&c2c); err != nil {
		ctx.JSON(http.StatusBadRequest, badresponse.SetResponseReason("非法参数"))
	}
	c2c.Category.ParentID = c2c.ParentCategoryID
	database.UpsertCategory(&c2c.Category)
	ctx.JSON(http.StatusOK, c2c.Category)

}
