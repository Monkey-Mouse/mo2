package mo2utils

import (
	"github.com/Monkey-Mouse/mo2/server/model"
	"github.com/Monkey-Mouse/mo2/services/mo2search"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
)

const (
	blogIndex = "blog"
)

func init() {
	mo2search.CreateOrLoadIndex(blogIndex)
}

type blogI struct {
	Title       string
	Description string
	KeyWords    []string
}

// IndexBlog index the blog
func IndexBlog(blog *model.Blog) {
	mo2search.Index(blogIndex, blog.ID.Hex(), blog)
}

// IndexBlogs index multiple blogs
func IndexBlogs(blog []model.Blog) {
	for _, v := range blog {
		IndexBlog(&v)
	}
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
func QueryPrefix(term string) search.DocumentMatchCollection {
	re := mo2search.Query(blogIndex, bleve.NewPrefixQuery(term), 0, 10, nil)
	return re.Hits
}
