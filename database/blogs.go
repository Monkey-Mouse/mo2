package database

import (
	"context"
	"log"
	"time"

	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/Monkey-Mouse/mo2/mo2utils"
	"github.com/Monkey-Mouse/mo2/server/model"

	"github.com/Monkey-Mouse/mo2/mo2utils/mo2errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	OperationKey     = "operation"
	OperationRecycle = "recycle"
	OperationRestore = "restore"
)
const (
	DurationBeforeBlogDelete = time.Hour * 24 * 30
)

// BlogCol blog col
var BlogCol *mongo.Collection = GetCollection("blog")

// DraftCol draft col
var DraftCol *mongo.Collection = GetCollection("draft")

func init() {
	createBlogIndexes(BlogCol)
	createBlogIndexes(DraftCol)
}

func createBlogIndexes(col *mongo.Collection) {
	col.Indexes().CreateMany(context.TODO(), append([]mongo.IndexModel{
		{Keys: bson.M{"ket_words": 1}},
		{Keys: bson.M{"author_id": 1}},
		{Keys: bson.M{"categories": 1}},
		{Keys: bson.M{"project_id": 1}},
	}, model.IndexModels...))
}
func chooseCol(isDraft bool) (col *mongo.Collection) {
	col = DraftCol
	if !isDraft {
		col = BlogCol
	}
	return
}

// InsertBlog insert
func insertBlog(b *model.Blog, isDraft bool) (mErr mo2errors.Mo2Errors) {
	col := chooseCol(isDraft)
	if res, err := col.InsertOne(context.TODO(), b); err != nil {
		log.Println(err)
		mErr.InitError(err)
	} else {
		mErr.InitNoError("insert %v", res.InsertedID)
	}
	return
}

func ScoreBlog(ctx context.Context, b *model.Blog, prev float64, now float64) (sum float64, num int) {
	col := chooseCol(false)
	sum = b.ScoreSum + now
	num = b.ScoreNum
	if prev < 0 {
		num++
	} else {
		sum -= prev
	}
	col.UpdateByID(ctx, b.ID, bson.D{{"$set", bson.M{"score_sum": sum, "score_num": num}}})
	return
}

// upsertBlog
func upsertBlog(b *model.Blog, isDraft bool) (mErr mo2errors.Mo2Errors) {
	col := chooseCol(isDraft)
	b.EntityInfo.Update()
	result, err := col.UpdateOne(
		context.TODO(),
		bson.D{{"_id", b.ID}, {"author_id", b.AuthorID}},
		bson.D{{"$set", bson.M{
			"entity_info": b.EntityInfo,
			"title":       b.Title,
			"description": b.Description,
			"content":     b.Content,
			"cover":       b.Cover,
			"key_words":   b.KeyWords,
			"categories":  b.CategoryIDs,
			"y_doc":       b.YDoc,
			"project_id":  b.ProjectID,
		}}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Println(err)
		mErr.InitError(err)
		return
	}
	if !isDraft {
		log.Println("发布时删除草稿" + b.ID.String())
		mErr = DeleteBlogs(true, b.ID)
	}
	if result.UpsertedCount != 0 {
		mErr.InitNoError("新建文章" + b.ID.String())
	}
	return
}

// DeleteBlogs 彻底删除
func DeleteBlogs(isDraft bool, blogIDs ...primitive.ObjectID) (mErr mo2errors.Mo2Errors) {
	if res, err := chooseCol(isDraft).DeleteMany(context.TODO(), bson.M{"_id": bson.M{"$in": blogIDs}}); err != nil {
		mErr.InitError(err)
	} else {
		mErr.InitNoError("delete %v %v(s)\n", res.DeletedCount, If(isDraft, "draft", "blog"))
	}
	log.Println(mErr)
	return
}

// UpsertBlog upsert blog or draft
func UpsertBlog(b *model.Blog, isDraft bool) (mErr mo2errors.Mo2Errors) {
	if b.ID == primitive.NilObjectID {
		b.Init()
		mErr = insertBlog(b, isDraft)
	} else {
		mErr = upsertBlog(b, isDraft)
	}
	if !isDraft {
		mo2utils.IndexBlog(b)
	}
	return
}

// ProcessBlog process blog or draft
// 根据operation对blog/draft的信息进行更新
// recycle:新增关于本blog/draft的recycleBin信息，且isDeleted字段置为true状态
// restore:将recycleBin中关于本blog/draft的信息进行删除，且isDeleted字段恢复为false状态
func ProcessBlog(isDraft bool, b *model.Blog, operation string) (mErr mo2errors.Mo2Errors) {

	item := model.RecycleItem{
		ID:         primitive.NewObjectID(),
		ItemID:     b.ID,
		CreateTime: time.Now(),
		DeleteTime: time.Now().Add(DurationBeforeBlogDelete),
		Handler:    If(isDraft, model.HandlerDraft, model.HandlerBlog).(string),
	}
	switch operation {
	case OperationRecycle:
		b.EntityInfo.IsDeleted = true
		mErr = UpsertRecycleItem(item)
	case OperationRestore:
		b.EntityInfo.IsDeleted = false
		mErr = DeleteByRecycleItemInfo(b.ID, item.Handler)
	default:
		log.Println("invalid operation")
		mErr.Init(mo2errors.Mo2NoExist, "invalid operation")
	}
	return
}

//find blog by user
func FindBlogsByUser(u dto.LoginUserInfo, filter model.Filter) (b []model.Blog) {
	return FindBlogsByValue("author_id", u.ID, filter)
}
func getBlogListQueryOption() *options.FindOptions {
	return options.Find().SetSort(bson.D{{"entity_info.update_time", -1}}).SetProjection(bson.D{{"content", 0}, {"y_doc", 0}})
}

//find blog by value
func FindBlogsByValue(field string, val interface{}, filter model.Filter) (b []model.Blog) {
	col := chooseCol(filter.IsDraft)
	opts := getBlogListQueryOption().SetSkip(int64(filter.Page * filter.PageSize)).SetLimit(int64(filter.PageSize))
	cursor, err := col.Find(context.TODO(), bson.M{field: val, "entity_info.isdeleted": filter.IsDeleted}, opts)
	if err != nil {
		panic(err)
	}
	err = cursor.All(context.TODO(), &b)
	if err != nil {
		panic(err)
	}
	return
}

//find blog by id
func FindBlogById(id primitive.ObjectID, isDraft bool, opts ...*options.FindOneOptions) (b model.Blog) {
	col := chooseCol(isDraft)
	err := col.FindOne(context.TODO(), bson.D{{"_id", id}}, opts...).Decode(&b)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}
		panic(err)
	}
	return
}

// FindBlogs find blog
func FindBlogs(filter model.Filter) (b []model.Blog) {
	col := chooseCol(filter.IsDraft)
	opts := getBlogListQueryOption().SetSkip(int64(filter.Page * filter.PageSize)).SetLimit(int64(filter.PageSize))
	f := bson.D{{"entity_info.isdeleted", filter.IsDeleted}}
	if len(filter.Ids) != 0 {
		f = bson.D{
			{"entity_info.isdeleted", filter.IsDeleted},
			{"_id", bson.M{"$in": filter.Ids}},
		}
	}
	cursor, err := col.Find(context.TODO(), f, opts)
	err = cursor.All(context.TODO(), &b)
	if err != nil {
		panic(err)
	}
	return
}

// FindBlogsByCategoryId 寻找包括categoryId的所有blogs的信息
func FindBlogsByCategoryId(catID primitive.ObjectID, filter model.Filter) (bs []model.Blog, mErr mo2errors.Mo2Errors) {
	var cursor *mongo.Cursor
	var err error
	col := chooseCol(filter.IsDraft)
	opts := getBlogListQueryOption() //.SetSkip(int64(filter.Page * filter.PageSize)).SetLimit(int64(filter.PageSize))

	cursor, err = col.Find(context.TODO(), bson.M{"categories": catID, "entity_info.isdeleted": filter.IsDeleted}, opts)

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
