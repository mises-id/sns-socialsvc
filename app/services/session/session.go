package session

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/config/env"
	"github.com/mises-id/sns-socialsvc/lib/codes"
	"github.com/mises-id/sns-socialsvc/lib/mises"
	"github.com/sirupsen/logrus"
)

var (
	secret      = env.Envs.JWTSecret
	misesClient mises.Client
)

func SignIn(ctx context.Context, auth string) (string, error) {
	fmt.Println("misesClient: ", misesClient)
	misesid, pubkey, err := misesClient.Auth(auth)
	if err != nil {
		logrus.Errorf("mises verify error: %v", err)
		return "", codes.ErrAuthorizeFailed
	}
	user, created, err := models.FindOrCreateUserByMisesid(ctx, misesid)
	if err != nil {
		return "", err
	}
	if created && len(pubkey) > 0 {
		_ = misesClient.Register(misesid, pubkey)
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":      user.UID,
		"misesid":  user.Misesid,
		"username": user.Username,
		"exp":      time.Now().Add(env.Envs.TokenDuration).Unix(),
	})
	return at.SignedString([]byte(secret))
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
		UID:      uint64(mapClaims["uid"].(float64)),
		Misesid:  mapClaims["misesid"].(string),
		Username: mapClaims["username"].(string),
	}, nil
}

func SetupMisesClient() {
	misesClient = mises.New()
}

func MockMisesClient(mock mises.Client) {
	misesClient = mock
}
