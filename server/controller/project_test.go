package controller

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/agiledragon/gomonkey"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		patches.ApplyFunc(database.UpsertProject,
			func(ctx context.Context,
				p *database.Project, update bson.M) (*mongo.UpdateResult, error) {
				return nil, nil
			},
		)
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
		p := b.(*database.Project)
		if len(p.ManagerIDs) > 0 {
			t.Errorf("insert project should not have managerids")
		}
		if len(p.MemberIDs) > 0 {
			t.Errorf("insert project should not have memberids")
		}
		if p.OwnerID != id {
			t.Errorf("ownerid %v should equal to userid %v", p.OwnerID, id)
		}

	})
}
