package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//搜索歌单
func main() {
	//分别设置 关键词 类型 页码 数量
	keyWords := "冰与火之舞"
	sType := "1"
	page := 1
	limit := 1
	//这里返回的就是搜索结果
	res, err := Begin(keyWords, sType, page, limit)
	if err != nil {
		fmt.Println("请求失败：", err)
		return
	}
	fmt.Println(res)
}

func Begin(keyWords string, sType string, page int, limit int) (string, error) {
	// 创建客户端
	client := &http.Client{}
	// 格式化参数
	Limit := strconv.Itoa(limit)
	//这里可以简单设置下偏移量，可以达到具体哪一页哪一条数据
	//但是offset需要啥limit的倍数
	Offset := strconv.Itoa(page)
	// 拿取数据需要设置body
	form := url.Values{}
	//搜索词
	form.Set("s", keyWords)
	//搜索类型
	form.Set("type", sType)
	//返回的数据条数
	form.Set("limit", Limit)
	//偏移量
	form.Set("offset", Offset)
	body := strings.NewReader(form.Encode())
	// 创建请求
	//这里是http://music.163.com/api/search/get/这个api获取搜索歌单的信息
	request, _ := http.NewRequest("POST", "http://music.163.com/api/search/get/", body)
	//设置头部
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Cookie", "appver=2.0.2")
	request.Header.Set("Referer", "http://music.163.com")
	request.Header.Set("Content-Length", strconv.Itoa(body.Len()))

	//发起请求
	response, err := client.Do(request)

	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	//读取数据，返回
	resBody, err1 := ioutil.ReadAll(response.Body)
	if err1 != nil {
		return "", err1
	}
	return string(resBody), nil
}
