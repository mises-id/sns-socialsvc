package message

import (
	"context"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
)

type ListMessageParams struct {
	models.ListMessageParams
}

type ReadMessageParams struct {
	*models.ReadMessageParams
}

func ListMessage(ctx context.Context, params *ListMessageParams) ([]*models.Message, pagination.Pagination, error) {
	return models.ListMessage(ctx, &params.ListMessageParams)
}

func ReadMessages(ctx context.Context, params *ReadMessageParams) error {
	return models.ReadMessages(ctx, params.ReadMessageParams)
}
