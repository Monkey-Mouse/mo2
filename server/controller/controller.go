package controller

import (
	"mo2/server/middleware"
	"mo2/server/model"
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
		api.GetWithRL("/logs", c.Log, 10)
		uploads := api.Group("", model.OrdinaryUser)
		{
			uploads.Get("/img/:filename", c.GenUploadToken)
			uploads.Post("/file", c.Upload)
		}
		relation := api.Group("relation", model.OrdinaryUser)
		{

			relation.Post("categories/:type", c.RelateCategories2Entity, model.GeneralAdmin)
			relation.Post("category/:type", c.RelateCategory2Entity, model.GeneralAdmin)
			relation.Get("category/:type/:ID", c.FindCategoriesByType)
			relation.Get("blogs/:type/:ID", c.FindBlogsByType)
		}
		blogs := api.Group("blogs", model.Anonymous, model.OrdinaryUser)
		{
			blogs.Get("query", c.QueryBlogs)

			user := blogs.Group("", model.OrdinaryUser)
			{
				user.Post("category", c.UpsertCategory)
				user.Get("category", c.FindAllCategories)

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
		accounts := api.Group("/accounts", model.Anonymous, model.OrdinaryUser)
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
