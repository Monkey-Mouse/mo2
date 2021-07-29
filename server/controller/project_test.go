package controller

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/Monkey-Mouse/mo2/mo2utils"
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

func TestController_GetProject(t *testing.T) {
	ctx := &gin.Context{
		Params: gin.Params{{Key: "id", Value: primitive.NewObjectID().Hex()}},
	}
	ctx1 := &gin.Context{}
	p := &database.Project{}
	database.UpsertProject(ctx, p, nil)
	defer database.DeleteProject(ctx, p.ID)
	ctx2 := &gin.Context{
		Params: gin.Params{{Key: "id", Value: p.ID.Hex()}},
	}
	type args struct {
		ctx *gin.Context
		u   dto.LoginUserInfo
	}
	tests := []struct {
		name       string
		c          *Controller
		args       args
		wantStatus int
		wantBody   interface{}
		wantErr    bool
	}{
		{name: "test 404", c: &Controller{},
			args: args{
				ctx: ctx,
			}, wantStatus: 404, wantBody: nil, wantErr: true},
		{name: "test 400", c: &Controller{},
			args: args{
				ctx: ctx1,
			}, wantStatus: 400, wantBody: nil, wantErr: true},
		{name: "test 200", c: &Controller{},
			args: args{
				ctx: ctx2,
			}, wantStatus: 200, wantBody: p, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{}
			gotStatus, gotBody, err := c.GetProject(tt.args.ctx, tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("Controller.GetProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("Controller.GetProject() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
			if gotBody != nil {
				p1 := gotBody.(*database.Project)
				p2 := tt.wantBody.(*database.Project)
				p1.EntityInfo = p2.EntityInfo
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("Controller.GetProject() gotBody = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}

func TestController_DeleteProject(t *testing.T) {
	ctx := &gin.Context{
		Params: gin.Params{{Key: "id", Value: primitive.NewObjectID().Hex()}},
	}
	ctx1 := &gin.Context{}
	uid := primitive.NewObjectID()
	p := &database.Project{
		OwnerID: uid,
	}
	database.UpsertProject(ctx, p, nil)
	ctx2 := &gin.Context{
		Params: gin.Params{{Key: "id", Value: p.ID.Hex()}},
	}
	defer database.DeleteProject(ctx, p.ID)
	type args struct {
		ctx *gin.Context
		u   dto.LoginUserInfo
	}
	tests := []struct {
		name       string
		c          *Controller
		args       args
		wantStatus int
		wantBody   interface{}
		wantErr    bool
	}{
		{name: "test 404", c: &Controller{},
			args: args{
				ctx: ctx,
				u:   dto.LoginUserInfo{ID: uid},
			}, wantStatus: 404, wantBody: nil, wantErr: true},
		{name: "test 400", c: &Controller{},
			args: args{
				ctx: ctx1,
			}, wantStatus: 400, wantBody: nil, wantErr: true},
		{name: "test 200", c: &Controller{},
			args: args{
				ctx: ctx2,
				u:   dto.LoginUserInfo{ID: uid},
			}, wantStatus: 200, wantBody: nil, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{}
			gotStatus, _, err := c.DeleteProject(tt.args.ctx, tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("Controller.DeleteProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("Controller.DeleteProject() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestController_JoinProject(t *testing.T) {
	email := "xxxx"
	ctx := &gin.Context{}
	uid := primitive.NewObjectID()
	p := &database.Project{
		OwnerID: uid,
	}
	database.UpsertProject(ctx, p, nil)
	ctx1 := &gin.Context{}

	t1 := mo2utils.GenerateJwtCode(dto.LoginUserInfo{Email: "xx", ID: p.ID, Name: "member_i_ds"})

	token := mo2utils.GenerateJwtCode(dto.LoginUserInfo{Email: email, ID: p.ID, Name: "member_i_ds"})
	ctx2 := &gin.Context{}

	defer database.DeleteProject(ctx, p.ID)
	type args struct {
		ctx   *gin.Context
		u     dto.LoginUserInfo
		token string
	}
	tests := []struct {
		name       string
		c          *Controller
		args       args
		wantStatus int
		wantBody   interface{}
		wantErr    bool
	}{
		{name: "test 400", c: &Controller{},
			args: args{
				ctx:   ctx,
				u:     dto.LoginUserInfo{Email: email},
				token: "",
			}, wantStatus: 400, wantBody: nil, wantErr: true},
		{name: "test 400 wrong email", c: &Controller{},
			args: args{
				ctx:   ctx1,
				u:     dto.LoginUserInfo{Email: email},
				token: t1,
			}, wantStatus: 400, wantBody: nil, wantErr: true},
		{name: "test 200", c: &Controller{},
			args: args{
				ctx:   ctx2,
				u:     dto.LoginUserInfo{Email: email},
				token: token,
			}, wantStatus: 200, wantBody: nil, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pat := gomonkey.ApplyMethod(reflect.TypeOf(tt.args.ctx), "Query", func(*gin.Context, string) string {
				return tt.args.token
			})
			defer pat.Reset()
			c := &Controller{}
			gotStatus, _, err := c.JoinProject(tt.args.ctx, tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("Controller.JoinProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("Controller.JoinProject() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestController_ListProject(t *testing.T) {
	ctx := &gin.Context{}
	patches := gomonkey.ApplyMethod(reflect.TypeOf(ctx), "BindQuery",
		func(_ *gin.Context, obj interface{}) error {
			filter := obj.(*listFilter)
			*filter = listFilter{
				Page:     0,
				PageSize: 10,
				Tags:     []string{},
				Uid:      primitive.NewObjectID().Hex(),
			}
			return nil
		},
	)
	defer patches.Reset()
	type args struct {
		ctx *gin.Context
		u   dto.LoginUserInfo
	}
	tests := []struct {
		name       string
		c          *Controller
		args       args
		wantStatus int
		wantBody   interface{}
		wantErr    bool
	}{
		{name: "test 200", c: &Controller{},
			args: args{
				ctx: ctx,
				u:   dto.LoginUserInfo{ID: primitive.NewObjectID()},
			}, wantStatus: 200, wantBody: []*database.Project{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{}
			gotStatus, gotBody, err := c.ListProject(tt.args.ctx, tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("Controller.ListProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("Controller.ListProject() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("Controller.ListProject() gotBody = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}
