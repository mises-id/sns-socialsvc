package channel_user

import (
	"context"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"github.com/mises-id/sns-socialsvc/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	PageChannelUserInput struct {
		PageParams *pagination.PageQuickParams
		Misesid    string
	}
)

//create channel user
func CreateChannelUser(ctx context.Context, uid uint64, channel_str string) error {

	channel_str = utils.RemoveChannelUrlProfix(channel_str)
	channel_id, err := primitive.ObjectIDFromHex(channel_str)
	if err != nil {
		return err
	}
	channel, err := models.FindChannelListByID(ctx, channel_id)
	return createChannelUser(ctx, uid, channel)
}

func createChannelUser(ctx context.Context, uid uint64, channel *models.ChannelList) error {
	channel_id := channel.ID

	//create channel user
	channelUser := &models.ChannelUser{
		ChannelID:      channel.ID,
		ChannelMisesid: channel.Misesid,
		ChannelUID:     channel.UID,
		UID:            uid,
		ValidState:     enum.UserValidDefalut,
		AirdropState:   enum.ChannelAirdropDefault,
	}
	_, err := models.CreateChannelUser(ctx, channelUser)
	if err != nil && err != models.ChannelUserExist {
		return err
	}
	return models.UpdateUserChannelIDByUID(ctx, uid, channel_id)

}

//page channel user
func PageChannelUser(ctx context.Context, in *PageChannelUserInput) ([]*models.ChannelUser, pagination.Pagination, error) {

	if in.Misesid == "" {
		return []*models.ChannelUser{}, &pagination.QuickPagination{}, nil
	}
	params := &models.PageChannelUserInput{
		PageParams: in.PageParams,
		Misesid:    in.Misesid,
	}
	res, page, err := models.PageChannelUser(ctx, params)
	if err != nil {
		return nil, nil, err
	}
	return res, page, nil
}
