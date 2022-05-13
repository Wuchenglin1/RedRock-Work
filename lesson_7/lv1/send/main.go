package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

const account = "小秘密哦"

func main() {
	sendMsg()
}

func sendMsg() {
	conn, err := amqp.Dial(account)
	if err != nil {
		log.Fatalf("连接rabbitmq失败 : %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("打开管道失败: %v", err)
	}
	defer ch.Close()

	chName := "ch1"
	q, err := ch.QueueDeclare(chName, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("声明队列失败: %v", err)
	}

	var body string
	fmt.Println("请输入想发送的内容")
	fmt.Scanln(&body)

	err = ch.Publish("", q.Name, false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte(body)})
	if err != nil {
		log.Fatalf("发布消息失败: %v", err)
	}
	log.Fatalf("向 %v 发送消息 %v 成功!", chName, body)
}
