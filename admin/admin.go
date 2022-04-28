package admin

import (
	"github.com/mises-id/sns-socialsvc/config/env"
	"github.com/mises-id/sns-socialsvc/lib/storage"
)

func init() {
	storage.SetupImageStorage(env.Envs.StorageHost, env.Envs.StorageKey, env.Envs.StorageSalt)
}

type AdminApi interface {
	User() UserApi
	Status() StatusApi
	ChannelUser() ChannelUserApi
}

type adminApi struct {
}

func New() AdminApi {
	return &adminApi{}
}

func (*adminApi) User() UserApi {
	return NewUserApi()
}

func (*adminApi) Status() StatusApi {
	return NewStatusApi()
}
func (*adminApi) ChannelUser() ChannelUserApi {
	return NewChannelUserApi()
}
