package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"mo2/server/model"
	"testing"
	"time"
)

func upsertRecycleItem4Test(num int) (ids []primitive.ObjectID) {
	ids = make([]primitive.ObjectID, num)
	for i := 0; i < num; i++ {
		ids[i] = primitive.NewObjectID()
		UpsertRecycleItem(model.RecycleItem{ID: ids[i]})
	}
	return
}

func upsertBlogs4Test(num int, isDraft bool) (ids []primitive.ObjectID) {
	ids = make([]primitive.ObjectID, num)
	for i := 0; i < num; i++ {
		ids[i] = primitive.NewObjectID()
		upsertBlog(&model.Blog{ID: ids[i]}, isDraft)
	}
	return
}

func TestDeleteExpireItems(t *testing.T) {

	ids := upsertRecycleItem4Test(3)
	blogID := primitive.NewObjectID()
	upsertBlog(&model.Blog{ID: blogID}, false)
	deleteOne := model.RecycleItem{ID: ids[0], DeleteTime: time.Now()}
	notDeleteOne := model.RecycleItem{ID: ids[1], DeleteTime: time.Date(3000, 5, 5, 5, 5, 5, 5, time.UTC)}
	deleteBlog := model.RecycleItem{ID: ids[2], ItemID: blogID, DeleteTime: time.Now(), Handler: model.HandlerBlog}

	defer func() {
		DeleteRecycleItems(ids...)
		DeleteBlogs(false, blogID)
	}()
	tests := []struct {
		name        string
		item        model.RecycleItem
		wantDeleted bool
	}{
		{"delete", deleteOne, true},
		{"deleteBlog", deleteBlog, true},
		{"notDelete", notDeleteOne, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UpsertRecycleItem(tt.item)
			if gotMErr := DeleteExpireItems(); gotMErr.IsError() {
				t.Errorf("DeleteExpireItems()  = %v", gotMErr)
			} else {

				if tt.wantDeleted && !tt.item.ItemID.IsZero() {
					var err error

					switch tt.item.Handler {
					case model.HandlerBlog:
						var res model.Blog
						err = blogCol.FindOne(context.TODO(), bson.M{"_id": tt.item.ItemID}).Decode(&res)
					case model.HandlerDraft:
						var res model.Blog
						err = draftCol.FindOne(context.TODO(), bson.M{"_id": tt.item.ItemID}).Decode(&res)
					default:
						log.Println("invalid handler")
					}
					if err != mongo.ErrNoDocuments {
						t.Errorf("%v not deleted", tt.item.Handler)
					}
				}

			}
		})
	}
}

func TestDeleteItems(t *testing.T) {
	num := 2
	ids := upsertRecycleItem4Test(num)
	defer DeleteRecycleItems(ids...)
	tests := []struct {
		name string
		ids  []primitive.ObjectID
	}{
		{"numEqual", ids},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var items []model.RecycleItem
			if cursor, err := binCol.Find(context.TODO(), bson.M{"_id": bson.M{"$in": tt.ids}}); err != nil {
				t.Error(err)
			} else {
				if err = cursor.All(context.TODO(), &items); err != nil {
					t.Error(err)
				} else {
					if len(ids) != len(items) {
						t.Errorf("upsert %v not equal to find %v", len(ids), len(items))
					}
				}
			}

			if gotMErr := DeleteRecycleItems(tt.ids...); gotMErr.IsError() {
				t.Errorf("DeleteRecycleItems() = %v", gotMErr)
			} else {

				if cursor, err := binCol.Find(context.TODO(), bson.M{"_id": bson.M{"$in": tt.ids}}); err != nil {
					t.Error(err)
				} else {
					if err = cursor.All(context.TODO(), &items); err != nil {
						t.Error(err)
					} else {
						if len(items) != 0 {
							t.Errorf("still exist %v after deleted", len(items))
						}
					}
				}

			}

		})
	}
}

func TestUpsertRecycleItem(t *testing.T) {
	ids := upsertRecycleItem4Test(1)
	newId := primitive.NewObjectID()
	defer DeleteRecycleItems(ids...)
	defer DeleteRecycleItems(newId)
	foo := model.RecycleItem{
		ID:         primitive.NewObjectID(),
		ItemID:     primitive.NewObjectID(),
		CreateTime: time.Now(),
		DeleteTime: time.Now(),
		Handler:    "foo",
	}
	defer DeleteRecycleItems(foo.ID)
	tests := []struct {
		name string
		item model.RecycleItem
	}{
		{"insert", model.RecycleItem{ID: newId}},
		{"upsert", model.RecycleItem{ID: ids[0]}},
		{"foo", foo},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if gotMErr := UpsertRecycleItem(tt.item); gotMErr.IsError() {
				t.Errorf("UpsertRecycleItem() = %v", gotMErr)
			}
			var result model.RecycleItem
			if err := binCol.FindOne(context.TODO(), bson.M{"_id": tt.item.ID}).Decode(&result); err != nil {
				t.Errorf("error: %v", err)
			} else {
				if result.ID != tt.item.ID {
					t.Error("not insert successfully")
				}
			}
		})
	}
}
