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
	lookupUserNum       = 10
	sendTweetNum        = 3
	followTwitterNum    = 5
	checkTwitterUserNum = 5
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
		if user_twitter.IsAirdrop == true || user_twitter.ValidState == 2 {
			continue
		}
		uid := user_twitter.UID
		twitter_user, err := getTwitterUserById(ctx, user_twitter.TwitterUserId)
		if err != nil {
			fmt.Printf("[%s]uid[%d] RunLookupTwitterUser GetTwitterUserById Error:%s \n", time.Now().Local().String(), uid, err.Error())
			user_twitter.FindTwitterUserState = 3
			if strings.Contains(err.Error(), "httpStatusCode=401") {
				//user_twitter.FindTwitterUserState = 4
				//delete
				models.DeleteUserTwitterAuthByID(ctx, user_twitter.ID)
				continue
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
			FollowingCount: uint64(*twitter_user.PublicMetrics.FollowingCount),
			TweetCount:     uint64(*twitter_user.PublicMetrics.TweetCount),
		}
		user_twitter.TwitterUser = TwitterUser
		//follow
		user_twitter.FollowState = 1
		channel_user, err := models.FindChannelUserByUID(ctx, uid)
		var amount int64
		var do_channeluser bool
		var valid_state enum.UserValidState
		user_twitter.ValidState = 3
		valid_state = enum.UserValidFailed
		if channel_user != nil && (channel_user.ValidState == enum.UserValidDefalut || channel_user.ValidState == enum.UserValidFailed) {
			do_channeluser = true
			fmt.Printf("RunLookupTwitterUser DoChannelUser True [%s] UID[%d]\n", time.Now().Local().String(), uid)
		}
		//check
		followers_count := user_twitter.TwitterUser.FollowersCount
		//is_valid
		if IsValidTwitterUser(user_twitter.TwitterUser) {
			if followers_count >= 500 && followers_count <= 10000 {
				user_twitter.ValidState = 4
				fmt.Printf("[%s]uid[%d] RunLookupTwitterUser CheckValidState FollowersCount[%d]", time.Now().Local().String(), uid, followers_count)
			} else {
				airdropData, err := createAirdrop(ctx, user_twitter)
				if err != nil {
					fmt.Printf("[%s]uid[%d] RunLookupTwitterUser CreateAirdrop Error:%s \n", time.Now().Local().String(), uid, err.Error())
					user_twitter.FindTwitterUserState = 3
					models.UpdateUserTwitterAuthFindState(ctx, user_twitter)
					continue
				}
				user_twitter.Amount = airdropData.Coin
				user_twitter.IsAirdrop = true
				user_twitter.ValidState = 2
				user_twitter.SendTweeState = 1
				//channel_user
				if do_channeluser {
					amount = user_twitter.Amount / 10
					valid_state = enum.UserValidSucessed
				}
			}
		}
		user_twitter.FindTwitterUserState = 2
		//update
		err = models.UpdateUserTwitterAuthTwitterUser(ctx, user_twitter)
		if err != nil {
			fmt.Printf("[%s]uid[%d] RunLookupTwitterUser UpdateUserTwitterAuthTwitterUser Error:%s \n", time.Now().Local().String(), uid, err.Error())
			continue
		}
		fmt.Printf("[%s]uid[%d] RunLookupTwitterUser Success \n", time.Now().Local().String(), uid)
		//do channel_user
		if do_channeluser {
			if err := channel_user.UpdateCreateAirdrop(ctx, valid_state, amount); err != nil {
				fmt.Printf("RunLookupTwitterUser UpdateChannelUser [%s] UID[%d] Error:%s\n", time.Now().Local().String(), uid, err.Error())
			} else {
				fmt.Printf("RunLookupTwitterUser UpdateChannelUser [%s] UID[%d] Success\n", time.Now().Local().String(), uid)
			}
		}
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
		SortBy:         followerSortOrIDAsc(),
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
		user_twitter.SendTweeState = 2
		/* mises := utils.UMisesToMises(uint64(GetTwitterAirdropCoin(ctx, user_twitter)))
		misesid := utils.RemoveMisesidProfix(user_twitter.Misesid)
		tweet := fmt.Sprintf("I have claimed %.2f $MIS airdrop by using Mises Browser @Mises001, which supports Web3 sites and extensions on mobile.\n\nhttps://www.mises.site/download?MisesID=%s\n\n#Mises #Browser #Wallet #web3 #extension", mises, misesid)
		if err := sendTweet(ctx, user_twitter, tweet); err != nil {
			fmt.Printf("[%s]uid[%d] Send Tweet Error:%s \n", time.Now().Local().String(), uid, err.Error())
			user_twitter.SendTweeState = 3
			if strings.Contains(err.Error(), "httpStatusCode=401") {
				user_twitter.SendTweeState = 4
			}
			if strings.Contains(err.Error(), "httpStatusCode=429") {
				return nil
			}
		} */
		if err := reTweet(ctx, user_twitter); err != nil {
			fmt.Printf("[%s]uid[%d] Send Tweet Error:%s \n", time.Now().Local().String(), uid, err.Error())
			user_twitter.SendTweeState = 3
			if strings.Contains(err.Error(), "httpStatusCode=401") {
				user_twitter.SendTweeState = 4
			}
			if strings.Contains(err.Error(), "httpStatusCode=429") {
				user_twitter.SendTweeState = 5
			}
		}
		//like tweet
		user_twitter.LikeTweeState = 2
		if user_twitter.SendTweeState == 4 {
			user_twitter.LikeTweeState = 4
		} else {
			if err := likeTweet(ctx, user_twitter); err != nil {
				fmt.Printf("[%s]uid[%d] Like Tweet Error:%s \n", time.Now().Local().String(), uid, err.Error())
				user_twitter.LikeTweeState = 3
				if strings.Contains(err.Error(), "httpStatusCode=401") {
					user_twitter.LikeTweeState = 4
				}
				if strings.Contains(err.Error(), "httpStatusCode=429") {
					user_twitter.LikeTweeState = 5
				}
			}
		}
		if err := models.UpdateUserTwitterAuthSendTweet(ctx, user_twitter); err != nil {
			fmt.Printf("[%s]uid[%d] RunSendTweet UpdateUserTwitterAuthSendTweet Error:%s\n ", time.Now().Local().String(), uid, err.Error())
			continue
		}
		if user_twitter.SendTweeState == 2 {
			fmt.Printf("[%s]uid[%d] RunSendTweet Success \n", time.Now().Local().String(), uid)
		}
		if user_twitter.LikeTweeState == 2 {
			fmt.Printf("[%s]uid[%d] LikeTweet Success \n", time.Now().Local().String(), uid)
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
		SortBy:      followerSortOrIDAsc(),
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
			if strings.Contains(err.Error(), "httpStatusCode=429") {
				user_twitter.FollowState = 5
			}
		}
		if err = models.UpdateUserTwitterAuthFollow(ctx, user_twitter); err != nil {
			fmt.Printf("[%s]uid[%d],RunFollowTwitter UpdateUserTwitterAuthFollow Error:%s\n", time.Now().String(), uid, err.Error())
			continue
		}
		if user_twitter.FollowState == 2 {
			fmt.Printf("[%s]uid[%d] RunFollowTwitter Success \n", time.Now().Local().String(), uid)
		}
	}
	return nil
}

//check TwitterUser
func PlanCheckTwitterUser(ctx context.Context) error {
	fmt.Printf("[%s]PlanCheckTwitterUser Start\n", time.Now().Local().String())
	err := runCheckTwitterUser(ctx)
	fmt.Printf("[%s]PlanCheckTwitterUser End\n", time.Now().Local().String())
	return err
}

func followerSortOrIDAsc() string {
	sort := "followers_count_sort"
	m := utils.GetRand(1, 100) % 3
	if m == 0 {
		sort = "id_asc"
	}
	return sort
}

func runCheckTwitterUser(ctx context.Context) error {
	//get list

	params := &search.UserTwitterAuthSearch{
		ValidState: 4,
		SortBy:     followerSortOrIDAsc(),
		ListNum:    int64(checkTwitterUserNum),
	}
	user_twitter_list, err := models.ListUserTwitterAuth(ctx, params)
	if err != nil {
		return err
	}
	num := len(user_twitter_list)
	if num <= 0 {
		return nil
	}
	fmt.Printf("[%s]PlanCheckTwitterUser %d \n", time.Now().Local().String(), num)
	//do list
	for _, user_twitter := range user_twitter_list {
		if user_twitter.ValidState != 4 {
			continue
		}
		uid := user_twitter.UID
		if user_twitter.TwitterUser == nil {
			fmt.Printf("[%s] uid[%d],Error PlanCheckTwitterUser TwitterUser is Null\n", time.Now().String(), uid)
			continue
		}
		followers, err := userFollowers(ctx, user_twitter)
		if err != nil {
			fmt.Printf("[%s] uid[%d],PlanCheckTwitterUser UserFollowers Error:%s\n", time.Now().String(), uid, err.Error())
			if strings.Contains(err.Error(), "httpStatusCode=429") {
				continue
			}
			user_twitter.ValidState = 3 //invalid
			updateUserTwitterAuthTwitterUser(ctx, user_twitter)
			continue
		}
		if followers == nil || len(followers.Data) == 0 {
			user_twitter.ValidState = 3 //invalid
			updateUserTwitterAuthTwitterUser(ctx, user_twitter)
			continue
		}
		//check followers
		followersNum := len(followers.Data)
		et := time.Now().UTC().AddDate(0, -3, 0)
		fmt.Printf("[%s] uid[%d],PlanCheckTwitterUser Check ET[%s]\n", time.Now().String(), uid, et.String())
		var zeroTweetNum, zeroFollowerNum, recentRegisterNum, totalFollowerNum int
		for _, follower := range followers.Data {
			followerUser := &models.TwitterUser{
				TwitterUserId:  *follower.ID,
				UserName:       *follower.Username,
				Name:           *follower.Name,
				CreatedAt:      *follower.CreatedAt,
				FollowersCount: uint64(*follower.PublicMetrics.FollowersCount),
				FollowingCount: uint64(*follower.PublicMetrics.FollowingCount),
				TweetCount:     uint64(*follower.PublicMetrics.TweetCount),
			}
			totalFollowerNum += *follower.PublicMetrics.FollowersCount
			if followerUser.TweetCount == 0 {
				zeroTweetNum++
			}
			if followerUser.FollowersCount == 0 {
				zeroFollowerNum++
			}
			if et.UTC().Unix() < followerUser.CreatedAt.UTC().Unix() {
				recentRegisterNum++
			}
			fmt.Println("followerUser: ", followerUser)
		}
		checkResult := &models.CheckResult{
			CheckNum:          followersNum,
			ZeroTweetNum:      zeroTweetNum,
			ZeroFollowerNum:   zeroFollowerNum,
			RecentRegisterNum: recentRegisterNum,
			TotalFollowerNum:  totalFollowerNum,
		}
		user_twitter.CheckResult = checkResult
		fmt.Printf("[%s] uid[%d],PlanCheckTwitterUser Check FollowersNum[%d],zeroTweetNum[%d],zeroFollowerNum[%d],recentRegisterNum[%d],totalFollowerNum[%d]\n", time.Now().String(), uid, followersNum, zeroTweetNum, zeroFollowerNum, recentRegisterNum, totalFollowerNum)
		channel_user, err := models.FindChannelUserByUID(ctx, uid)
		var amount int64
		var do_channeluser bool
		valid_state := enum.UserValidFailed
		if channel_user != nil && (channel_user.ValidState == enum.UserValidDefalut || channel_user.ValidState == enum.UserValidFailed) {
			do_channeluser = true
			fmt.Printf("PlanCheckTwitterUser DoChannelUser True [%s] UID[%d]\n", time.Now().Local().String(), uid)
		}
		airdropData, err := createAirdrop(ctx, user_twitter)
		if err != nil {
			fmt.Printf("[%s]uid[%d] PlanCheckTwitterUser CreateAirdrop Error:%s \n", time.Now().Local().String(), uid, err.Error())
			user_twitter.ValidState = 5
			updateUserTwitterAuthTwitterUser(ctx, user_twitter)
			continue
		}
		user_twitter.Amount = airdropData.Coin
		user_twitter.IsAirdrop = true
		user_twitter.ValidState = 2
		user_twitter.SendTweeState = 1
		//channel_user
		if do_channeluser {
			amount = user_twitter.Amount / 10
			valid_state = enum.UserValidSucessed
		}
		//update
		err = updateUserTwitterAuthTwitterUser(ctx, user_twitter)
		if err == nil {
			fmt.Printf("[%s] uid[%d],coin[%d] PlanCheckTwitterUser Success \n", time.Now().Local().String(), uid, user_twitter.Amount)
		}
		//do channel_user
		if do_channeluser {
			if err := channel_user.UpdateCreateAirdrop(ctx, valid_state, amount); err != nil {
				fmt.Printf("PlanCheckTwitterUser UpdateChannelUser [%s] UID[%d] Error:%s\n", time.Now().Local().String(), uid, err.Error())
			} else {
				fmt.Printf("PlanCheckTwitterUser UpdateChannelUser [%s] UID[%d] Success\n", time.Now().Local().String(), uid)
			}
		}
	}
	return nil
}

func updateUserTwitterAuthTwitterUser(ctx context.Context, user_twitter *models.UserTwitterAuth) error {
	err := models.UpdateUserTwitterAuthTwitterUser(ctx, user_twitter)
	if err != nil {
		fmt.Printf("[%s] uid[%d] PlanCheckTwitterUser UpdateUserTwitterAuthTwitterUser Error:%s \n", time.Now().Local().String(), user_twitter.UID, err.Error())
	}
	return err
}
