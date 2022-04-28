package tool

import "RedRock-Work/lesson_6/lv1/user/server/grpc/user"

func GetRsp(rsp *user.Rsp, message string, status int32) *user.Rsp {
	rsp.Message = message
	rsp.Status = status
	return rsp
}
