package main

import (
	"RedRock-Work/lesson_6/lv1.00000000000000001/server/api"
	"RedRock-Work/lesson_6/lv1.00000000000000001/server/dao"
)

func main() {
	dao.InitMysql()
	dao.InitRedis()
	api.InitRouter()
}
