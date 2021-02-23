package controller

import (
	"github.com/gin-gonic/gin"
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
// @Description add by json
// @Tags category
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.Category
// @Router /api/blogs/findAllCategories [get]
func (c *Controller) FindAllCategories(ctx *gin.Context) {
	cats := database.FindAllCategories()
	ctx.JSON(http.StatusOK, cats)
}
