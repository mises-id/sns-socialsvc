package enum

import "github.com/mises-id/sns-socialsvc/lib/codes"

type RelationType uint32

const (
	Following RelationType = iota + 1
	Fan
	Friend
)

var (
	relationTypeMap = map[RelationType]string{
		Following: "following",
		Fan:       "fan",
		Friend:    "friend",
	}
	relationTypeStringMap = map[string]RelationType{}
)

func init() {
	for key, val := range relationTypeMap {
		relationTypeStringMap[val] = key
	}
}

func RelationTypeFromString(relationType string) (RelationType, error) {
	result, ok := relationTypeStringMap[relationType]
	if !ok {
		return Following, codes.ErrInvalidArgument.Newf("invalid relation type: %s", relationType)
	}
	return result, nil
}

func (relationType RelationType) String() string {
	return relationTypeMap[relationType]
}
