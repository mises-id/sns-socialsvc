package enum

import "github.com/mises-id/sns-socialsvc/lib/codes"

type TagType string

const (
	TagBlank           TagType = ""
	TagStarUser        TagType = "star_user"
	TagProblemUser     TagType = "problem_user"
	TagRecommendUser   TagType = "recommend_user"
	TagRecommendStatus TagType = "recommend_status"
)

var (
	tagTypeMap = map[TagType]string{
		TagStarUser:        "star_user",
		TagProblemUser:     "problem_user",
		TagRecommendUser:   "recommend_user",
		TagRecommendStatus: "recommend_status",
	}
	tagTypeStringMap = map[string]TagType{}
)

func init() {
	for key, val := range tagTypeMap {
		tagTypeStringMap[val] = key
	}
}

func TagTypeFromString(tp string) (TagType, error) {
	tagType, ok := tagTypeStringMap[tp]
	if !ok {
		return TagBlank, codes.ErrInvalidArgument.Newf("invalid status type: %s", tp)
	}
	return tagType, nil
}
