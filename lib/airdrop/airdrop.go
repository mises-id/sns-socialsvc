package airdrop

import (
	"fmt"

	"github.com/mises-id/sdk"
	"github.com/mises-id/sdk/types"
	"github.com/mises-id/sns-socialsvc/config/env"
)

var (
	AirdropClient IClient
)

type (
	IClient interface {
		SetListener(listener types.MisesAppCmdListener)
		RunSync(uid string, pubkey string, coin int64) error
		RunAsync(uid string, pubkey string, coin int64, opts ...Options) error
		SetTrackID(id string) Options
	}

	Client struct {
		app types.MApp
	}
	Options func(cmd types.MisesAppCmd) types.MisesAppCmd
)

func (c Client) SetListener(listener types.MisesAppCmdListener) {
	c.app.SetListener(listener)
}
func (c Client) RunSync(uid string, pubkey string, coin int64) error {
	return c.app.RunSync(c.app.NewFaucetCmd(uid, pubkey, coin))
}
func (c Client) RunAsync(uid string, pubkey string, coin int64, opts ...Options) error {
	appcmd := c.app.NewFaucetCmd(uid, pubkey, coin)
	for _, opt := range opts {
		appcmd = opt(appcmd)
	}
	return c.app.RunAsync(appcmd, false)
}
func (c Client) SetTrackID(id string) Options {
	return func(cmd types.MisesAppCmd) types.MisesAppCmd {
		cmd.SetTrackID(id)
		return cmd
	}
}

func New() IClient {
	if env.Envs.DebugAirdropPrefix != "" {
		return &Client{
			app: nil,
		}
	}
	mo := types.MSdkOption{
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

func SetAirdropClient() {
	AirdropClient = New()
}
