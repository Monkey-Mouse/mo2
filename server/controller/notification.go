package controller

import (
	"fmt"
	"mo2/mo2utils"
	"mo2/server/controller/badresponse"
	"mo2/services/loghelper"
	"net/http"

	"github.com/Monkey-Mouse/mo2log/logmodel"
	"github.com/Monkey-Mouse/mo2log/service/logservice"
	"github.com/gin-gonic/gin"
)

// GetNotificationNum godoc
// @Summary get notification num
// @Description get notification num
// @Tags notification
// @Produce  json
// @Success 200
// @Router /api/notification/num [get]
func (c *Controller) GetNotificationNum(ctx *gin.Context) {
	info, _ := mo2utils.GetUserInfo(ctx)
	num, err := notifyLogClient.Client.GetUserNewMsgNum(ctx, &logservice.UserID{UserId: info.ID[:]})
	if err != nil {
		fmt.Println(err)
	}

	ctx.JSON(200, gin.H{"num": num.Num})
}

// GetNotifications godoc
// @Summary get notifications
// @Description get notifications
// @Tags notification
// @Param page query int false "int 0" 0
// @Param pagesize query int false "int 5" 5
// @Produce  json
// @Success 200
// @Router /api/notification [get]
func (c *Controller) GetNotifications(ctx *gin.Context) {
	info, _ := mo2utils.GetUserInfo(ctx)
	page, pagesize, err := mo2utils.ParsePagination(ctx)
	if err != nil {
		badresponse.SetErrResponse(ctx, http.StatusUnprocessableEntity, "格式错误")
		return
	}
	msgs, err := notifyLogClient.Client.GetUserMsgs(ctx, &logservice.ListRequest{
		UserId:   info.ID[:],
		Page:     int32(page),
		Pagesize: int32(pagesize),
	})
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(500, badresponse.SetResponseReason("Internal error"))
		return
	}
	logs := make([]*logmodel.LogModel, 0, len(msgs.Logs))
	for _, v := range msgs.Logs {
		logs = append(logs, loghelper.ProtoToLog(v))
	}
	ctx.JSON(200, logs)
}
