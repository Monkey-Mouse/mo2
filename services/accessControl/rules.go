package accessControl

import (
	"github.com/Monkey-Mouse/go-abac/abac"
	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/Monkey-Mouse/mo2/server/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AllowOwn struct {
	userInfo dto.LoginUserInfo
	id       primitive.ObjectID `json:"id"`
	filter   model.Filter
}

func (r *AllowOwn) ProcessContext(ctx abac.ContextType) {
	r.userInfo = ctx.Value("userInfo").(dto.LoginUserInfo)
	r.id = ctx.Value("id").(primitive.ObjectID)
	r.filter = ctx.Value("filter").(model.Filter)
}
func (r *AllowOwn) JudgeRule() (bool, error) {
	blog := database.FindBlogById(r.id, r.filter.IsDraft)
	if blog.AuthorID == r.userInfo.ID {
		return true, nil
	} else {
		return false, nil
	}
}
