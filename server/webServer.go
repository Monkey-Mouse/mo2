package server

import (
	_ "mo2/docs"
	"mo2/mo2utils"
	"mo2/server/controller"
	"mo2/server/middleware"
	"mo2/server/model"

	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupHandlers(c *controller.Controller) {
	api := middleware.H.Group("/api")
	{
		api.Get("/logs", c.Log)
		api.Get("/img/:filename", c.GenUploadToken, model.OrdinaryUser)
		blogs := api.Group("blogs")
		{
			blogs.Get("query", c.QueryBlogs)
			blogs.Post("addCategory", c.UpsertCategory)
			blogs.Get("findAllCategories", c.FindAllCategories)
			blogs.Post("addBlogs2Categories", c.AddBlogs2Categories)
			blogs.Get("findCategoryByUserId", c.FindCategoryByUserId)
			blogs.Post("addCategory2User", c.AddCategory2User)
			blogs.Get("findCategoriesByUserId", c.FindCategoriesByUserId)
			blogs.Post("addCategory2Category", c.AddCategory2Category)
		}
		accounts := api.Group("/accounts")
		{
			//accounts.GET(":id", c.ShowAccount)
			//accounts.POST("addUser",c.AddMo2User)
			accounts.Post("", c.AddAccount)
			accounts.Post("login", c.LoginAccount)
			accounts.Get("logout", c.LogoutAccount)
			accounts.Get("detail/:id", c.ShowAccount)
			accounts.Get("listBrief", c.ListAccountsInfo)

			/*accounts.GET("", c.ListAccounts)
			accounts.POST("", c.AddAccount)
			accounts.DELETE(":id", c.DeleteAccount)
			accounts.PATCH(":id", c.UpdateAccount)
			accounts.POST(":id/images", c.UploadAccountImage)*/
		}
	}
}

func RunServer() {

	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("dist", true)))
	r.Use(middleware.AuthMiddlware)
	r.GET("/sayHello", controller.SayHello)
	c := controller.NewController()
	setupHandlers(c)
	middleware.H.RegisterMapedHandlers(r)
	v1 := r.Group("/api")
	{
		blogs := v1.Group("/blogs")
		{
			blogs.POST("publish", c.UpsertBlog)

			find := blogs.Group("/find")
			{
				find.GET("own", c.FindBlogsByUser)
				find.GET("userId", c.FindBlogsByUserId)

				find.GET("id", c.FindBlogById)

			}

		}

	}
	auth := r.Group("/auth", mo2utils.BasicAuth())
	{
		auth.GET("home", func(ctx *gin.Context) {

			//TODO change the info generate way
			user, err := ctx.Cookie("jwtToken")
			if err != nil {
				ctx.JSON(http.StatusForbidden, "login first!")
			} else {
				ctx.JSON(http.StatusOK, gin.H{"home": user + " Welcome to your home"})

			}
		})
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.NoRoute(func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "dist/index.html")
	})
	// r.GET("/", func(c *gin.Context) {
	// 	http.ServeFile(c.Writer, c.Request, "dist/index.html")
	// })
	// r.Static("/static", "dist/static")
	r.Run(":5000")
}
