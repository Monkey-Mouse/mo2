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
var catUserCol = GetCollection("userCategory")

func UpsertCategory(c *model.Category) {
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
func FindSubCategories(ID primitive.ObjectID) (cs []model.Category, mErr mo2errors.Mo2Errors) {
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
func FindCategories(ids []primitive.ObjectID) (cs []model.Category) {
	var c model.Category
	for _, id := range ids {
		err := catCol.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&c)
		if err != nil {
			if err != mongo.ErrNoDocuments {
				log.Fatal(err)
			}
		}
		cs = append(cs, c)
	}
	return
}

func SortCategories(c model.Category, m map[string][]model.Category) {
	var cs []model.Category
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

//find category by id
func FindCategoryById(id primitive.ObjectID) (c model.Category) {
	err := catCol.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&c)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
	return
}
func FindAllCategories() (cs []model.Category) {
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
func AddBlogs2Categories(ab2cs dto.AddBlogs2Categories) (result []model.Blog) {
	// 将所有满足条件的blog/draft进行更新
	blogCol.UpdateMany(context.TODO(), bson.D{{"_id", bson.D{{"$in", ab2cs.BlogIDs}}}}, bson.D{{"$addToSet", bson.M{"categories": bson.M{"$each": ab2cs.CategoryIDs}}}})
	draftCol.UpdateMany(context.TODO(), bson.D{{"_id", bson.D{{"$in", ab2cs.BlogIDs}}}}, bson.D{{"$addToSet", bson.M{"categories": bson.M{"$each": ab2cs.CategoryIDs}}}})

	cursor, _ := blogCol.Find(context.TODO(), bson.D{{"_id", bson.M{"_id": bson.M{"$in": ab2cs.BlogIDs}}}}, options.Find().SetProjection(bson.M{"content": 0}))
	cursor.All(context.TODO(), &result)
	return
}
func AddCategories2Category(parCategoryID primitive.ObjectID, subCategoryIDs ...primitive.ObjectID) {
	catCol.UpdateOne(context.TODO(), bson.M{"_id": bson.M{"$in": subCategoryIDs}}, bson.M{"$set": bson.M{"parent_id": parCategoryID}})
}
func AddBlog2Categories(blog *model.Blog, categoryIDs []primitive.ObjectID, isDraft bool) {
	blog.CategoryIDs = append(blog.CategoryIDs, categoryIDs...)
	upsertBlog(blog, isDraft)
}
func AddCategory2User(category model.Category, userId primitive.ObjectID) {
	if !category.IsValid() {
		category.Init()
	}
	AddCategoryId2User(category.ID, userId)
}
func AddCategoryIdStr2User(catIdStr string, userId primitive.ObjectID) {
	catId, err := primitive.ObjectIDFromHex(catIdStr)
	var c model.Category
	if err != nil {
		// not exist, create one category first
		UpsertCategory(&c)
		catId = c.ID
	}
	AddCategoryId2User(catId, userId)
}
func AddCategoryId2User(catId primitive.ObjectID, userIds ...primitive.ObjectID) (mErr mo2errors.Mo2Errors) {
	//todo check if valid
	res, err := catCol.UpdateMany(context.TODO(), bson.M{"_id": catId}, bson.M{"$addToSet": bson.M{"owner_ids": bson.M{"$each": userIds}}})
	if err != nil {
		mErr.Init(mo2errors.Mo2Error, err.Error())
		return
	}
	mErr.Init(mo2errors.Mo2NoError, fmt.Sprintf("%v modified", res.ModifiedCount))
	return
}

//find category by userid
func FindCategoryByUserId(id primitive.ObjectID) (c model.Category) {
	var cu model.CategoryUser
	err := catUserCol.FindOne(context.TODO(), bson.D{{"user_id", id}}).Decode(&cu)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
	if err = catCol.FindOne(context.TODO(), bson.D{{"_id", cu.CategoryID}}).Decode(&c); err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
	return
}

//find categories by user id
func FindCategoriesByUserId(userId ...primitive.ObjectID) (cs []model.Category, mErr mo2errors.Mo2Errors) {
	// disable sort in backend
	//m = make(map[string][]model.Category)
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
