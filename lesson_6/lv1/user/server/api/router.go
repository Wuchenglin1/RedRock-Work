package api

import (
	"RedRock-Work/lesson_6/lv1/user/server/grpc/user"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func InitRouter() {
	const port = ":8888"
	ls, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("listen error : %s", err)
	}
	defer ls.Close()

	s := grpc.NewServer()
	user.RegisterVerifyUserServer(s, &server{})
	fmt.Printf("listening on %v ", port)
	if err = s.Serve(ls); err != nil {
		log.Fatalf("server error : %s", err)
	}
}
