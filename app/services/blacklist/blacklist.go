package blacklist

import (
	"context"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
)

type ListBlacklistParams struct {
	UID        uint64
	PageParams *pagination.PageQuickParams
}

func ListBlacklist(ctx context.Context, params *ListBlacklistParams) ([]*models.Blacklist, pagination.Pagination, error) {
	return models.ListBlacklist(ctx, params.UID, params.PageParams)
}

func CreateBlacklist(ctx context.Context, uid, targetUID uint64) (*models.Blacklist, error) {
	// check user
	targetUser, err := models.FindUser(ctx, targetUID)
	if err != nil {
		return nil, err
	}
	blacklist, err := models.CreateBlacklist(ctx, uid, targetUser.UID)
	if err == nil {
		return blacklist, nil
	}
	// delete follow
	if err = models.EnsureDeleteFollow(ctx, uid, targetUID); err != nil {
		return nil, err
	}
	return blacklist, nil
}

func DeleteBlacklist(ctx context.Context, uid, targetUID uint64) error {
	return models.DeleteBlacklist(ctx, uid, targetUID)
}
