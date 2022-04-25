package main

import (
	"RedRock-Work/lesson_6/lv0/proto/login"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

const port = ":6666"

type server struct {
	login.UnimplementedLoginServer
}

func (*server) Login(ctx context.Context, req *login.ReqMsg) (*login.RspMsg, error) {
	rsp := new(login.RspMsg)
	user := req.GetUser()
	rsp.OK = CheckPassword(user.UserName, user.Password)
	fmt.Println(user, rsp.OK)
	if rsp.OK {
		rsp.Message = "登录成功!"
		return rsp, nil
	} else {
		rsp.Message = "账号或密码错误!请重新输入账号密码!"
		return rsp, nil
	}
}

func main() {
	ls, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("listen error : %s", err)
	}
	defer ls.Close()

	s := grpc.NewServer()
	login.RegisterLoginServer(s, &server{})
	fmt.Printf("listening on %v\n", port)
	if err = s.Serve(ls); err != nil {
		log.Fatalf("server error : %s", err)
	}
}

func CheckPassword(userName string, password string) bool {
	return password == userName+"123456"
}
