package models

import (
	"context"
	"testing"

	"github.com/mises-id/sns-socialsvc/lib/db"
)

func TestCreateStatus(t *testing.T) {
	db.SetupMongo(context.TODO())
	defer db.DB().Collection("statuses").Drop(context.TODO())
	status, err := CreateStatus(context.TODO(), &CreateStatusParams{
		UID:     uint64(1),
		Content: "test status",
	})
	if err != nil {
		t.Error(err)
	}
	if status.Content != "test status" {
		t.Errorf("status content = %s; expected %s", status.Content, "test status")
	}
}
