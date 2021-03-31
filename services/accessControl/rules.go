package accessControl

import (
	"context"
	"github.com/Monkey-Mouse/go-abac/abac"
	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/Monkey-Mouse/mo2/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

type AllowOwn struct {
	UserInfo dto.LoginUserInfo
	ID       primitive.ObjectID `json:"ID"`
	Filter   model.Filter
}

func (r *AllowOwn) ProcessContext(ctx abac.ContextType) {
	*r = ctx.Value("allowOwn").(AllowOwn)
}
func (r *AllowOwn) JudgeRule() (bool, error) {
	blog := database.FindBlogById(r.ID, r.Filter.IsDraft)
	if blog.AuthorID == r.UserInfo.ID {
		return true, nil
	} else {
		return false, nil
	}
}

const accessManagerStr = "accessManager"

var AccessManagerCol = database.GetCollection(accessManagerStr)

type AccessFilter struct {
	VisitorID primitive.ObjectID `json:"visitor_id" bson:"visitor_id"`
	ManagerID primitive.ObjectID `json:"manager_id" bson:"manager_id"`
	RoleList  []string           `json:"role_list,omitempty" example:"'admin':xxxxx 'write':xxxxx" bson:"role_map,omitempty"`
}

// JudgeRule
// 判断id是否在manager的某个role(s)之内，逻辑“与”，必须满足列表内的所有role
func (a *AccessFilter) JudgeRule() (bool, error) {
	var res model.AccessManager
	opt := options.FindOne().SetProjection(bson.M{"_id": 1})
	for _, role := range a.RoleList {
		key := strings.Join([]string{"role_map", role}, ".")
		err := AccessManagerCol.FindOne(context.TODO(), bson.M{"$and": []bson.M{{"_id": a.ManagerID}, {key: bson.M{"$eq": a.VisitorID}}}}, opt).Decode(&res)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return false, nil
			}
			return false, err
		}
	}
	return true, nil
}
func (a *AccessFilter) ProcessContext(ctx abac.ContextType) {
	*a = ctx.Value("accessManager").(AccessFilter)
}
