package enum

type (
	State               int
	SortType            int
	UserValidState      int32
	ChannelAirdropState int32
)

const (
	ChannelAirdropDefault ChannelAirdropState = 0
	ChannelAirdropPending ChannelAirdropState = 1
	ChannelAirdropSuccess ChannelAirdropState = 2
	ChannelAirdropFailed  ChannelAirdropState = 3

	StateClose        State          = 0
	StateOpen         State          = 1
	UserValidDefalut  UserValidState = 0
	UserValidSucessed UserValidState = 1
	UserValidFailed   UserValidState = 2

	SortAsc  SortType = 1
	SortDesc SortType = -1
)
