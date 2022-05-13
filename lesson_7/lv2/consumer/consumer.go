package main

import (
	mq "RedRock-Work/lesson_7/lv2/broker"
	"bytes"
	"fmt"
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
		Type:  2,
		Topic: []byte("我是主题"),
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
	var res [2048]byte

	_, err = conn.Read(res[:])
	if err != nil {
		log.Fatalf("read data err : %v", err)
		return
	}

	buffer := bytes.NewBuffer(res[:])
	message, err := mq.B2M(buffer)
	if err != nil {
		log.Fatalf("b2m error : %v", err)
		return
	}
	fmt.Println(message)
}
