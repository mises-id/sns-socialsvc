package message

type FansMeta struct {
	UID         uint64 `bson:"uid,omitempty"`
	FanUsername string `bson:"fan_username,omitempty"`
}
