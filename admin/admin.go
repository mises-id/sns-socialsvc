package admin

type AdminApi interface {
	User() UserApi
	Status() StatusApi
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
