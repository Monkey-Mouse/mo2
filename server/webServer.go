package server

import (
	_ "mo2/docs"
	"mo2/server/controller"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RunServer() {

	r := gin.Default()
	r.GET("/sayHello", controller.SayHello)

	c := controller.NewController()
	v1 := r.Group("/api")
	{
		accounts := v1.Group("/accounts")
		{
			accounts.GET(":id", c.ShowAccount)
			//accounts.POST("addUser",c.AddMo2User)
			accounts.POST("", c.AddAccount)

			/*accounts.GET("", c.ListAccounts)
			accounts.POST("", c.AddAccount)
			accounts.DELETE(":id", c.DeleteAccount)
			accounts.PATCH(":id", c.UpdateAccount)
			accounts.POST(":id/images", c.UploadAccountImage)*/
		}
		blogs := v1.Group("/blogs")
		{
			blogs.POST("", c.PublishBlog)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run()
}
