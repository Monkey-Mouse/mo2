package accessControl

import (
	"context"
	"github.com/Monkey-Mouse/go-abac/abac"
	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/Monkey-Mouse/mo2/mo2utils/mo2errors"
	"github.com/Monkey-Mouse/mo2/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

const (
	RuleAllowOwn     = "allowOwn"
	RuleAccessFilter = "accessFilter"
)

const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleMember = "member"
	RoleReader = "reader"
)

type AllowOwn struct {
	UserInfo dto.LoginUserInfo
	ID       primitive.ObjectID `json:"id"`
	Filter   model.Filter
	Resource string
}

func (r *AllowOwn) ProcessContext(ctx abac.ContextType) {
	if val := ctx.Value(RuleAllowOwn); val != nil {
		*r = val.(AllowOwn)
	}
}
func (r *AllowOwn) JudgeRule() (bool, error) {
	if r.UserInfo.ID.IsZero() {
		return false, mo2errors.Init(mo2errors.Mo2NoLogin, "not login")
	}
	switch r.Resource {
	case ResourceBlog:
		{
			blog := database.FindBlogById(r.ID, r.Filter.IsDraft)
			if blog.AuthorID == r.UserInfo.ID {
				return true, nil
			} else {
				return false, nil
			}
		}
	case ResourceGroup:
		if group, mErr := database.FindGroup(r.ID); mErr.IsError() {
			return false, mErr
		} else {
			if group.OwnerID == r.UserInfo.ID {
				return true, nil
			} else {
				return false, nil
			}
		}
	default:
		return false, mo2errors.Init(mo2errors.Mo2NoExist, "source not available")
	}

}

type AccessFilter struct {
	VisitorID primitive.ObjectID `json:"visitor_id" bson:"visitor_id"`
	GroupID   primitive.ObjectID `json:"group_id" bson:"group_id"`
	RoleList  [][]string         `json:"role_list,omitempty" example:"'admin':xxxxx 'write':xxxxx" bson:"role_map,omitempty"`
}

// JudgeRule
// 判断id是否在manager的某个role(s)之内，
// 第一层，逻辑“或”，满足其中一个role的组合即可
// 第二层，逻辑“与”，必须满足列表内的所有role
func (a *AccessFilter) JudgeRule() (bool, error) {
	var res model.AccessManager
	opt := options.FindOne().SetProjection(bson.M{"_id": 1})

	for _, allowRoles := range a.RoleList {
		pass := true
		for _, role := range allowRoles {
			key := strings.Join([]string{"access_manager", "role_map", role}, ".")
			err := database.GroupCol.FindOne(context.TODO(), bson.M{"$and": []bson.M{{"_id": a.GroupID}, {key: bson.M{"$eq": a.VisitorID}}}}, opt).Decode(&res)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					pass = false
				} else {
					return false, err
				}

			}
		}
		if pass {
			return true, nil
		}
	}
	return false, nil
}
func (a *AccessFilter) ProcessContext(ctx abac.ContextType) {
	if val := ctx.Value(RuleAccessFilter); val != nil {
		*a = val.(AccessFilter)
	}
}
