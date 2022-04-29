package tool

import "RedRock-Work/lesson_6/lv1.00000000000000001/server/model"

func GetRsp(status int32, message string) model.RspMsg {
	return model.RspMsg{
		Status:  status,
		Message: message,
	}
}
