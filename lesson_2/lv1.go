package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
)

func main() {

	//连接redis
	dialOption := redis.DialPassword("root")
	c, err := redis.Dial("tcp", "110.42.165.192:6379", dialOption)
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
