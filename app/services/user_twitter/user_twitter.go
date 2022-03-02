package user_twitter

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/fields"
	"github.com/michimani/gotwi/tweets"
	"github.com/michimani/gotwi/tweets/types"
	"github.com/mises-id/sns-socialsvc/app/models"
)

var twiClient *gotwi.GotwiClient

type (
	TweetsIn struct {
		Query      string
		StartTime  *time.Time
		EndTime    *time.Time
		NextToken  string
		MaxResults int
	}
)

func UserTwitterAuth() {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://127.0.0.1:10809")
	}
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{Transport: transport}
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
	dateNow := time.Now().UTC().AddDate(0, 0, -1)
	startTime := time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 0, 0, 0, 0, dateNow.Location())
	endTime := time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 23, 59, 59, 0, dateNow.Location())
	tweetIn := &TweetsIn{
		Query:     "#GrowWithGoogle",
		StartTime: &startTime,
		EndTime:   &endTime,
	}
	GetTwitter(context.TODO(), tweetIn)
}

func GetTwitter(ctx context.Context, in *TweetsIn) {
	res, err := GetTwitterList(ctx, in)
	if err != nil {
		fmt.Println("err:", err.Error())
	}
	//do
	TwitterAuth(ctx, res)
	//
	if res.Meta.NextToken != nil {
		in.NextToken = *res.Meta.NextToken
		GetTwitter(ctx, in)
	}
	return

}

func TwitterAuth(ctx context.Context, tweets *types.SearchTweetsRecentResponse) {
	num := gotwi.IntValue(tweets.Meta.ResultCount)
	fmt.Println("total: ", num)
	fmt.Println("meta: ", tweets.Meta)
	misesids := make([]string, 0)
	twitterUserIds := make([]string, 0)
	userTwitters := make([]*models.UserTwitterAuth, 0)
	//check tweet text
	for k, v := range tweets.Data {
		text := gotwi.StringValue(v.Text)
		misesid, err := getMisesIdByTweetText(text)
		if err != nil {
			continue
		}
		user := tweets.Includes.Users[k]
		tweetId := gotwi.StringValue(v.ID)
		twitterUserId := gotwi.StringValue(user.ID)
		misesids = append(misesids, misesid)
		twitterUserIds = append(twitterUserIds, twitterUserId)
		tweetUser := &models.TwitterUser{
			Name:           gotwi.StringValue(user.Name),
			UserName:       gotwi.StringValue(user.Username),
			FollowersCount: uint64(gotwi.IntValue(user.PublicMetrics.FollowersCount)),
			TweetCount:     uint64(gotwi.IntValue(user.PublicMetrics.TweetCount)),
			CreatedAt:      *user.CreatedAt,
		}
		userTwitter := &models.UserTwitterAuth{
			Misesid:       misesid,
			AuthTweetID:   tweetId,
			TwitterUserId: twitterUserId,
			TwitterUser:   tweetUser,
			CreatedAt:     time.Now(),
		}
		userTwitters = append(userTwitters, userTwitter)
		fmt.Println("tweets id: ", gotwi.StringValue(v.ID))
		fmt.Println("tweets create_at: ", v.CreatedAt)
	}
	//insert
	err := models.CreateUserTwitterAuthMany(ctx, userTwitters)
	if err != nil {
		fmt.Println("insert error: ", err.Error())
	}
}

func getMisesIdByTweetText(text string) (string, error) {

	return "", nil
}

func GetTwitterList(ctx context.Context, in *TweetsIn) (*types.SearchTweetsRecentResponse, error) {

	params := &types.SearchTweetsRecentParams{
		Query:     in.Query,
		StartTime: in.StartTime,
		EndTime:   in.EndTime,
		Expansions: fields.ExpansionList{
			fields.ExpansionAuthorID,
		},
		UserFields: fields.UserFieldList{
			fields.UserFieldCreatedAt,
		},
		TweetFields: fields.TweetFieldList{
			fields.TweetFieldCreatedAt,
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
func GetShareTwitterUrl(ctx context.Context, uid uint64) (string, error) {
	//find user
	user, err := models.FindUser(ctx, uid)
	if err != nil {
		return "", err
	}
	var tweetUrl string
	misesid := user.Misesid
	twitterUrl := "https://twitter.com/intent/tweet?text="
	text := `Welcome to Mises @Mises001 #mises

https://home.mises.site/home/me?misesid=` + misesid
	tweetUrl = twitterUrl + url.QueryEscape(text)

	return tweetUrl, nil
}
