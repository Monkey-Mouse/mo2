package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/Monkey-Mouse/mo2/server/model"
	"github.com/Monkey-Mouse/mo2/services/loghelper"
	"github.com/Monkey-Mouse/mo2log/service/logservice"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	likedResp       = gin.H{"liked": true}
	notLikedResp    = gin.H{"liked": false}
	errWrongType    = errors.New("wrong type")
	okResp          = gin.H{"status": "ok"}
	errMicroService = errors.New("internal microservice error")
	errBadidFormat  = errors.New("bad id format")
)

func (c *Controller) Liked(ctx *gin.Context, u dto.LoginUserInfo) (status int, body interface{}, err error) {
	hex := ctx.Param("id")
	t := ctx.Param("type")
	id, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		return http.StatusUnprocessableEntity, nil, err
	}
	switch t {
	case "blog":
		_, err = notifyLogClient.Client.Exist(ctx,
			&logservice.ExtRequest{
				Operator:        u.ID[:],
				Operation:       loghelper.LIKE_BLOG,
				OperationTarget: id[:],
			},
		)
	default:
		return 422, nil, errWrongType
	}
	if err != nil {
		return 200, notLikedResp, nil
	}
	return 200, likedResp, nil
}

func (c *Controller) Like(ctx *gin.Context, u dto.LoginUserInfo) (status int, body interface{}, err error) {
	t := ctx.Param("type")
	var data struct {
		ID primitive.ObjectID `json:"id"`
	}
	ctx.ShouldBindJSON(&data)

	switch t {
	case "blog":
		var b model.Blog
		_, err := notifyLogClient.Client.Exist(ctx,
			&logservice.ExtRequest{
				Operator:        u.ID[:],
				Operation:       loghelper.LIKE_BLOG,
				OperationTarget: data.ID[:],
			},
		)
		if err == nil {
			bs, err := json.Marshal(bson.D{
				{"operator_id", u.ID},
				{"operation", loghelper.LIKE_BLOG},
				{"operation_target_id", data.ID},
			})
			if err != nil {
				return http.StatusUnprocessableEntity, nil, err
			}
			notifyLogClient.Client.Delete(ctx, &logservice.Query{
				Query: bs,
			})
		} else {
			database.BlogCol.FindOne(ctx,
				bson.M{"_id": data.ID},
				options.FindOne().SetProjection(bson.D{{"author_id", 1}, {"title", 1}})).Decode(&b)
			notifyLogClient.LogInfo(
				loghelper.Log{
					Operator:             u.ID,
					Operation:            loghelper.LIKE_BLOG,
					OperationTarget:      data.ID,
					OperationTargetOwner: b.AuthorID,
					ExtraMessage:         fmt.Sprintf(`点赞了你的文章<a href="/article/%s">%s</a>`, b.ID.Hex(), b.Title),
				},
			)
		}
	default:
		return 422, nil, errWrongType
	}
	body = okResp
	return
}
func (c *Controller) LikeNum(ctx *gin.Context) (status int, body interface{}, err error) {
	t := ctx.Param("type")
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		return http.StatusUnprocessableEntity, nil, errBadidFormat
	}
	num := 0
	switch t {
	case "blog":
		var numP *logservice.Num
		numP, err = notifyLogClient.Client.SendCountQuery(ctx,
			bson.D{
				{"operation", loghelper.LIKE_BLOG},
				{"operation_target_id", id},
			},
		)
		if err != nil {
			return 500, nil, errMicroService
		}
		num = int(numP.Num)
	default:
		return 422, nil, errWrongType
	}
	body = gin.H{"num": num}
	return
}
