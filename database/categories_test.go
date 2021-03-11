package database

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"mo2/dto"
	"mo2/server/model"
	"time"
)

//插入n个blog/draft，并返回它们的id列表
func InsertBlogs4Test(isDraft bool, num int) (ids []primitive.ObjectID) {
	ids = make([]primitive.ObjectID, num)
	for i := 0; i < num; i++ {
		ids[i] = primitive.NewObjectID()
		if mErr := insertBlog(&model.Blog{ID: ids[i]}, isDraft); mErr.IsError() {
			log.Println(mErr)
		}
	}
	return
}

func ExampleDeleteCategory() {
	testNum := 100
	isDraft := true
	catID := primitive.NewObjectID()
	UpsertCategory(&model.Directory{ID: catID})
	ids := InsertBlogs4Test(isDraft, testNum)

	defer func() {
		if mErr := deleteBlogs(isDraft, ids...); mErr.IsError() {
			log.Println(mErr)
		}
	}()

	RelateCategories2Blogs(dto.RelateEntitySet2EntitySet{
		RelatedIDs:  []primitive.ObjectID{catID},
		RelateToIDs: ids,
	})
	bs, mErr := FindBlogsByCategoryId(catID, isDraft)
	if mErr.IsError() {
		fmt.Println(mErr)
	}
	fmt.Println("before delete: ", len(bs))
	t1 := time.Now()
	if mErr := DeleteCategory(catID); mErr.IsError() {
		fmt.Println(mErr)
	}
	log.Println(time.Since(t1))
	bs, mErr = FindBlogsByCategoryId(catID, isDraft)
	if mErr.IsError() {
		fmt.Println(mErr)

	}
	fmt.Println("after delete: ", len(bs))
	fmt.Println(FindCategoryById(catID))

	// Output:
	// before delete:  100
	// after delete:  0
	// {ObjectID("000000000000000000000000") ObjectID("000000000000000000000000")  { } []}

}
