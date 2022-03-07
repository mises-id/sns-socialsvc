package enum

type AirdropType string

type AirdropStatus int32

const (
	AirdropDefault AirdropStatus = iota
	AirdropPending
	AirdropSuccess
	AirdropFailed
	AirdropTwitter = "twitter"
)
