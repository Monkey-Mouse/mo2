package controller

import (
	"mo2/database"
	"mo2/dto"
	"mo2/mo2utils"
	"mo2/mo2utils/mo2errors"
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
// @Param account body model.Directory true "Add category"
// @Success 200 {object} model.Directory
// @Router /api/blogs/category [post]
func (c *Controller) UpsertCategory(ctx *gin.Context) {
	var cat model.Directory
	if err := ctx.ShouldBindJSON(&cat); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
		return
	}
	database.UpsertCategory(&cat)
	ctx.JSON(http.StatusOK, cat)
}

// DeleteCategory godoc
// @Summary delete category
// @Description 根据id删除，并解除与之相关实体的所有关联
// @Tags category
// @Accept  json
// @Produce  json
// @Param ids body []primitive.ObjectID true "category id to delete"
// @Success 200 {object} model.Directory
// @Router /api/directories/category [delete]
func (c *Controller) DeleteCategory(ctx *gin.Context) {
	var catIDs []primitive.ObjectID
	if err := ctx.ShouldBindJSON(&catIDs); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
		return
	}
	if userInfo, exist := mo2utils.GetUserInfo(ctx); exist {
		if allowIDs, mErr := database.RightFilter(userInfo.ID, catIDs...); mErr.IsError() {
			ctx.AbortWithStatusJSON(http.StatusConflict, badresponse.SetResponseError(mErr))
			return
		} else {
			if mErr = database.DeleteCategoryCompletely(allowIDs...); mErr.IsError() {
				ctx.AbortWithStatusJSON(http.StatusConflict, badresponse.SetResponseError(mErr))
				return
			} else {
				ctx.Status(http.StatusOK)
			}
		}
	} else {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseReason("请先登录"))
		return
	}
}

// FindAllCategories godoc
// @Summary find categories
// @Description 若id为空，返回所有categories；若id不为空，返回该id的category
// @Tags category
// @Produce  json
// @Param id query string false "string ObjectID" ""
// @Success 200 {object} []model.Directory
// @Router /api/blogs/category [get]
func (c *Controller) FindAllCategories(ctx *gin.Context) {
	idStr := ctx.Query("id")
	var cats []model.Directory
	if len(idStr) == 0 {
		cats = database.FindAllCategories()
	} else {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
			return
		}
		var ids []primitive.ObjectID
		var mErr mo2errors.Mo2Errors
		ids = append(ids, id)
		if cats, mErr = database.FindCategories(ids); mErr.IsError() {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseError(mErr))
			return
		}
	}
	ctx.JSON(http.StatusOK, cats)

}

// RelateCategory2Entity godoc
// @Summary relate category to given type
// @Description
// @Tags relation
// @Accept  json
// @Produce  json
// @Param type path string true "relatedTo user/blog/category/userMain"
// @Param e2e body dto.RelateEntity2Entity true "category id to be related"
// @Success 200 {object} dto.RelateEntity2Entity
// @Router /api/relation/category/{type} [post]
func (c *Controller) RelateCategory2Entity(ctx *gin.Context) {
	var e2e dto.RelateEntity2Entity
	if err := ctx.ShouldBindJSON(&e2e); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法参数"))
		return
	}
	switch ctx.Param(typeKey) {
	case typeUser:
		database.RelateCategory2User(e2e.RelatedID, e2e.RelateToID)
	case typeUserMain:
		database.RelateMainCategory2User(e2e)
	case typeCategory:
		database.RelateSubCategory2Category(e2e)
	case typeBlog:
		database.RelateCategory2Blog(e2e)
	}
	ctx.JSON(http.StatusOK, e2e)

}

// FindCategoriesByType godoc
// @Summary find categories by given type
// @Description  根据type返回不同结果：[user] 个人的所有category|[sub] 所有子category
// @Tags relation
// @Accept  json
// @Produce  json
// @Param type path string true "find by user/sub"
// @Param ID path string true "ID"
// @Success 200 {object} []model.Directory
// @Failure 400 {object} badresponse.ResponseError
// @Failure 404 {object} badresponse.ResponseError
// @Router /api/relation/category/{type}/{ID} [get]
func (c *Controller) FindCategoriesByType(ctx *gin.Context) {
	idStr := ctx.Param("ID")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法输入"))
		return
	}
	var cats []model.Directory
	var mErr mo2errors.Mo2Errors
	switch ctx.Param(typeKey) {
	case typeUser:
		cats, mErr = database.FindCategoriesByUserId(id)
	case typeSubCategories:
		cats, mErr = database.FindSubCategories(id)
	}
	if mErr.IsError() {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseError(mErr))
		return
	}
	ctx.JSON(http.StatusOK, cats)
}

// RelateCategories2Entity godoc
// @Summary 将列表内的子categories关联到单个实体上
// @Description （根据path中提供的关联类型选择对应方法）目前有：父category
// @Tags relation
// @Accept  json
// @Produce  json
// @Param type path string true "types to relate"
// @Param id body dto.RelateEntitySet2Entity true "sub category id and parent id"
// @Success 200 {object} model.Directory
// @Router /api/relation/categories/{type} [post]
func (c *Controller) RelateCategories2Entity(ctx *gin.Context) {

	var s2e dto.RelateEntitySet2Entity
	if err := ctx.ShouldBindJSON(&s2e); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason("非法参数"))
		return
	}
	switch ctx.Param(typeKey) {
	case typeCategory:
		database.RelateSubCategories2Category(s2e)
	case typeBlog:
		database.RelateCategories2Blog(s2e)
	}
	ctx.JSON(http.StatusOK, s2e)
}
