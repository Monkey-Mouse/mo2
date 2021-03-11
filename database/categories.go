package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mo2/dto"
	"mo2/mo2utils/mo2errors"
	"mo2/server/model"
	"sync"
)

var catCol = GetCollection("category")

// UpsertCategory 更新、插入category
func UpsertCategory(c *model.Directory) (mErr mo2errors.Mo2Errors) {
	if c.ID.IsZero() {
		c.Init()
	}
	res, err := catCol.UpdateOne(context.TODO(), bson.D{{"_id", c.ID}},
		bson.M{"$set": bson.M{"parent_id": c.ParentID, "name": c.Name, "info": c.Info}},
		options.Update().SetUpsert(true))
	if err != nil {
		mErr.InitError(err)
		log.Println(mErr)
	} else {
		mErr.InitNoError("update %v ", res.ModifiedCount)
	}
	return
}

// DeleteCategoryCompletely 删除category，并删除blog中的冗余数据
func DeleteCategoryCompletely(ids ...primitive.ObjectID) (mErr mo2errors.Mo2Errors) {
	var cat model.Directory
	var removeAll sync.WaitGroup
	removeAll.Add(1)
	removeAll.Add(len(ids))
	errSignal := make(chan mo2errors.Mo2Errors)
	for _, id := range ids {
		go func(id primitive.ObjectID) {
			defer removeAll.Done()
			err := catCol.FindOneAndDelete(context.TODO(), bson.M{"_id": id}).Decode(&cat)
			if err != nil {
				mErr.InitError(err)
				//errSignal<-mErr
				log.Println(id)
				log.Println(err)
			} else {
			}
			UpdateSubCategories(id, cat.ParentID)
		}(id)
	}

	go func() {
		defer removeAll.Done()
		defer close(errSignal)
		errSignal <- RemoveCategoriesInAllBlogs(ids...)
	}()
	for errSig := range errSignal {
		mErr = errSig
		log.Println(len(ids), errSig.Error())
	}
	removeAll.Wait()
	return
}

// DeleteCategory
func DeleteCategory(ids ...primitive.ObjectID) (mErr mo2errors.Mo2Errors) {
	if res, err := catCol.DeleteMany(context.TODO(), bson.M{"_id": bson.M{"$in": ids}}); err != nil {
		mErr.InitError(err)
	} else {
		mErr.InitNoError("delete %v categories", res.DeletedCount)
	}
	log.Println(mErr)
	return
}

func UpdateSubCategories(catID primitive.ObjectID, parentID primitive.ObjectID) (mErr mo2errors.Mo2Errors) {
	res, err := catCol.UpdateMany(context.TODO(), bson.M{"parent_id": catID}, bson.M{"$set": bson.M{"parent_id": parentID}})
	if err != nil {
		mErr.InitError(err)
		log.Println(mErr)
		return
	}
	mErr.InitNoError("update %v subCategories", res.ModifiedCount)
	return
}

// RemoveCategoriesInAllBlogs 删除所有blog中的category存在id
func RemoveCategoriesInAllBlogs(catIDs ...primitive.ObjectID) (mErr mo2errors.Mo2Errors) {
	resBlog, err := blogCol.UpdateMany(context.TODO(), bson.M{"categories": bson.M{"$in": catIDs}}, bson.M{"$pullAll": bson.M{"categories": catIDs}})
	if err != nil {
		mErr.InitError(err)
		log.Println(err)
		return
	}
	resDraft, err := draftCol.UpdateMany(context.TODO(), bson.M{"categories": bson.M{"$in": catIDs}}, bson.M{"$pullAll": bson.M{"categories": catIDs}})
	if err != nil {
		mErr.InitError(err)
		log.Println(err)
		return
	}
	mErr.Init(mo2errors.Mo2NoError, fmt.Sprintf("update %v blog(s),%v draft(s)", resBlog.MatchedCount, resDraft.MatchedCount))
	return
}

// FindSubCategories 寻找一个category id的所有子category的详细信息
func FindSubCategories(ID primitive.ObjectID) (cs []model.Directory, mErr mo2errors.Mo2Errors) {
	results, err := catCol.Find(context.TODO(), bson.M{"parent_id": ID})
	if err != nil {
		mErr.InitError(err)
		return
	}
	if err = results.All(context.TODO(), &cs); err != nil {
		mErr.InitError(err)
		return
	}
	return
}

// FindCategories 寻找位于ids列表中所有categories的详细信息
func FindCategories(ids []primitive.ObjectID) (cs []model.Directory, mErr mo2errors.Mo2Errors) {
	cursor, err := catCol.Find(context.TODO(), bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		mErr.InitError(err)
		return
	}
	if err = cursor.All(context.TODO(), &cs); err != nil {
		mErr.InitError(err)
		return
	}
	return
}

// 根据用户id，返回请求的ids中有权限进行某种操作的过滤id列表
func RightFilter(userID primitive.ObjectID, requestIDs ...primitive.ObjectID) (allowIDs []primitive.ObjectID, mErr mo2errors.Mo2Errors) {
	cursor, err := catCol.Find(context.TODO(), bson.M{"$and": []bson.M{bson.M{"_id": bson.M{"$in": requestIDs}}, bson.M{"owner_ids": userID}}}, options.Find().SetProjection(bson.M{"_id": 1}))
	var allowDirects []model.Directory
	if err != nil {
		mErr.InitError(err)
		log.Println(mErr)
	} else {
		if err = cursor.All(context.TODO(), &allowDirects); err != nil {
			mErr.InitError(err)
			log.Println(mErr)
		} else {
			mErr.InitNoError("%v of %v are allowed", len(allowDirects), len(requestIDs))
			log.Println(mErr)
		}
	}
	for _, allowDirect := range allowDirects {
		allowIDs = append(allowIDs, allowDirect.ID)
	}
	return
}

// FindBlogsByCategoryId 寻找包括categoryId的所有blogs的信息
func FindBlogsByCategoryId(id primitive.ObjectID, isDraft bool) (bs []model.Blog, mErr mo2errors.Mo2Errors) {
	var cursor *mongo.Cursor
	var err error
	if isDraft {
		cursor, err = draftCol.Find(context.TODO(), bson.M{"categories": id})
	} else {
		cursor, err = blogCol.Find(context.TODO(), bson.M{"categories": id})
	}
	if err != nil {
		mErr.InitError(err)
		return
	}
	if err = cursor.All(context.TODO(), &bs); err != nil {
		mErr.InitError(err)
		return
	}
	return
}

// SortCategories 递归建立categories的树形结构
func SortCategories(c model.Directory, m map[string][]model.Directory) {
	var cs []model.Directory
	cs, _ = FindSubCategories(c.ID)
	if len(cs) == 0 {
		return
	} else {
		m[c.ID.Hex()] = cs
		for _, category := range cs {
			SortCategories(category, m)
		}
	}
}

//FindCategoryById find category by id
func FindCategoryById(id primitive.ObjectID) (c model.Directory) {
	err := catCol.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&c)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
	return
}

//FindAllCategories find all categories
func FindAllCategories() (cs []model.Directory) {
	results, err := catCol.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	err = results.All(context.TODO(), &cs)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// RelateMainCategory2User 为用户创建主归档目录
func RelateMainCategory2User(e2e dto.RelateEntity2Entity) (mErr mo2errors.Mo2Errors) {
	if _, err := catCol.UpdateOne(context.TODO(), bson.M{"_id": e2e.RelatedID}, bson.M{"$set": bson.M{"parent_id": e2e.RelateToID}}); err != nil {
		mErr.InitError(err)
	} else {
		mErr = RelateCategory2User(e2e.RelatedID, e2e.RelateToID)
	}
	return
}

// RelateSubCategory2Category 将子归档subCategory的parent_id设为category
func RelateSubCategory2Category(e2e dto.RelateEntity2Entity) {
	catCol.UpdateOne(context.TODO(), bson.M{"_id": e2e.RelatedID}, bson.M{"$set": bson.M{"parent_id": e2e.RelateToID}})
}

// RelateSubCategories2Category 将子归档subCategories的parent_id设为category
func RelateSubCategories2Category(s2e dto.RelateEntitySet2Entity) {
	catCol.UpdateMany(context.TODO(), bson.M{"_id": bson.M{"$in": s2e.RelatedIDs}}, bson.M{"$set": bson.M{"parent_id": s2e.RelateToID}})
}

// RelateCategory2Blog 在blog及draft的categories中添加category的id
func RelateCategory2Blog(e2e dto.RelateEntity2Entity) {
	blogCol.UpdateOne(context.TODO(), bson.M{"_id": e2e.RelateToID}, bson.M{"$addToSet": bson.M{"categories": e2e.RelatedID}})
	draftCol.UpdateOne(context.TODO(), bson.M{"_id": e2e.RelateToID}, bson.M{"$addToSet": bson.M{"categories": e2e.RelatedID}})
}

// RelateCategory2Blog 在blog及draft的categories中添加categories的id
func RelateCategories2Blog(s2e dto.RelateEntitySet2Entity) {
	blogCol.UpdateOne(context.TODO(), bson.M{"_id": s2e.RelateToID}, bson.M{"$addToSet": bson.M{"categories": bson.M{"$each": s2e.RelatedIDs}}})
	draftCol.UpdateOne(context.TODO(), bson.M{"_id": s2e.RelateToID}, bson.M{"$addToSet": bson.M{"categories": bson.M{"$each": s2e.RelatedIDs}}})
}

// RelateCategories2Blogs 在blogs及drafts的categories中添加categories的id
// Todo find nice way to make of use
func RelateCategories2Blogs(s2s dto.RelateEntitySet2EntitySet) (result []model.Blog) {
	// 将所有满足条件的blog/draft进行更新
	blogCol.UpdateMany(context.TODO(), bson.D{{"_id", bson.D{{"$in", s2s.RelateToIDs}}}}, bson.D{{"$addToSet", bson.M{"categories": bson.M{"$each": s2s.RelatedIDs}}}})
	draftCol.UpdateMany(context.TODO(), bson.D{{"_id", bson.D{{"$in", s2s.RelateToIDs}}}}, bson.D{{"$addToSet", bson.M{"categories": bson.M{"$each": s2s.RelatedIDs}}}})
	// todo delete useless result
	cursor, _ := blogCol.Find(context.TODO(), bson.D{{"_id", bson.M{"_id": bson.M{"$in": s2s.RelateToIDs}}}}, options.Find().SetProjection(bson.M{"content": 0}))
	cursor.All(context.TODO(), &result)
	return
}

// RelateCategory2User 在category的ownerIds中添加userIDs
func RelateCategory2User(catID primitive.ObjectID, userIds ...primitive.ObjectID) (mErr mo2errors.Mo2Errors) {
	//todo check if valid
	res, err := catCol.UpdateMany(context.TODO(), bson.M{"_id": catID}, bson.M{"$addToSet": bson.M{"owner_ids": bson.M{"$each": userIds}}})
	if err != nil {
		mErr.Init(mo2errors.Mo2Error, err.Error())
		return
	}
	mErr.Init(mo2errors.Mo2NoError, fmt.Sprintf("%v modified", res.ModifiedCount))
	return
}

// RelateCategories2User 在列表category的ownerIds中添加userIDs
func RelateCategories2User(catIDs []primitive.ObjectID, userIds ...primitive.ObjectID) (mErr mo2errors.Mo2Errors) {
	res, err := catCol.UpdateMany(context.TODO(), bson.M{"_id": bson.M{"$in": catIDs}}, bson.M{"$addToSet": bson.M{"owner_ids": bson.M{"$each": userIds}}})
	if err != nil {
		mErr.Init(mo2errors.Mo2Error, err.Error())
		return
	}
	mErr.Init(mo2errors.Mo2NoError, fmt.Sprintf("%v modified", res.ModifiedCount))
	return
}

//FindCategoriesByUserId find categories by user id
func FindCategoriesByUserId(userId ...primitive.ObjectID) (cs []model.Directory, mErr mo2errors.Mo2Errors) {
	// disable sort in backend
	//m = make(map[string][]model.Directory)
	//SortCategories(c, m)
	cursor, err := catCol.Find(context.TODO(), bson.M{"owner_ids": bson.M{"$in": userId}})
	if err != nil {
		mErr.InitError(err)
		return
	}
	if err = cursor.All(context.TODO(), &cs); err != nil {
		mErr.InitError(err)
		return
	}
	return
}
