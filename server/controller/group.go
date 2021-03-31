package controller

import (
	"github.com/Monkey-Mouse/go-abac/abac"
	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/mo2utils"
	"github.com/Monkey-Mouse/mo2/server/controller/badresponse"
	"github.com/Monkey-Mouse/mo2/server/model"
	"github.com/Monkey-Mouse/mo2/services/accessControl"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	accountStr = "account"
	groupStr   = "group"
)

// UpsertGroup godoc
// @Summary create group or update info
// @Description add by json
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param draft query bool false "bool true" true
// @Param account body model.Group true "Add blog"
// @Success 201 {object} model.Group
// @Success 204
// @Failure 400 {object} badresponse.ResponseError
// @Failure 401 {object} badresponse.ResponseError
// @Router /api/group [put]
func (c *Controller) UpsertGroup(ctx *gin.Context) {
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
				RoleList:  []string{"admin"},
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
			}
		}
	} else {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, badresponse.SetResponseReason(badresponse.UnauthorizeReason))
		return
	}
}
