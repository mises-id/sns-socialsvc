//go:build tests
// +build tests

package nft

import (
	"context"
	"fmt"
	"testing"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/search"
	service "github.com/mises-id/sns-socialsvc/app/services/nft"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"github.com/mises-id/sns-socialsvc/tests"
	"github.com/stretchr/testify/suite"
)

type NftServerSuite struct {
	tests.BaseTestSuite
	nfts []*models.NftAsset
}

func (suite *NftServerSuite) SetupSuite() {
	suite.Collections = []string{"users", "nftassets", "likes"}

	suite.BaseTestSuite.SetupSuite()
	suite.CreateTestUsers(10)

	nftMap := map[uint64][]int64{
		1: {2, 3, 4, 5, 6, 7, 8},
		2: {4},
		4: {3},
	}
	for uid, targetIDs := range nftMap {
		for _, targetID := range targetIDs {

			contract := models.AssetContract{
				Address: fmt.Sprintf("0xaddress%d", targetID),
			}
			ass := models.Asset{
				UID:           uid,
				AssetId:       targetID,
				AssetContract: &contract,
			}
			nft := &models.NftAsset{
				Asset: ass,
			}
			db.ODM(context.Background()).Create(nft)
			suite.nfts = append(suite.nfts, nft)
		}
	}

}

func (suite *NftServerSuite) TearDownSuite() {
	suite.BaseTestSuite.TearDownSuite()
}

func TestNftServer(t *testing.T) {
	suite.Run(t, &NftServerSuite{})
}

func (suite *NftServerSuite) TestListNft() {

	suite.T().Run("list nft first page", func(t *testing.T) {
		nfts, page, err := service.PageNftAsset(context.TODO(), 1, &service.NftAssetInput{NftAssetSearch: &search.NftAssetSearch{
			UID: 1,
			PageParams: &pagination.PageQuickParams{
				Limit: 5,
			},
		}})
		suite.Nil(err)
		suite.EqualValues(len(nfts), 5)
		suite.EqualValues(page.GetPageSize(), 5)
		quickPage := page.BuildJSONResult().(*pagination.QuickPagination)
		suite.NotEmpty(quickPage.NextID)
	})

	suite.T().Run("list nft with page", func(t *testing.T) {
		_, page, err := service.PageNftAsset(context.TODO(), 1, &service.NftAssetInput{NftAssetSearch: &search.NftAssetSearch{
			UID: 1,
			PageParams: &pagination.PageQuickParams{
				Limit: 5,
			},
		}})
		suite.Nil(err)
		quickPage := page.BuildJSONResult().(*pagination.QuickPagination)
		suite.NotEmpty(quickPage.NextID)
		_, page, err = service.PageNftAsset(context.TODO(), 1, &service.NftAssetInput{NftAssetSearch: &search.NftAssetSearch{
			UID: 1,
			PageParams: &pagination.PageQuickParams{
				Limit:  5,
				NextID: quickPage.NextID,
			},
		}})
		suite.Nil(err)
		quickPage = page.BuildJSONResult().(*pagination.QuickPagination)
		suite.Empty(quickPage.NextID)
	})
}

func (suite *NftServerSuite) TestLikeNft() {
	suite.T().Run("like unlike nft", func(t *testing.T) {
		_, err := service.LikeNftAsset(context.TODO(), 2, suite.nfts[0].ID)
		suite.Nil(err)
		likes, _, err := service.PageLike(context.TODO(), 2, &service.PageLikeParams{LikeSearch: &models.LikeSearch{TargetID: suite.nfts[0].ID, TargetType: enum.LikeNft}})
		suite.Nil(err)
		suite.EqualValues(1, len(likes))
		asset, err := service.FindNftAsset(context.TODO(), 2, suite.nfts[0].ID)
		suite.Nil(err)
		suite.EqualValues(1, asset.LikesCount)

		err = service.UnlikeNftAsset(context.TODO(), 2, suite.nfts[0].ID)
		suite.Nil(err)
		likes1, _, err := service.PageLike(context.TODO(), 2, &service.PageLikeParams{LikeSearch: &models.LikeSearch{TargetID: suite.nfts[0].ID, TargetType: enum.LikeNft}})
		suite.Nil(err)
		suite.EqualValues(0, len(likes1))

		asset1, err := service.FindNftAsset(context.TODO(), 2, suite.nfts[0].ID)
		suite.Nil(err)
		suite.EqualValues(0, asset1.LikesCount)
	})
}
