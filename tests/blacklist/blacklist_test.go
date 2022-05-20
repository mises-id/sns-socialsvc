//go:build tests
// +build tests

package blacklist

import (
	"context"
	"testing"

	"github.com/mises-id/sns-socialsvc/app/models"
	service "github.com/mises-id/sns-socialsvc/app/services/blacklist"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"github.com/mises-id/sns-socialsvc/tests"
	"github.com/stretchr/testify/suite"
)

type BlacklistServerSuite struct {
	tests.BaseTestSuite
	statuses []*models.Status
}

func (suite *BlacklistServerSuite) SetupSuite() {
	suite.Collections = []string{"users", "blacklists"}

	suite.BaseTestSuite.SetupSuite()
	suite.CreateTestUsers(10)

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

}

func (suite *BlacklistServerSuite) TearDownSuite() {
	suite.BaseTestSuite.TearDownSuite()
}

func TestBlacklistServer(t *testing.T) {
	suite.Run(t, &BlacklistServerSuite{})
}

func (suite *BlacklistServerSuite) TestListBlacklist() {

	suite.T().Run("list blacklist first page", func(t *testing.T) {
		blacklists, page, err := service.ListBlacklist(context.TODO(), &service.ListBlacklistParams{UID: 1, PageParams: &pagination.PageQuickParams{
			Limit: 5,
		}})
		suite.Nil(err)
		suite.EqualValues(len(blacklists), 5)
		suite.EqualValues(page.GetPageSize(), 5)
		quickPage := page.BuildJSONResult().(*pagination.QuickPagination)
		suite.NotEmpty(quickPage.NextID)
	})

	suite.T().Run("list blacklist with page", func(t *testing.T) {
		_, page, err := service.ListBlacklist(context.TODO(), &service.ListBlacklistParams{UID: 1, PageParams: &pagination.PageQuickParams{
			Limit: 5,
		}})
		suite.Nil(err)
		quickPage := page.BuildJSONResult().(*pagination.QuickPagination)
		suite.NotEmpty(quickPage.NextID)
		_, page, err = service.ListBlacklist(context.TODO(), &service.ListBlacklistParams{UID: 1, PageParams: &pagination.PageQuickParams{
			Limit:  5,
			NextID: quickPage.NextID,
		}})
		suite.Nil(err)
		quickPage = page.BuildJSONResult().(*pagination.QuickPagination)
		suite.Empty(quickPage.NextID)
	})
}

func (suite *BlacklistServerSuite) TestCreateBlacklist() {

	suite.T().Run("create exsist blacklist", func(t *testing.T) {
		_, err := service.CreateBlacklist(context.TODO(), 1, 2)
		suite.Nil(err)
	})

	suite.T().Run("create new blacklist", func(t *testing.T) {
		_, err := service.CreateBlacklist(context.TODO(), 2, 3)
		suite.Nil(err)
	})
}

func (suite *BlacklistServerSuite) TestDeleteBlacklist() {
	suite.T().Run("delete exsist blacklist", func(t *testing.T) {

	})

	suite.T().Run("delete non blacklist", func(t *testing.T) {

	})
}
