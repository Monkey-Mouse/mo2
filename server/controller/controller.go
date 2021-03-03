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

const imgGenToken = "/img/:filename"

// SetupHandlers set up
func SetupHandlers(c *Controller) {
	api := middleware.H.Group("/api")
	{
		api.Get("/logs", c.Log)
		api.Get(imgGenToken, c.GenUploadToken)
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
