package mo2utils

import (
	"mo2/server/model"
	"mo2/services/mo2search"
)

const (
	blogIndex = "blog"
)

func init() {
	mo2search.CreateOrLoadIndex(blogIndex)
}

func indexBlog(blog model.Blog) {
	mo2search.Indexes[blogIndex].Index(blog.ID.Hex(), blog)
}
func indexBlogs(blog []model.Blog) {
	batch := mo2search.Indexes[blogIndex].NewBatch()
	for _, v := range blog {
		batch.Index(v.ID.Hex(), v)
	}
	mo2search.Indexes[blogIndex].Batch(batch)
}
