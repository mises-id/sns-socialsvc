package enum

import "github.com/mises-id/sns-socialsvc/lib/codes"

type MessageType uint32

const (
	NewComment MessageType = iota
	NewLike
	NewFans
	NewForward
)

var (
	messageTypeMap = map[MessageType]string{
		NewComment: "new_comment",
		NewLike:    "new_like",
		NewFans:    "new_fans",
		NewForward: "new_fowards",
	}
	messageTypeStringMap = map[string]MessageType{}
)

func init() {
	for key, val := range messageTypeMap {
		messageTypeStringMap[val] = key
	}
}

func (tp MessageType) String() string {
	return messageTypeMap[tp]
}

func MessageTypeFromString(tp string) (MessageType, error) {
	messageType, ok := messageTypeStringMap[tp]
	if !ok {
		return NewComment, codes.ErrInvalidArgument.Newf("invalid message type: %s", tp)
	}
	return messageType, nil
}
