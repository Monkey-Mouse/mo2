package server

import (
	_ "mo2/docs"
	"mo2/server/controller"

	"mo2/server/middleware"
	"net/http"

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
			//accounts.GET(":id", c.ShowAccount)
			//accounts.POST("addUser",c.AddMo2User)
			accounts.POST("", c.AddAccount)
			accounts.POST("login", c.LoginAccount)
			accounts.GET("logout", c.LogoutAccount)

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
	auth := r.Group("/auth", middleware.BasicAuth())
	{
		auth.GET("home", func(ctx *gin.Context) {
			user, err := ctx.Cookie("user")
			if err != nil {
				ctx.JSON(http.StatusForbidden, "login first!")
			} else {
				ctx.JSON(http.StatusOK, gin.H{"home": user + " Welcome to your home"})

			}
		})
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":5000")
}
