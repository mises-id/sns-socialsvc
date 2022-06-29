package user_twitter

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/fields"
	"github.com/michimani/gotwi/tweets"
	"github.com/michimani/gotwi/tweets/types"
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/config/env"
)

var (
	twiClient         *gotwi.GotwiClient
	tweeTtag          string
	misesidPrefix     = "did:mises:"
	validRegisterDate string
)

type (
	TweetsIn struct {
		Query      string
		StartTime  *time.Time
		EndTime    *time.Time
		NextToken  string
		MaxResults int
	}
	misesTweet struct {
		AuthorID  string
		TweetID   string
		TweetText string
		CreatedAt time.Time
	}
)

func init() {
	tweeTtag = env.Envs.TWEET_TAG
	validRegisterDate = env.Envs.VALID_TWITTER_REGISTER_DATE
}

func TwitterAuth(ctx context.Context) {
	if !models.GetAirdropStatus(ctx) {
		fmt.Println("airdrop status end")
		return
	}
	/* 	proxy := func(_ *http.Request) (*url.URL, error) {
	   		return url.Parse("http://127.0.0.1:1087")
	   	}
	   	transport := &http.Transport{Proxy: proxy}
	   	client := &http.Client{Transport: transport} */
	client := &http.Client{}
	in := &gotwi.NewGotwiClientInput{
		HTTPClient:           client,
		AuthenticationMethod: gotwi.AuthenMethodOAuth2BearerToken,
	}

	c, err := gotwi.NewGotwiClient(in)
	if err != nil {
		fmt.Println(err)
		return
	}
	twiClient = c
	//dateNow := time.Now().UTC().AddDate(0, 0, -3)
	//startTime := time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 0, 0, 0, 0, dateNow.Location())
	//endTime := time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 23, 59, 59, 0, dateNow.Location())
	sh, _ := time.ParseDuration("-7m")
	startTime := time.Now().UTC().Add(sh)
	tweetIn := &TweetsIn{
		Query:     tweeTtag,
		StartTime: &startTime,
		//EndTime:   &endTime,
	}
	getTwitter(context.TODO(), tweetIn)
}

func getTwitter(ctx context.Context, in *TweetsIn) {
	res, err := getTwitterList(ctx, in)
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	twitterAuth(ctx, res)
	if res.Meta.NextToken != nil {
		in.NextToken = *res.Meta.NextToken
		getTwitter(ctx, in)
	}
	return
}

func twitterAuth(ctx context.Context, tweets *types.SearchTweetsRecentResponse) {

	misesids := make([]string, 0)
	twitterUserIds := make([]string, 0)
	userTwitters := make([]*models.UserTwitterAuth, 0)
	tweetAuthors := make(map[string]*models.TwitterUser, 0)
	misesTweetsMap := make(map[string][]misesTweet, 0)
	//misesTweetsMap := make(map[string]misesTweet, 0)
	for _, v := range tweets.Includes.Users {
		userMetrics := v.PublicMetrics
		twitterUserId := gotwi.StringValue(v.ID)
		tweetUser := &models.TwitterUser{
			TwitterUserId:  twitterUserId,
			Name:           gotwi.StringValue(v.Name),
			UserName:       gotwi.StringValue(v.Username),
			FollowersCount: uint64(gotwi.IntValue(userMetrics.FollowersCount)),
			TweetCount:     uint64(gotwi.IntValue(userMetrics.TweetCount)),
			CreatedAt:      *v.CreatedAt,
		}
		twitterUserIds = append(twitterUserIds, twitterUserId)
		tweetAuthors[twitterUserId] = tweetUser

	}
	//check tweet text
	for _, v := range tweets.Data {
		if v.Entities == nil || len(v.Entities.URLs) == 0 {
			continue
		}
		urls := v.Entities.URLs
		url := gotwi.StringValue(urls[0].ExpandedURL)
		text := gotwi.StringValue(v.Text)
		misesid, err := getMisesIdByTweetText(url)
		if err != nil {
			continue
		}
		twitterUserId := gotwi.StringValue(v.AuthorID)
		tweetId := gotwi.StringValue(v.ID)
		misesids = append(misesids, misesid)
		mises_tweet := misesTweet{
			AuthorID:  twitterUserId,
			TweetID:   tweetId,
			TweetText: text,
			CreatedAt: *v.CreatedAt,
		}
		//misesTweetsMap[misesid] = mises_tweet
		_, ok := misesTweetsMap[misesid]
		if !ok {
			misesTweetsMap[misesid] = []misesTweet{mises_tweet}
		} else {
			misesTweetsMap[misesid] = append(misesTweetsMap[misesid], mises_tweet)
		}

	}
	//find users by misesids
	misesUserNum := len(misesids)
	if misesUserNum == 0 {
		return
	}
	users, err := models.FindUserByMisesids(ctx, misesids...)
	if err != nil {
		fmt.Println("find user by misesids error: ", err.Error())
		return
	}
	//fitler:  find user twitter auth by misesids or twitter user id
	existUserTwitterAuths, err := models.ListUserTwitterAuthByMisesidsOrTwitterUserIds(ctx, misesids, twitterUserIds)
	if err != nil {
		fmt.Println("find exists user twitter auth error: ", err.Error())
	}
	authorIDs := make([]string, 0)
	for _, user := range users {
		mises_tweet := getMisesTweet(user.Misesid, misesTweetsMap, existUserTwitterAuths)
		//mises_tweet := misesTweetsMap[user.Misesid]
		twitter_user_id := mises_tweet.AuthorID
		if checkMisesidOrTwitterUserIdIsExists(user.Misesid, tweetAuthors[twitter_user_id], existUserTwitterAuths) {
			continue
		}
		userTwitter := &models.UserTwitterAuth{
			UID:     user.UID,
			Misesid: user.Misesid,
			TweetInfo: &models.TweetInfo{
				TweetID:   mises_tweet.TweetID,
				Text:      mises_tweet.TweetText,
				CreatedAt: mises_tweet.CreatedAt,
			},
			TwitterUserId: twitter_user_id,
			TwitterUser:   tweetAuthors[twitter_user_id],
			CreatedAt:     time.Now(),
		}
		userTwitters = append(userTwitters, userTwitter)
		authorIDs = append(authorIDs, twitter_user_id)
	}
	if len(userTwitters) == 0 {
		return
	}
	//insert
	err1 := models.CreateUserTwitterAuthMany(ctx, userTwitters)
	if err1 != nil {
		fmt.Println("insert user twitter auth error: ", err1.Error())
	}
}
func isUseAuthorID(twitter_user_id string, authorIDs []string) bool {
	is_use := false
	for _, authorID := range authorIDs {
		if authorID == twitter_user_id {
			return true
		}
	}
	return is_use
}

func getMisesTweet(misesid string, misesTweetsMap map[string][]misesTweet, existUserTwitterAuths []*models.UserTwitterAuth) misesTweet {
	mises_tweets := misesTweetsMap[misesid]
	num := len(mises_tweets)
	res := mises_tweets[0]
	if num > 1 {
		for _, mises_tweet := range mises_tweets {
			twitter_user_id := mises_tweet.AuthorID
			is_exist := false
			for _, exists := range existUserTwitterAuths {
				if twitter_user_id == exists.TwitterUserId {
					is_exist = true
					break
				}
			}
			if !is_exist {
				res = mises_tweet
			}
		}

	}
	return res
}

func checkMisesidOrTwitterUserIdIsExists(misesid string, twitter_user *models.TwitterUser, existUserTwitterAuths []*models.UserTwitterAuth) bool {
	for _, exists := range existUserTwitterAuths {
		if (misesid == exists.Misesid && IsValidTwitterUser(exists.TwitterUser)) || twitter_user.TwitterUserId == exists.TwitterUserId {
			return true
		}
	}
	return false
}

func IsValidTwitterUser(twitter_user *models.TwitterUser) (is_valid bool) {
	validRegisterDate = env.Envs.VALID_TWITTER_REGISTER_DATE
	timeFormat := "2006-01-02"
	st, _ := time.Parse(timeFormat, validRegisterDate)
	vt := st.Unix()
	twitterUserCreatedAt := twitter_user.CreatedAt.Unix()
	if vt >= twitterUserCreatedAt {
		is_valid = true
	}
	return is_valid
}

func getMisesIdByTweetText(text string) (string, error) {
	sep := "?MisesID="
	arr := strings.Split(text, sep)
	if len(arr) < 2 {
		return "", errors.New("invalid misesid")
	}
	return addMisesidProfix(arr[1]), nil
}

func addMisesidProfix(misesid string) string {
	if !strings.HasPrefix(misesid, misesidPrefix) {
		return misesidPrefix + misesid
	}
	return misesid
}
func removeMisesidProfix(misesid string) string {
	if strings.HasPrefix(misesid, misesidPrefix) {
		return strings.TrimPrefix(misesid, misesidPrefix)
	}
	return misesid
}

func getTwitterList(ctx context.Context, in *TweetsIn) (*types.SearchTweetsRecentResponse, error) {

	params := &types.SearchTweetsRecentParams{
		Query:     in.Query,
		StartTime: in.StartTime,
		EndTime:   in.EndTime,
		Expansions: fields.ExpansionList{
			fields.ExpansionAuthorID,
		},
		UserFields: fields.UserFieldList{
			fields.UserFieldCreatedAt,
			fields.UserFieldPublicMetrics,
		},
		TweetFields: fields.TweetFieldList{
			fields.TweetFieldCreatedAt,
			fields.TweetFieldEntities,
		},
		MaxResults: 20,
	}
	if in.NextToken != "" {
		params.NextToken = in.NextToken
	}
	tr, err := tweets.SearchTweetsRecent(context.Background(), twiClient, params)
	if err != nil {
		fmt.Println("search tweet recent error: ", err.Error())
		return nil, err
	}
	return tr, nil
}

//
func GetShareTweetUrl(ctx context.Context, uid uint64) (string, error) {
	//find user
	user, err := models.FindUser(ctx, uid)
	if err != nil {
		return "", err
	}
	var tweetUrl string
	misesid := removeMisesidProfix(user.Misesid)
	twitterUrl := "https://twitter.com/intent/tweet?text="
	text := "join us and get 3% airdrop! \n\nhttps://www.mises.site/download?MisesID=" + misesid + " \n\n" + tweeTtag + " #Decentralized  #SocialMedia"
	tweetUrl = twitterUrl + url.QueryEscape(text)

	return tweetUrl, nil
}
