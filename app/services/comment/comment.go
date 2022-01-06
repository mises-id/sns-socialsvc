package comment

import (
	"context"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	_, err := models.FindStatus(ctx, params.StatusID)
	if err != nil {
		return nil, err
	}
	if params.ParentID != primitive.NilObjectID {
		parent, err := models.FindComment(ctx, params.ParentID)
		if err != nil {
			return nil, err
		}
		if parent.GroupID == primitive.NilObjectID {
			commentParams.GroupID = parent.ID
		} else {
			commentParams.GroupID = parent.GroupID
		}
		commentParams.OpponentID = parent.UID
	}
	return models.CreateComment(ctx, commentParams)
}
