package controller

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mo2/database"
	"mo2/server/controller/badresponse"
	"net/http"
)

// ListAccountsInfo godoc
// @Summary List directories brief info
// @Description from a list of directory ids [usage]:/api/directories/{collection}?id=60223d4042d6febff9f276f0&id=60236866d2a68483adaccc38
// @Tags directory
// @Accept  json
// @Produce  json
// @Param collection path string true "category/..."
// @Param id query array true "directory IDs list"
// @Success 200 {object} []model.Directory
// @Router /api/directories/{collection} [get]
func (c *Controller) ListDirectoriesInfo(ctx *gin.Context) {
	directoryIDstrs, exist := ctx.GetQueryArray("id")
	var directoryIDs []primitive.ObjectID
	for _, idStr := range directoryIDstrs {
		if id, err := primitive.ObjectIDFromHex(idStr); err != nil {

		} else {
			directoryIDs = append(directoryIDs, id)
		}
	}
	col := ctx.Param("collection")
	if exist {
		directories, mErr := database.FindDirectoryInfo(col, directoryIDs...)
		if mErr.IsError() {
			ctx.AbortWithStatusJSON(http.StatusConflict, badresponse.SetResponseError(mErr))
		} else {
			ctx.JSON(http.StatusOK, directories)
		}
	}

}
