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

func checkOwnerOrManager(p *database.Project, u dto.LoginUserInfo) bool {
	if p.OwnerID == u.ID {
		return true
	}
	for _, v := range p.ManagerIDs {
		if v == u.ID {
			return true
		}
	}
	return false
}

func (c *Controller) Score(ctx *gin.Context, u dto.LoginUserInfo) (status int, body interface{}, err error) {
	score := &dto.ScoreReq{}
	err = ctx.BindJSON(score)
	if err != nil {
		return 400, nil, err
	}
	b := database.FindBlogById(score.Target, false, &options.FindOneOptions{
		Projection: bson.D{{"content", 0}, {"y_doc", 0}},
	})
	if !b.ProjectID.IsZero() {
		p, err := database.GetProject(ctx, bson.M{"_id": b.ProjectID, "tags": "课程"})
		if err == nil && !checkOwnerOrManager(p, u) {
			return 403, nil, fmt.Errorf("该文章只有所属项目的管理员有权利打分")
		}
	}

	v, err := notifyLogClient.Client.Exist(ctx, &logservice.ExtRequest{
		Operator:        u.ID[:],
		Operation:       loghelper.SCORE,
		OperationTarget: score.Target[:],
	})
	reScore := "重新"
	if err != nil {
		reScore = ""
	}
	notifyLogClient.LogInfo(loghelper.Log{
		Operator:             u.ID,
		Operation:            loghelper.SCORE,
		OperationTarget:      score.Target,
		OperationTargetOwner: b.AuthorID,
		LogLevel:             logservice.LogModel_INFO,
		ExtraMessage:         fmt.Sprintf(`给你的文章<a href="/article/%s">%s</a>%s打分：%.1f`, score.Target.Hex(), b.Title, reScore, score.Score),
	})
	sum, num := 0.0, 0
	if err != nil {
		sum, num = database.ScoreBlog(ctx, &b, -1, float64(score.Score))
	} else {
		ss := strings.Split(v.ExtraMessage, "打分：")
		f, _ := strconv.ParseFloat(ss[1], 64)
		sum, num = database.ScoreBlog(ctx, &b, f, score.Score)
	}
	return 200, gin.H{"sum": sum, "num": num}, nil
}
func (c *Controller) IsScored(ctx *gin.Context, u dto.LoginUserInfo) (status int, body interface{}, err error) {
	score := &dto.ScoreReq{}
	err = ctx.BindJSON(score)
	if err != nil {
		return 400, nil, err
	}
	_, err = notifyLogClient.Client.Exist(ctx, &logservice.ExtRequest{
		Operator:        u.ID[:],
		Operation:       loghelper.SCORE,
		OperationTarget: score.Target[:],
	})
	return 200, err == nil, nil
}
