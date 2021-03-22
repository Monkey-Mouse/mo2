package database

import (
	"context"

	"github.com/Monkey-Mouse/mo2/mo2utils/mo2errors"
	"github.com/Monkey-Mouse/mo2/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CategoryCol = "category"
)

//ChooseCol 选择数据集合
// 根据collection名称选取相应名称的集合
func ChooseCol(collection string) (col *mongo.Collection, mErr mo2errors.Mo2Errors) {
	switch collection {
	case CategoryCol:
		col = catCol
	default:
		mErr.Init(mo2errors.Mo2NoExist, "collection %v not exist", collection)
	}
	return
}

// FindDirectoryInfo 根据collection选取集合，并在其中寻找在id在ids列表中的所有directory类型的信息
// 现在返回的数据有id，info与name
func FindDirectoryInfo(collection string, ids ...primitive.ObjectID) (info []model.Directory, mErr mo2errors.Mo2Errors) {
	col, mErr := ChooseCol(collection)
	if mErr.IsError() {
		return
	}
	cursor, err := col.Find(context.TODO(), bson.M{"_id": bson.M{"$in": ids}}, options.Find().SetProjection(bson.M{"_id": 1, "info": 1, "name": 1}))
	if err != nil {
		mErr.InitError(err)
		return
	}
	if err = cursor.All(context.TODO(), &info); err != nil {
		mErr.InitError(err)
		return
	}
	return
}
