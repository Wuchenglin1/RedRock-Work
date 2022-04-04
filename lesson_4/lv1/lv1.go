package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type web struct {
	Name string
	URL  string
}

func main() {
	startTime := time.Now().Unix()
	url := "http://xiaodiaodaya.cn"
	//开始爬取主页的所有url
	res, err := Begin(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	//正则表达式查找到每一个段子的url
	reg := regexp.MustCompile("<a href=\"(.*?)</a>")
	webUrl := reg.FindAllStringSubmatch(res, -1)

	//将url和name存进map中
	m := make(map[int]web)
	index := 0
	for i := 0; i < len(webUrl); i++ {
		w := web{}
		arr := strings.Split(webUrl[i][1], "\">")
		//去除掉无用的一个网址
		k := strings.Index(arr[0], "target")
		if k != -1 {
			continue
		}
		FindStr := strings.Index(arr[0], ".html")
		if FindStr == -1 {
			str, err1 := Begin(url + arr[0])
			if err1 != nil {
				fmt.Println(err1)
				return
			}

			reg = regexp.MustCompile("<!--listS-->(.*?)<!--listE-->")
			allUrl := reg.FindAllStringSubmatch(str, -1)

			reg = regexp.MustCompile("<a href=\"(.*?)</a>")
			arr2 := reg.FindAllStringSubmatch(allUrl[0][1], -1)
			for j := range arr2 {
				arr3 := strings.Split(arr2[j][1], "\">")
				FindStr = strings.Index(arr3[0], ".html")
				if FindStr == -1 {
					continue
				}
				newWeb := web{}
				newWeb.URL = url + arr3[0]
				newWeb.Name = arr3[1]
				m[index] = newWeb
				index++
			}
			continue
		}
		w.URL = url + arr[0]
		w.Name = arr[1]
		m[index] = w
		index++
	}
	//遍历map，爬下所有的内容
	index = 0
	//Jock := make(map[int]string)

	file, err1 := os.OpenFile("jock.text", os.O_RDWR|os.O_CREATE, 0766)
	if err1 != nil {
		log.Println("io error :", err1)
		return
	}

	defer file.Close()

	for i := range m {
		//这里有个 http://xiaodiaodaya.cn17.html 莫名的网址钻出来捣乱我的程序，去死去死
		if m[i].URL == "http://xiaodiaodaya.cn17.html" {
			continue
		}
		res, err = Begin(m[i].URL)
		if err != nil {
			fmt.Println(err)
		}
		//获取详情内容
		reg = regexp.MustCompile("<!--listS-->(.*?)<!--listE-->")
		arr1 := reg.FindAllStringSubmatch(res, -1)
		strArr1 := strings.Split(arr1[0][1], "<br/><br/>")
		_, _ = file.WriteString(m[i].URL)
		_, _ = file.WriteString("\n")
		for k := range strArr1 {
			str := DeleteString1(strArr1[k])
			str = DeleteString2(str)
			str = DeleteString3(str)
			fmt.Println(str)
			_, err = file.WriteString(str)
			_, err1 = file.WriteString("\n")
			if err != nil || err1 != nil {
				fmt.Println(err, err1)
				return
			}
		}
		fmt.Println("这是第", i+1, "个网址:", m[i].URL)
	}
	endTime := time.Now().Unix()
	fmt.Println(endTime - startTime)
}

func Begin(url string) (string, error) {
	//创建http客户端
	client := &http.Client{}

	//发起GET请求
	req, _ := http.NewRequest("GET", url, nil)

	//设置请求头的UserAgent
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Host", "<calculated when request is sent>")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	req.Header.Add("Accept", "")
	req.Header.Add("Accept-Encoding", "")
	req.Header.Add("Connection", "keep-alive")

	//执行请求,获取数据
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	result := string(body)
	return result, nil

}

// DeleteString1 SplitString 去除掉内容中的<br/>
func DeleteString1(inStr string) (outStr string) {
	strArr := strings.Split(inStr, "<br/>")
	for i := range strArr {
		outStr += strArr[i]
	}
	return outStr
}

// DeleteString2 DetailString2 去除内容中的<font color='#FFFFFF'>xiaodiaodaya.cn</font>
func DeleteString2(inStr string) (outStr string) {
	strArr := strings.Split(inStr, "<font color='#FFFFFF'>xiaodiaodaya.cn</font>")
	for i := range strArr {
		outStr += strArr[i]
	}
	return outStr
}

// DeleteString3 删除内容中的图片
func DeleteString3(inStr string) (outStr string) {
	reg := regexp.MustCompile("<img src=\"(.*?)\"/>")
	strArr := reg.FindAllStringSubmatch(inStr, -1)
	if len(strArr) == 0 {
		return inStr
	}
	arr := strings.Split(inStr, strArr[0][0])
	for i := range arr {
		outStr += arr[i]
	}
	return
}
