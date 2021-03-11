package database

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"mo2/dto"
	"mo2/server/model"
	"time"
)

// InsertBlogs4Test 插入n个blog/draft，并返回它们的id列表
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

// InsertCategories4Test 插入n个category，并返回它们的id列表
func InsertCategories4Test(num int) (ids []primitive.ObjectID) {
	ids = make([]primitive.ObjectID, num)
	for i := 0; i < num; i++ {
		ids[i] = primitive.NewObjectID()
		if mErr := UpsertCategory(&model.Directory{ID: ids[i]}); mErr.IsError() {
			log.Println(mErr)
		}
	}
	return
}

func ExampleDeleteCategoryCompletely() {
	testNum := 100
	catNum := 10
	isDraft := true
	catIDs := InsertCategories4Test(catNum)
	ids := InsertBlogs4Test(isDraft, testNum)

	defer func() {
		if mErr := deleteBlogs(isDraft, ids...); mErr.IsError() {
			log.Println(mErr)
		}
		DeleteCategory(catIDs...)
	}()

	RelateCategories2Blogs(dto.RelateEntitySet2EntitySet{
		RelatedIDs:  catIDs,
		RelateToIDs: ids,
	})

	if bs, mErr := FindBlogsByCategoryId(catIDs[0], isDraft); mErr.IsError() {
		log.Println(mErr)
	} else {
		fmt.Println("before delete: ", len(bs))
	}
	t1 := time.Now()
	if mErr := DeleteCategoryCompletely(catIDs...); mErr.IsError() {
		fmt.Println(mErr)
	}
	log.Println(time.Since(t1))

	if bs, mErr := FindBlogsByCategoryId(catIDs[0], isDraft); mErr.IsError() {
		log.Println(mErr)
	} else {
		fmt.Println("after delete: ", len(bs))
	}
	fmt.Println(FindCategoryById(catIDs[0]))

	// Output:
	// before delete:  100
	// after delete:  0
	// {ObjectID("000000000000000000000000") ObjectID("000000000000000000000000")  { } []}

}

func ExampleRelateSubCategories2Category() {
	num := 10
	ids := InsertCategories4Test(num)
	parentID := InsertCategories4Test(1)[0]
	defer func() {
		DeleteCategory(ids...)
		DeleteCategory(parentID)
	}()
	RelateSubCategories2Category(dto.RelateEntitySet2Entity{RelatedIDs: ids, RelateToID: parentID})
	if res, mErr := FindSubCategories(parentID); mErr.IsError() {
		log.Fatal(mErr)
	} else {
		fmt.Println("after relate: ", len(res))
	}
	//Output:
	//after relate:  10
}

func ExampleUpdateSubCategories() {
	num := 10
	ids := InsertCategories4Test(num)
	parentID := InsertCategories4Test(1)[0]
	newParentID := primitive.NewObjectID()
	defer func() {
		DeleteCategory(ids...)
		DeleteCategory(parentID)
	}()

	RelateSubCategories2Category(dto.RelateEntitySet2Entity{RelatedIDs: ids, RelateToID: parentID})
	res, _ := FindSubCategories(parentID)
	fmt.Println("relate to old: ", len(res))
	res, _ = FindSubCategories(newParentID)
	fmt.Println("relate to new: ", len(res))
	mErr := UpdateSubCategories(parentID, newParentID)
	fmt.Println(mErr)
	res, _ = FindSubCategories(parentID)
	fmt.Println("relate to old: ", len(res))
	res, _ = FindSubCategories(newParentID)
	fmt.Println("relate to new: ", len(res))

	//Output:
	//relate to old:  10
	//relate to new:  0
	//200: update [10] subCategories
	//relate to old:  0
	//relate to new:  10
}
