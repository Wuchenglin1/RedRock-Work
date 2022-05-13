package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

const account = "小秘密哦"

func main() {
	receiveMsg()
}

func receiveMsg() {
	conn, err := amqp.Dial(account)
	if err != nil {
		log.Fatalf("连接rabbitmq失败: %v", err)
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
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("注册失败: %v", err)
	}

	exitCh := make(chan bool)

	go func() {
		for d := range msgs {
			fmt.Printf("从 %v 接收到消息 %v \n", chName, string(d.Body))
		}
	}()
	fmt.Printf("正在从 %v 接收消息\n", chName)
	<-exitCh
}
