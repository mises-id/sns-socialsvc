package user_twitter

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/search"
	"github.com/mises-id/sns-socialsvc/lib/utils"
)

const (
	lookupUserNum    = 10
	sendTweetNum     = 3
	followTwitterNum = 3
)

func PlanLookupTwitterUser(ctx context.Context) error {
	fmt.Printf("[%s]RunLookupTwitterUser Start\n", time.Now().Local().String())
	err := runLookupTwitterUser(ctx)
	fmt.Printf("[%s]RunLookupTwitterUser End\n", time.Now().Local().String())
	return err
}

func runLookupTwitterUser(ctx context.Context) error {

	//get list
	params := &search.UserTwitterAuthSearch{
		FindTwitterUserState: 1,
		SortType:             enum.SortAsc,
		SortKey:              "_id",
		ListNum:              int64(lookupUserNum),
	}
	user_twitter_list, err := models.ListUserTwitterAuth(ctx, params)
	if err != nil {
		return err
	}
	num := len(user_twitter_list)
	if num <= 0 {
		return nil
	}
	fmt.Printf("[%s]RunLookupTwitterUser %d \n", time.Now().Local().String(), num)
	//do list
	for _, user_twitter := range user_twitter_list {
		if user_twitter.IsAirdrop == true {
			continue
		}
		uid := user_twitter.UID
		twitter_user, err := getTwitterUserById(ctx, user_twitter.TwitterUserId)
		if err != nil {
			fmt.Printf("[%s]uid[%d] RunLookupTwitterUser GetTwitterUserById Error:%s \n", time.Now().Local().String(), uid, err.Error())
			user_twitter.FindTwitterUserState = 3
			if strings.Contains(err.Error(), "httpStatusCode=401") {
				user_twitter.FindTwitterUserState = 4
			}
			models.UpdateUserTwitterAuthFindState(ctx, user_twitter)
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
		//follow
		user_twitter.FollowState = 1
		//is_valid
		if IsValidTwitterUser(user_twitter.TwitterUser) {
			if err := createAirdrop(ctx, user_twitter); err != nil {
				fmt.Printf("[%s]uid[%d] RunLookupTwitterUser CreateAirdrop Error:%s \n", time.Now().Local().String(), uid, err.Error())
				user_twitter.FindTwitterUserState = 3
				models.UpdateUserTwitterAuthFindState(ctx, user_twitter)
				continue
			}
			user_twitter.IsAirdrop = true
			user_twitter.SendTweeState = 1
		}
		user_twitter.FindTwitterUserState = 2
		//update
		err = models.UpdateUserTwitterAuthTwitterUser(ctx, user_twitter)
		if err != nil {
			fmt.Printf("[%s]uid[%d] RunLookupTwitterUser UpdateUserTwitterAuthTwitterUser Error:%s \n", time.Now().Local().String(), uid, err.Error())
			continue
		}
		fmt.Printf("[%s]uid[%d] RunLookupTwitterUser Success \n", time.Now().Local().String(), uid)
	}
	return nil
}
func PlanSendTweet(ctx context.Context) error {
	fmt.Printf("[%s]RunSendTweet Start\n", time.Now().Local().String())
	err := runSendTweet(ctx)
	fmt.Printf("[%s]RunSendTweet End\n", time.Now().Local().String())
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
	fmt.Printf("[%s]RunSendTweet %d \n", time.Now().Local().String(), num)
	//do list
	for _, user_twitter := range user_twitter_list {
		uid := user_twitter.UID
		mises := utils.UMisesToMises(uint64(GetTwitterAirdropCoin(ctx, user_twitter)))
		misesid := utils.RemoveMisesidProfix(user_twitter.Misesid)
		tweet := fmt.Sprintf("I have claimed $%.2f $MIS airdrop by using Mises Browser @Mises001, which supports Web3 sites and extensions on mobile.\n\nhttps://www.mises.site/download?MisesID=%s\n\n#Mises #Browser #web3 #extension", mises, misesid)
		user_twitter.SendTweeState = 2
		if err := sendTweet(ctx, user_twitter, tweet); err != nil {
			fmt.Printf("[%s]uid[%d] Send Tweet Error:%s \n", time.Now().Local().String(), uid, err.Error())
			user_twitter.SendTweeState = 3
			if strings.Contains(err.Error(), "httpStatusCode=401") {
				user_twitter.SendTweeState = 4
			}
		}
		if err := models.UpdateUserTwitterAuthSendTweet(ctx, user_twitter); err != nil {
			fmt.Printf("[%s]uid[%d] RunSendTweet UpdateUserTwitterAuthSendTweet Error:%s\n ", time.Now().Local().String(), uid, err.Error())
			continue
		}
		if user_twitter.SendTweeState == 2 {
			fmt.Printf("[%s]uid[%d] RunSendTweet Success \n", time.Now().Local().String(), uid)
		}
	}
	return nil
}

//follow twitter
func FollowTwitter(ctx context.Context) error {
	fmt.Printf("[%s]RunFollowTwitter Start\n", time.Now().Local().String())
	err := runFollowTwitter(ctx)
	fmt.Printf("[%s]RunFollowTwitter End\n", time.Now().Local().String())
	return err
}

func runFollowTwitter(ctx context.Context) error {
	//get list
	params := &search.UserTwitterAuthSearch{
		FollowState: 1,
		SortType:    enum.SortAsc,
		SortKey:     "_id",
		ListNum:     int64(followTwitterNum),
	}
	user_twitter_list, err := models.ListUserTwitterAuth(ctx, params)
	if err != nil {
		return err
	}
	num := len(user_twitter_list)
	if num <= 0 {
		return nil
	}
	fmt.Printf("[%s]RunFollowTwitter %d \n", time.Now().Local().String(), num)
	//do list
	for _, user_twitter := range user_twitter_list {
		uid := user_twitter.UID
		//to follow
		user_twitter.FollowState = 2
		if err := apiFollowTwitterUser(ctx, user_twitter, targetTwitterId); err != nil {
			fmt.Printf("[%s]uid[%d],RunFollowTwitter ApiFollowTwitterUser error:%s\n", time.Now().String(), uid, err.Error())
			user_twitter.FollowState = 3
			if strings.Contains(err.Error(), "httpStatusCode=401") {
				user_twitter.FollowState = 4
			}
		}
		if err = models.UpdateUserTwitterAuthFollew(ctx, user_twitter); err != nil {
			fmt.Printf("[%s]uid[%d],RunFollowTwitter UpdateUserTwitterAuthFollew Error:%s\n", time.Now().String(), uid, err.Error())
			continue
		}
		if user_twitter.FollowState == 2 {
			fmt.Printf("[%s]uid[%d] RunFollowTwitter Success \n", time.Now().Local().String(), uid)
		}
	}

	return nil
}
