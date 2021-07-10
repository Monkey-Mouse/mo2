package controller

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/Monkey-Mouse/mo2/services/loghelper"
	"github.com/Monkey-Mouse/mo2log/service/logservice"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (c *Controller) Score(ctx *gin.Context, u dto.LoginUserInfo) (status int, body interface{}, err error) {
	score := &dto.ScoreReq{}
	err = ctx.BindJSON(score)
	if err != nil {
		return 400, nil, err
	}
	b := database.FindBlogById(score.Target, false, &options.FindOneOptions{
		Projection: bson.D{{"content", 0}, {"y_doc", 0}},
	})

	v, err := notifyLogClient.Client.Exist(ctx, &logservice.ExtRequest{
		Operator:        u.ID[:],
		Operation:       loghelper.SCORE,
		OperationTarget: score.Target[:],
	})
	reScore := ""
	if err != nil {
		reScore = "重新"
	}
	notifyLogClient.LogInfo(loghelper.Log{
		Operator:             u.ID,
		Operation:            loghelper.SCORE,
		OperationTarget:      score.Target,
		OperationTargetOwner: b.AuthorID,
		LogLevel:             logservice.LogModel_INFO,
		ExtraMessage:         fmt.Sprintf(`给你的文章<a href="/article/%s">%s</a>%s打分：%f`, score.Target.Hex(), b.Title, reScore, score.Score),
	})
	if err != nil {
		database.ScoreBlog(ctx, &b, -1, float64(score.Score))
	} else {
		ss := strings.Split(v.ExtraMessage, "打分：")
		f, _ := strconv.ParseFloat(ss[1], 64)
		database.ScoreBlog(ctx, &b, f, score.Score)
	}
	return 200, gin.H{"success": true}, nil
}
