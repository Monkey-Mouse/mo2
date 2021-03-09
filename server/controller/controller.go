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
	apiImgGenToken = "/img/:filename"
	apiLogs        = "/logs"
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
		blogs := api.Group("blogs", model.Anonymous, model.OrdinaryUser)
		{
			blogs.Get("query", c.QueryBlogs)

			user := blogs.Group("", model.OrdinaryUser)
			{
				user.Post("category", c.UpsertCategory)
				user.Get("category", c.FindAllCategories)
				user.Get("category/parent", c.FindSubCategories)
				user.Post("addBlogs2Categories", c.AddBlogs2Categories)
				user.Get("findCategoryByUserId", c.FindCategoryByUserId)
				user.Post("category/user/:userID", c.AddCategory2User)
				user.Get("category/user/:userID", c.FindCategoriesByUserId)
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
