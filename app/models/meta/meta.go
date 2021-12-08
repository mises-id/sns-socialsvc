package meta

import (
	"encoding/json"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/codes"
)

type MetaData interface {
	isMetaData()
}

func BuildStatusMeta(statusType enum.StatusType, metaData json.RawMessage) (MetaData, error) {
	if metaData == nil {
		return nil, nil
	}
	var data MetaData
	switch statusType {
	default:
		return &TextMeta{}, codes.ErrInvalidArgument.New("invalid status type")
	case enum.TextStatus:
		data = &TextMeta{}
	case enum.LinkStatus:
		data = &LinkMeta{}
	}
	return data, json.Unmarshal(metaData, data)
}
