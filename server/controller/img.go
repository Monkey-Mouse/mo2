package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/Monkey-Mouse/mo2/mo2img"
	"github.com/Monkey-Mouse/mo2/mo2utils"

	"github.com/gin-gonic/gin"
)

// GenUploadToken generate img upload token
// GenUploadToken godoc
// @Summary Gen img token
// @Description add by json
// @Tags img
// @Produce  json
// @Param filename path string true "file name"
// @Success 200 {object} dto.ImgUploadToken
// @Router /api/img/{filename} [get]
func (c *Controller) GenUploadToken(ctx *gin.Context) {
	user, _ := mo2utils.GetUserInfo(ctx)
	fileKey := ctx.Param("filename")
	savekey := fmt.Sprintf("%s/%v%v", user.ID.Hex(), time.Now().UnixNano(), fileKey)
	token := mo2img.GenerateUploadToken(savekey)
	ctx.JSON(http.StatusOK, dto.ImgUploadToken{Token: token, FileKey: savekey})

}
