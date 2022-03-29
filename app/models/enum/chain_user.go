package enum

type ChainUserStatus int32

const (
	ChainUserDefault ChainUserStatus = iota
	ChainUserPending
	ChainUserSuccess
	ChainUserFailed
)
