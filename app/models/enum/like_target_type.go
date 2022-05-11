package enum

import "github.com/mises-id/sns-socialsvc/lib/codes"

type LikeTargetType uint32

const (
	LikeStatus LikeTargetType = iota
	LikeComment
	LikeNft
	LikeNftComment
)

var (
	likeTargetTypeMap = map[LikeTargetType]string{
		LikeStatus:     "status",
		LikeComment:    "comment",
		LikeNft:        "nft",
		LikeNftComment: "nft_comment",
	}
	likeTargetTypeStringMap = map[string]LikeTargetType{}
)

func init() {
	for key, val := range likeTargetTypeMap {
		likeTargetTypeStringMap[val] = key
	}
}

func (tp LikeTargetType) String() string {
	return likeTargetTypeMap[tp]
}

func LikeTargetTypeFromString(tp string) (LikeTargetType, error) {
	likeTargetType, ok := likeTargetTypeStringMap[tp]
	if !ok {
		return LikeStatus, codes.ErrInvalidArgument.Newf("invalid like target type: %s", tp)
	}
	return likeTargetType, nil
}
