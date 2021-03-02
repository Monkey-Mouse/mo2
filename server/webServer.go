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
		uploads := api.Group("", model.OrdinaryUser)
		{
			uploads.Get("/img/:filename", c.GenUploadToken)
			uploads.Post("/file", c.Upload)

		}
		blogs := api.Group("blogs")
		{
			open := blogs.Group("", model.Anonymous, model.OrdinaryUser)
			{
				open.Get("query", c.QueryBlogs)

			}
			user := blogs.Group("", model.OrdinaryUser)
			{
				user.Post("addCategory", c.UpsertCategory)
				user.Get("findAllCategories", c.FindAllCategories)
				user.Post("addBlogs2Categories", c.AddBlogs2Categories)
				user.Get("findCategoryByUserId", c.FindCategoryByUserId)
				user.Post("addCategory2User", c.AddCategory2User)
				user.Get("findCategoriesByUserId", c.FindCategoriesByUserId)
				user.Post("addCategory2Category", c.AddCategory2Category)
				user.Post("publish", c.UpsertBlog)
				user.Delete(":id", c.DeleteBlog)
				user.Put(":id", c.RestoreBlog)
			}

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
			accounts.Delete("", c.DeleteAccount, model.OrdinaryUser)
			accounts.Put("", c.UpdateAccount, model.OrdinaryUser)
			accounts.Get("verify", c.VerifyEmail, model.Anonymous, model.OrdinaryUser)
			accounts.Post("role", c.AddAccountRole, model.GeneralAdmin, model.OrdinaryUser)
			accounts.Post("login", c.LoginAccount)
			accounts.Post("logout", c.LogoutAccount)
			accounts.Get("detail/:id", c.ShowAccount)
			accounts.Get("listBrief", c.ListAccountsInfo)
		}
	}
}

// RunServer start web server
func RunServer() {

	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("dist", true)))
	r.GET("/sayHello", controller.SayHello)
	c := controller.NewController()
	setupHandlers(c)
	middleware.H.RegisterMapedHandlers(r, func(ctx *gin.Context) (userInfo middleware.RoleHolder, err error) {
		str, err := ctx.Cookie("jwtToken")
		if err != nil {
			return
		}
		userInfo, err = mo2utils.ParseJwt(str)
		return
	}, mo2utils.UserInfoKey)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.NoRoute(func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "dist/index.html")
	})
	r.Run(":5001")
}
