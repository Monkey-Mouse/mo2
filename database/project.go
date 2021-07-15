package database

import (
	"context"
	"fmt"

	"github.com/Monkey-Mouse/mo2/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const projStr = "group"

var ProjCol = GetCollection(projStr)

type Project struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Name        string               `bson:"name"`
	Tags        []string             `bson:"tags"`
	OwnerID     primitive.ObjectID   `bson:"owner_id"`
	ManagerIDs  []primitive.ObjectID `bson:"manager_i_ds"`
	MemberIDs   []primitive.ObjectID `bson:"member_i_ds"`
	BlogIDs     []primitive.ObjectID `bson:"blog_i_ds"`
	Description string               `bson:"description"`
	Avatar      string               `bson:"avatar"`
	EntityInfo  model.Entity         `bson:"entity_info"`
}

func UpsertProject(ctx context.Context, p *Project, update bson.M) (*mongo.UpdateResult, error) {
	if p.ID.IsZero() {
		p.ID = primitive.NewObjectID()
		p.EntityInfo = model.Entity{}
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
	p.ID = re.UpsertedID.(primitive.ObjectID)
	return re, nil
}

func ListProject(ctx context.Context, page, pageSize int64, query interface{}) ([]*Project, error) {
	pageSize = sanitizePagesize(pageSize)
	c, err := ProjCol.Find(ctx, query, getPaginationOption(page, pageSize))
	if err != nil {
		return nil, err
	}
	re := make([]*Project, pageSize)
	err = c.All(ctx, &re)
	return re, err
}
func GetProject(ctx context.Context, query interface{}) (p *Project, err error) {
	p = &Project{}
	err = ProjCol.FindOne(ctx, query).Decode(p)
	return
}

func AddBlogToProj(ctx context.Context, pid primitive.ObjectID, bid primitive.ObjectID) error {
	if pid.IsZero() {
		return nil
	}
	proj, err := UpsertProject(ctx, &Project{ID: pid}, bson.M{"$push": bson.M{"blog_i_ds": bid}})
	if err != nil {
		return err
	}
	if proj.UpsertedCount == 0 {
		return fmt.Errorf("project not found")
	}
	return nil
}
