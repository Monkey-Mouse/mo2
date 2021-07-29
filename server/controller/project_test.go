package controller

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/Monkey-Mouse/mo2/mo2utils/mo2errors"
	"github.com/Monkey-Mouse/mo2/server/model"
	emailservice "github.com/Monkey-Mouse/mo2/services/emailService"
	"github.com/agiledragon/gomonkey"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func buildUpsertProjectCtx(
	proj database.Project, err error) (*gomonkey.Patches, *gin.Context) {
	ctx := &gin.Context{}

	patches := gomonkey.ApplyMethod(reflect.TypeOf(ctx), "BindJSON", func(_ *gin.Context, obj interface{}) error {
		p := obj.(*database.Project)
		*p = proj
		return err
	})
	return patches, ctx
}

func TestController_UpsertProject(t *testing.T) {
	t.Run("test 400", func(t *testing.T) {
		patches, ctx := buildUpsertProjectCtx(
			database.Project{},
			fmt.Errorf("test err"))
		defer patches.Reset()
		s, _, err := NewController().UpsertProject(ctx, dto.LoginUserInfo{})
		if err == nil {
			t.Error("should throw err but not")
		}
		if s != 400 {
			t.Errorf("should produce status 400, but %v", s)
		}
	})
	t.Run("test insert", func(t *testing.T) {
		patches, ctx := buildUpsertProjectCtx(
			database.Project{
				ManagerIDs: []primitive.ObjectID{primitive.NewObjectID()},
				MemberIDs:  []primitive.ObjectID{primitive.NewObjectID()},
			},
			nil)
		defer patches.Reset()
		id := primitive.NewObjectID()
		s, b, err := NewController().UpsertProject(ctx, dto.LoginUserInfo{
			ID: id,
		})
		if err != nil {
			t.Error("should not throw err but throw " + err.Error())
		}
		if s != 200 {
			t.Errorf("should produce status 200, but %v", s)
		}
		m := b.(gin.H)
		p := m["project"].(*database.Project)
		if len(p.ManagerIDs) > 0 {
			t.Errorf("insert project should not have managerids")
		}
		if len(p.MemberIDs) > 0 {
			t.Errorf("insert project should not have memberids")
		}
		if p.OwnerID != id {
			t.Errorf("ownerid %v should equal to userid %v", p.OwnerID, id)
		}
		_, err = database.DeleteProject(ctx, p.ID)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("test 404", func(t *testing.T) {
		patches, ctx := buildUpsertProjectCtx(
			database.Project{
				ID: primitive.NewObjectID(),
			},
			nil)
		defer patches.Reset()
		id := primitive.NewObjectID()
		s, _, err := NewController().UpsertProject(ctx, dto.LoginUserInfo{
			ID: id,
		})
		if err == nil {
			t.Error("should throw err but didn't ")
		}
		if s != 404 {
			t.Errorf("should produce status 404, but %v", s)
		}
	})
	t.Run("test update  group", func(t *testing.T) {
		database.AccCol.Drop(context.TODO())
		database.ProjCol.Drop(context.TODO())
		uid := primitive.NewObjectID()
		proj := &database.Project{
			Tags:    []string{"xxx", "课x"},
			OwnerID: uid,
		}
		email := "xxx@xx"
		err := database.CreateActiveAccounts(context.TODO(), []model.Account{
			{ID: uid, Email: email, UserName: "xxx", HashedPwd: "xxxxx"},
		})
		if err != nil {
			t.Errorf("insert user data err %v", err)
		}
		_, err = database.UpsertProject(context.TODO(), proj, nil)
		if err != nil {
			t.Errorf("upsert data %v err %v", proj, err)
		}
		proj.ManagerIDs = []primitive.ObjectID{}
		patches, ctx := buildUpsertProjectCtx(
			database.Project{
				ID:         proj.ID,
				Tags:       []string{"xxx", "课程"},
				ManagerIDs: []primitive.ObjectID{uid},
				MemberIDs:  []primitive.ObjectID{primitive.NewObjectID()},
				OwnerID:    uid,
			},
			nil)
		patches.ApplyMethod(reflect.TypeOf(ctx), "ClientIP", func(*gin.Context) string {
			return "xxxx"
		})
		receivers := []string{}
		patches.ApplyFunc(emailservice.SendEmail, func(mail *emailservice.Mo2Email,
			senderAddr string) (err *mo2errors.Mo2Errors) {
			receivers = mail.Receivers
			return nil
		})
		defer patches.Reset()
		s, b, err := NewController().UpsertProject(ctx, dto.LoginUserInfo{
			ID: uid,
		})
		if err != nil {
			t.Error("should not throw err but throw " + err.Error())
		}
		if s != 200 {
			t.Errorf("should produce status 200, but %v", s)
		}
		m := b.(gin.H)
		p := m["project"].(*database.Project)
		if len(p.ManagerIDs) > 0 {
			t.Errorf("insert project should not have managerids")
		}
		if len(p.MemberIDs) > 0 {
			t.Errorf("insert project should not have memberids")
		}
		if len(receivers) != 1 {
			t.Errorf("receivers should contain 1 email, but %v", len(receivers))
		}
		if receivers[0] != email {
			t.Errorf("receivers should only contain %v, but %v", email, receivers)
		}
	})
}
