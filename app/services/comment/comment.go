package comment

import (
	"context"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ListCommentParams struct {
	models.ListCommentParams
}

type CreateCommentParams struct {
	*models.CreateCommentParams
}

func ListComment(ctx context.Context, params *ListCommentParams) ([]*models.Comment, pagination.Pagination, error) {
	return models.ListComment(ctx, &params.ListCommentParams)
}

func CreateComment(ctx context.Context, params *CreateCommentParams) (*models.Comment, error) {
	commentParams := params.CreateCommentParams
	// check status exsist
	status, err := models.FindStatus(ctx, params.StatusID)
	if err != nil {
		return nil, err
	}
	var groupComment *models.Comment
	if params.ParentID != primitive.NilObjectID {
		parent, err := models.FindComment(ctx, params.ParentID)
		if err != nil {
			return nil, err
		}
		if parent.GroupID == primitive.NilObjectID {
			commentParams.GroupID = parent.ID
			groupComment = parent
		} else {
			commentParams.GroupID = parent.GroupID
			groupComment, err = models.FindComment(ctx, commentParams.GroupID)
			if err != nil {
				return nil, err
			}
		}
		commentParams.OpponentID = parent.UID
	}
	comment, err := models.CreateComment(ctx, commentParams)
	if err != nil {
		return nil, err
	}
	return comment, incrCommentCounter(ctx, status, groupComment)
}

func LikeComment(ctx context.Context, uid uint64, commentID primitive.ObjectID) (*models.Like, error) {
	comment, err := models.FindComment(ctx, commentID)
	if err != nil {
		return nil, err
	}
	like, err := models.FindLike(ctx, uid, commentID, enum.LikeComment)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if err == nil {
		return like, nil
	}
	like, err = models.CreateLike(ctx, uid, comment.UID, commentID, enum.LikeComment)
	if err != nil {
		return nil, err
	}
	return like, comment.IncCommentCounter(ctx, "likes_count")
}

func UnlikeComment(ctx context.Context, uid uint64, commentID primitive.ObjectID) error {
	like, err := models.FindLike(ctx, uid, commentID, enum.LikeComment)
	if err != nil {
		return err
	}
	comment, err := models.FindComment(ctx, commentID)
	if err != nil {
		return err
	}
	if err = models.DeleteLike(ctx, like.ID); err != nil {
		return err
	}
	return comment.IncCommentCounter(ctx, "likes_count", -1)
}

func incrCommentCounter(ctx context.Context, status *models.Status, comment *models.Comment) error {
	if err := status.IncStatusCounter(ctx, "comments_count"); err != nil {
		return err
	}
	return comment.IncCommentCounter(ctx, "comments_count")
}
