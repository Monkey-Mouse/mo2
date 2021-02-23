package controller

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mo2/database"
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
// @Summary find all categories
// @Description 若id为空，返回所有categories；若id不为空，返回该id的category
// @Tags category
// @Accept  json
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
