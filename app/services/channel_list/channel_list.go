package channel_list

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/search"
	"github.com/mises-id/sns-socialsvc/config/env"
	"github.com/mises-id/sns-socialsvc/lib/codes"
	"github.com/mises-id/sns-socialsvc/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	invalidMisesid = errors.New("misesid invalid")
	googlePlayUrl  = "https://play.google.com/store/apps/details?id="
	playAppUrl     = "https://play.app.goo.gl/?link="
)

type (
	ChannelUrlInput struct {
		Misesid string
	}
	ChannelUrlOutput struct {
		Url              string
		PosterUrl        string
		TotalChannelUser uint64
		AirdropAmount    float64 //mises
	}
)

//channel info
func ChannelInfo(ctx context.Context, in *ChannelUrlInput) (*ChannelUrlOutput, error) {

	out := &ChannelUrlOutput{}
	if in.Misesid == "" {
		return out, codes.ErrInvalidArgument
	}
	//find channel by misesid
	misesid := utils.AddMisesidProfix(in.Misesid)
	channel, err := models.FindChannelListByMisesid(ctx, misesid)
	//if not exist to create
	if err == mongo.ErrNoDocuments {
		channel, err = models.CreateChannelByMisesid(ctx, misesid)
		if err == models.ChannelExist {
			channel, err = models.FindChannelListByMisesid(ctx, misesid)
		}
	}
	if channel == nil {
		return out, codes.ErrInvalidArgument.Newf(err.Error())
	}
	url := getChannelUrl(ctx, channel)
	out.Url = url
	out.PosterUrl = getChannelPosterUrl(ctx, channel)
	out.TotalChannelUser = countChannelTotalUser(ctx, channel.ID)
	out.AirdropAmount = getChannelAirdropAmount(ctx, channel.UID)
	return out, nil
}

//get channel airdrop amount
func getChannelAirdropAmount(ctx context.Context, uid uint64) float64 {
	var mises float64
	user_ext, err := models.FindOrCreateUserExt(ctx, uid)
	if err != nil {
		fmt.Printf("uid[%d],count channel user error: %s \n", uid, err.Error())
	} else {
		umises := user_ext.ChannelAirdropCoin
		mises = utils.UMisesToMises(umises)
	}
	return mises
}

//count channel user
func countChannelTotalUser(ctx context.Context, channel_id primitive.ObjectID) uint64 {

	params := &search.ChannelUserSearch{
		ChannelID: channel_id,
	}
	c, err := models.CountChannelUser(ctx, params)
	if err != nil {
		c = 0
		fmt.Printf("channel_id[%s],count channel user error: %s \n", channel_id.Hex(), err.Error())
	}
	return uint64(c)
}

//get channel url
func getChannelUrl(ctx context.Context, ch *models.ChannelList) string {

	appid := env.Envs.GooglePlayAppID
	referrer := "utm_source=" + utils.AddChannelUrlProfix(ch.ID.Hex())
	googlePlay := googlePlayUrl + appid + "&referrer=" + url.QueryEscape(referrer)
	return playAppUrl + url.QueryEscape(googlePlay)
}
func getChannelPosterUrl(ctx context.Context, ch *models.ChannelList) string {

	appid := env.Envs.GooglePlayAppID
	referrer := "utm_source=" + utils.AddChannelUrlProfix(ch.ID.Hex()) + "&utm_medium=poster"
	googlePlay := googlePlayUrl + appid + "&referrer=" + url.QueryEscape(referrer)
	return playAppUrl + url.QueryEscape(googlePlay)
}
