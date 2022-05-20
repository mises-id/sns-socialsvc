//go:build tests
// +build tests

package comment

import (
	"context"
	"testing"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	service "github.com/mises-id/sns-socialsvc/app/services/comment"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/tests"
	"github.com/stretchr/testify/suite"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentServerSuite struct {
	tests.BaseTestSuite
	statuses []*models.Status
}

func (suite *CommentServerSuite) SetupSuite() {
	suite.Collections = []string{"users", "statuses", "comments"}

	suite.BaseTestSuite.SetupSuite()
	suite.CreateTestUsers(10)

}

func (suite *CommentServerSuite) TearDownSuite() {
	suite.BaseTestSuite.TearDownSuite()
}

func TestCommentServer(t *testing.T) {
	suite.Run(t, &CommentServerSuite{})
}

func (suite *CommentServerSuite) TestCreateComment() {

	status := &models.Status{
		UID:           5,
		FromType:      enum.FromPost,
		StatusType:    enum.TextStatus,
		CommentsCount: 2,
		Content:       "hello",
		IsPublic:      true,
	}
	db.ODM(context.Background()).Create(status)
	firstComment := &models.Comment{
		StatusID: status.ID,
		UID:      4,
		Content:  "first comment",
	}
	db.ODM(context.Background()).Create(firstComment)
	secondComment := &models.Comment{
		StatusID:   status.ID,
		UID:        5,
		ParentID:   firstComment.ID,
		GroupID:    firstComment.ID,
		OpponentID: 4,
		Content:    "first comment",
	}
	db.ODM(context.Background()).Create(secondComment)
	suite.T().Run("create first level comment", func(t *testing.T) {
		s := &models.Status{}
		db.ODM(context.TODO()).Where(bson.M{"_id": status.ID}).First(s)
		suite.EqualValues(2, s.CommentsCount)

		comment, err := service.CreateComment(context.TODO(), &service.CreateCommentParams{CreateCommentParams: &models.CreateCommentParams{
			StatusID: status.ID,
			UID:      1,
			Content:  "test content",
		}})
		suite.Nil(err)
		suite.EqualValues(0, comment.CommentsCount)
		suite.Equal(primitive.NilObjectID, comment.GroupID)
		suite.NotNil(comment.User)
		s1 := &models.Status{}
		db.ODM(context.TODO()).Where(bson.M{"_id": status.ID}).First(s1)
		suite.EqualValues(3, s1.CommentsCount)
	})

	suite.T().Run("create second level comment", func(t *testing.T) {
		comment, err := service.CreateComment(context.TODO(), &service.CreateCommentParams{CreateCommentParams: &models.CreateCommentParams{
			StatusID: status.ID,
			UID:      1,
			ParentID: firstComment.ID,
			Content:  "test content",
		}})
		suite.Nil(err)
		suite.Equal(comment.GroupID.Hex(), firstComment.ID.Hex())
		suite.Nil(comment.Opponent)
		s := &models.Status{}
		db.ODM(context.TODO()).Where(bson.M{"_id": status.ID}).First(s)
		suite.EqualValues(s.CommentsCount, 4)
		c := &models.Comment{}
		db.ODM(context.TODO()).Where(bson.M{"_id": firstComment.ID}).First(c)
		suite.EqualValues(c.CommentsCount, 1)
		models.PreloadCommentData(context.TODO(), c)
		suite.GreaterOrEqual(len(c.Comments), 1)
	})

	suite.T().Run("comment for second level comment", func(t *testing.T) {
		comment, err := service.CreateComment(context.TODO(), &service.CreateCommentParams{CreateCommentParams: &models.CreateCommentParams{
			StatusID: status.ID,
			UID:      1,
			ParentID: secondComment.ID,
			Content:  "test content",
		}})
		suite.Nil(err)
		suite.Equal(comment.GroupID.Hex(), firstComment.ID.Hex())
		suite.NotNil(comment.User)
		suite.NotNil(comment.Opponent)
		suite.EqualValues(comment.OpponentID, 5)
		c1 := &models.Comment{}
		db.ODM(context.TODO()).Where(bson.M{"_id": firstComment.ID}).First(c1)
		suite.EqualValues(c1.CommentsCount, 2)

		models.PreloadCommentData(context.TODO(), c1)
		suite.GreaterOrEqual(len(c1.Comments), 2)

		c2 := &models.Comment{}
		db.ODM(context.TODO()).Where(bson.M{"_id": secondComment.ID}).First(c2)

		suite.EqualValues(c2.CommentsCount, 0)
	})

}
