package stream

import (
	"context"
	"fmt"

	streamLib "github.com/mises-id/sns-socialsvc/lib/stream"
	tmtypes "github.com/tendermint/tendermint/types"
)

func Run(ctx context.Context) error {
	callback := &EventStreamingCallback{}
	callback.done = make(chan bool)
	callback.maxCount = 10
	err := streamLib.StreamClient.StartEventStreaming(callback)
	if err != nil {
		fmt.Println("StartEventStreaming error: ", err.Error())
		return nil
	}
	callback.wait()

	resp, err := streamLib.StreamClient.ParseEvent(callback.header, callback.tx)
	if err != nil {
		fmt.Println("ParseEvent error: ", err.Error())
		return nil
	}
	fmt.Printf("ParseEvent %s", resp.TxHash)
	return nil
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
		cb.done <- true
	}
}
func (cb *EventStreamingCallback) OnNewBlockHeaderEvent(h *tmtypes.EventDataNewBlockHeader) {
	fmt.Printf("OnNewBlockHeaderEvent\n")
	cb.eventCount++
	cb.header = h
	if cb.eventCount > cb.maxCount || (cb.tx != nil && cb.header != nil) {
		cb.done <- true
	}
}
func (cb *EventStreamingCallback) OnEventStreamingTerminated() {
	fmt.Printf("OnEventStreamingTerminated")
	cb.done <- true
}
func (cb *EventStreamingCallback) wait() {
	<-cb.done
}
