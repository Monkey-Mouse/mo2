package controller

import (
	"github.com/Monkey-Mouse/mo2/mo2utils/adapter"
	"github.com/Monkey-Mouse/mo2/server/middleware"
	"github.com/Monkey-Mouse/mo2/server/model"
)

// Controller example
type Controller struct {
}

// NewController example
func NewController() *Controller {
	return &Controller{}
}

// Message example
type Message struct {
	Message string `json:"message" example:"message"`
}

const (
	apiImgGenToken    = "/img/:filename"
	apiLogs           = "/logs"
	typeKey           = "type"
	typeCategory      = "category"
	typeCategories    = "categories"
	typeUser          = "user"
	typeUserMain      = "userMain"
	typeBlog          = "blog"
	typeUsers         = "users"
	typeBlogs         = "blogs"
	typeSubCategories = "sub"
)

// SetupHandlers set up
func SetupHandlers(c *Controller) {
	api := middleware.H.Group("/api")
	{
		api.GetWithRL("accounts/verify", c.VerifyEmail, 3)
		api.GetWithRL("/logs", c.Log, 10)
		oau := api.Group("oauth")
		{
			oau.GetWithRL("github", c.GithubOauth, 3)
		}
		noti := api.Group("notification", model.OrdinaryUser)
		{
			noti.GET("num", c.GetNotificationNum)
			noti.GET("", c.GetNotifications)
		}
		admin := api.Group("/admin")
		{
			admin.POST("index", c.IndexAll)
		}
		uploads := api.Group("", model.OrdinaryUser)
		{
			uploads.GET("/img/:filename", adapter.ReAdapter(c.GenUploadToken))
			uploads.POST("/file", c.Upload)
		}
		relation := api.Group("relation", model.OrdinaryUser, model.Anonymous)
		{

			relation.POST("categories/:type", c.RelateCategories2Entity, model.GeneralAdmin)
			relation.POST("category/:type", c.RelateCategory2Entity, model.GeneralAdmin)
			relation.GET("category/:type/:ID", c.FindCategoriesByType)
			relation.GET("blogs/:type/:ID", c.FindBlogsByType)
		}
		like := api.Group("like", model.OrdinaryUser, model.Anonymous)
		{
			like.PostWithRL(":type", adapter.ReAdapterWithUinfo(c.Like), 5, model.OrdinaryUser)
			like.GetWithRL("num/:type/:id", adapter.ReAdapter(c.LikeNum), 5)
			like.GetWithRL("ext/:type/:id", adapter.ReAdapterWithUinfo(c.Liked), 5, model.OrdinaryUser)
		}
		directories := api.Group("directories", model.OrdinaryUser, model.Anonymous)
		{
			user := directories.Group("", model.OrdinaryUser)
			{
				user.DELETE("category", c.DeleteCategory)
			}
			directories.GET(":collection", c.ListDirectoriesInfo)
		}

		blogs := api.Group("blogs", model.Anonymous, model.OrdinaryUser)
		{
			blogs.GET("query", c.QueryBlogs)

			user := blogs.Group("", model.OrdinaryUser)
			{
				user.POST("score", adapter.ReAdapterWithUinfo(c.Score))
				user.POST("isscored", adapter.ReAdapterWithUinfo(c.IsScored))
				user.POST("category", c.UpsertCategory)
				user.GET("category", c.FindAllCategories)

				user.POST("publish", c.UpsertBlog)
				user.DELETE(":id", c.DeleteBlog)
				user.PUT(":operation/:id", c.ProcessBlog)
				user.POST("doctype", adapter.ReAdapterWithUinfo(c.SetDocType))
			}

			find := blogs.Group("/find")
			{
				find.GET("own", c.FindBlogsByUser, model.OrdinaryUser)
				find.GET("", c.FindBlogsByID)
				find.GET("id", c.FindBlogById)
			}
		}
		accounts := api.Group("/accounts", model.Anonymous, model.OrdinaryUser)
		{
			accounts.POST("", c.AddAccount)
			accounts.DELETE("", c.DeleteAccount, model.OrdinaryUser)
			accounts.PUT("", c.UpdateAccount, model.OrdinaryUser)
			accounts.POST("role", c.AddAccountRole, model.GeneralAdmin, model.OrdinaryUser)
			accounts.POST("login", c.LoginAccount)
			accounts.POST("logout", c.LogoutAccount)
			accounts.GET("detail/:id", c.ShowAccount)
			accounts.GET("listBrief", c.ListAccountsInfo)
			accounts.GET("", adapter.ReAdapterWithUinfo(c.SearchAccount))
		}
		comment := api.Group("/comment", model.Anonymous, model.OrdinaryUser)
		{
			comment.GET(":id", c.GetComment)
			comment.POST("", c.PostComment, model.OrdinaryUser)
			comment.POST(":id", c.PostSubComment, model.OrdinaryUser)
		}
		api.GET("commentcount/:id", c.GetCommentNum)
		group := api.Group("/group")
		{
			user := group.Group("", model.OrdinaryUser)
			{
				user.PUT("", c.UpdateGroup)
				user.POST("", c.InsertGroup)
				user.DELETE(":id", c.DeleteGroup)
			}
			group.GET(":id", c.FindGroup)

		}
		proj := api.Group("project")
		{
			proj.POST("", adapter.ReAdapterWithUinfo(c.UpsertProject), model.OrdinaryUser)
			proj.GET("", adapter.ReAdapterWithUinfo(c.ListProject))
			proj.DELETE(":id", adapter.ReAdapterWithUinfo(c.DeleteProject), model.OrdinaryUser)
			proj.GET(":id", adapter.ReAdapterWithUinfo(c.GetProject))
		}
	}
}
