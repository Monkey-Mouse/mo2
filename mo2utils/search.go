package mo2utils

import (
	"github.com/Monkey-Mouse/mo2/server/model"
	"github.com/Monkey-Mouse/mo2/services/mo2search"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
)

const (
	blogIndex    = "blog"
	userIndex    = "user"
	projectIndex = "project"
)

func init() {
	mo2search.CreateOrLoadIndex(blogIndex)
	mo2search.CreateOrLoadIndex(userIndex)
	mo2search.CreateOrLoadIndex(projectIndex)
}

// IndexBlog index the blog
func IndexBlog(blog *model.Blog) {
	mo2search.Index(blogIndex, blog.ID.Hex(), blog)
}

func IndexAccount(account *model.Account) {
	account.Infos = nil
	account.HashedPwd = ""
	mo2search.Index(userIndex, account.ID.Hex(), account)
}
func IndexProject(project *model.Project) {
	mo2search.Index(projectIndex, project.ID.Hex(), project)
}

// IndexBlogs index multiple blogs
func IndexBlogs(blog []model.Blog) {
	for _, v := range blog {
		IndexBlog(&v)
	}
}

func IndexAccounts(accounts []model.Account) {
	for _, v := range accounts {
		IndexAccount(&v)
	}
}
func QueryUser(search string, page int, pagesize int) search.DocumentMatchCollection {
	queryT := bleve.NewMatchQuery(search)
	queryT.SetField("userName")
	queryT.SetBoost(5)
	queryD := bleve.NewMatchQuery(search)
	queryD.SetField("email")
	queryD.SetBoost(5)

	query := bleve.NewDisjunctionQuery(queryT, queryD)
	re := mo2search.Query(userIndex, query, page, pagesize, []string{"*"})
	return re.Hits
}
func QueryProject(search string, page int, pagesize int) search.DocumentMatchCollection {
	queryN := bleve.NewMatchQuery(search)
	queryN.SetField("name")
	queryN.SetBoost(5)
	queryD := bleve.NewMatchQuery(search)
	queryD.SetField("description")
	queryD.SetBoost(1)
	queryT := bleve.NewMatchQuery(search)
	queryT.SetField("tags")
	queryT.SetBoost(5)

	query := bleve.NewDisjunctionQuery(queryT, queryD, queryN)
	re := mo2search.Query(projectIndex, query, page, pagesize, []string{"*"})
	return re.Hits
}
func QueryBlog(search string, page int, pagesize int) search.DocumentMatchCollection {
	queryT := bleve.NewMatchQuery(search)
	queryT.SetField("title")
	queryT.SetBoost(5)
	queryD := bleve.NewMatchQuery(search)
	queryD.SetField("description")
	queryD.SetBoost(1)

	query := bleve.NewDisjunctionQuery(queryT, queryD)
	re := mo2search.Query(blogIndex, query, page, pagesize, []string{"*"})
	return re.Hits
}
func DeleteBlogIndex(id string) {
	mo2search.Delete(blogIndex, id)
}
func DeleteAccountIndex(id string) {
	mo2search.Delete(userIndex, id)
}
func DeleteProjectIndex(id string) {
	mo2search.Delete(projectIndex, id)
}
func QueryAccountPrefix(term string) search.DocumentMatchCollection {
	re := mo2search.Query(userIndex, bleve.NewPrefixQuery(term), 0, 10, nil)
	return re.Hits
}
func QueryBlogPrefix(term string) search.DocumentMatchCollection {
	re := mo2search.Query(blogIndex, bleve.NewPrefixQuery(term), 0, 10, nil)
	return re.Hits
}
