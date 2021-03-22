package database

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/Monkey-Mouse/mo2/mo2utils/mo2errors"
	"github.com/Monkey-Mouse/mo2/server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var catCol = GetCollection(CategoryCol)

func init() {
	createCategoryIndexes(catCol)
}

func createCategoryIndexes(col *mongo.Collection) {
	col.Indexes().CreateMany(context.TODO(), append([]mongo.IndexModel{
		{Keys: bson.M{"_id": 1}},
		{Keys: bson.M{"parent_id": 1}},
	}, model.IndexModels...))
}

// UpsertCategory 更新、插入category
func UpsertCategory(c *model.Directory) (mErr mo2errors.Mo2Errors) {
	update := bson.M{"$set": bson.M{"parent_id": c.ParentID, "name": c.Name, "info": c.Info, "owner_ids": c.OwnerIDs}}
	res, err := catCol.UpdateOne(context.TODO(), bson.D{{"_id", c.ID}}, update, options.Update().SetUpsert(true))
	if err != nil {
		mErr.InitError(err)
		log.Println(mErr)
	} else {
		mErr.InitNoError("update %v ", res.UpsertedID)
	}
	return
}

// FindOrCreateRoot4User
// 寻找用户的根目录，若不存在，新建并返回
func FindOrCreateRoot4User(userID primitive.ObjectID) (cat model.Directory, mErr mo2errors.Mo2Errors) {
	if err := catCol.FindOne(context.TODO(), bson.M{"parent_id": userID}).Decode(&cat); err != nil {
		if err == mongo.ErrNoDocuments {
			//create new root for user
			cat = model.Directory{ID: primitive.NewObjectID(), ParentID: userID, OwnerIDs: []primitive.ObjectID{userID}}
			mErr = UpsertCategory(&cat)
			return
		}
		mErr.InitError(err)
		log.Println(mErr)
	} else {
		mErr.InitNoError("have already exists")
	}
	return
}

// idExistInList 检查id是否存在与给的的id列表中，若存在，返回true；否则返回false
func idExistInList(checkId primitive.ObjectID, ids []primitive.ObjectID) (exist bool) {
	exist = false
	for _, id := range ids {
		if id == checkId {
			exist = true
		}
	}
	return
}

// UpsertDirectoryByUser 用户id为userID的用户更新、插入category
// 若插入的ownerIDs不存在，则将请求方加入
// 		若存在，则检查请求方id是否存在访问权限，若无则直接返回
// 若插入的parentID不存在，则检查请求用户是否有root
// 				若有root，则使用root的id做为parentID插入
// 				若无root，则创建并使用新建的root的id做为parentID插入
func UpsertDirectoryByUser(c *model.Directory, userID primitive.ObjectID) (mErr mo2errors.Mo2Errors) {
	if len(c.OwnerIDs) == 0 {
		c.OwnerIDs = append(c.OwnerIDs, userID)
	} else {
		if idExistInList(userID, c.OwnerIDs) == false {
			mErr.Init(mo2errors.Mo2Unauthorized, "无访问权限")
			return
		}
	}
	if c.ID.IsZero() {
		c.Init()
	}
	if c.ParentID.IsZero() {
		// find add people's root directory(category)
		var root model.Directory
		if root, mErr = FindOrCreateRoot4User(userID); mErr.IsError() {
			return
		} else {
			c.ParentID = root.ID
		}
	}
	mErr = UpsertCategory(c)
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
	resBlog, err := BlogCol.UpdateMany(context.TODO(), bson.M{"categories": bson.M{"$in": catIDs}}, bson.M{"$pullAll": bson.M{"categories": catIDs}})
	if err != nil {
		mErr.InitError(err)
		log.Println(err)
		return
	}
	resDraft, err := DraftCol.UpdateMany(context.TODO(), bson.M{"categories": bson.M{"$in": catIDs}}, bson.M{"$pullAll": bson.M{"categories": catIDs}})
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

// RightFilter
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
	BlogCol.UpdateOne(context.TODO(), bson.M{"_id": e2e.RelateToID}, bson.M{"$addToSet": bson.M{"categories": e2e.RelatedID}})
	DraftCol.UpdateOne(context.TODO(), bson.M{"_id": e2e.RelateToID}, bson.M{"$addToSet": bson.M{"categories": e2e.RelatedID}})
}

// RelateCategory2Blog 在blog及draft的categories中添加categories的id
func RelateCategories2Blog(s2e dto.RelateEntitySet2Entity) {
	BlogCol.UpdateOne(context.TODO(), bson.M{"_id": s2e.RelateToID}, bson.M{"$addToSet": bson.M{"categories": bson.M{"$each": s2e.RelatedIDs}}})
	DraftCol.UpdateOne(context.TODO(), bson.M{"_id": s2e.RelateToID}, bson.M{"$addToSet": bson.M{"categories": bson.M{"$each": s2e.RelatedIDs}}})
}

// RelateCategories2Blogs 在blogs及drafts的categories中添加categories的id
// Todo find nice way to make of use
func RelateCategories2Blogs(s2s dto.RelateEntitySet2EntitySet) (result []model.Blog) {
	// 将所有满足条件的blog/draft进行更新
	BlogCol.UpdateMany(context.TODO(), bson.D{{"_id", bson.D{{"$in", s2s.RelateToIDs}}}}, bson.D{{"$addToSet", bson.M{"categories": bson.M{"$each": s2s.RelatedIDs}}}})
	DraftCol.UpdateMany(context.TODO(), bson.D{{"_id", bson.D{{"$in", s2s.RelateToIDs}}}}, bson.D{{"$addToSet", bson.M{"categories": bson.M{"$each": s2s.RelatedIDs}}}})
	// todo delete useless result
	cursor, _ := BlogCol.Find(context.TODO(), bson.D{{"_id", bson.M{"_id": bson.M{"$in": s2s.RelateToIDs}}}}, options.Find().SetProjection(bson.M{"content": 0}))
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
	} else {
		mErr.Init(mo2errors.Mo2NoError, fmt.Sprintf("%v modified", res.ModifiedCount))
	}
	log.Println(mErr)
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
