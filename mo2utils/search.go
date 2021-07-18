package mo2utils

import (
	"github.com/Monkey-Mouse/mo2/server/model"
	"github.com/Monkey-Mouse/mo2/services/mo2search"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
)

const (
	blogIndex = "blog"
	userIndex = "user"
)

func init() {
	mo2search.CreateOrLoadIndex(blogIndex)
	mo2search.CreateOrLoadIndex(userIndex)
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
func QueryAccountPrefix(term string) search.DocumentMatchCollection {
	re := mo2search.Query(userIndex, bleve.NewPrefixQuery(term), 0, 10, nil)
	return re.Hits
}
func QueryBlogPrefix(term string) search.DocumentMatchCollection {
	re := mo2search.Query(blogIndex, bleve.NewPrefixQuery(term), 0, 10, nil)
	return re.Hits
}
