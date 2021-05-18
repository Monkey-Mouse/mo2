package controller

import (
	"net/http"

	"github.com/Monkey-Mouse/go-abac/abac"
	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/mo2utils"
	"github.com/Monkey-Mouse/mo2/server/controller/badresponse"
	"github.com/Monkey-Mouse/mo2/server/model"
	"github.com/Monkey-Mouse/mo2/services/accessControl"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	accountStr = "account"
	groupStr   = "group"
)

// InsertGroup godoc
// @Summary insert group
// @Description add by json
// @Tags group
// @Accept  json
// @Produce  json
// @Param group body model.Group true "Add group"
// @Success 201 {object} model.Group
// @Success 204
// @Failure 400 {object} badresponse.ResponseError
// @Failure 401 {object} badresponse.ResponseError
// @Router /api/group [post]
func (c *Controller) InsertGroup(ctx *gin.Context) {
	var group model.Group
	if err := ctx.ShouldBindJSON(&group); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason(badresponse.BadRequestReason))
		return
	}
	if userInfo, exist := mo2utils.GetUserInfo(ctx); exist {
		if pass, err := accessControl.Ctrl.CanAnd(abac.IQueryInfo{
			Subject:  accountStr,
			Action:   abac.ActionCreate,
			Resource: accessControl.ResourceGroup,
			Context:  abac.DefaultContext{},
		}); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseError(err))
			return
		} else if !pass {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason(badresponse.NoAccessReason))
			return
		} else if pass {
			if group.ID.IsZero() {
				group.ID = primitive.NewObjectID()
			}
			group.OwnerID = userInfo.ID
			if mErr := database.UpsertGroup(group); mErr.IsError() {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseError(mErr))
				return
			} else {
				ctx.JSON(http.StatusCreated, group)
				return
			}
		}
	} else {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseReason(badresponse.UnauthorizeReason))
		return
	}
}

// UpdateGroup godoc
// @Summary update group
// @Description add by json
// @Tags group
// @Accept  json
// @Produce  json
// @Param group body model.Group true "update group"
// @Success 201 {object} model.Group
// @Success 204
// @Failure 400 {object} badresponse.ResponseError
// @Failure 401 {object} badresponse.ResponseError
// @Router /api/group [put]
func (c *Controller) UpdateGroup(ctx *gin.Context) {
	var group model.Group
	if err := ctx.ShouldBindJSON(&group); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason(badresponse.BadRequestReason))
		return
	}
	if userInfo, exist := mo2utils.GetUserInfo(ctx); exist {
		if pass, err := accessControl.Ctrl.CanOr(abac.IQueryInfo{
			Subject:  accountStr,
			Action:   abac.ActionUpdate,
			Resource: accessControl.ResourceGroup,
			Context: abac.DefaultContext{accessControl.RuleAllowOwn: accessControl.AllowOwn{
				UserInfo: userInfo,
				ID:       group.ID,
				Resource: accessControl.ResourceGroup,
			}, accessControl.RuleAccessFilter: accessControl.AccessFilter{
				VisitorID: userInfo.ID,
				GroupID:   group.ID,
				RoleList:  [][]string{{accessControl.RoleAdmin}},
			}},
		}); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseError(err))
			return
		} else if !pass {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason(badresponse.NoAccessReason))
			return
		} else if pass {
			if mErr := database.UpsertGroup(group); mErr.IsError() {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseError(mErr))
				return
			} else {
				ctx.JSON(http.StatusCreated, group)
				return
			}
		}
	} else {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseReason(badresponse.UnauthorizeReason))
		return
	}
}

// DeleteGroup godoc
// @Summary delete group
// @Description delete by id in path
// @Tags group
// @Accept  json
// @Produce  json
// @Param id path primitive.ObjectID true "group id to delete"
// @Success 202
// @Failure 400 {object} badresponse.ResponseError
// @Failure 401 {object} badresponse.ResponseError
// @Router /api/group/{id} [delete]
func (c *Controller) DeleteGroup(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if id.IsZero() || err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason(badresponse.BadRequestReason))
		return
	}
	if userInfo, exist := mo2utils.GetUserInfo(ctx); exist {
		if pass, err := accessControl.Ctrl.CanOr(abac.IQueryInfo{
			Subject:  accountStr,
			Action:   abac.ActionDelete,
			Resource: accessControl.ResourceGroup,
			Context: abac.DefaultContext{accessControl.RuleAllowOwn: accessControl.AllowOwn{
				UserInfo: userInfo,
				ID:       id,
				Resource: accessControl.ResourceGroup,
			}, accessControl.RuleAccessFilter: accessControl.AccessFilter{
				VisitorID: userInfo.ID,
				GroupID:   id,
				RoleList:  [][]string{{accessControl.RoleAdmin}},
			}},
		}); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseError(err))
			return
		} else if !pass {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason(badresponse.NoAccessReason))
			return
		} else if pass {
			if mErr := database.DeleteGroupByID(id); mErr.IsError() {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseError(mErr))
				return
			} else {
				ctx.Status(http.StatusAccepted)
				return
			}
		}
	} else {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseReason(badresponse.UnauthorizeReason))
		return
	}
}

// FindGroup godoc
// @Summary get group
// @Description delete by id in path
// @Tags group
// @Accept  json
// @Produce  json
// @Param id path primitive.ObjectID true "group id to find"
// @Success 200 {object} model.Group
// @Success 202
// @Failure 400 {object} badresponse.ResponseError
// @Failure 401 {object} badresponse.ResponseError
// @Router /api/group/{id} [get]
func (c *Controller) FindGroup(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil || id.IsZero() {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseReason(badresponse.BadRequestReason))
		return
	}
	if group, mErr := database.FindGroup(id); mErr.IsError() {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, badresponse.SetResponseError(mErr))
		return
	} else {
		ctx.JSON(http.StatusOK, group)
		return
	}
}
