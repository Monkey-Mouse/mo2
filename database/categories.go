package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mo2/dto"
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
func FindSubCategories(c model.Category) (cs []model.Category) {
	results, err := catCol.Find(context.TODO(), bson.M{"parent_id": c.ID})
	if err != nil {
		log.Fatal(err)
	}
	if err = results.All(context.TODO(), &cs); err != nil {
		log.Fatal(err)
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
	cs = FindSubCategories(c)
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
func AddBlogs2Categories(ab2cs dto.AddBlogs2Categories) (results []dto.QueryBlog) {
	var categories []model.Category
	for _, categoryID := range ab2cs.CategoryIDs {
		var category model.Category
		category = FindCategoryById(categoryID)
		if category.IsValid() {
			categories = append(categories, category)
		}
	}
	var blog, draft model.Blog
	if len(categories) > 0 {
		for _, blogID := range ab2cs.BlogIDs {
			blog = FindBlogById(blogID, false)
			if blog.IsValid() {
				AddBlog2Categories(&blog, categories, false)
				results = append(results, dto.MapBlog2QueryBlog(blog))
			}
			draft = FindBlogById(blogID, true)
			if draft.IsValid() {
				AddBlog2Categories(&draft, categories, true)
				results = append(results, dto.MapBlog2QueryBlog(draft))
			}
		}
	}
	return results
}
func AddBlog2Categories(blog *model.Blog, categories []model.Category, isDraft bool) {
	blog.Categories = append(blog.Categories, categories...)
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
func AddCategoryId2User(catId primitive.ObjectID, userId primitive.ObjectID) {
	catUser := model.CategoryUser{
		UserID:     userId,
		CategoryID: catId,
	}
	//todo check if valid
	catUserCol.InsertOne(context.TODO(), catUser)
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
func FindCategoriesByUserId(id primitive.ObjectID) (m map[string][]model.Category) {
	c := FindCategoryByUserId(id)
	if c.ID.IsZero() {
		return
	}
	m = make(map[string][]model.Category)
	SortCategories(c, m)

	return

}