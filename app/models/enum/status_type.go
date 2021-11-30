package enum

import "github.com/mises-id/socialsvc/lib/codes"

type StatusType uint8

const (
	TextStatus StatusType = iota
	LinkStatus
)

var (
	statusTypeMap = map[StatusType]string{
		TextStatus: "text",
		LinkStatus: "link",
	}
	statusTypeStringMap = map[string]StatusType{}
)

func init() {
	for key, val := range statusTypeMap {
		statusTypeStringMap[val] = key
	}
}

func (tp StatusType) String() string {
	return statusTypeMap[tp]
}

func StatusTypeFromString(tp string) (StatusType, error) {
	statusType, ok := statusTypeStringMap[tp]
	if !ok {
		return TextStatus, codes.ErrInvalidArgument.Newf("invalid status type: %s", tp)
	}
	return statusType, nil
}
