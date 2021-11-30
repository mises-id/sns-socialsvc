package enum

import (
	"github.com/mises-id/socialsvc/lib/codes"
)

type Gender uint8

const (
	GenderOther Gender = iota
	GenderMale
	GenderFemale
)

var (
	genderMap = map[Gender]string{
		GenderOther:  "other",
		GenderMale:   "male",
		GenderFemale: "female",
	}
	genderStringMap = map[string]Gender{}
)

func init() {
	for key, val := range genderMap {
		genderStringMap[val] = key
	}
}

func GenderFromString(gender string) (Gender, error) {
	result, ok := genderStringMap[gender]
	if !ok {
		return GenderOther, codes.ErrInvalidArgument.Newf("invalid gender: %s", gender)
	}
	return result, nil
}

func (gender Gender) String() string {
	return genderMap[gender]
}
