package comment

import (
	"context"
	"testing"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func cleanTables(names ...string) {
	for _, name := range names {
		db.DB().Collection(name).Drop(context.TODO())
	}
}

func TestCreateComment(t *testing.T) {
	db.SetupMongo(context.TODO())
	storage.SetupImageStorage("127.0.0.1", "xxx", "xx")
	defer cleanTables("users", "statuses", "comments")
	for i := 0; i < 10; i++ {
		db.ODM(context.Background()).Create(&models.User{
			UID: uint64(i + 1),
		})
	}
	status := &models.Status{
		UID:           5,
		FromType:      enum.FromPost,
		StatusType:    enum.TextStatus,
		CommentsCount: 2,
		Content:       "hello",
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
	t.Run("create first level comment", func(t *testing.T) {
		comment, err := CreateComment(context.TODO(), &CreateCommentParams{
			CreateCommentParams: &models.CreateCommentParams{
				StatusID: status.ID,
				UID:      1,
				Content:  "test content",
			},
		})
		if err != nil {
			t.Error(err)
		}
		if comment.CommentsCount != 0 {
			t.Errorf("wrong comment.comments_count after created: %d", comment.CommentsCount)
		}
		if comment.GroupID != primitive.NilObjectID {
			t.Errorf("wrong group id for first level comment: %s", comment.GroupID.Hex())
		}
		if comment.User == nil {
			t.Errorf("blank comment.user after created")
		}
		s := &models.Status{}
		db.ODM(context.TODO()).Where(bson.M{"_id": status.ID}).First(s)
		if s.CommentsCount != 3 {
			t.Errorf("wrong comments count after comment to status: %d", s.CommentsCount)
		}
	})

	t.Run("create second level comment", func(t *testing.T) {
		comment, err := CreateComment(context.TODO(), &CreateCommentParams{
			CreateCommentParams: &models.CreateCommentParams{
				StatusID: status.ID,
				UID:      1,
				ParentID: firstComment.ID,
				Content:  "test content",
			},
		})
		if err != nil {
			t.Error(err)
		}
		if comment.GroupID.Hex() != firstComment.ID.Hex() {
			t.Error("wrong group id for second level comment")
		}
		if comment.Opponent != nil {
			t.Errorf("unexpect comment.opponent after created")
		}
		s := &models.Status{}
		db.ODM(context.TODO()).Where(bson.M{"_id": status.ID}).First(s)
		if s.CommentsCount != 4 {
			t.Errorf("wrong status.comments_count after create a second level comment: %d", s.CommentsCount)
		}
		c := &models.Comment{}
		db.ODM(context.TODO()).Where(bson.M{"_id": firstComment.ID}).First(c)
		if c.CommentsCount != 1 {
			t.Errorf("wrong comment.comments_count after create a second level comment: %d", s.CommentsCount)
		}
		models.PreloadCommentData(context.TODO(), c)
		if len(c.Comments) < 1 {
			t.Errorf("wrong comment.comments after create a second level comment: %v", len(c.CommentIDs))
		}
	})

	t.Run("comment for second level comment", func(t *testing.T) {
		comment, err := CreateComment(context.TODO(), &CreateCommentParams{
			CreateCommentParams: &models.CreateCommentParams{
				StatusID: status.ID,
				UID:      1,
				ParentID: secondComment.ID,
				Content:  "test content",
			},
		})
		if err != nil {
			t.Error(err)
		}
		if comment.GroupID.Hex() != firstComment.ID.Hex() {
			t.Errorf("wrong group id for second level comment %s , %s", comment.GroupID.Hex(), firstComment.ID.Hex())
		}
		if comment.User == nil {
			t.Error("blank user for comment")
		}
		if comment.Opponent == nil {
			t.Errorf("blank comment.opponent after created")
		}
		if comment.OpponentID != 5 {
			t.Errorf("wrong comment.opponent_id after created")
		}
		c1 := &models.Comment{}
		db.ODM(context.TODO()).Where(bson.M{"_id": firstComment.ID}).First(c1)
		if c1.CommentsCount != 2 {
			t.Errorf("wrong first comment.comments_count after create a comment for second level comment: %d", c1.CommentsCount)
		}

		models.PreloadCommentData(context.TODO(), c1)
		if len(c1.Comments) < 2 {
			t.Errorf("wrong comment.comments after comment for a second level comment: %v", len(c1.CommentIDs))
		}

		c2 := &models.Comment{}
		db.ODM(context.TODO()).Where(bson.M{"_id": secondComment.ID}).First(c2)
		if c2.CommentsCount != 0 {
			t.Errorf("wrong second comment.comments_count after create a comment for second level comment: %d", c2.CommentsCount)
		}
	})

}
