package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

var mux sync.Mutex
var room sync.Map

type Client struct {
	Connection   *websocket.Conn
	IsRegister   bool
	IsLogin      bool
	Password     string
	Channel      int //频道
	NextSendTime int64
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Broad struct {
	MsgCh      chan Msg
	Register   chan Msg
	UnRegister chan Msg
	CountCh    chan websocket.Conn
}

type Msg struct {
	Time     time.Time
	UserName string
	MsgType  int
	Content  []byte
	Channel  int
}

var B = Broad{
	MsgCh:      make(chan Msg, 50),
	Register:   make(chan Msg, 50),
	UnRegister: make(chan Msg, 50),
}

func WsHandler(c *gin.Context) {
	userName, ch := getQuery(c)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	//报错
	if err != nil {
		log.Println("coon.ReadMessage() err1 :", err)
		return
	}
	channel, err := strconv.Atoi(ch)
	if err != nil {
		channel = 1
	}
	//加入聊天室的信息
	B.Register <- Msg{
		Time:     time.Now(),
		UserName: userName,
		MsgType:  websocket.TextMessage,
		Content:  nil,
		Channel:  channel,
	}
	//创建一个client
	client := Client{
		Connection:   conn,
		IsRegister:   false,
		IsLogin:      false,
		Channel:      channel,
		NextSendTime: time.Now().Add(time.Second * 5).Unix(),
	}

	//检测用户是否重名
	_, ok := room.Load(userName)
	if ok {
		conn.WriteMessage(websocket.TextMessage, S2B("该用户已存在!"))
		conn.Close()
		return
	}
	//进入房间即开始计时挂机时间
	err = conn.SetReadDeadline(time.Now().Add(time.Minute * 10))

	if err != nil {
		fmt.Println("send error message err3:", err)
		room.Delete(userName)
		conn.Close()
		return
	}
	//把client存入map
	room.Store(userName, &client)

	go WriteMsg(userName, channel, conn)
}

// WriteMsg Rec 读消息
func WriteMsg(userName string, channel int, conn *websocket.Conn) {
	defer conn.Close()
	for {
		msgTyp, content, err := conn.ReadMessage()
		if err != nil {
			//这里报错一般是客户端关闭的时候发来的close消息
			//fmt.Println("send message err:", err)
			B.UnRegister <- Msg{
				Time:     time.Now(),
				UserName: userName,
				MsgType:  websocket.TextMessage,
				Content:  S2B("退出了聊天室"),
			}
			room.Delete(userName)
			conn.Close()
			return
		}
		//这里都统一处理二进制数据类型和文本数据类型了
		if msgTyp == websocket.BinaryMessage || msgTyp == websocket.TextMessage {
			//先刷新挂机踢出群聊时间
			conn.SetReadDeadline(time.Now().Add(time.Minute * 1))
			//拿取client
			value, _ := room.Load(userName)
			client := value.(*Client)

			//检测注册状态
			if !client.IsRegister {
				strContent := B2S(content)
				n := strings.Index(strContent, "/register ")
				//指令不正确
				if n == -1 {
					conn.WriteMessage(websocket.TextMessage, S2B(time.Now().Format(time.Stamp)+"您还没有注册！"+"\n请输入/register 密码 密码 来注册"))
					continue
				}
				//格式不正确
				if n != 0 {
					conn.WriteMessage(websocket.TextMessage, S2B(time.Now().Format(time.Stamp)+"您还没有注册！"+"\n请输入/register 密码 密码 来注册"))
					continue
				}
				password := strContent[10:]
				bt := strings.Split(password, " ")
				if len(bt) != 2 {
					client.Connection.WriteMessage(websocket.TextMessage, S2B(time.Now().Format(time.Stamp)+"格式不正确!"))
					continue
				}
				if bt[0] != bt[1] {
					client.Connection.WriteMessage(websocket.TextMessage, S2B(time.Now().Format(time.Stamp)+"密码不相同!"))
					continue
				}
				client.IsRegister = true
				//账号密码存储在map中
				client.Password = bt[0]
				client.Connection.WriteMessage(websocket.TextMessage, S2B(time.Now().Format(time.Stamp)+"您注册成功"))
				continue
			}
			//检测是否登录
			if !client.IsLogin {
				userPassword := client.Password
				strContent := B2S(content)
				n := strings.Index(strContent, "/login ")
				//指令有误
				if n == -1 {
					client.Connection.WriteMessage(websocket.TextMessage, S2B(time.Now().Format(time.Stamp)+"您输入的指令有误!"+"\n请输入/login password来登录"))
					continue
				}
				//格式错误
				if n != 0 {
					client.Connection.WriteMessage(websocket.TextMessage, S2B(time.Now().Format(time.Stamp)+"您输入的指令有误!"+"\n请输入/login password来登录"))
					continue
				}
				password := strContent[7:]
				if n = strings.Index(password, " "); n != -1 {
					client.Connection.WriteMessage(websocket.TextMessage, S2B(time.Now().Format(time.Stamp)+"密码不能包含空格!"))
					continue
				}
				if password != userPassword {
					client.Connection.WriteMessage(websocket.TextMessage, S2B(time.Now().Format(time.Stamp)+"密码错误!"))
					continue
				}
				client.IsLogin = true
				client.Connection.WriteMessage(websocket.TextMessage, S2B(time.Now().Format(time.Stamp)+"登录成功!\n请尽情发言吧~!"))
				continue
			}
			//已注册登录
			//先检测用户是否频繁发消息
			if time.Now().Unix() >= client.NextSendTime {
				client.Connection.WriteMessage(websocket.TextMessage, S2B("喝杯水歇歇吧,不要频繁发消息喔~"))
				//同样也要刷新发消息间隔时间
				client.NextSendTime = time.Now().Add(time.Second * 5).Unix()
				continue
			}
			//刷新检测发消息间隔时间
			client.NextSendTime = time.Now().Add(time.Second * 5).Unix()
			//塞入消息
			B.MsgCh <- Msg{
				Time:     time.Now(),
				UserName: userName,
				MsgType:  msgTyp,
				Content:  content,
				Channel:  channel,
			}
		}

		fmt.Println("message from ", userName, " : ", string(content))
	}
}

// Broadcast 广播消息
func Broadcast() {
	for {
		select {
		case msg := <-B.MsgCh:

			room.Range(func(k, v interface{}) bool {
				client := v.(*Client)
				content := append(S2B(msg.Time.Format(time.Stamp)+msg.UserName+" 说 : \n"), msg.Content...)
				if client.Channel == msg.Channel {
					err := client.Connection.WriteMessage(msg.MsgType, content)
					if err != nil {
						log.Println("conn.WriteMessage err: ", err)
					}
					return true
				}
				return true
			})
		case msg := <-B.Register:
			room.Range(func(k, v interface{}) bool {
				client := v.(*Client)
				content := append(S2B(msg.Time.Format(time.Stamp) + " " + msg.UserName + " 进入了聊天室\n"))
				if client.Channel == msg.Channel {
					err := client.Connection.WriteMessage(msg.MsgType, content)
					if err != nil {
						log.Println("conn.WriteMessage err: ", err)
					}
					return true
				}
				return true
			})
		case msg := <-B.UnRegister:
			room.Range(func(k, v interface{}) bool {
				client := v.(*Client)
				content := append(S2B(msg.Time.Format(time.Stamp)+" "+msg.UserName+" "), msg.Content...)
				if client.Channel == msg.Channel {
					err := client.Connection.WriteMessage(msg.MsgType, content)
					if err != nil {
						log.Println("conn.WriteMessage err: ", err)
					}
					return true
				}
				return true
			})
		}
	}
}

func main() {
	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": 200,
			"data":   "successful",
		})
	})
	go Broadcast()
	r.GET("/ws", WsHandler)

	r.Run()

}

var num = 1

//读取参数
func getQuery(c *gin.Context) (string, string) {
	mux.Lock()
	userName := c.DefaultQuery("name", "newUser"+strconv.Itoa(num))
	num++
	mux.Unlock()
	//获取频道
	ch := c.DefaultQuery("ch", "1")
	return userName, ch
}

// S2B 高效string->[]byte
func S2B(str string) (bytes []byte) {
	x := *(*[2]uintptr)(unsafe.Pointer(&str))
	bytes = *(*[]byte)(unsafe.Pointer(&[3]uintptr{x[0], x[1], x[1]}))
	return
}
func B2S(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
