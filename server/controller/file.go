package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	_ "mime/multipart" //in godoc comment
	"net/http"
)

// Upload godoc
// @Summary simple test
// @Description say something
// @Accept multipart/form-data
// @Produce  json
// @Param form body string true "file"
// @Success 200 {string} json
// @Router /file [post]
func (c *Controller) Upload(ctx *gin.Context) {
	// Multipart form
	form, _ := ctx.MultipartForm()
	files := form.File["upload[]"]

	for _, file := range files {
		log.Println(file.Filename)
		//importService.Transform(file)
		// Upload the file to specific dst.
		// ctx.SaveUploadedFile(file, "\\foo")
	}
	ctx.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}
