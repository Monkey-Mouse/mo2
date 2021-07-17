package database

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateProject(t *testing.T) {
	ProjCol.Drop(context.TODO())
	p := &Project{}
	_, err := UpsertProject(context.TODO(), p, nil)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	if p.ID.IsZero() {
		t.Errorf("id of return value should not be nil")
	}
}

func TestListProject(t *testing.T) {
	_, err := ListProject(context.TODO(), 0, 10, bson.D{})
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
}
