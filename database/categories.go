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
)

var catCol = GetCollection("category")

// upsertCategory 更新、插入category
func UpsertCategory(c *model.Directory) {
	if c.ID.IsZero() {
		c.Init()
	}
	_, err := catCol.UpdateOne(context.TODO(), bson.D{{"_id", c.ID}},
		bson.M{"$set": bson.M{"parent_id": c.ParentID, "name": c.Name}},
		options.Update().SetUpsert(true))
	if err != nil {
		log.Fatal(err)
	}
}

// FindSubCategories 寻找一个categoryid的所有子category的详细信息
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
func RelateMainCategory2User(e2e dto.RelateEntity2Entity) {
	catCol.UpdateOne(context.TODO(), bson.M{"_id": e2e.RelatedID}, bson.M{"$set": bson.M{"parent_id": e2e.RelateToID}})
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
	blogCol.UpdateMany(context.TODO(), bson.D{{"_id", bson.D{{"$in", s2s.RelatedIDs}}}}, bson.D{{"$addToSet", bson.M{"categories": bson.M{"$each": s2s.RelateToIDs}}}})
	draftCol.UpdateMany(context.TODO(), bson.D{{"_id", bson.D{{"$in", s2s.RelatedIDs}}}}, bson.D{{"$addToSet", bson.M{"categories": bson.M{"$each": s2s.RelateToIDs}}}})
	// todo delete useless result
	cursor, _ := blogCol.Find(context.TODO(), bson.D{{"_id", bson.M{"_id": bson.M{"$in": s2s.RelatedIDs}}}}, options.Find().SetProjection(bson.M{"content": 0}))
	cursor.All(context.TODO(), &result)
	return
}

// RelateCategory2User 在category的ownerIds中添加userIDs
func RelateCategory2User(catId primitive.ObjectID, userIds ...primitive.ObjectID) (mErr mo2errors.Mo2Errors) {
	//todo check if valid
	res, err := catCol.UpdateMany(context.TODO(), bson.M{"_id": catId}, bson.M{"$addToSet": bson.M{"owner_ids": bson.M{"$each": userIds}}})
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
