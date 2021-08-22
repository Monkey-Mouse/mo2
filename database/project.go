package database

import (
	"context"

	"github.com/Monkey-Mouse/mo2/mo2utils"
	"github.com/Monkey-Mouse/mo2/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const projStr = "group"

var ProjCol = GetCollection(projStr)

func init() {
	ProjCol.Indexes().CreateMany(context.TODO(), append([]mongo.IndexModel{
		{Keys: bson.M{"tags": 1}},
		{Keys: bson.M{"owner_id": 1}},
	}, model.IndexModels...))
}

// UpsertProject 插入或更新project
// - 如果插入，ManagerIDs和MemberIDs会被设置为空数组
func UpsertProject(ctx context.Context, p *model.Project, update bson.M) (*mongo.UpdateResult, error) {
	if p.ID.IsZero() {
		p.ID = primitive.NewObjectID()
		p.EntityInfo = model.Entity{}
		p.BlogIDs = []primitive.ObjectID{}
		p.ManagerIDs = []primitive.ObjectID{}
		p.MemberIDs = []primitive.ObjectID{}
		p.EntityInfo.Create()
	} else {
		p.EntityInfo.Update()
	}
	if update == nil {
		update = bson.M{"$set": p}
	}
	re, err := ProjCol.UpdateOne(ctx, bson.M{"_id": p.ID}, update, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}
	mo2utils.IndexProject(p)
	if re.UpsertedID != nil {
		p.ID = re.UpsertedID.(primitive.ObjectID)
	}
	return re, nil
}

func ListProject(ctx context.Context, page, pageSize int64, query interface{}) ([]*model.Project, error) {
	pageSize = sanitizePagesize(pageSize)
	c, err := ProjCol.Find(ctx, query, getPaginationOption(page, pageSize))
	if err != nil {
		return nil, err
	}
	re := make([]*model.Project, pageSize)
	err = c.All(ctx, &re)
	return re, err
}
func GetProject(ctx context.Context, query interface{}) (p *model.Project, err error) {
	p = &model.Project{}
	err = ProjCol.FindOne(ctx, query).Decode(p)
	return
}
func DeleteProject(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {

	re, err := ProjCol.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	mo2utils.DeleteProjectIndex(id.Hex())
	return re, nil
}
