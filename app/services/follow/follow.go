package follow

import (
	"context"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/codes"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"github.com/mises-id/sns-socialsvc/lib/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func LatestFollowing(ctx context.Context, uid uint64) ([]*models.Follow, error) {
	follows, err := models.LatestFollowing(ctx, uid)
	if err != nil {
		return nil, err
	}
	return follows, nil
}

func ListFriendship(ctx context.Context, uid uint64, relationType enum.RelationType, pageParams *pagination.QuickPagination) ([]*models.Follow, pagination.Pagination, error) {
	// check user exsit
	_, err := models.FindUser(ctx, uid)
	if err != nil {
		return nil, nil, err
	}
	follows, page, err := models.ListFollow(ctx, uid, relationType, pageParams)
	if err != nil {
		return nil, nil, err
	}
	currentUID, ok := ctx.Value(utils.CurrentUIDKey{}).(uint64)
	if ok && currentUID == uid {
		err = models.ReadNewFans(ctx, uid)
		if err != nil {
			return nil, nil, err
		}
	}
	return follows, page, nil
}

func Follow(ctx context.Context, fromUID, toUID uint64) (*models.Follow, error) {
	if fromUID == toUID {
		return nil, codes.ErrInvalidArgument
	}
	// check user exsist
	_, err := models.FindUser(ctx, fromUID)
	if err != nil {
		return nil, err
	}
	_, err = models.FindUser(ctx, toUID)
	if err != nil {
		return nil, err
	}
	blocked, err := models.IsBlocked(ctx, fromUID, toUID)
	if err != nil {
		return nil, err
	}
	if blocked {
		return nil, codes.ErrUserInBlacklist
	}
	isFriend := false
	follow, err := models.GetFollow(ctx, fromUID, toUID)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	fansFollow, err := models.GetFollow(ctx, toUID, fromUID)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if err == nil {
		isFriend = true
		if !fansFollow.IsFriend {
			if err = fansFollow.SetFriend(ctx, true); err != nil {
				return nil, err
			}
		}
	}
	if follow != nil {
		return follow, follow.SetFriend(ctx, isFriend)
	}
	return models.CreateFollow(ctx, fromUID, toUID, isFriend)
}

func Unfollow(ctx context.Context, fromUID, toUID uint64) error {
	follow, err := models.GetFollow(ctx, fromUID, toUID)
	if err != nil {
		return nil
	}
	fansFollow, err := models.GetFollow(ctx, toUID, fromUID)
	if err == nil && fansFollow.IsFriend {
		if err = fansFollow.SetFriend(ctx, false); err != nil {
			return err
		}
	} else if err != mongo.ErrNoDocuments {
		return err
	}
	return follow.Delete(ctx)
}
