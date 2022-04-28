package tool

import (
	"RedRock-Work/lesson_6/lv1/user/server/model"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

var _cfg *model.Config

func GetConfig() *model.Config {
	file, err := os.Open("../config/config.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&_cfg)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return _cfg
}
