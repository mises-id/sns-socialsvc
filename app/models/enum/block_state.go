package enum

import "github.com/mises-id/sns-socialsvc/lib/codes"

type BlockState uint32

const (
	Normal BlockState = iota
	Block
	Blocked
	BlockAndBlocked
)

var (
	blockStateMap = map[BlockState]string{
		Normal:          "normal",
		Block:           "block",
		Blocked:         "blocked",
		BlockAndBlocked: "block_and_blocked",
	}
	blockStateStringMap = map[string]BlockState{}
)

func init() {
	for key, val := range blockStateMap {
		blockStateStringMap[val] = key
	}
}

func BlockStateFromString(st string) (BlockState, error) {
	result, ok := blockStateStringMap[st]
	if !ok {
		return Normal, codes.ErrInvalidArgument.Newf("invalid block state: %s", st)
	}
	return result, nil
}

func (st BlockState) String() string {
	return blockStateMap[st]
}
