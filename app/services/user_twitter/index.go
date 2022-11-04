package user_twitter

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/fields"
	"github.com/michimani/gotwi/resources"
	"github.com/michimani/gotwi/tweets"
	tweetsType "github.com/michimani/gotwi/tweets/types"
	"github.com/michimani/gotwi/users"
	"github.com/michimani/gotwi/users/types"
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/search"
	"github.com/mises-id/sns-socialsvc/config/env"
	"github.com/mises-id/sns-socialsvc/lib/codes"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	callbackBase                 = "https://api.alb.mises.site/api/v1/twitter/callback"
	RequestTokenEndpoint         = "https://api.twitter.com/oauth/request_token"
	AccessTokenEndpoint          = "https://api.twitter.com/oauth/access_token"
	AuthEndpoint                 = "https://api.twitter.com/oauth/authorize"
	OAuthVersion10               = "1.0"
	OAuthSignatureMethodHMACSHA1 = "HMAC-SHA1"
	oauth1header                 = `OAuth oauth_callback="%s",oauth_consumer_key="%s",oauth_nonce="%s",oauth_signature="%s",oauth_signature_method="%s",oauth_timestamp="%s",oauth_token="%s",oauth_version="%s"`
)

var (
	OAuthConsumerKey    = ""
	OAuthConsumerSecret = ""
	OAuthToken          = ""
	OAuthTokenSecret    = ""
	targetTwitterId     = "1442753558311424001"
)

const (
	followTwitterNum = 45
)

type (
	CreateOAuthSignatureInput struct {
		HTTPMethod       string
		RawEndpoint      string
		OAuthConsumerKey string
		OAuthToken       string
		SigningKey       string
		ParameterMap     map[string]string
	}
	CreateOAuthSignatureOutput struct {
		OAuthNonce           string
		OAuthSignatureMethod string
		OAuthTimestamp       string
		OAuthVersion         string
		OAuthSignature       string
	}
	Endpoint     string
	EndpointInfo struct {
		Raw                      string
		Base                     string
		EncodedQueryParameterMap map[string]string
	}
	AirdropInfoOutput struct {
		Twitter *models.UserTwitterAuth
		Airdrop *models.Airdrop
	}
)

func init() {
	OAuthConsumerKey = env.Envs.GOTWI_API_KEY
	OAuthConsumerSecret = env.Envs.GOTWI_API_KEY_SECRET
}

//get twitter auth url
func GetTwitterAuthUrl(ctx context.Context, uid uint64) (string, error) {
	callback := fmt.Sprintf("%s?uid=%d", callbackBase, uid)
	auth, err := RequestToken(ctx, callback)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s?%s", AuthEndpoint, auth)
	return url, nil
}

//get airdrop info
func GetAirdropInfo(ctx context.Context, uid uint64) (*AirdropInfoOutput, error) {
	user_twitter, err := models.FindUserTwitterAuthByUid(ctx, uid)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if user_twitter != nil {
		user_twitter.IsValid = IsValidTwitterUser(user_twitter.TwitterUser)
		if user_twitter.IsValid {
			user_twitter.Amount = GetTwitterAirdropCoin(ctx, user_twitter)
		}
	}
	airdrop, err := models.FindAirdropByUid(ctx, uid)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	res := &AirdropInfoOutput{
		Twitter: user_twitter,
		Airdrop: airdrop,
	}
	return res, nil

}

//receive airdrop
func ReceiveAirdrop(ctx context.Context, uid uint64, tweet string) error {
	//check twitter auth
	user_twitter, err := models.FindUserTwitterAuthByUid(ctx, uid)
	if err != nil {
		return codes.ErrForbidden.Newf("Twitter is unauthorized")
	}
	if !IsValidTwitterUser(user_twitter.TwitterUser) {
		return codes.ErrForbidden.Newf("Twitter is invalid.")
	}
	if !models.GetAirdropStatus(ctx) {
		return codes.ErrForbidden.Newf("Airdrop end")
	}
	airdrop, _ := models.FindAirdropByUid(ctx, uid)
	if airdrop != nil {
		return codes.ErrForbidden.Newf("Repeat to receive")
	}
	if tweet == "" {
		tweet = "hi"
	}
	//send tweet
	if err := sendTweet(ctx, user_twitter, tweet); err != nil {
		fmt.Printf("uid[%d] send tweet err:%s ", uid, err.Error())
		//return codes.ErrForbidden.Newf("Send twitter failed.")
	}
	//to follow twitter
	user_twitter.IsFollowed = false
	if err := models.UpdateUserTwitterAuthFollew(ctx, user_twitter); err != nil {
		fmt.Printf("[%s]uid[%d] UpdateUserTwitterAuthFollew err:%s ", time.Now().String(), uid, err.Error())
	}
	//create airdrop order
	if err := createAirdrop(ctx, user_twitter); err != nil {
		return err
	}
	//cancel auth token
	return nil
}

//send tweet
func sendTweet(ctx context.Context, user_twitter *models.UserTwitterAuth, tweet string) error {

	if user_twitter.OauthToken == "" || user_twitter.OauthTokenSecret == "" {
		return codes.ErrForbidden.Newf("OAuthToken and OAuthTokenSecret is required")
	}
	transport := &http.Transport{Proxy: setProxy()}
	client := &http.Client{Transport: transport}
	in := &gotwi.NewGotwiClientInput{
		HTTPClient:           client,
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           user_twitter.OauthToken,
		OAuthTokenSecret:     user_twitter.OauthTokenSecret,
	}
	twitter_client, err := gotwi.NewGotwiClient(in)
	if err != nil {
		return err
	}
	params := &tweetsType.ManageTweetsPostParams{
		Text: &tweet,
	}
	_, err = tweets.ManageTweetsPost(ctx, twitter_client, params)

	return err
}

//follow twitter
func FollowTwitter(ctx context.Context) error {
	return runFollowTwitter(ctx)
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
	fmt.Println("runFollowTwitter num: ", num)
	//do list
	for _, user_twitter := range user_twitter_list {
		//to follow
		if err := apiFollowTwitterUser(ctx, user_twitter, targetTwitterId); err != nil {
			fmt.Printf("[%s]twitter_user_id[%s],apiFollowTwitterUser error:%s", time.Now().String(), user_twitter.TwitterUserId, err.Error())
		}
		user_twitter.IsFollowed = true
		if err = models.UpdateUserTwitterAuthFollew(ctx, user_twitter); err != nil {
			fmt.Printf("[%s]twitter_user_id[%s],UpdateUserTwitterAuthFollew error:%s", time.Now().String(), user_twitter.TwitterUserId, err.Error())
		}
	}

	return nil
}

//apiFollowTwitterUser
func apiFollowTwitterUser(ctx context.Context, user_twitter *models.UserTwitterAuth, target_user_id string) error {
	if user_twitter == nil {
		return errors.New("user_twitter is null")
	}
	if user_twitter.OauthToken == "" || user_twitter.OauthTokenSecret == "" {
		return codes.ErrForbidden.Newf("OAuthToken and OAuthTokenSecret is required")
	}
	transport := &http.Transport{Proxy: setProxy()}
	client := &http.Client{Transport: transport}
	in := &gotwi.NewGotwiClientInput{
		HTTPClient:           client,
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           user_twitter.OauthToken,
		OAuthTokenSecret:     user_twitter.OauthTokenSecret,
	}
	twitter_client, err := gotwi.NewGotwiClient(in)
	if err != nil {
		return err
	}
	params := &types.FollowsFollowingPostParams{
		ID:           user_twitter.TwitterUserId,
		TargetUserID: &target_user_id,
	}
	_, err = users.FollowsFollowingPost(ctx, twitter_client, params)

	return err
}

func createAirdrop(ctx context.Context, user_twitter *models.UserTwitterAuth) error {
	airdropAdd := &models.Airdrop{
		UID:       user_twitter.UID,
		Misesid:   user_twitter.Misesid,
		Status:    enum.AirdropDefault,
		Type:      enum.AirdropTwitter,
		Coin:      GetTwitterAirdropCoin(ctx, user_twitter),
		TxID:      "",
		CreatedAt: time.Now(),
	}
	_, err := models.CreateAirdrop(ctx, airdropAdd)
	if err != nil {
		return err
	}
	return nil
}

func getTwitterCallbackUrl(code, username, misesid string) string {
	return env.Envs.TwitterAuthSuccessCallback + "?code=" + code + "&username=" + username + "&misesid=" + misesid
}

//twitter auth callback
func TwitterCallback(ctx context.Context, uid uint64, oauth_token, oauth_verifier string) string {

	var (
		callback0 string = getTwitterCallbackUrl("0", "", "")
		callback1 string = getTwitterCallbackUrl("1", "", "")
		callback2 string = getTwitterCallbackUrl("2", "", "")
	)
	if oauth_token == "" || oauth_verifier == "" {
		fmt.Printf("Oauth_token[%s],oauth_verifier[%s] err", oauth_token, oauth_verifier)
		return callback2
	}
	user, err := models.FindUser(ctx, uid)
	if err != nil {
		fmt.Println("Twitter callback find user err: ", err.Error())
		return callback2
	}
	userMisesid := user.Misesid
	callback2 = getTwitterCallbackUrl("2", "", userMisesid)
	//find twitter user
	access_token, err := AccessToken(ctx, oauth_token, oauth_verifier)
	if err != nil {
		fmt.Println("Twitter callback access token err: ", err.Error())
		return callback2
	}
	params, _ := url.ParseQuery(access_token)
	user_ids, ok := params["user_id"]
	if !ok || len(user_ids) <= 0 {
		fmt.Println("Twitter callback user_id err: ", err.Error())
		return callback2
	}
	oauth_tokens, ok := params["oauth_token"]
	oauth_token_secrets, ok := params["oauth_token_secret"]
	twitter_user_id := user_ids[0]
	oauth_token_new := oauth_tokens[0]
	oauth_token_secret := oauth_token_secrets[0]
	//check twitter_user_id
	twitter_auth, err := models.FindUserTwitterAuthByTwitterUserId(ctx, twitter_user_id)

	if twitter_auth != nil && twitter_auth.UID != uid {
		callback1 = getTwitterCallbackUrl("1", twitter_auth.TwitterUser.UserName, userMisesid)
		fmt.Println("FindUserTwitterAuthByTwitterUserId exist ")
		return callback1
	}
	//check uid
	user_twitter, err := models.FindUserTwitterAuthByUid(ctx, uid)
	if err != nil && err != mongo.ErrNoDocuments {
		fmt.Println("Twitter callback FindUserTwitterAuthByUid err: ", err.Error())
		return callback2
	}
	twitter_user, err := getTwitterUserById(ctx, twitter_user_id)
	if err != nil {
		fmt.Println("Twitter callback getTwitterUserById err: ", err.Error())
		return callback2
	}
	callback0 = getTwitterCallbackUrl("0", *twitter_user.Username, userMisesid)
	TwitterUser := &models.TwitterUser{
		TwitterUserId:  *twitter_user.ID,
		UserName:       *twitter_user.Username,
		Name:           *twitter_user.Name,
		CreatedAt:      *twitter_user.CreatedAt,
		FollowersCount: uint64(*twitter_user.PublicMetrics.FollowersCount),
		TweetCount:     uint64(*twitter_user.PublicMetrics.TweetCount),
	}
	//check airdrop
	airdrop, err := models.FindAirdropByUid(ctx, uid)

	if user_twitter == nil {
		//create
		if airdrop != nil {
			fmt.Println("Twitter callback airdrop exist")
			return callback0
		}
		add := &models.UserTwitterAuth{
			UID:              uid,
			Misesid:          user.Misesid,
			TwitterUserId:    twitter_user_id,
			TwitterUser:      TwitterUser,
			OauthToken:       oauth_token_new,
			OauthTokenSecret: oauth_token_secret,
		}
		err = models.CreateUserTwitterAuth(ctx, add)

	} else {
		//update
		user_twitter.OauthToken = oauth_token_new
		user_twitter.OauthTokenSecret = oauth_token_secret
		if airdrop == nil {
			user_twitter.TwitterUser = TwitterUser
			user_twitter.TwitterUserId = twitter_user_id
		}
		err = models.UpdateUserTwitterAuth(ctx, user_twitter)
	}
	if err != nil {
		fmt.Println("Twitter callback save err: ", err.Error())
	}
	return callback0
}
func TwitterCallbackOld(ctx context.Context, uid uint64, oauth_token, oauth_verifier string) string {

	var (
		callback0 string = getTwitterCallbackUrl("0", "", "")
		callback1 string = getTwitterCallbackUrl("1", "", "")
		callback2 string = getTwitterCallbackUrl("2", "", "")
	)
	if oauth_token == "" || oauth_verifier == "" {
		fmt.Printf("Oauth_token[%s],oauth_verifier[%s] err", oauth_token, oauth_verifier)
		return callback2
	}
	user, err := models.FindUser(ctx, uid)
	if err != nil {
		fmt.Println("Twitter callback find user err: ", err.Error())
		return callback2
	}
	userMisesid := user.Misesid
	callback2 = getTwitterCallbackUrl("2", "", userMisesid)
	//find twitter user
	access_token, err := AccessToken(ctx, oauth_token, oauth_verifier)
	if err != nil {
		fmt.Println("Twitter callback access token err: ", err.Error())
		return callback2
	}
	params, _ := url.ParseQuery(access_token)
	user_ids, ok := params["user_id"]
	if !ok || len(user_ids) <= 0 {
		fmt.Println("Twitter callback user_id err: ", err.Error())
		return callback2
	}
	oauth_tokens, ok := params["oauth_token"]
	oauth_token_secrets, ok := params["oauth_token_secret"]
	twitter_user_id := user_ids[0]
	oauth_token_new := oauth_tokens[0]
	oauth_token_secret := oauth_token_secrets[0]
	//check twitter_user_id
	twitter_auth, err := models.FindUserTwitterAuthByTwitterUserId(ctx, twitter_user_id)

	if twitter_auth != nil && twitter_auth.UID != uid {
		callback1 = getTwitterCallbackUrl("1", twitter_auth.TwitterUser.UserName, userMisesid)
		fmt.Println("FindUserTwitterAuthByTwitterUserId exist ")
		return callback1
	}
	//check uid
	user_twitter, err := models.FindUserTwitterAuthByUid(ctx, uid)
	if err != nil && err != mongo.ErrNoDocuments {
		fmt.Println("Twitter callback FindUserTwitterAuthByUid err: ", err.Error())
		return callback2
	}
	twitter_user, err := getTwitterUserById(ctx, twitter_user_id)
	if err != nil {
		fmt.Println("Twitter callback getTwitterUserById err: ", err.Error())
		return callback2
	}
	callback0 = getTwitterCallbackUrl("0", *twitter_user.Username, userMisesid)
	TwitterUser := &models.TwitterUser{
		TwitterUserId:  *twitter_user.ID,
		UserName:       *twitter_user.Username,
		Name:           *twitter_user.Name,
		CreatedAt:      *twitter_user.CreatedAt,
		FollowersCount: uint64(*twitter_user.PublicMetrics.FollowersCount),
		TweetCount:     uint64(*twitter_user.PublicMetrics.TweetCount),
	}
	//check airdrop
	airdrop, err := models.FindAirdropByUid(ctx, uid)

	if user_twitter == nil {
		//create
		if airdrop != nil {
			fmt.Println("Twitter callback airdrop exist")
			return callback0
		}
		add := &models.UserTwitterAuth{
			UID:              uid,
			Misesid:          user.Misesid,
			TwitterUserId:    twitter_user_id,
			TwitterUser:      TwitterUser,
			OauthToken:       oauth_token_new,
			OauthTokenSecret: oauth_token_secret,
		}
		err = models.CreateUserTwitterAuth(ctx, add)

	} else {
		//update
		user_twitter.OauthToken = oauth_token_new
		user_twitter.OauthTokenSecret = oauth_token_secret
		if airdrop == nil {
			user_twitter.TwitterUser = TwitterUser
			user_twitter.TwitterUserId = twitter_user_id
		}
		err = models.UpdateUserTwitterAuth(ctx, user_twitter)
	}
	if err != nil {
		fmt.Println("Twitter callback save err: ", err.Error())
	}
	return callback0
}

func getTwitterUserById(ctx context.Context, twitter_user_id string) (*resources.User, error) {
	transport := &http.Transport{Proxy: setProxy()}
	client := &http.Client{Transport: transport}
	in := &gotwi.NewGotwiClientInput{
		HTTPClient:           client,
		AuthenticationMethod: gotwi.AuthenMethodOAuth2BearerToken,
	}
	twitter_client, err := gotwi.NewGotwiClient(in)
	if err != nil {
		return nil, err
	}
	params := &types.UserLookupIDParams{
		ID: twitter_user_id,
		UserFields: fields.UserFieldList{
			fields.UserFieldCreatedAt,
			fields.UserFieldPublicMetrics,
		},
	}
	tr, err := users.UserLookupID(ctx, twitter_client, params)
	if err != nil {
		fmt.Println("User look up id error: ", err.Error())
		return nil, err
	}
	return &tr.Data, nil
}

func setProxy() func(*http.Request) (*url.URL, error) {
	return func(_ *http.Request) (*url.URL, error) {
		return nil, nil
		return url.Parse("http://127.0.0.1:1087")
	}
}

//get twitter auth request_token
func RequestToken(ctx context.Context, callback string) (string, error) {

	api := fmt.Sprintf("%s?oauth_callback=%s", RequestTokenEndpoint, callback)
	transport := &http.Transport{Proxy: setProxy()}
	client := &http.Client{Transport: transport}
	req, _ := http.NewRequest("POST", api, nil)
	ParameterMap := map[string]string{
		"oauth_callback": callback,
	}
	in := &CreateOAuthSignatureInput{
		HTTPMethod:       req.Method,
		RawEndpoint:      req.URL.String(),
		OAuthConsumerKey: OAuthConsumerKey,
		OAuthToken:       OAuthToken,
		SigningKey:       getSignKey(),
		ParameterMap:     ParameterMap,
	}

	out, err := CreateOAuthSignature(in)
	if err != nil {
		return "", err
	}
	auth := fmt.Sprintf(oauth1header,
		url.QueryEscape(callback),
		url.QueryEscape(OAuthConsumerKey),
		url.QueryEscape(out.OAuthNonce),
		url.QueryEscape(out.OAuthSignature),
		url.QueryEscape(out.OAuthSignatureMethod),
		url.QueryEscape(out.OAuthTimestamp),
		url.QueryEscape(OAuthToken),
		url.QueryEscape(out.OAuthVersion),
	)
	req.Header.Add("Authorization", auth)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", errors.New(res.Status)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	return string(body), nil
}
func AccessToken(ctx context.Context, oauth_token, oauth_verifier string) (string, error) {

	api := fmt.Sprintf("%s?oauth_token=%s&oauth_verifier=%s", AccessTokenEndpoint, oauth_token, oauth_verifier)

	transport := &http.Transport{Proxy: setProxy()}
	client := &http.Client{Transport: transport}
	req, _ := http.NewRequest("POST", api, nil)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", errors.New(res.Status)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	return string(body), nil
}

func getSignKey() string {
	return fmt.Sprintf("%s&%s", OAuthConsumerSecret, OAuthTokenSecret)
}

func CreateOAuthSignature(in *CreateOAuthSignatureInput) (*CreateOAuthSignatureOutput, error) {
	out := CreateOAuthSignatureOutput{
		OAuthSignatureMethod: OAuthSignatureMethodHMACSHA1,
		OAuthVersion:         OAuthVersion10,
	}
	nonce, err := generateOAthNonce()
	if err != nil {
		return nil, err
	}
	out.OAuthNonce = nonce

	ts := fmt.Sprintf("%d", time.Now().Unix())
	out.OAuthTimestamp = ts
	endpointBase := endpointBase(in.RawEndpoint)

	parameterString := createParameterString(nonce, ts, in)
	sigBase := createSignatureBase(in.HTTPMethod, endpointBase, parameterString)
	sig, err := calculateSignature(sigBase, in.SigningKey)
	if err != nil {
		return nil, err
	}
	out.OAuthSignature = sig

	return &out, nil
}
func generateOAthNonce() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	nonce := base64.StdEncoding.EncodeToString(key)
	symbols := []string{"+", "/", "="}
	for _, s := range symbols {
		nonce = strings.Replace(nonce, s, "", -1)
	}
	return nonce, nil
}

func endpointBase(e string) string {
	queryIdx := strings.Index(e, "?")
	if queryIdx < 0 {
		return e
	}

	return e[:queryIdx]
}

func (e Endpoint) String() string {
	return string(e)
}

func (e Endpoint) Detail() (*EndpointInfo, error) {
	d := EndpointInfo{
		Raw:                      e.String(),
		EncodedQueryParameterMap: map[string]string{},
	}

	queryIdx := strings.Index(e.String(), "?")
	if queryIdx < 0 {
		d.Base = string(e)
		return &d, nil
	}

	d.Base = e.String()[:queryIdx]
	queryPart := e.String()[queryIdx+1:]
	paramsPairs := strings.Split(queryPart, "&")
	for _, pp := range paramsPairs {
		keyValue := strings.Split(pp, "=")
		var err error
		v := ""
		if len(keyValue) == 2 {
			v, err = url.QueryUnescape(keyValue[1])
			if err != nil {
				return nil, err
			}
		}
		d.EncodedQueryParameterMap[keyValue[0]] = v
	}

	return &d, nil
}

func createParameterString(nonce, ts string, in *CreateOAuthSignatureInput) string {
	qv := url.Values{}
	for k, v := range in.ParameterMap {
		qv.Add(k, v)
	}

	qv.Add("oauth_consumer_key", in.OAuthConsumerKey)
	qv.Add("oauth_nonce", nonce)
	qv.Add("oauth_signature_method", OAuthSignatureMethodHMACSHA1)
	qv.Add("oauth_timestamp", ts)
	qv.Add("oauth_token", in.OAuthToken)
	qv.Add("oauth_version", OAuthVersion10)

	encoded := qv.Encode()
	encoded = regexp.MustCompile(`([^%])(\+)`).ReplaceAllString(encoded, "$1%20")
	return encoded
}

func createSignatureBase(method, endpointBase, parameterString string) string {
	return fmt.Sprintf(
		"%s&%s&%s",
		url.QueryEscape(strings.ToUpper(method)),
		url.QueryEscape(endpointBase),
		url.QueryEscape(parameterString),
	)
}

func calculateSignature(base, key string) (string, error) {
	b := []byte(key)
	h := hmac.New(sha1.New, b)
	_, err := io.WriteString(h, base)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
