package mises

import (
	"log"

	"github.com/mises-id/sdk"
	"github.com/mises-id/sdk/misesid"
	"github.com/mises-id/sdk/types"
	"github.com/mises-id/sns-socialsvc/config/env"
)

func init() {
	if env.Envs.MisesEndpoint != "" {
		if err := misesid.SetTestEndpoint(env.Envs.MisesEndpoint); err != nil {
			log.Fatal("init mises sdk test endpoint error")
		}
	}
}

type User struct {
	ID string
}

type Client interface {
	Auth(auth string) (string, string, error)
	Register(misesUID string, pubKey string) error
}

type ClientImpl struct {
	client types.MSdk
	app    types.MApp
}

func (c *ClientImpl) Register(misesUID string, pubKey string) error {

	return c.app.RegisterUserAsync(misesUID, pubKey)
}
func (c *ClientImpl) Auth(auth string) (string, string, error) {
	// just for staging environment

	return c.client.VerifyLogin(auth)
}

func New() Client {
	opt := sdk.MSdkOption{
		ChainID: env.Envs.MisesChainID,
	}
	appinfo := misesid.NewMisesAppInfoReadonly(
		"Mises Discover'",
		"https://www.mises.site",
		"https://home.mises.site",
		[]string{"mises.site"},
		"Mises Network",
	)
	sdk, app := sdk.NewSdkForApp(opt, appinfo)
	return &ClientImpl{
		client: sdk,
		app:    app,
	}
}
