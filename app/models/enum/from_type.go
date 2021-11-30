package enum

import "github.com/mises-id/socialsvc/lib/codes"

type FromType uint8

const (
	FromPost FromType = iota
	FromForward
	FromComment
	FromLike
)

type FromTypeFilter struct {
	FromType FromType
}

var (
	fromTypeMap = map[FromType]string{
		FromPost:    "post",
		FromForward: "forward",
		FromComment: "comment",
		FromLike:    "like",
	}
	fromTypeStringMap  = map[string]FromType{}
	fromTypeCounterMap = map[FromType]string{
		FromForward: "forwards_count",
		FromComment: "comments_count",
		FromLike:    "likes_count",
	}
)

func init() {
	for key, val := range fromTypeMap {
		fromTypeStringMap[val] = key
	}
}

func (tp FromType) String() string {
	return fromTypeMap[tp]
}

func (tp FromType) CounterKey() string {
	return fromTypeCounterMap[tp]
}

func FromTypeFromString(tp string) (FromType, error) {
	fromType, ok := fromTypeStringMap[tp]
	if !ok {
		return FromPost, codes.ErrInvalidArgument.Newf("invalid from type: %s", tp)
	}
	return fromType, nil
}
