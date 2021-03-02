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
		api.Get("/logs", c.Log, model.Anonymous, model.OrdinaryUser)
		api.Get("/img/:filename", c.GenUploadToken, model.OrdinaryUser)
		blogs := api.Group("blogs", model.OrdinaryUser)
		{
			blogs.Get("query", c.QueryBlogs, model.Anonymous)
			blogs.Post("addCategory", c.UpsertCategory)
			blogs.Get("findAllCategories", c.FindAllCategories)
			blogs.Post("addBlogs2Categories", c.AddBlogs2Categories)
			blogs.Get("findCategoryByUserId", c.FindCategoryByUserId)
			blogs.Post("addCategory2User", c.AddCategory2User)
			blogs.Get("findCategoriesByUserId", c.FindCategoriesByUserId)
			blogs.Post("addCategory2Category", c.AddCategory2Category)
			blogs.Post("publish", c.UpsertBlog)
			blogs.Delete(":id", c.DeleteBlog)
			blogs.Put(":id", c.RestoreBlog)
			find := blogs.Group("/find", model.OrdinaryUser)
			{
				find.Get("own", c.FindBlogsByUser)
				find.Get("userId", c.FindBlogsByUserId, model.Anonymous)
				find.Get("id", c.FindBlogById, model.Anonymous)
			}
		}
		accounts := api.Group("/accounts")
		{
			accounts.Post("", c.AddAccount, model.Anonymous) //todo: whether add ordinaryUser
			accounts.Delete("", c.DeleteAccount, model.OrdinaryUser)
			accounts.Put("", c.UpdateAccount, model.OrdinaryUser)
			accounts.Get("verify", c.VerifyEmail, model.Anonymous, model.OrdinaryUser)
			accounts.Post("role", c.AddAccountRole, model.GeneralAdmin, model.OrdinaryUser) //todo update after have generalAdmin
			accounts.Post("login", c.LoginAccount, model.Anonymous)                         //todo: whether add ordinaryUser
			accounts.Post("logout", c.LogoutAccount, model.OrdinaryUser)
			accounts.Get("detail/:id", c.ShowAccount, model.Anonymous, model.OrdinaryUser)
			accounts.Get("listBrief", c.ListAccountsInfo, model.Anonymous, model.OrdinaryUser)
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
