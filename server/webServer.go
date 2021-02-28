package server

import (
	_ "mo2/docs"
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
			blogs.Post("addCategory", c.UpsertCategory, model.OrdinaryUser)
			blogs.Get("findAllCategories", c.FindAllCategories, model.OrdinaryUser)
			blogs.Post("addBlogs2Categories", c.AddBlogs2Categories, model.OrdinaryUser)
			blogs.Get("findCategoryByUserId", c.FindCategoryByUserId, model.OrdinaryUser)
			blogs.Post("addCategory2User", c.AddCategory2User, model.OrdinaryUser)
			blogs.Get("findCategoriesByUserId", c.FindCategoriesByUserId, model.OrdinaryUser)
			blogs.Post("addCategory2Category", c.AddCategory2Category, model.OrdinaryUser)
			blogs.Post("publish", c.UpsertBlog, model.OrdinaryUser)
			blogs.Delete(":id", c.DeleteBlog, model.OrdinaryUser)
			blogs.Put(":id", c.RestoreBlog, model.OrdinaryUser)
			find := blogs.Group("/find")
			{
				find.Get("own", c.FindBlogsByUser, model.OrdinaryUser)
				find.Get("userId", c.FindBlogsByUserId)
				find.Get("id", c.FindBlogById)
			}
		}
		accounts := api.Group("/accounts")
		{
			accounts.Post("", c.AddAccount)
			accounts.Delete("", c.DeleteAccount)
			accounts.Put("", c.UpdateAccount)
			accounts.Get("verify", c.VerifyEmail)
			accounts.Post("role", c.AddAccountRole)
			accounts.Post("login", c.LoginAccount)
			accounts.Post("logout", c.LogoutAccount)
			accounts.Get("detail/:id", c.ShowAccount)
			accounts.Get("listBrief", c.ListAccountsInfo)
		}
		auth := api.Group("/auth")
		{
			auth.Get("home", func(ctx *gin.Context) {
				//TODO change the info generate way
				user, err := ctx.Cookie("jwtToken")
				if err != nil {
					ctx.JSON(http.StatusForbidden, "login first!")
				} else {
					ctx.JSON(http.StatusOK, gin.H{"home": user + " Welcome to your home"})
				}
			})
		}
	}
}

func RunServer() {

	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("dist", true)))
	r.Use(middleware.AuthMiddleware)
	r.GET("/sayHello", controller.SayHello)
	c := controller.NewController()
	setupHandlers(c)
	middleware.H.RegisterMapedHandlers(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.NoRoute(func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "dist/index.html")
	})
	// r.GET("/", func(c *gin.Context) {
	// 	http.ServeFile(c.Writer, c.Request, "dist/index.html")
	// })
	// r.Static("/static", "dist/static")
	r.Run(":5001")
}
