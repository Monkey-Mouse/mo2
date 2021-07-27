package controller

import (
	"fmt"

	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/Monkey-Mouse/mo2/mo2utils"
	emailservice "github.com/Monkey-Mouse/mo2/services/emailService"
	"github.com/gin-gonic/gin"
	"github.com/willf/bloom"
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
		p.ManagerIDs = []primitive.ObjectID{}
		p.MemberIDs = []primitive.ObjectID{}
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
	invite := false
	if !mo2utils.Contains(prev.Tags, "课程") {
		filter := bloom.NewWithEstimates(100, 0.01)
		managers := make([]primitive.ObjectID, 0, len(prev.ManagerIDs))
		members := make([]primitive.ObjectID, 0, len(prev.MemberIDs))
		invitations := make([]primitive.ObjectID, 0)
		meminvitations := make([]primitive.ObjectID, 0)

		for _, v := range prev.ManagerIDs {
			filter.Add(v[:])
		}
		for _, v := range prev.MemberIDs {
			filter.Add(v[:])
		}
		for _, v := range p.ManagerIDs {
			if !filter.Test(v[:]) {
				invitations = append(invitations, v)
				invite = true
			} else {
				managers = append(managers, v)
			}
		}
		p.ManagerIDs = managers
		for _, v := range p.MemberIDs {
			if !filter.Test(v[:]) {
				meminvitations = append(meminvitations, v)
				invite = true
			} else {
				members = append(members, v)
			}
		}
		p.MemberIDs = members
		if len(invitations) > 0 {
			infos := database.ListAccountsBrief(invitations)
			if len(infos) > 0 {
				for _, v := range infos {
					url := "https://www.motwo.cn/project/" + p.ID.Hex()
					token := mo2utils.GenerateJwtCode(dto.LoginUserInfo{Email: v.Email, ID: p.ID, Name: "manager_i_ds"})
					emailservice.SendEmail(emailservice.InvitationMessage(url+"?token="+token+"&email="+v.Email, p.Name, []string{v.Email}), ctx.ClientIP())
				}
			}
		}
		if len(meminvitations) > 0 {
			infos := database.ListAccountsBrief(meminvitations)
			if len(infos) > 0 {
				for _, v := range infos {
					url := "https://www.motwo.cn/project/" + p.ID.Hex()
					token := mo2utils.GenerateJwtCode(dto.LoginUserInfo{Email: v.Email, ID: p.ID, Name: "member_i_ds"})
					emailservice.SendEmail(emailservice.InvitationMessage(url+"?token="+token+"&email="+v.Email, p.Name, []string{v.Email}), ctx.ClientIP())
				}
			}
		}
	}
	_, err = database.UpsertProject(ctx, p, nil)
	if err != nil {
		return 500, nil, err
	}
	return 200, gin.H{"invite": invite, "project": p}, nil
}

func (c *Controller) JoinProject(ctx *gin.Context, u dto.LoginUserInfo) (status int, body interface{}, err error) {
	token := ctx.Query("token")
	info, err := mo2utils.ParseJwt(token)
	if err != nil {
		return 400, nil, err
	}
	if info.Email != u.Email {
		return 400, nil, fmt.Errorf("wrong email")
	}
	_, err = database.UpsertProject(ctx, &database.Project{ID: info.ID}, bson.M{"$addToSet": bson.M{info.Name: u.ID}})
	if err != nil {
		return 500, nil, err
	}
	p, err := database.GetProject(ctx, bson.M{"_id": info.ID})
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
