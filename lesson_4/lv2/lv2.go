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
	//本来这里应该是设置偏移量的，但是感觉好像没啥需要
	Offset := strconv.Itoa(page)
	// 设置body
	form := url.Values{}
	form.Set("s", keyWords)
	form.Set("type", sType)
	form.Set("limit", Limit)
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
	// 发起请求
	response, reqErr := client.Do(request)
	// 错误处理
	if reqErr != nil {
		return "", reqErr
	}
	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	return string(resBody), nil
}
