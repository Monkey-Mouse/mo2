package mo2utils

import (
	"mo2/server/model"
	"mo2/services/mo2search"

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
}

// IndexBlog index the blog
func IndexBlog(blog *model.Blog) {
	mo2search.Indexes[blogIndex].Index(blog.ID.Hex(), blogI{blog.Title, blog.Content})
}

// IndexBlogs index multiple blogs
func IndexBlogs(blog []model.Blog) {
	batch := mo2search.Indexes[blogIndex].NewBatch()
	for _, v := range blog {
		batch.Index(v.ID.Hex(), blogI{v.Title, v.Content})
	}
	mo2search.Indexes[blogIndex].Batch(batch)
}

func QueryBlogs(search string) search.DocumentMatchCollection {
	re := mo2search.Query(blogIndex, bleve.NewQueryStringQuery(search))
	return re.Hits
}
