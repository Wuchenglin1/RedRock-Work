package api

import (
	"RedRock-Work/lesson_6/lv1.00000000000000001/server/model"
	"RedRock-Work/lesson_6/lv1.00000000000000001/server/service"
)

type server struct{}

func (*server) Register(user model.InputUser) model.RspMsg {
	return service.Register(user)
}

func (*server) Login(user model.InputUser) model.RspMsg {
	return service.Login(user)
}

func (*server) ChangePassword(user model.InputUser) model.RspMsg {
	return service.ChangePassword(user)
}
