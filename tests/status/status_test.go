//go:build tests
// +build tests

package status

import (
	"context"
	"testing"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models"
	service "github.com/mises-id/sns-socialsvc/app/services/status"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"github.com/mises-id/sns-socialsvc/tests"
	"github.com/stretchr/testify/suite"
)

type StatusServerSuite struct {
	tests.BaseTestSuite

	statuses []*models.Status
}

func (suite *StatusServerSuite) SetupSuite() {
	suite.Collections = []string{"counters", "users", "follows", "statuses", "likes"}

	suite.BaseTestSuite.SetupSuite()
	suite.CreateTestUsers(10)
	for i := 0; i < 20; i++ {
		db.ODM(context.Background()).Create(&models.Status{
			UID:     uint64(1),
			Content: "test status 1",
		})
	}
	hideTime := time.Now().Add(time.Second * -10)
	for i := 0; i < 20; i++ {
		db.ODM(context.Background()).Create(&models.Status{
			UID:      uint64(1),
			Content:  "test status 2",
			HideTime: &hideTime,
		})
	}
	hideTime = hideTime.AddDate(0, 0, 1)
	for i := 0; i < 20; i++ {
		db.ODM(context.Background()).Create(&models.Status{
			UID:      uint64(1),
			Content:  "test status 3",
			HideTime: &hideTime,
		})
	}

}

func (suite *StatusServerSuite) TearDownSuite() {
	suite.BaseTestSuite.TearDownSuite()
}

func TestStatusServer(t *testing.T) {
	suite.Run(t, &StatusServerSuite{})
}

func (suite *StatusServerSuite) TestCreateStatus() {
	status, err := service.CreateStatus(context.TODO(), 2, &service.CreateStatusParams{
		StatusType: "text",
		Content:    "test status",
	})
	suite.Nil(err)
	suite.Equal(status.Content, "test status")
}

func (suite *StatusServerSuite) TestListShowStatus() {
	statuses, _, err := models.ListStatus(context.TODO(), &models.ListStatusParams{
		UIDs:     []uint64{1},
		OnlyShow: true,
		PageParams: &pagination.PageQuickParams{
			Limit: 50,
		},
	})
	suite.Nil(err)
	suite.Equal(40, len(statuses))
}
