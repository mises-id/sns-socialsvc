package blacklist

import (
	"context"
	"testing"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"github.com/mises-id/sns-socialsvc/lib/storage"
)

func cleanTables(names ...string) {
	for _, name := range names {
		db.DB().Collection(name).Drop(context.TODO())
	}
}

func TestListBlacklist(t *testing.T) {
	db.SetupMongo(context.TODO())
	storage.SetupImageStorage("127.0.0.1", "xxx", "xx")
	defer cleanTables("users", "blacklists")
	for i := 0; i < 10; i++ {
		db.ODM(context.Background()).Create(&models.User{
			UID: uint64(i + 1),
		})
	}
	blacklistMap := map[uint64][]uint64{
		1: {2, 3, 4, 5, 6, 7, 8},
		2: {4},
		4: {3},
	}
	for uid, targetIDs := range blacklistMap {
		for _, targetID := range targetIDs {
			db.ODM(context.Background()).Create(&models.Blacklist{
				UID:       uid,
				TargetUID: targetID,
			})
		}
	}

	t.Run("list blacklist first page", func(t *testing.T) {
		blacklists, page, err := ListBlacklist(context.TODO(), &ListBlacklistParams{UID: 1, PageParams: &pagination.PageQuickParams{
			Limit: 5,
		}})
		if err != nil {
			t.Error(err)
		}
		if len(blacklists) != 5 {
			t.Errorf("list blacklist wrong result count: %d ", len(blacklists))
		}
		if page.GetPageSize() != 5 {
			t.Errorf("list blacklist wrong page size: %d ", page.GetPageSize())
		}
		quickPage := page.BuildJSONResult().(*pagination.QuickPagination)

		if quickPage.NextID == "" {
			t.Errorf("list blacklist wrong next id: %s ", quickPage.NextID)
		}
	})

	t.Run("list blacklist with page", func(t *testing.T) {
		_, page, err := ListBlacklist(context.TODO(), &ListBlacklistParams{UID: 1, PageParams: &pagination.PageQuickParams{
			Limit: 5,
		}})
		if err != nil {
			t.Error(err)
		}
		quickPage := page.BuildJSONResult().(*pagination.QuickPagination)
		_, page, err = ListBlacklist(context.TODO(), &ListBlacklistParams{UID: 1, PageParams: &pagination.PageQuickParams{
			Limit:  5,
			NextID: quickPage.NextID,
		}})
		if err != nil {
			t.Error(err)
		}
		quickPage = page.BuildJSONResult().(*pagination.QuickPagination)
		if quickPage.NextID != "" {
			t.Errorf("list blacklist wrong next id: %s ", quickPage.NextID)
		}
	})
}

func TestCreateBlacklist(t *testing.T) {
	db.SetupMongo(context.TODO())
	storage.SetupImageStorage("127.0.0.1", "xxx", "xx")
	defer cleanTables("users", "blacklists")
	for i := 0; i < 10; i++ {
		db.ODM(context.Background()).Create(&models.User{
			UID: uint64(i + 1),
		})
	}
	blacklistMap := map[uint64][]uint64{
		1: {2, 3, 4, 5, 6, 7, 8},
		2: {4},
		4: {3},
	}
	for uid, targetIDs := range blacklistMap {
		for _, targetID := range targetIDs {
			db.ODM(context.Background()).Create(&models.Blacklist{
				UID:       uid,
				TargetUID: targetID,
			})
		}
	}

	t.Run("create exsist blacklist", func(t *testing.T) {
		_, err := CreateBlacklist(context.TODO(), 1, 2)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("create new blacklist", func(t *testing.T) {
		_, err := CreateBlacklist(context.TODO(), 2, 3)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestDeleteBlacklist(t *testing.T) {
	t.Run("delete exsist blacklist", func(t *testing.T) {

	})

	t.Run("delete non blacklist", func(t *testing.T) {

	})
}
