package service

import (
	"RedRock-Work/lesson_6/lv1/user/client/model"
	"RedRock-Work/lesson_6/lv1/user/server/grpc/user"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Register(u model.UserInfo) (int32, string, error) {
	newCredentials := insecure.NewCredentials()
	conn, err := grpc.Dial("localhost:8888", grpc.WithTransportCredentials(newCredentials))
	if err != nil {
		return 0, "", err
	}
	defer conn.Close()

	client := user.NewVerifyUserClient(conn)
	reqUser := user.ReqUser{
		UserName: u.Username,
		Password: u.Password,
	}
	rsp, err := client.Register(context.Background(), &reqUser)
	if err != nil {
		return 0, "", err
	}
	return rsp.GetStatus(), rsp.GetMessage(), nil
}

func Login(u model.UserInfo) (int32, string, error) {
	newCredentials := insecure.NewCredentials()
	conn, err := grpc.Dial("localhost:8888", grpc.WithTransportCredentials(newCredentials))
	if err != nil {
		return 0, "", err
	}
	defer conn.Close()

	client := user.NewVerifyUserClient(conn)
	reqUser := user.ReqUser{
		UserName: u.Username,
		Password: u.Password,
	}
	rsp, err := client.Login(context.Background(), &reqUser)
	if err != nil {
		return 0, "", err
	}
	return rsp.GetStatus(), rsp.GetMessage(), nil
}

func ChangePassword(u model.UserInfo) (int32, string, error) {
	newCredentials := insecure.NewCredentials()
	conn, err := grpc.Dial("localhost:8888", grpc.WithTransportCredentials(newCredentials))
	if err != nil {
		return 0, "", err
	}

	client := user.NewVerifyUserClient(conn)
	reqUser := user.ReqUser{
		UserName: u.Username,
		Password: u.Password,
	}
	rsp, err := client.ChangePassword(context.Background(), &reqUser)
	if err != nil {
		return 0, "", err
	}
	return rsp.GetStatus(), rsp.GetMessage(), nil
}
