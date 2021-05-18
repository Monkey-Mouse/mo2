package controller

import (
	"errors"
	"fmt"
	"time"

	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/Monkey-Mouse/mo2/mo2img"
	"github.com/Monkey-Mouse/mo2/mo2utils"
	"github.com/Monkey-Mouse/mo2/services/loghelper"
	"github.com/Monkey-Mouse/mo2log/service/logservice"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

var imgLogClient = &loghelper.LogClient{}

func init() {
	imgLogClient.Init("IMG_LOG")
}

// GenUploadToken generate img upload token
// GenUploadToken godoc
// @Summary Gen img token
// @Description add by json
// @Tags img
// @Produce  json
// @Param filename path string true "file name"
// @Success 200 {object} dto.ImgUploadToken
// @Router /api/img/{filename} [get]
func (c *Controller) GenUploadToken(ctx *gin.Context) (status int, json interface{}, err error) {
	user, _ := mo2utils.GetUserInfo(ctx)
	n, err := imgLogClient.Client.Count(ctx, &logservice.UserID{UserId: user.ID[:]})
	if err != nil {
		return 500, nil, errors.New("Internal micro service error")
	}
	if n.Num >= 50 {
		return 429, nil, errors.New("上传次数于24h内达到上限50，暂时无法上传文件")
	}
	fileKey := ctx.Param("filename")
	savekey := fmt.Sprintf("%s/%v%v", user.ID.Hex(), time.Now().UnixNano(), fileKey)
	token := mo2img.GenerateUploadToken(savekey)
	defer imgLogClient.LogInfo(loghelper.Log{
		Operator:             user.ID,
		Operation:            1,
		OperationTarget:      primitive.NilObjectID,
		OperationTargetOwner: primitive.NilObjectID,
		ExtraMessage:         "",
	})
	return 200, dto.ImgUploadToken{Token: token, FileKey: savekey}, nil
}
