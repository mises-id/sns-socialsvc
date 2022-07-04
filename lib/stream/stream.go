package stream

import (
	"fmt"

	cosmos_types "github.com/cosmos/cosmos-sdk/types"
	"github.com/mises-id/sdk"
	"github.com/mises-id/sdk/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	StreamClient IStreamClient
)

func init() {

}

type IStreamClient interface {
	SetListener(listener types.MisesAppCmdListener)
	StartEventStreaming(listener types.MisesEventStreamingListener) error
	ParseEvent(header *tmtypes.EventDataNewBlockHeader, tx *tmtypes.EventDataTx) (*cosmos_types.TxResponse, error)
}

type StreamClientImpl struct {
	app types.MApp
}

func (c *StreamClientImpl) SetListener(listener types.MisesAppCmdListener) {
	c.app.SetListener(listener)
}
func (c *StreamClientImpl) StartEventStreaming(listener types.MisesEventStreamingListener) error {
	return c.app.StartEventStreaming(listener)
}
func (c *StreamClientImpl) ParseEvent(header *tmtypes.EventDataNewBlockHeader, tx *tmtypes.EventDataTx) (*cosmos_types.TxResponse, error) {
	return c.app.ParseEvent(header, tx)
}

func NewStream() IStreamClient {
	/*
		if env.Envs.DebugMisesPrefix != "" {
			return &ClientImpl{
				client: nil,
				app:    nil,
			}
		} */
	/* opt := types.MSdkOption{
		ChainID: env.Envs.MisesChainID,
	} */
	opt := types.MSdkOption{
		ChainID: "mainnet",
		Debug:   true,
		//RpcURI:  "http://mises.ihuaj.com:26657",
	}
	appinfo := types.NewMisesAppInfoReadonly(
		"Mises Stream",
		"https://www.mises.site",
		"https://home.mises.site",
		[]string{"mises.site"},
		"Mises Network",
	)
	_, app := sdk.NewSdkForApp(opt, appinfo)
	fmt.Println("new sdk for app stream success")
	return &StreamClientImpl{
		app: app,
	}
}

func SetStreamClient() {
	StreamClient = NewStream()
}
