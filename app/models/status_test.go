package models

import (
	"context"
	"testing"
	"time"

	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"github.com/mises-id/sns-socialsvc/lib/storage"
)

func TestCreateStatus(t *testing.T) {
	storage.SetupImageStorage("127.0.0.1", "xxx", "xx")
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

func TestListShowStatus(t *testing.T) {
	storage.SetupImageStorage("127.0.0.1", "xxx", "xx")
	db.SetupMongo(context.TODO())
	for i := 0; i < 20; i++ {
		db.ODM(context.Background()).Create(&Status{
			UID:     uint64(1),
			Content: "test status 1",
		})
	}
	hideTime := time.Now().Add(time.Second * -10)
	for i := 0; i < 20; i++ {
		db.ODM(context.Background()).Create(&Status{
			UID:      uint64(1),
			Content:  "test status 2",
			HideTime: &hideTime,
		})
	}
	hideTime = hideTime.AddDate(0, 0, 1)
	for i := 0; i < 20; i++ {
		db.ODM(context.Background()).Create(&Status{
			UID:      uint64(1),
			Content:  "test status 3",
			HideTime: &hideTime,
		})
	}
	defer db.DB().Collection("statuses").Drop(context.TODO())
	statuses, _, err := ListStatus(context.TODO(), &ListStatusParams{
		UIDs:     []uint64{1},
		OnlyShow: true,
		PageParams: &pagination.PageQuickParams{
			Limit: 50,
		},
	})

	if err != nil {
		t.Error(err)
	}
	if len(statuses) != 40 {
		t.Errorf("list status wrong result count: %d ", len(statuses))
	}
}
