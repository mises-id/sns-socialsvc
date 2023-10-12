package session

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/services/channel_user"
	"github.com/mises-id/sns-socialsvc/config/env"
	"github.com/mises-id/sns-socialsvc/lib/codes"
	"github.com/mises-id/sns-socialsvc/lib/mises"
	"github.com/sirupsen/logrus"
)

var (
	secret      = env.Envs.JWTSecret
	misesClient mises.Client
)

type (
	RegisterCallback struct {
	}
	referrerData struct {
		utm_source string
	}

	SignInParams struct {
		Auth      string
		Referrer  string
		UserAgent *models.UserAgent
	}
)

func SignIn(ctx context.Context, params *SignInParams) (string, bool, error) {
	misesid, pubkey, err := misesClient.Auth(params.Auth)
	if err != nil {
		logrus.Errorf("mises verify error: %v", err)
		return "", false, codes.ErrAuthorizeFailed
	}
	user, created, err := models.FindOrCreateUserByMisesid(ctx, misesid, pubkey)
	if err != nil {
		return "", created, err
	}
	//signin after
	signinAfter(ctx, user, params)
	if !user.OnChain && len(pubkey) > 0 {
		chainUserRegister(ctx, misesid, pubkey)
	}
	referrer := params.Referrer
	//referrer not empty
	if referrer != "" && user.ChannelID.IsZero() && (user.CreatedAt.Unix()+24*60*60-time.Now().UTC().Unix()) > 0 {
		err := models.InsertReferrer(ctx, user.UID, referrer)
		if err != nil {
			fmt.Println("insert referrer error: ", err.Error())
		}
		ref, err := handleReferrer(referrer)
		if err != nil {
			fmt.Printf("uid[%d], referrer[%s], error:%s\n ", user.UID, referrer, err.Error())
		} else {
			err = channel_user.CreateChannelUser(ctx, user.UID, ref.utm_source)
			if err != nil {
				fmt.Printf("uid[%d], referrer[%s], error:%s\n ", user.UID, referrer, err.Error())
			}
		}
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":         user.UID,
		"misesid":     user.Misesid,
		"username":    user.Username,
		"eth_address": strings.ToLower(user.EthAddress),
		"exp":         time.Now().Add(env.Envs.TokenDuration).Unix(),
	})
	token, err := at.SignedString([]byte(secret))
	return token, created, err
}

func signinAfter(ctx context.Context, user *models.User, params *SignInParams) error {
	return models.CreateUserLoginLog(ctx, user.UID, params.UserAgent)
}

func Auth(ctx context.Context, authToken string) (*models.User, error) {
	claim, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		if err.Error() == "Token is expired" {
			return nil, codes.ErrTokenExpired
		}
		return nil, err
	}
	mapClaims := claim.Claims.(jwt.MapClaims)
	return &models.User{
		UID:        uint64(mapClaims["uid"].(float64)),
		Misesid:    mapClaims["misesid"].(string),
		Username:   mapClaims["username"].(string),
		EthAddress: mapClaims["eth_address"].(string),
	}, nil
}

func handleReferrer(referrer string) (*referrerData, error) {
	referrer, _ = url.QueryUnescape(referrer)
	ref := &referrerData{}
	params, _ := url.ParseQuery(referrer)
	utm_sources, ok := params["utm_source"]
	if ok && len(utm_sources) > 0 {
		ref.utm_source = utm_sources[0]
	} else {
		return nil, errors.New("invalid referrer")
	}
	return ref, nil
}

func chainUserRegister(ctx context.Context, misesid, pubkey string) {
	chainUser := &models.ChainUser{
		Misesid: misesid,
		Pubkey:  pubkey,
	}
	err := models.CreateChainUser(ctx, chainUser)
	if err != nil {
		fmt.Printf("mises[%s],pubkey[%s] user register chain error:%s \n", misesid, pubkey, err.Error())
	}
}

func SetupMisesClient() {
	misesClient = mises.New()
}

func MockMisesClient(mock mises.Client) {
	misesClient = mock
}
