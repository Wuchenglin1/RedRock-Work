package main

import (
	"RedRock-Work/lesson_6/lv1/user/server/api"
	"RedRock-Work/lesson_6/lv1/user/server/dao"
)

func main() {
	dao.InitMysql()
	dao.InitRedis()
	api.InitRouter()
}
