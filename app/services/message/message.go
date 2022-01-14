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

type MessageSummary struct {
	LatestMessage      *models.Message
	Total              uint32
	NotificationsCount uint32
	UsersCount         uint32
}

func ListMessage(ctx context.Context, params *ListMessageParams) ([]*models.Message, pagination.Pagination, error) {
	return models.ListMessage(ctx, &params.ListMessageParams)
}

func ReadMessages(ctx context.Context, params *ReadMessageParams) error {
	return models.ReadMessages(ctx, params.ReadMessageParams)
}

func GetMessageSummary(ctx context.Context, uid uint64) (*MessageSummary, error) {
	var err error
	summary := &MessageSummary{}
	summary.NotificationsCount, err = models.UnreadMessagesCount(ctx, uid)
	if err != nil {
		return nil, err
	}
	summary.UsersCount, err = models.UnreadUsersCount(ctx, uid)
	if err != nil {
		return nil, err
	}
	summary.LatestMessage, err = models.LatestUnreadMessage(ctx, uid)
	if err != nil {
		return nil, err
	}
	summary.Total = summary.NotificationsCount + summary.UsersCount
	return summary, nil
}
