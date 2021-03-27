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
	mo2search.Indexes[blogIndex].Index(blog.ID.Hex(), blogI{blog.Title, blog.Description, blog.KeyWords})
}

// IndexBlogs index multiple blogs
func IndexBlogs(blog []model.Blog) {
	batch := mo2search.Indexes[blogIndex].NewBatch()
	for _, v := range blog {
		batch.Index(v.ID.Hex(), blogI{v.Title, v.Description, v.KeyWords})
	}
	mo2search.Indexes[blogIndex].Batch(batch)
}

func QueryBlog(search string) search.DocumentMatchCollection {
	queryT := bleve.NewMatchQuery(search)
	queryT.SetField("Title")
	queryT.SetBoost(5)
	queryD := bleve.NewMatchQuery(search)
	queryD.SetField("Description")
	queryD.SetBoost(1)
	query := bleve.NewDisjunctionQuery(queryT, queryD)
	re := mo2search.Query(blogIndex, query)
	return re.Hits
}
func DeleteBlogIndex(id string) {
	mo2search.Indexes[blogIndex].Delete(id)
}
func QueryPrefix(term string) search.DocumentMatchCollection {
	re := mo2search.Query(blogIndex, bleve.NewPrefixQuery(term))
	return re.Hits
}
