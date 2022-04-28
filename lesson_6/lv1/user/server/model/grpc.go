package model

import (
	"RedRock-Work/lesson_6/lv1/user/server/grpc/user"
)

type Server struct {
	user.UnimplementedVerifyUserServer
}
