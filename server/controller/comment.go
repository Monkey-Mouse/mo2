package controller

import (
	"mo2/database"
	"mo2/mo2utils"
	"mo2/server/controller/badresponse"
	"mo2/server/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetComment get comments
// @Summary get comments
// @Description get json comments
// @Tags comments
// @Accept  json
// @Produce  json
// @Param id path string true "article id"
// @Param page query int false "int 0" 0
// @Param pagesize query int false "int 5" 5
// @Success 200 {object} []model.Comment
// @Failure 422 {object} badresponse.ResponseError
// @Router /api/comment/{id} [get]
func (c *Controller) GetComment(ctx *gin.Context) {
	sid := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(sid)
	if err != nil {
		badresponse.SetErrResponse(ctx, http.StatusUnprocessableEntity, "格式错误")
		return
	}
	page, pagesize, err := mo2utils.ParsePagination(ctx)
	if err != nil {
		badresponse.SetErrResponse(ctx, http.StatusUnprocessableEntity, "格式错误")
		return
	}
	cs := database.GetComments(id, page, pagesize)
	ctx.JSON(200, cs)
}

// PostComment godoc
// @Summary upsert comments
// @Description upsert json comments
// @Tags comments
// @Accept  json
// @Produce  json
// @Param comment body model.Comment true "comment"
// @Success 200 {object} model.Comment
// @Failure 422 {object} badresponse.ResponseError
// @Router /api/comment [post]
func (c *Controller) PostComment(ctx *gin.Context) {
	var cmt model.Comment
	if ctx.ShouldBindJSON(&cmt) != nil {
		badresponse.SetErrResponse(ctx, http.StatusUnprocessableEntity, "格式错误")
		return
	}
	u, _ := mo2utils.GetUserInfo(ctx)
	cmt.Aurhor = u.ID
	database.UpsertComment(&cmt)
	ctx.JSON(200, &cmt)
}

// PostSubComment post subcomments
// @Summary upsert subcomments
// @Description upsert json comments
// @Tags comments
// @Accept  json
// @Produce  json
// @Param id path string true "comment id"
// @Param comment body model.Subcomment true "subcomment"
// @Success 200 {object} model.Subcomment
// @Failure 422 {object} badresponse.ResponseError
// @Router /api/comment/{id} [post]
func (c *Controller) PostSubComment(ctx *gin.Context) {
	sid := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(sid)
	if err != nil {
		badresponse.SetErrResponse(ctx, http.StatusUnprocessableEntity, "格式错误")
		return
	}
	var cmt model.Subcomment
	if ctx.ShouldBindJSON(&cmt) != nil {
		badresponse.SetErrResponse(ctx, http.StatusUnprocessableEntity, "格式错误")
		return
	}
	u, _ := mo2utils.GetUserInfo(ctx)
	cmt.Aurhor = u.ID
	database.UpsertSubComment(id, &cmt)
	ctx.JSON(200, &cmt)
}

// DeleteComment delete comments
func (c *Controller) DeleteComment(ctx *gin.Context) {

}

// DeleteSubComment delete subcomments
func (c *Controller) DeleteSubComment(ctx *gin.Context) {

}
