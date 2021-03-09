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
// @Router /api/blogs/category [post]
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
// @Router /api/blogs/category [get]
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

// FindAllCategories godoc
// @Summary find categories
// @Description id不为空，返回该id的子目录subCategories
// @Tags category
// @Produce  json
// @Param id query string true "string ObjectID" ""
// @Success 200 {object} []model.Category
// @Router /api/blogs/category/{parentID} [get]
func (c *Controller) FindSubCategories(ctx *gin.Context) {
	idStr := ctx.Query("parentID")
	var cats []model.Category

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
	}
	var ids []primitive.ObjectID
	ids = append(ids, id)
	cats = database.FindCategories(ids)

	ctx.JSON(http.StatusOK, cats)

}

// AddBlogs2Categories godoc
// @Summary add blogs to chosen categories
// @Description blogs 与 categories皆为id列表，方便批量操作
// @Tags category
// @Accept  json
// @Produce  json
// @Param id body dto.AddBlogs2Categories true "dto.AddBlogs2Categories"
// @Success 200 {object} []model.Blog
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
// @Param userID path string true "user id"
// @Param id body primitive.ObjectID true "category ids to be added"
// @Success 200 {object} dto.AddCategory2User
// @Router /api/blogs/category/user/{userID} [post]
func (c *Controller) AddCategory2User(ctx *gin.Context) {
	var c2u primitive.ObjectID
	if err := ctx.ShouldBindJSON(&c2u); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法参数"))
		return
	}
	userID, err := primitive.ObjectIDFromHex(ctx.Param("userID"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法参数"))
		return
	}
	// todo 唯一性
	database.AddCategoryId2User(c2u, userID)
	ctx.JSON(http.StatusOK, c2u)

}

// FindCategoriesByUserId godoc
// @Summary find categories by user id
// @Description  return (main category)个人的主存档 于前端不可见，用于后端存储
// @Tags category
// @Accept  json
// @Produce  json
// @Param userID path string false "user ID"
// @Success 200 {object} []model.Category
// @Failure 400 {object} badresponse.ResponseError
// @Failure 404 {object} badresponse.ResponseError
// @Router /api/blogs/category/user/{userID} [get]
func (c *Controller) FindCategoriesByUserId(ctx *gin.Context) {
	idStr := ctx.Param("userID")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
		return
	}
	cs, mErr := database.FindCategoriesByUserId(id)
	if mErr.IsError() {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseError(mErr))
		return
	}
	ctx.JSON(http.StatusOK, cs)
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
