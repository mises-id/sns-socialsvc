package enum

import "github.com/mises-id/sns-socialsvc/lib/codes"

type MessageState uint32

const (
	Unread MessageState = iota
	Read
)

var (
	messageStateMap = map[MessageState]string{
		Unread: "unread",
		Read:   "read",
	}
	messageStateStringMap = map[string]MessageState{}
)

func init() {
	for key, val := range messageStateMap {
		messageStateStringMap[val] = key
	}
}

func (tp MessageState) String() string {
	return messageStateMap[tp]
}

func MessageStateFromString(tp string) (MessageState, error) {
	messageState, ok := messageStateStringMap[tp]
	if !ok {
		return Unread, codes.ErrInvalidArgument.Newf("invalid message state: %s", tp)
	}
	return messageState, nil
}
