package enum

import "github.com/mises-id/sns-socialsvc/lib/codes"

type MessageType uint32

const (
	NewComment MessageType = iota
	NewLikeStatus
	NewLikeComment
	NewFans
	NewForward
	NewNftComment
	NewLikeNft
	NewLikeNftComment
)

var (
	messageTypeMap = map[MessageType]string{
		NewComment:        "new_comment",
		NewLikeStatus:     "new_like_status",
		NewLikeComment:    "new_like_comment",
		NewFans:           "new_fans",
		NewForward:        "new_fowards",
		NewNftComment:     "new_nft_comment",
		NewLikeNft:        "new_like_nft",
		NewLikeNftComment: "new_like_nft_comment",
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
