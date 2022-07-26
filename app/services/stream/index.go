package stream

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankcodec "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributioncodec "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashcodec "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakecodec "github.com/cosmos/cosmos-sdk/x/staking/types"
	dbm "github.com/tendermint/tm-db"

	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/std"

	misestypes "github.com/mises-id/mises-tm/x/misestm/types"
	"github.com/mises-id/sdk/misesid"

	"github.com/tendermint/tendermint/store"
	tmtypes "github.com/tendermint/tendermint/types"

	streamLib "github.com/mises-id/sns-socialsvc/lib/stream"
)

func Run(ctx context.Context) {
	callback := &EventStreamingCallback{}
	callback.done = make(chan bool)
	callback.maxCount = 10000
	err := streamLib.StreamClient.StartEventStreaming(callback)
	if err != nil {
		fmt.Println("StartEventStreaming error: ", err.Error())

	}
	//callback.wait()

	for i := range callback.done {
		resp, err := streamLib.StreamClient.ParseEvent(callback.header, callback.tx)
		if err != nil {
			fmt.Println(i)
			fmt.Println("ParseEvent error: ", err.Error())

		} else {
			fmt.Println("h: ", callback.header)
			fmt.Println("tx: ", callback.tx)
			fmt.Println("txx: ", callback.tx.Tx)
			fmt.Printf("ParseEvent %s", resp.GetTx())
		}
	}
}

type EventStreamingCallback struct {
	done       chan bool
	eventCount int
	maxCount   int
	header     *tmtypes.EventDataNewBlockHeader
	tx         *tmtypes.EventDataTx
}

func (cb *EventStreamingCallback) OnTxEvent(t *tmtypes.EventDataTx) {
	fmt.Printf("OnTxEvent\n")
	cb.eventCount++
	cb.tx = t
	if cb.eventCount > cb.maxCount || (cb.tx != nil && cb.header != nil) {
		fmt.Printf("done1")
		cb.done <- true
	}
}
func (cb *EventStreamingCallback) OnNewBlockHeaderEvent(h *tmtypes.EventDataNewBlockHeader) {
	fmt.Printf("OnNewBlockHeaderEvent\n")
	cb.eventCount++
	cb.header = h
	if cb.eventCount > cb.maxCount || (cb.tx != nil && cb.header != nil) {
		fmt.Printf("done2")
		cb.done <- true
	}
}
func (cb *EventStreamingCallback) OnEventStreamingTerminated() {
	fmt.Printf("OnEventStreamingTerminated")
	fmt.Printf("done3")
	cb.done <- true
}
func (cb *EventStreamingCallback) wait() {
	<-cb.done
}

type intoAny interface {
	AsAny() *codectypes.Any
}

type amount struct {
	Denom  string `json:"denom"`
	Amount int64  `json:"amount"`
}
type msg struct {
	FromAddress         string    `json:"from_address"`
	ToAddress           string    `json:"to_address"`
	ValidatorSrcAddress string    `json:"validator_src_address" bson:"validator_src_address"`
	ValidatorDstAddress string    `json:"validator_dst_address" bson:"validator_dst_address"`
	Creator             string    `json:"creator"`
	Amount              []*amount `json:"amount"`
}

func Test(ctx context.Context) error {

	misesid.SetConfig()
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	authcodec.RegisterInterfaces(interfaceRegistry)
	bankcodec.RegisterInterfaces(interfaceRegistry)
	stakecodec.RegisterInterfaces(interfaceRegistry)
	distributioncodec.RegisterInterfaces(interfaceRegistry)
	slashcodec.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	misestypes.RegisterInterfaces(interfaceRegistry)
	codec := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(codec, tx.DefaultSignModes)

	//txs := "CqQBCqEBCiovbWlzZXNpZC5taXNlc3RtLnYxYmV0YTEuTXNnVXBkYXRlVXNlckluZm8ScwosbWlzZXMxOHV1bmh6Nzg0d2gzbWpmZWd1emF4dWFxamQyNnk4YWN3Y3Z5Y2MSNmRpZDptaXNlczptaXNlczE4dXVuaHo3ODR3aDNtamZlZ3V6YXh1YXFqZDI2eThhY3djdnljYxoHEgVvdGhlciIAKAESkAEKTgpGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQKnPhwuCsctw/ewM06ytLVKxJjHe+wcQiXrBfCV3CYfARIECgIIARI+CgoKBHVtaXMSAjYxEMCaDCIsbWlzZXMxdjQ5ZGp1OXZkcXkwOXp4N2hsc2tzZjB1N2FnNW1qNDU3OW10c2saQAqSxQA8oRExDD6Up/KW/Fo9gBQ6WrxZl5OIS+XMB1gTMAJR9rTb9pjBZJv1+o+7hBTlmYByphKCCvvFu+WC7/w="
	//txs := "Cp4BCo0BChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEm0KLG1pc2VzMW55a2Nndzd2MjRoeHFzZmFlZzNoc3V4dmZ4ajBnNmRqOW4yM2E0EixtaXNlczFnOW16ZHNxcnkzOXUyNHA4d3hxM3Y1ZHBoemZzcWs1dXI0bTBmZBoPCgR1bWlzEgc0MzYwMDAwEgxtaXNlcyBnbyBzZGsSZQpSCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAitqIBq5hNaK+A7RFvv1b48vH+COX6TJ67dMVCVlDe2oEgQKAggBGMPnAxIPCgkKBHVtaXMSATcQy4gEGkCy2j+kmkyYXu9NcvEG/uQ/QNpoK5nc68/qfpTN13SxNFrQ3CbLSw+jIiEXxcrdcCo+KOLrNt5E8APiiPMBX5Pk"
	//txs := "CvQJCp4BCjcvY29zbW9zLmRpc3RyaWJ1dGlvbi52MWJldGExLk1zZ1dpdGhkcmF3RGVsZWdhdG9yUmV3YXJkEmMKLG1pc2VzMTdlZGZkNXByenJrYTdjNTltOWFkbXdzenl3czk5cGE0bWhrN2E4EjNtaXNlc3ZhbG9wZXIxOWZzbmQzcnBjOXBkajVtN3dlbG54dmo2a2xyeGV0czAwdm03aGgKngEKNy9jb3Ntb3MuZGlzdHJpYnV0aW9uLnYxYmV0YTEuTXNnV2l0aGRyYXdEZWxlZ2F0b3JSZXdhcmQSYwosbWlzZXMxN2VkZmQ1cHJ6cmthN2M1OW05YWRtd3N6eXdzOTlwYTRtaGs3YTgSM21pc2VzdmFsb3BlcjE1Y3JxZ21uNGUyYzQwZGhlajBzMGZoMDJwOHlyYW1kNTVuN2N3bQqeAQo3L2Nvc21vcy5kaXN0cmlidXRpb24udjFiZXRhMS5Nc2dXaXRoZHJhd0RlbGVnYXRvclJld2FyZBJjCixtaXNlczE3ZWRmZDVwcnpya2E3YzU5bTlhZG13c3p5d3M5OXBhNG1oazdhOBIzbWlzZXN2YWxvcGVyMWEzbmp2cjM5OGw5ZmtrMms1dnpnYWtmMG52YWdxd3A3bTJrczdqCp4BCjcvY29zbW9zLmRpc3RyaWJ1dGlvbi52MWJldGExLk1zZ1dpdGhkcmF3RGVsZWdhdG9yUmV3YXJkEmMKLG1pc2VzMTdlZGZkNXByenJrYTdjNTltOWFkbXdzenl3czk5cGE0bWhrN2E4EjNtaXNlc3ZhbG9wZXIxanNzcTRjdzhxNTk1N2VmdnV2bGF4bXhtdmp5MjJ2MnNrbmF3ZnkKmgEKIy9jb3Ntb3Muc3Rha2luZy52MWJldGExLk1zZ0RlbGVnYXRlEnMKLG1pc2VzMTdlZGZkNXByenJrYTdjNTltOWFkbXdzenl3czk5cGE0bWhrN2E4EjNtaXNlc3ZhbG9wZXIxOWZzbmQzcnBjOXBkajVtN3dlbG54dmo2a2xyeGV0czAwdm03aGgaDgoEdW1pcxIGNTYwMDUzCpkBCiMvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dEZWxlZ2F0ZRJyCixtaXNlczE3ZWRmZDVwcnpya2E3YzU5bTlhZG13c3p5d3M5OXBhNG1oazdhOBIzbWlzZXN2YWxvcGVyMTVjcnFnbW40ZTJjNDBkaGVqMHMwZmgwMnA4eXJhbWQ1NW43Y3dtGg0KBHVtaXMSBTgwMjM4CpkBCiMvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dEZWxlZ2F0ZRJyCixtaXNlczE3ZWRmZDVwcnpya2E3YzU5bTlhZG13c3p5d3M5OXBhNG1oazdhOBIzbWlzZXN2YWxvcGVyMWEzbmp2cjM5OGw5ZmtrMms1dnpnYWtmMG52YWdxd3A3bTJrczdqGg0KBHVtaXMSBTE3NDEyCpgBCiMvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dEZWxlZ2F0ZRJxCixtaXNlczE3ZWRmZDVwcnpya2E3YzU5bTlhZG13c3p5d3M5OXBhNG1oazdhOBIzbWlzZXN2YWxvcGVyMWpzc3E0Y3c4cTU5NTdlZnZ1dmxheG14bXZqeTIydjJza25hd2Z5GgwKBHVtaXMSBDYwMzASZQpQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohArhNEKNgnWeTN4GlPscLlPTFb40zyfxuASem9F9xvRqEEgQKAggBGFESEQoLCgR1bWlzEgMyMDAQgIl6GkCwZb+YsjHWZ56fcPKHeG5o9X94QOZV0zhnuIZ1vPv95knd2QDid9sQ0WgWvaG2hHYUzHrzdjr0AxUmcm88Ow5D"
	txs := "Ct0BCtoBCiovY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dCZWdpblJlZGVsZWdhdGUSqwEKLG1pc2VzMWNxd25leDZ4MGF3NGV3YWV2MnUyOHU5cHRmdG5hbmc4ZnZ4OXNyEjNtaXNlc3ZhbG9wZXIxanMwOWo1N25nZTZ5cXBoZTl3dnI4a3A3MGxlM3dxZzI4cXUwZ2EaM21pc2VzdmFsb3BlcjFtdG0yeTJ2ZzRrcHFsZ3M0aDAwMmEyajA3ZnI4eDg0NXF1Z3NlNiIRCgR1bWlzEgkyMDAwMDAwMDASZgpRCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohA/lcGdMNTOcNWt6vmQE6HhpjEEB/Uc0wd7qjZfBOHfArEgQKAggBGJQDEhEKCwoEdW1pcxIDMjAwEICJehpAjxwjieG3H1q/oliDblczRtg9vdUnRKohG1oCI7SbeA9PC5g59l5Z5ky5NhQG9LKW9ifFcuv/OPxS8uZqrELsBA=="
	//txs := "CqEBCp4BCjcvY29zbW9zLmRpc3RyaWJ1dGlvbi52MWJldGExLk1zZ1dpdGhkcmF3RGVsZWdhdG9yUmV3YXJkEmMKLG1pc2VzMW52eW5kNzRtZmZ4NDZ6OW5mcnNjczkycG4zbWw0bmFtY3czaHBsEjNtaXNlc3ZhbG9wZXIxanNzcTRjdzhxNTk1N2VmdnV2bGF4bXhtdmp5MjJ2MnNrbmF3ZnkSZQpQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAvGcKzN1wb4kzTrXWaZ0er+voSNGFuKP+87vG8l/yp5kEgQKAggBGBISEQoLCgR1bWlzEgMyMDAQgIl6GkBY5Fd0sSBQwPnu8FYWmG338+0GlEL4cVO2RZlddGLcK1E4s6KWmQRoq+cgyIoP2zIlj5uQt4ZKQ6Vhc5VgHiD4"
	txc, _ := base64.StdEncoding.DecodeString(txs)
	hs := bytes.HexBytes(tmtypes.Tx(txc).Hash()).String()
	fmt.Println("hs: ", hs)
	txb, err := txCfg.TxDecoder()(txc)
	if err != nil {
		fmt.Println("tx decoders error: ", err.Error())
		return err
	}
	if err := txb.ValidateBasic(); err != nil {
		return err
	}
	p, ok := txb.(intoAny)
	if !ok {
		return fmt.Errorf("expecting a type implementing intoAny, got: %T", txb)
	}
	m := &msg{}
	anyTx := p.AsAny()
	fmt.Println(txb.GetMsgs())
	js, _ := codec.MarshalJSON(txb.GetMsgs()[0])
	json.Unmarshal(js, m)
	msgType := cosmostypes.MsgTypeURL(txb.GetMsgs()[0])
	fmt.Println("txbmsg: ", m, "js", string(js), "msgType", msgType)
	fmt.Println("txbany: ", anyTx.TypeUrl)
	fmt.Println("txbany: ", anyTx.String())
	return nil
}

func cg() ([]byte, error) {
	config := cfg.DefaultBaseConfig()
	config.DBPath = "/Users/cg/.misestm/data"
	if err := config.ValidateBasic(); err != nil {
		return nil, err
	}
	dbType := dbm.BackendType(config.DBBackend)
	// Get BlockStore
	blockStoreDB, err := dbm.NewDB("blockstore", dbType, config.DBDir())
	if err != nil {
		return nil, err
	}
	fmt.Println("dbdir: ", config.DBDir())
	blockStore := store.NewBlockStore(blockStoreDB)
	block := blockStore.LoadBlock(1737167)
	err = block.ValidateBasic()
	if err != nil {
		return nil, err
	}
	//fmt.Println(block.Header)
	fmt.Println(block.Data.Txs[0].Hash())
	return block.Data.Txs.Hash(), nil
}
