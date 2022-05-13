package mq

import (
	"bufio"
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type Msg struct {
	//1:读消息  2:拉消息  3:ack
	Type      byte
	MessageId byte

	TopicLen int64
	Topic    []byte
	MsgLen   int64
	Message  []byte
}

func B2M(read io.Reader) (Msg, error) {
	msg := Msg{}
	buf := make([]byte, 256)

	n, err := read.Read(buf)
	if err != nil {
		log.Fatalf("Read error : %v", err)
		return msg, err
	}
	fmt.Println("已读取数据字节:", n)

	//type
	msg.Type = buf[0]
	//id
	msg.MessageId = buf[1]

	buffer := bytes.NewBuffer(buf[2:6])
	err = binary.Read(buffer, binary.LittleEndian, &msg.TopicLen)
	if err != nil {
		log.Fatalf("read topiclen error : %v", err)
		return msg, err
	}

	buffer = bytes.NewBuffer(buf[6+msg.TopicLen : 10+msg.TopicLen])
	err = binary.Read(buffer, binary.LittleEndian, &msg.Topic)
	if err != nil {
		log.Fatalf("read topic error : %v", err)
		return msg, err
	}

	buffer = bytes.NewBuffer(buf[6+msg.TopicLen : 10+msg.TopicLen+msg.MsgLen])
	err = binary.Read(buffer, binary.LittleEndian, &msg.Message)
	if err != nil {
		log.Fatalf("read topic error : %v", err)
		return msg, err
	}
	return msg, nil
}

func M2B(msg Msg) ([]byte, error) {
	data := make([]byte, 256)

	buffer := bytes.NewBuffer([]byte{})

	data = append(data, msg.Type, msg.MessageId)

	err := binary.Write(buffer, binary.LittleEndian, msg.TopicLen)
	if err != nil {
		log.Fatalf("write topiclen error : %v", err)
		return data, err
	}

	data = append(data, buffer.Bytes()...)
	data = append(data, msg.Topic...)

	err = binary.Write(buffer, binary.LittleEndian, msg.MsgLen)
	if err != nil {
		log.Fatalf("write msgLen error : %v", err)
		return data, err
	}
	return data, nil
}

type Broker struct {
	que map[string]*list.List
	l   sync.Mutex
}

func NewBroker() *Broker {
	return &Broker{}
}

func (b *Broker) Run() {
	listen, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		log.Fatalf("listen prot error : %v", err)
	}

	b.que = make(map[string]*list.List)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalf("accept error : %v", err)
			return
		}

		reader := bufio.NewReader(conn)

		msg, err := B2M(reader)
		if err != nil {
			log.Fatalf("B2M error : %v", err)
			return
		}

		switch msg.Type {
		case 1:
			b.l.Lock()
			if b.que[string(msg.Topic)] == nil {
				b.que[string(msg.Topic)] = &list.List{}
			}
			b.que[string(msg.Topic)].PushBack(msg)
			b.l.Unlock()
		case 2:
			b.l.Lock()
			if b.que[string(msg.Topic)] == nil {
				b.que[string(msg.Topic)] = &list.List{}
			}
			if b.que[string(msg.Topic)].Len() == 0 {
				//conn写入错误
				continue
			}
			res := b.que[string(msg.Topic)].Front()

			data, err := M2B(res.Value.(Msg))
			if err != nil {
				log.Fatalf("M2B error : %v", err)
				return
			}
			_, err = conn.Write(data)
			if err != nil {
				log.Fatalf("m2b error : %v", err)
				return
			}

			b.que[string(msg.Topic)].Remove(res)
			b.l.Unlock()
		}
	}
}
