package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var conn redis.Conn

type User struct {
	UserName string
	Password string
	Nano     int64
}

type Config struct {
	Ip       string `json:"ip"`
	Password string `json:"password"`
}

type Goods struct {
	Id    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Label string  `json:"label"`
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
	conn = c

	//开启gin框架
	r := gin.Default()

	user := r.Group("/user")
	{
		user.POST("/register", Register)
		user.POST("/login", Login)
	}

	store := r.Group("/goods")
	{
		store.GET("/get", GetGoodsInfo)
		store.POST("/browse", BrowseGoods)
		store.POST("/buy", BuyGoods)
		//推荐商品
		store.GET("/recommended", Recommended)

	}

	r.Run()
}

func Register(c *gin.Context) {
	u := User{}
	//模拟注册
	u.UserName = c.PostForm("userName")
	if u.UserName == "" {
		Response(c, false, "用户名不能为空")
		return
	}
	u.Password = c.PostForm("password")
	if u.Password == "" {
		Response(c, false, "密码不能为空")
		return
	}
	t := time.Now().UnixNano()
	conn.Do("zadd", t, 0, "phone")
	conn.Do("zadd", t, 0, "fleece")
	conn.Do("zadd", t, 0, "book")
	conn.Do("zadd", t, 0, "food")
	conn.Do("zadd", t, 0, "cooker")
	conn.Do("zadd", t, 0, "medical")
	conn.Do("hset", u.UserName, t, u.Password)

	c.JSON(200, "注册成功")

}

func Login(c *gin.Context) {
	u := User{}
	//模拟注册
	u.UserName = c.PostForm("userName")
	if u.UserName == "" {
		Response(c, false, "用户名不能为空")
		return
	}
	u.Password = c.PostForm("password")
	if u.Password == "" {
		Response(c, false, "密码不能为空")
		return
	}
	data, err := redis.Values(conn.Do("hgetall", u.UserName))
	if err != nil {
		fmt.Println(err)
		return
	}
	var arr = make([]string, 2)
	for k, v := range data {
		arr[k] = string(v.([]byte))
	}
	if arr[0] == "" {
		Response(c, false, "账号不存在！")
		return
	}

	if u.Password != arr[1] {
		Response(c, false, "您输入的密码有误")
		return
	}

	c.SetCookie("userName", u.UserName, 3600, "/", "", false, true)
	c.SetCookie("nano", arr[0], 3600, "/", "", false, true)
	Response(c, true, "登录成功")
}

func GetGoodsInfo(c *gin.Context) {
	var u User
	userName, err := c.Cookie("userName")
	if err != nil || userName == "" {
		Response(c, false, "您还没有登录")
		fmt.Println(err)
		return
	}
	u.UserName = userName

	key, err1 := redis.Values(conn.Do("keys", "*"))
	if err1 != nil {
		Response(c, false, "服务器错误")
		fmt.Println(err1)
		return
	}
	var g Goods
	m := make(map[int]Goods)
	for k, v := range key {
		value := string(v.([]byte))
		reg := regexp.MustCompile("[^\\d]{1}")
		result := reg.FindAllStringSubmatch(value, -1)
		if len(result) != 0 {
			continue
		}
		res, err3 := redis.Bytes(conn.Do("get", value))
		if err3 != nil {
			fmt.Println(err)
			return
		}
		arr := strings.Split(string(res), ";")
		price, err4 := strconv.ParseFloat(arr[1], 64)
		if err4 != nil {
			fmt.Println(err4)
			Response(c, false, "服务器错误")
			return
		}
		g.Name = arr[0]
		g.Price = price
		g.Id = value
		g.Label = arr[2]
		m[k] = g
	}
	Response(c, false, m)
}

func BrowseGoods(c *gin.Context) {
	var u User
	userName, err := c.Cookie("userName")
	nano, err1 := c.Cookie("nano")
	if err1 != nil || userName == "" || nano == "" {
		Response(c, false, "您还没有登录")
		fmt.Println(err1)
		return
	}
	u.UserName = userName
	u.Nano, err = strconv.ParseInt(nano, 10, 64)
	if err != nil {
		Response(c, false, "您还没有登录")
		fmt.Println(err)
		return
	}

	gid := c.PostForm("gid")
	if gid == "" {
		Response(c, false, "您还没有输入gid")
		return
	}

	res, err1 := redis.Bytes(conn.Do("get", gid))
	if err1 != nil {
		Response(c, false, "该gid不存在！")
		return
	}
	var g Goods
	arr := strings.Split(string(res), ";")
	price, err2 := strconv.ParseFloat(arr[1], 64)
	if err2 != nil {
		fmt.Println(err2)
		Response(c, false, "服务器错误")
		return
	}
	g.Id = gid
	g.Name = arr[0]
	g.Price = price
	g.Label = arr[2]

	//这里就给这个gid的商品加分
	_, err = conn.Do("zincrby", u.Nano, 1, g.Label)
	if err != nil {
		fmt.Println(err)
		Response(c, false, "服务器错误")
		return
	}
	Response(c, true, g)
}

func BuyGoods(c *gin.Context) {
	var u User
	userName, err := c.Cookie("userName")
	nano, err1 := c.Cookie("nano")
	if err1 != nil || userName == "" || nano == "" {
		Response(c, false, "您还没有登录")
		fmt.Println(err1)
		return
	}
	u.UserName = userName
	u.Nano, err = strconv.ParseInt(nano, 10, 64)
	if err != nil {
		Response(c, false, "您还没有登录")
		fmt.Println(err)
		return
	}

	gid := c.PostForm("gid")
	if gid == "" {
		Response(c, false, "gid不能为空")
		return
	}

	res, err3 := redis.Bytes(conn.Do("get", gid))
	if err3 != nil {
		Response(c, false, "该gid不存在！")
		return
	}
	var g Goods
	arr := strings.Split(string(res), ";")
	price, err2 := strconv.ParseFloat(arr[1], 64)
	if err2 != nil {
		fmt.Println(err2)
		Response(c, false, "服务器错误")
		return
	}
	g.Id = gid
	g.Name = arr[0]
	g.Price = price
	g.Label = arr[2]

	//这里就给这个gid的商品加分
	_, err = conn.Do("zincrby", u.Nano, 5, g.Label)
	if err != nil {
		fmt.Println(err)
		Response(c, false, "服务器错误")
		return
	}

	Response(c, false, "购买"+gid+"成功")
}

func Recommended(c *gin.Context) {
	account := 0
	var u User
	userName, err := c.Cookie("userName")
	nano, err1 := c.Cookie("nano")
	if err1 != nil || userName == "" || nano == "" {
		Response(c, false, "您还没有登录")
		fmt.Println(err1)
		return
	}
	u.UserName = userName
	u.Nano, err = strconv.ParseInt(nano, 10, 64)
	if err != nil {
		Response(c, false, "您还没有登录")
		fmt.Println(err)
		return
	}

	value, err2 := redis.Values(conn.Do("zrevrange", "1648209935170079100", 0, -1))
	if err != nil {
		fmt.Println(err2)
		return
	}
	arr := make([]string, 6)
	for k, v := range value {
		arr[k] = string(v.([]byte))
	}
	//先声明一个装goodsInfo的map
	goodsMap := make(map[int]Goods)
	//然后就按照分数排列来推送不同分类的商品
	for i := len(arr) - 1; i >= 0; i-- {
		//先将每个类别的商品存储到一个map里面
		v, err3 := redis.Values(conn.Do("zrangebyscore", arr[i], "-inf", "+inf"))
		if err3 != nil {
			fmt.Println(err3)
			Response(c, false, "服务器错误")
			return
		}
		//存取该个类别的所有商品
		m1 := make(map[int]string)
		for k, v2 := range v {
			m1[k] = string(v2.([]byte))
		}
		//再推送i个该类别商品商品
		v1, err4 := redis.Values(conn.Do("zrangebyscore", arr[i], "-inf", "+inf"))
		if err4 != nil {
			fmt.Println(err4)
			Response(c, false, "服务器错误")
			return
		}
		//存取该个类别的所有商品
		m := make(map[int]string)
		for k, v2 := range v1 {
			m[k] = string(v2.([]byte))
		}

		//随机抽该类别商品
		num := random(0, len(m)-1, i)

		for j := range num {

			res, err := redis.Bytes(conn.Do("get", m[j]))
			if err != nil {
				fmt.Println(err)
				Response(c, false, "服务器错误")
				return
			}
			arr1 := strings.Split(string(res), ";")
			price, err := strconv.ParseFloat(arr1[1], 64)
			if err != nil {
				fmt.Println(err)
				Response(c, false, "服务器错误")
				return
			}
			fmt.Println(arr1[0], price, arr1[2])

			g := Goods{
				Name:  arr1[0],
				Price: price,
				Label: arr1[2],
			}
			goodsMap[account] = g
			account++
		}
	}
	Response(c, true, goodsMap)
}

func Response(c *gin.Context, status bool, data interface{}) {
	c.JSON(200, gin.H{
		"status": status,
		"data":   data,
	})
}

func GetConfig() *Config {
	var cfg *Config
	file, err := os.Open("../config.json")
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

func random(start int, end int, count int) []int {
	//范围检查
	if end < start || (end-start) < count {
		return nil
	}

	//存放结果的slice
	nums := make([]int, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		//生成随机数
		num := r.Intn(end-start) + start

		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
		}
	}

	return nums
}
