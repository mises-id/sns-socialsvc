package mises

import (
	"log"
	"strings"

	"github.com/mises-id/sdk"
	"github.com/mises-id/sdk/types"
	"github.com/mises-id/sns-socialsvc/config/env"
)

func init() {

}

type User struct {
	ID string
}

type Client interface {
	Auth(auth string) (string, string, error)
	Register(misesUID string, pubKey string) error
	SetListener(listener types.MisesAppCmdListener)
}

type ClientImpl struct {
	client types.MSdk
	app    types.MApp
}

func (c *ClientImpl) SetListener(listener types.MisesAppCmdListener) {
	c.app.SetListener(listener)
}

func (c *ClientImpl) Register(misesUID string, pubKey string) error {

	return c.app.RunAsync(
		c.app.NewRegisterUserCmd(
			misesUID,
			pubKey,
			1000000,
		), false,
	)
}
func (c *ClientImpl) Auth(auth string) (string, string, error) {
	// just for staging environment
	if env.Envs.DebugMisesPrefix != "" {
		arr := strings.Split(auth, ":")
		if len(arr) > 1 {
			return arr[1], "", nil
		}
	}

	return c.client.VerifyLogin(auth)
}
func New() Client {

	if env.Envs.DebugMisesPrefix != "" {
		return &ClientImpl{
			client: nil,
			app:    nil,
		}
	}
	opt := sdk.MSdkOption{
		ChainID: env.Envs.MisesChainID,
	}
	appinfo := types.NewMisesAppInfoReadonly(
		"Mises Discover",
		"https://www.mises.site",
		"https://home.mises.site",
		[]string{"mises.site"},
		"Mises Network",
	)
	sdk, app := sdk.NewSdkForApp(opt, appinfo)
	if env.Envs.MisesEndpoint != "" {
		if err := sdk.SetEndpoint(env.Envs.MisesEndpoint); err != nil {
			log.Fatal("init mises sdk test endpoint error")
		}
	}
	return &ClientImpl{
		client: sdk,
		app:    app,
	}
}
