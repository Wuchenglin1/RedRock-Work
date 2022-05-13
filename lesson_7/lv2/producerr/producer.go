package main

import (
	mq "RedRock-Work/lesson_7/lv2/broker"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		log.Printf("connection failed : %v", err)
		return
	}
	defer conn.Close()

	msg := mq.Msg{
		Type:      1,
		MessageId: 0,
		Topic:     []byte("我是主题"),
		Message:   []byte("我是消息"),
	}

	msg.TopicLen = int64(len(msg.Topic))
	msg.MsgLen = int64(len(msg.Message))

	_bytes, err := mq.M2B(msg)
	if err != nil {
		return
	}

	_, err = conn.Write(_bytes)
	if err != nil {
		log.Fatalf("conn write error : %v", err)
		return
	}
}
