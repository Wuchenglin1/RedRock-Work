package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
)

type Config struct {
	Ip       string `json:"ip"`
	Password string `json:"password"`
}

func main() {

	cfg := GetConfig()

	//连接redis
	fmt.Println(cfg.Ip, cfg.Password)
	dialOption := redis.DialPassword(cfg.Password)
	c, err := redis.Dial("tcp", cfg.Ip, dialOption)
	if err != nil {
		log.Println("fail:", err)
		return
	}
	defer c.Close()

	//获取PubSubConn类型
	psc := redis.PubSubConn{Conn: c}

	//订阅 test 频道
	err = psc.Subscribe("test")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		//需要进行断言才能获取数据
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			fmt.Println(v.Error())
			break
		}
	}

}

func GetConfig() *Config {
	var cfg *Config
	file, err := os.Open("./config.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return cfg
}
