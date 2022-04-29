package api

import (
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app/option"
	"time"
)

func InitRouter() {
	s := begonia.NewServer(option.Addr(":12306"))
	//注册服务
	s.Register("userServer", &server{})
	for {
		time.Sleep(1 * time.Hour)
	}
}
