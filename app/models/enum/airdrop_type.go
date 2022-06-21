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

var (
	AirdropTypeMap = map[AirdropStatus]string{
		AirdropDefault: "default",
		AirdropPending: "pending",
		AirdropSuccess: "success",
		AirdropFailed:  "failed",
	}
)

func (tp AirdropStatus) String() string {
	return AirdropTypeMap[tp]
}
