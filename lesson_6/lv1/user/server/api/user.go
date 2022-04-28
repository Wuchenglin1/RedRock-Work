package api

import (
	"RedRock-Work/lesson_6/lv1/user/server/grpc/user"
	"RedRock-Work/lesson_6/lv1/user/server/model"
	"RedRock-Work/lesson_6/lv1/user/server/service"
	"context"
)

type server model.Server

func (*server) Register(ctx context.Context, userInfo *user.ReqUser) (*user.Rsp, error) {
	return service.Register(userInfo)
}

func (*server) Login(ctx context.Context, userInfo *user.ReqUser) (*user.Rsp, error) {
	return service.Login(userInfo)
}

func (*server) ChangePassword(ctx context.Context, userInfo *user.ReqUser) (*user.Rsp, error) {
	return service.ChangePassword(userInfo)
}
