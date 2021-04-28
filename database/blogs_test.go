package database

import (
	"github.com/Monkey-Mouse/mo2/server/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"testing"
)

func Test_upsertBlog(t *testing.T) {
	isDraft := true
	id := primitive.NewObjectID()

	uID := primitive.NewObjectID()
	defer DeleteBlogs(isDraft, id)

	type args struct {
		b       *model.Blog
		isDraft bool
	}
	tests := []struct {
		name         string
		args         args
		wantAuthorID primitive.ObjectID
	}{
		// first insert

		{"first insert", args{isDraft: isDraft, b: &model.Blog{ID: id, AuthorID: uID}}, uID},
		{"second update", args{isDraft: isDraft, b: &model.Blog{ID: id, AuthorID: primitive.NewObjectID()}}, uID},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//fmt.Println(tt.args.b)
			if mErr := upsertBlog(tt.args.b, tt.args.isDraft); mErr.IsError() {
				t.Errorf("upsertBlog() got Error= %v", mErr.ErrorTip)

			}

			if blog := FindBlogById(id, isDraft); !reflect.DeepEqual(blog.AuthorID, tt.wantAuthorID) {
				t.Errorf("upsertBlog() = %v, want %v", blog.AuthorID, tt.wantAuthorID)
			}
		})
	}
}
