package tool

import (
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app/client"
	"github.com/MashiroC/begonia/app/option"
)

func GetService(funcName string) (client.RemoteFunSync, error) {

	c := begonia.NewClient(option.Addr(":12306"))
	s, err := c.Service("userServer")
	if err != nil {
		return nil, err
	}
	loginFunc, err := s.FuncSync(funcName)
	if err != nil {
		return nil, err
	}
	return loginFunc, nil
}
