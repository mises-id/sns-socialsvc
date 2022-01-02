package comment

import (
	"context"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
)

type ListCommentParams struct {
	models.ListCommentParams
}

func ListComment(ctx context.Context, params *ListCommentParams) ([]*models.Comment, pagination.Pagination, error) {
	return models.ListComment(ctx, &params.ListCommentParams)
}
