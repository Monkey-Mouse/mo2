package controller

import (
	"fmt"

	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *Controller) UpsertProject(ctx *gin.Context, u dto.LoginUserInfo) (status int, body interface{}, err error) {
	p := &database.Project{}
	err = ctx.BindJSON(&p)
	if err != nil {
		return 400, nil, err
	}
	if p.ID.IsZero() {
		p.OwnerID = u.ID
		_, err = database.UpsertProject(ctx, p, nil)
		if err != nil {
			return 500, nil, err
		}
		return 200, p, nil
	}
	// 更新鉴权
	prev, err := database.GetProject(ctx, bson.M{"_id": p.ID, "$or": []bson.M{{"manager_i_ds": u.ID}, {"owner_id": u.ID}}})
	if err != nil {
		return 403, nil, fmt.Errorf("access denied")
	}
	p.OwnerID = prev.OwnerID
	_, err = database.UpsertProject(ctx, p, nil)
	if err != nil {
		return 500, nil, err
	}
	return 200, p, nil
}

func (c *Controller) ListProject(ctx *gin.Context, u dto.LoginUserInfo) (status int, body interface{}, err error) {
	var filter struct {
		Page     int64
		PageSize int64
		Tags     []string
		Uid      string
	}
	err = ctx.BindQuery(&filter)
	if err != nil {
		return 400, nil, err
	}
	id, ierr := primitive.ObjectIDFromHex(filter.Uid)
	exfilter := bson.M{}
	if ierr == nil {
		exfilter = bson.M{
			"$or": []bson.M{
				{"manager_i_ds": id},
				{"owner_id": id},
				{"member_i_ds": id},
			},
		}
	}
	if len(filter.Tags) != 0 {
		exfilter["tags"] = bson.M{
			"$all": filter.Tags,
		}
	}
	ps, err := database.ListProject(ctx, filter.Page, filter.PageSize, exfilter)
	if err != nil {
		return 500, nil, err
	}
	return 200, ps, nil

}

func (c *Controller) GetProject(ctx *gin.Context, u dto.LoginUserInfo) (status int, body interface{}, err error) {
	sid, _ := ctx.Params.Get("id")
	id, err := primitive.ObjectIDFromHex(sid)
	if err != nil {
		return 400, nil, err
	}
	p, err := database.GetProject(ctx, bson.M{"_id": id})
	if err != nil {
		return 404, nil, fmt.Errorf("not found")
	}
	return 200, p, nil
}
func (c *Controller) DeleteProject(ctx *gin.Context, u dto.LoginUserInfo) (status int, body interface{}, err error) {
	sid, _ := ctx.Params.Get("id")
	id, err := primitive.ObjectIDFromHex(sid)
	if err != nil {
		return 400, nil, err
	}
	_, err = database.GetProject(ctx, bson.M{"_id": id, "owner_id": u.ID})
	if err != nil {
		return 403, nil, fmt.Errorf("access denied")
	}
	p, err := database.DeleteProject(ctx, id)
	if err != nil {
		return 404, nil, fmt.Errorf("not found")
	}
	return 200, p, nil

}
