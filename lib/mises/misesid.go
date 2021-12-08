package mises

import (
	"log"
	"strings"

	"github.com/mises-id/sdk"
	"github.com/mises-id/sdk/types"
	"github.com/mises-id/sdk/user"
	"github.com/mises-id/sns-socialsvc/config/env"
)

func init() {
	if env.Envs.MisesTestEndpoint != "" {
		if err := user.SetTestEndpoint(env.Envs.MisesTestEndpoint); err != nil {
			log.Fatal("init mises sdk test endpoint error")
		}
	}
}

type User struct {
	ID string
}

type Client interface {
	Auth(auth string) (string, error)
}

type ClientImpl struct {
	client types.MSdk
}

func (c *ClientImpl) Auth(auth string) (string, error) {
	// just for staging environment
	if env.Envs.DebugMisesPrefix != "" {
		arr := strings.Split(auth, ":")
		if len(arr) > 1 {
			return arr[1], nil
		}
	}
	return c.client.VerifyLogin(auth)
}

func New() Client {
	return &ClientImpl{
		client: sdk.NewSdkForApp(sdk.MSdkOption{}),
	}
}
