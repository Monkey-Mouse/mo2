package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mo2/mo2utils/mo2errors"
	"mo2/server/model"
	"sync"
	"time"
)

var binCol = GetCollection("recycleBin")

// UpsertRecycleItem 新增或更新recycleBin中的item
func UpsertRecycleItem(item model.RecycleItem) (mErr mo2errors.Mo2Errors) {
	res, err := binCol.UpdateOne(context.TODO(), bson.M{"_id": item.ID}, bson.M{"$set": item}, options.Update().SetUpsert(true))
	if err != nil {
		mErr.InitError(err)
	} else {
		mErr.InitNoError("match %v rec item, upsert %v", res.MatchedCount, res.UpsertedCount)
	}
	log.Println(mErr)
	return
}

// DeleteRecycleItems
// 根据id删除Items
func DeleteRecycleItems(ids ...primitive.ObjectID) (mErr mo2errors.Mo2Errors) {
	res, err := binCol.DeleteMany(context.TODO(), bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		mErr.InitError(err)
	} else {
		mErr.InitNoError("delete %v item(s)", res.DeletedCount)
	}
	log.Printf("delete %v rec item(s)", res.DeletedCount)
	return
}

// DeleteByRecycleItemID
// 根据itemID删除Item
func DeleteByRecycleItemID(ID primitive.ObjectID) (mErr mo2errors.Mo2Errors) {
	res, err := binCol.DeleteMany(context.TODO(), bson.M{"item_id": ID})
	if err != nil {
		mErr.InitError(err)
	} else {
		mErr.InitNoError("delete %v item(s)", res.DeletedCount)
	}
	log.Printf("delete %v rec item(s)", res.DeletedCount)
	return
}

// DeleteByRecycleItemInfo
// 根据itemID删除Item
func DeleteByRecycleItemInfo(ID primitive.ObjectID, handler string) (mErr mo2errors.Mo2Errors) {
	res, err := binCol.DeleteMany(context.TODO(), bson.M{"item_id": ID, "handler": handler})
	if err != nil {
		mErr.InitError(err)
	} else {
		mErr.InitNoError("delete %v item(s)", res.DeletedCount)
	}
	log.Printf("delete %v rec item(s)", res.DeletedCount)
	return
}

// DeleteExpireItems
// 删除已到规定删除时间的Items
func DeleteExpireItems() (mErr mo2errors.Mo2Errors) {
	cursor, err := binCol.Find(context.TODO(), bson.M{"delete_time": bson.M{"$lte": time.Now()}})

	if err != nil {
		mErr.InitError(err)
	} else {
		var expItems []model.RecycleItem
		if err = cursor.All(context.TODO(), &expItems); err != nil {
			mErr.InitError(err)
		} else {
			var proc sync.WaitGroup
			proc.Add(len(expItems))
			// process each item
			for _, item := range expItems {
				go func(item model.RecycleItem) { deleteMarkItem(item); defer proc.Done() }(item)
			}
			proc.Wait()
			mErr.InitNoError("finish one clean recycle bin task")
		}
	}
	return
}

// 根据bin中提供的item信息，进行删除操作
func deleteMarkItem(item model.RecycleItem) {
	var mErr mo2errors.Mo2Errors
	switch item.Handler {
	case model.HandlerBlog:
		if mErr = DeleteBlogs(false, item.ItemID); !mErr.IsError() {
			DeleteRecycleItems(item.ID)
		}
	case model.HandlerDraft:
		if mErr = DeleteBlogs(true, item.ItemID); !mErr.IsError() {
			DeleteRecycleItems(item.ID)
		}
	default:
		log.Println("invalid handler")
		DeleteRecycleItems(item.ID)
	}
}
