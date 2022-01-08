package comment

import (
	"context"
	"testing"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/storage"
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
		UID:        5,
		FromType:   enum.FromPost,
		StatusType: enum.TextStatus,
		Content:    "hello",
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
		if comment.GroupID != primitive.NilObjectID {
			t.Errorf("wrong group id for first level comment: %s", comment.GroupID.Hex())
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
	})

}
