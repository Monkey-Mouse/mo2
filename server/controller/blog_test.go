package controller

import (
	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/server/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"testing"
)

func findBlog(t *testing.T, isDraft string, id string) *http.Request {
	return get(t, "/api/blogs/find/id", map[string]string{"draft": isDraft, "id": id})
}
func TestController_FindBlogByID(t *testing.T) {
	bID := primitive.NewObjectID()
	id := primitive.NewObjectID()
	req := findBlog(t, "true", bID.Hex())
	req1 := findBlog(t, "true", bID.Hex())
	req2 := findBlog(t, "true", bID.Hex())
	req3 := findBlog(t, "true", bID.Hex())
	req4 := findBlog(t, "false", bID.Hex())
	req5 := findBlog(t, "false", bID.Hex())
	req6 := findBlog(t, "false", bID.Hex())
	req7 := findBlog(t, "false", bID.Hex())

	addCookie(req1)
	addCookieWithID(req2, primitive.NewObjectID())
	addCookieWithID(req3, id)
	addCookie(req5)
	addCookieWithID(req6, primitive.NewObjectID())
	addCookieWithID(req7, id)

	database.UpsertBlog(&model.Blog{
		ID:       bID,
		AuthorID: id,
	}, false)

	database.UpsertBlog(&model.Blog{
		ID:       bID,
		AuthorID: id,
	}, true)
	testHTTP(t,
		tests{name: "Test find draft without login", req: req, wantCode: 403}, // processed in middleware
		tests{name: "Test find draft without author id", req: req1, wantCode: 204},
		tests{name: "Test find draft with wrong author id", req: req2, wantCode: 204},
		tests{name: "Test find draft with right author id", req: req3, wantCode: 200},
		tests{name: "Test find blog without login", req: req4, wantCode: 403}, // processed in middleware
		tests{name: "Test find blog without author id", req: req5, wantCode: 200},
		tests{name: "Test find blog with wrong author id", req: req6, wantCode: 200},
		tests{name: "Test find blog with right author id", req: req7, wantCode: 200},
	)
	database.DeleteBlogs(true, bID)
	database.DeleteBlogs(false, bID)
}
