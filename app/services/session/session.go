package session

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mises-id/sdk/types"
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

type (
	RegisterCallback struct {
		retryNum int
	}
)

func SignIn(ctx context.Context, auth string) (string, bool, error) {
	fmt.Println("misesClient: ", misesClient)
	misesid, pubkey, err := misesClient.Auth(auth)
	if err != nil {
		logrus.Errorf("mises verify error: %v", err)
		return "", false, codes.ErrAuthorizeFailed
	}
	user, created, err := models.FindOrCreateUserByMisesid(ctx, misesid)
	if err != nil {
		return "", created, err
	}
	if created && len(pubkey) > 0 {
		misesChainRegister(misesid, pubkey, 0)
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":      user.UID,
		"misesid":  user.Misesid,
		"username": user.Username,
		"exp":      time.Now().Add(env.Envs.TokenDuration).Unix(),
	})
	token, err := at.SignedString([]byte(secret))
	return token, created, err
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

func misesChainRegister(misesid, pubkey string, num int) {
	fmt.Printf("mises[%s] user register chain \n", misesid)
	misesClient.SetListener(&RegisterCallback{num})
	err1 := misesClient.Register(misesid, pubkey)
	if err1 != nil {
		fmt.Printf("mises[%s] user register chain error:%s \n", misesid, err1.Error())
	}
}

func (cb *RegisterCallback) OnTxGenerated(cmd types.MisesAppCmd) {
	misesid := cmd.MisesUID()
	fmt.Printf("Mises[%s] User Register OnTxGenerated\n", misesid)

}
func (cb *RegisterCallback) OnSucceed(cmd types.MisesAppCmd) {
	misesid := cmd.MisesUID()
	fmt.Printf("Mises[%s] User Register OnSucceed\n", misesid)

}
func (cb *RegisterCallback) OnFailed(cmd types.MisesAppCmd) {
	misesid := cmd.MisesUID()
	fmt.Printf("Mises[%s] User Register OnFailed\n", misesid)
	pubkey := cmd.PubKey()
	if cb.retryNum < 5 && misesid != "" && pubkey != "" {
		cb.retryNum++
		fmt.Printf("Mises[%s] User Register Retry\n", misesid)
		misesChainRegister(misesid, pubkey, cb.retryNum)
	}
}

func SetupMisesClient() {
	misesClient = mises.New()
}

func MockMisesClient(mock mises.Client) {
	misesClient = mock
}
