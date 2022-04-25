package main

import (
	"RedRock-Work/lesson_6/lv0/proto/login"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

const address = "localhost:6666"

func main() {
	newCredentials := insecure.NewCredentials()
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(newCredentials))
	if err != nil {
		log.Fatalf("connect faild : %s", err)
	}
	defer conn.Close()

	client := login.NewLoginClient(conn)
	user := login.User{}
	req := login.ReqMsg{
		User: &user,
	}
	fmt.Println("请输入账号名和密码")
	for {
		fmt.Scan(&user.UserName, &user.Password)
		rsp, err1 := client.Login(context.Background(), &req)
		if err1 != nil {
			log.Fatalf("response error : %s", err1)
		}
		if !rsp.GetOK() {
			fmt.Printf("来自服务器的消息: %v", rsp.Message)
			continue
		} else {
			fmt.Printf("来自服务器的消息: %v\n", rsp.Message)
			break
		}
	}
}
