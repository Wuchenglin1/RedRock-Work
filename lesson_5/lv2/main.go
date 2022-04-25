package main

import (
	"net"
	"net/http"
)

const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

type Handler func(conn *MyConn)
type Msg struct {
	Typ     int
	content []byte
}
type MyConn struct {
	conn net.Conn
	//关于读写缓存啥的

	//比如当遇到ping消息时触发
	pingHandler Handler
}

func Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (conn MyConn, err error) {
	//大致流程：
	//检查请求头 Connection
	//检查请求头 Upgrade
	//检查请求方式
	//检查请求头 Sec-Websocket-Version 是否为13
	//检查请求头 Origin是否是允许的
	//检查请求头 Sec-Websocket-Key
	//处理 Sec-Websocket-Protocol 子协议字段

	//关于hijacker，就是从http.ResponseWriter重新拿到conn
	//调用 http.Hijacker 拿到这个连接现在开始就可以使用websocket通信了
	//大概就是
	h, ok := w.(http.Hijacker)
	if !ok {
		//寄
	}
	conn.conn, _, err = h.Hijack()
	//手写回复报文
	resp := []byte{}
	//HTTP/1.1 101 Switching Protocols
	//一系列请求头
	//Upgrade: websocket
	//Connection: Upgrade
	//Sec-WebSocket-Accept：
	//Sec-WebSocket-Protocol:
	//请求头写完别忘了换行
	//将请求报文写入
	conn.conn.Write(resp)
	return
}

//读取一次数据
//根据协议解析消息
/*
0                   1                   2                   3
0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-------+-+-------------+-------------------------------+
|F|R|R|R| opcode|M| Payload len |    Extended payload length    |
|I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |
|N|V|V|V|       |S|             |   (if payload len==126/127)   |
| |1|2|3|       |K|             |                               |
+-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
|     Extended payload length continued, if payload len == 127  |
+ - - - - - - - - - - - - - - - +-------------------------------+
|                               |Masking-key, if MASK set to 1  |
+-------------------------------+-------------------------------+
| Masking-key (continued)       |          Payload Data         |
+-------------------------------- - - - - - - - - - - - - - - - +
:                     Payload Data continued ...                :
+ - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
|                     Payload Data continued ...                |
+---------------------------------------------------------------+
*/

func (conn *MyConn) ReadMsg(m Msg) {
	//conn.Read()
	//首先读取第一个字节
	//第一个字节第一位
	//同理把后面读出来

	//掩码处理

	//这里假设这是你读到的类型
	t := 0
	//ping pong 消息是心跳消息
	if t == PingMessage {
		//这种消息一般不返回
		//一般是要设置函数
		//比如说收到ping消息回复pong消息之类的
		if conn.pingHandler != nil {
			conn.pingHandler(conn)
		}
	} else if t == PongMessage {
		//同上
	} else if t == TextMessage {
		//收到文本消息返回
		m.Typ = t
		//m.content=读取到的内容
	} else if t == BinaryMessage {
		//收到二进制消息返回
		m.Typ = t
		//m.content=读取到的内容
	} else if t == CloseMessage {
		//收到关闭消息
		//向上返回，看你怎么处理
		m.Typ = t
		//m.content=读取到的内容
	} else {
		//读取到其它的
		//可能是你和客户端自定义的
		//可能是非法
	}
	return
}

func (conn *MyConn) WriteMsg(m Msg) (err error) {
	//按照数据帧写出数据
	p := []byte{}
	_, err = conn.conn.Write(p)
	if err != nil {
		//写不进去，咋办呢
	}
	return
}
