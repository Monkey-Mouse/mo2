package controller

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mo2/database"
	"mo2/dto"
	"mo2/server/model"
	"net/http"
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("非法输入"))
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
			ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("非法输入"))
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
		ctx.JSON(http.StatusBadRequest, SetResponseReason("非法参数"))
	}
	results := database.AddBlogs2Categories(ab2c)
	ctx.JSON(http.StatusOK, results)

}
