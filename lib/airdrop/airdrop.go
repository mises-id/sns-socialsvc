package airdrop

import (
	"fmt"

	"github.com/mises-id/sdk"
	"github.com/mises-id/sdk/types"
	"github.com/mises-id/sns-socialsvc/config/env"
)

type (
	IClient interface {
		SetListener(listener types.MisesAppCmdListener)
		RunSync(uid string, pubkey string, coin int64) error
		RunAsync(uid string, pubkey string, coin int64) error
	}

	Client struct {
		app types.MApp
	}
)

func (c Client) SetListener(listener types.MisesAppCmdListener) {
	c.app.SetListener(listener)
}
func (c Client) RunSync(uid string, pubkey string, coin int64) error {
	return c.app.RunSync(c.app.NewFaucetCmd(uid, pubkey, coin))
}
func (c Client) RunAsync(uid string, pubkey string, coin int64) error {
	return c.app.RunAsync(c.app.NewFaucetCmd(uid, pubkey, coin), false)
}

func New() IClient {
	if env.Envs.DebugAirdropPrefix != "" {
		return &Client{
			app: nil,
		}
	}
	mo := sdk.MSdkOption{
		ChainID:    env.Envs.MisesChainID,
		Debug:      true,
		PassPhrase: "mises.site",
	}
	appinfo := types.NewMisesAppInfoReadonly(
		"Mises Faucet",
		"https://www.mises.site",
		"https://home.mises.site",
		[]string{"mises.site"},
		"Mises Network",
	)
	_, app := sdk.NewSdkForApp(mo, appinfo)
	fmt.Println("new sdk for app airdrop success")
	client := &Client{
		app: app,
	}
	return client
}
