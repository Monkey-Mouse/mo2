package database

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"mo2/dto"
	"mo2/mo2utils/mo2errors"
	"mo2/server/model"
	"reflect"
	"testing"
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

func TestUpsertCategory(t *testing.T) {
	fooID := primitive.NewObjectID()
	parentID := primitive.NewObjectID()
	defer DeleteCategory(fooID)
	foo := model.Directory{
		ID:       fooID,
		ParentID: primitive.NewObjectID(),
		Name:     "FOO",
		Info: model.DirectoryInfo{
			Description: "TEST",
			Cover:       "HTTP://FOO.COM",
		},
		OwnerIDs: []primitive.ObjectID{primitive.NewObjectID()},
	}
	tests := []struct {
		name          string
		args          model.Directory
		wantFindCat   model.Directory
		wantErrorType int
	}{
		{"foo", foo, foo, mo2errors.Mo2NoError},
		{"addParentID", model.Directory{
			ID:       fooID,
			Name:     foo.Name,
			Info:     foo.Info,
			ParentID: parentID,
			OwnerIDs: foo.OwnerIDs,
		}, model.Directory{
			ID:       fooID,
			ParentID: parentID,
			Name:     foo.Name,
			Info:     foo.Info,
			OwnerIDs: foo.OwnerIDs,
		}, mo2errors.Mo2NoError},
		{"addOwnerID", model.Directory{
			ID:       fooID,
			ParentID: parentID,
			Name:     foo.Name,
			Info:     foo.Info,
			OwnerIDs: []primitive.ObjectID{parentID},
		}, model.Directory{
			ID:       fooID,
			ParentID: parentID,
			Name:     foo.Name,
			Info:     foo.Info,
			OwnerIDs: []primitive.ObjectID{parentID},
		}, mo2errors.Mo2NoError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer DeleteCategory(test.args.ID)
			if mErr := UpsertCategory(&test.args); mErr.ErrorCode != test.wantErrorType {
				t.Errorf("UpsertCategory() gotError = %v, want %v", mErr.ErrorCode, test.args)
			}
			findCat := FindCategoryById(test.args.ID)

			if !reflect.DeepEqual(findCat, test.wantFindCat) {
				t.Errorf("UpsertCategory() gotFindCat = %v, want %v", findCat, test.wantFindCat)
			}
		})
	}
}

func TestFindOrCreateRoot4User(t *testing.T) {
	userID := primitive.NewObjectID()
	noRootUserID := primitive.NewObjectID()
	root := model.Directory{
		ID:       primitive.NewObjectID(),
		ParentID: userID,
		Name:     "",
		Info:     model.DirectoryInfo{},
		OwnerIDs: nil,
	}
	defer DeleteCategory(root.ID)
	UpsertCategory(&root)

	tests := []struct {
		name   string
		userID primitive.ObjectID
	}{
		{"hasRoot", userID},
		{"notHasRoot", noRootUserID},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if cat, mErr := FindOrCreateRoot4User(test.userID); mErr.IsError() {
				t.Error(mErr)
			} else {
				defer DeleteCategory(cat.ID)
				if cat.ParentID != test.userID {
					t.Errorf("FindOrCreateRoot4User() want parentID %v get %v \n", test.userID, cat.ParentID)
				}
			}
		})
	}
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
	//200: update 10 subCategories
	//relate to old:  0
	//relate to new:  10
}

func ExampleRelateCategories2User() {
	ownCatNum := 100
	ownCatIDs := InsertCategories4Test(ownCatNum)
	userID := primitive.NewObjectID()
	defer func() {
		DeleteCategoryCompletely(ownCatIDs...)
	}()

	RelateCategories2User(ownCatIDs, userID)
	cs, _ := FindCategoriesByUserId(userID)
	if len(cs) == ownCatNum {
		fmt.Println("after relate, len(cs) equals to ownCatIDs")
	}

	//Output:
	//after relate, len(cs) equals to ownCatIDs
}

func ExampleRightFilter() {
	ownCatNum := 100
	notOwnCatNum := 200
	ownCatIDs := InsertCategories4Test(ownCatNum)
	notOwnCatIDs := InsertCategories4Test(notOwnCatNum)
	userID := primitive.NewObjectID()
	defer func() {
		DeleteCategoryCompletely(ownCatIDs...)
		DeleteCategoryCompletely(notOwnCatIDs...)
	}()

	RelateCategories2User(ownCatIDs, userID)

	allowIDs, _ := RightFilter(userID, append(ownCatIDs, notOwnCatIDs...)...)
	if len(allowIDs) == ownCatNum {
		fmt.Println("success")
	} else {
		fmt.Printf("allow: %v; own: %v", len(allowIDs), ownCatNum)
	}

	//Output:
	//success
}
