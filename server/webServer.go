package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "mo2/docs"
	"mo2/server/controller"
)

func RunServer() {

	r := gin.Default()
	r.GET("/sayHello", controller.SayHello)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run()
}
