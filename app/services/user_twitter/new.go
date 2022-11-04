package user_twitter

import (
	"context"
	"fmt"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/search"
)

const (
	lookupUserNum = 15
	sendTweetNum  = 3
)

func PlanLookupTwitterUser(ctx context.Context) error {
	err := runLookupTwitterUser(ctx)
	return err
}

func runLookupTwitterUser(ctx context.Context) error {
	//get list
	params := &search.UserTwitterAuthSearch{
		IsFindTwitterUserState: 1,
		SortType:               enum.SortAsc,
		SortKey:                "_id",
		ListNum:                int64(lookupUserNum),
	}
	user_twitter_list, err := models.ListUserTwitterAuth(ctx, params)
	if err != nil {
		return err
	}
	num := len(user_twitter_list)
	if num <= 0 {
		return nil
	}
	//do list
	for _, user_twitter := range user_twitter_list {
		if user_twitter.IsAirdrop == true {
			continue
		}
		twitter_user, err := getTwitterUserById(ctx, user_twitter.TwitterUserId)
		if err != nil {
			fmt.Println("runLookupTwitterUser getTwitterUserById err: ", err.Error())
			continue
		}
		TwitterUser := &models.TwitterUser{
			TwitterUserId:  *twitter_user.ID,
			UserName:       *twitter_user.Username,
			Name:           *twitter_user.Name,
			CreatedAt:      *twitter_user.CreatedAt,
			FollowersCount: uint64(*twitter_user.PublicMetrics.FollowersCount),
			TweetCount:     uint64(*twitter_user.PublicMetrics.TweetCount),
		}
		user_twitter.TwitterUser = TwitterUser
		//is_valid
		if IsValidTwitterUser(user_twitter.TwitterUser) {
			if err := createAirdrop(ctx, user_twitter); err != nil {
				fmt.Println("runLookupTwitterUser createAirdrop err: ", err.Error())
				continue
			}
			user_twitter.IsAirdrop = true
			user_twitter.SendTweeState = 1
		}
		user_twitter.IsFindTwitterUser = true
		//update
		err = models.UpdateUserTwitterAuthTwitterUser(ctx, user_twitter)
		if err != nil {
			fmt.Println("runLookupTwitterUser UpdateUserTwitterAuthTwitterUser err: ", err.Error())
			continue
		}
	}
	return nil
}
func PlanSendTweet(ctx context.Context) error {
	err := runSendTweet(ctx)
	return err
}

func runSendTweet(ctx context.Context) error {
	//get list
	params := &search.UserTwitterAuthSearch{
		SendTweetState: 1,
		SortType:       enum.SortAsc,
		SortKey:        "_id",
		ListNum:        int64(sendTweetNum),
	}
	user_twitter_list, err := models.ListUserTwitterAuth(ctx, params)
	if err != nil {
		return err
	}
	num := len(user_twitter_list)
	if num <= 0 {
		return nil
	}
	//do list
	/* for _, user_twitter := range user_twitter_list {


	} */
	return nil
}

func getTweetTemp(ctx context.Context) string {

	return ""
}
