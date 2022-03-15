package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/big"
	"net/http"
	"time"
)

type GitHubConfig struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
	Scope        string
	State        string
	Code         string
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type UserInfo struct {
	Login             string      `json:"login"`
	Id                int         `json:"id"`
	NodeId            string      `json:"node_id"`
	AvatarUrl         string      `json:"avatar_url"`
	GravatarId        string      `json:"gravatar_id"`
	Url               string      `json:"url"`
	HtmlUrl           string      `json:"html_url"`
	FollowersUrl      string      `json:"followers_url"`
	FollowingUrl      string      `json:"following_url"`
	GistsUrl          string      `json:"gists_url"`
	StarredUrl        string      `json:"starred_url"`
	SubscriptionsUrl  string      `json:"subscriptions_url"`
	OrganizationsUrl  string      `json:"organizations_url"`
	ReposUrl          string      `json:"repos_url"`
	EventsUrl         string      `json:"events_url"`
	ReceivedEventsUrl string      `json:"received_events_url"`
	Type              string      `json:"type"`
	SiteAdmin         bool        `json:"site_admin"`
	Name              string      `json:"name"`
	Company           interface{} `json:"company"`
	Blog              string      `json:"blog"`
	Location          interface{} `json:"location"`
	Email             interface{} `json:"email"`
	Hireable          interface{} `json:"hireable"`
	Bio               interface{} `json:"bio"`
	TwitterUsername   interface{} `json:"twitter_username"`
	PublicRepos       int         `json:"public_repos"`
	PublicGists       int         `json:"public_gists"`
	Followers         int         `json:"followers"`
	Following         int         `json:"following"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

var config = GitHubConfig{
	ClientId:     "15c711801e9cb68ad8d1",
	ClientSecret: "55ad2a951a7d1261155a3c1974799832afbcaa92",
	RedirectUrl:  "http://110.42.165.192:8090/authorization",
	Scope:        "user",
}

func main() {

	engine := gin.Default()
	//开启路由
	engine.GET("/", func(c *gin.Context) {
		state, err := CreateRandomString(20)
		if err != nil {
			fmt.Println(err)
			return
		}
		config.State = state
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s", config.ClientId, config.RedirectUrl, config.State))
		//c.Redirect(http.StatusMovedPermanently, "https://github.com/login/oauth/authorize?client_id=15c711801e9cb68ad8d1&scope=user&redirect_uri=http://127.0.0.1:8080/authorization")
	})
	engine.GET("/authorization", func(c *gin.Context) {
		code, flag := c.GetQuery("code")
		config.Code = code
		//打印code
		fmt.Println(code, flag)

		//获取token
		token, err := GetToken(code)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("token:", token)
		//这里进行一次打印token(为了好看)
		fmt.Println("accessToken:", token.AccessToken)
		fmt.Println("tokenType:", token.TokenType)
		fmt.Println("scope:", token.Scope)
		//通过accessToken拿取用户的信息
		Info, err := GetUserInfo(token.AccessToken)
		if err != nil {
			fmt.Println("Get UserInfo error:", err)
			return
		}
		userInfo := *Info

		fmt.Println(userInfo)

		c.JSON(200, gin.H{
			"Notice": "success",
			"Hello":  userInfo.Name,
			"Info:":  userInfo,
		})
	})
	_ = engine.Run(":8090")

}

func GetToken(code string) (TokenResponse, error) {
	token := new(TokenResponse)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%v&client_secret=%v&code=%v", config.ClientId, config.ClientSecret, code), nil)

	//可以改变header来改变传入的值形式,可以尝试绑定为json对象
	req.Header.Set("accept", "application/json")
	resp, err1 := client.Do(req)

	if err1 != nil {
		return *token, err1
	}
	defer resp.Body.Close()
	buf := make([]byte, 4096)
	for {
		//读取数据
		n, err2 := resp.Body.Read(buf)
		//读完就退出循环
		if n == 0 {
			break
		}
		if err2 != nil && err2 != io.EOF {
			return *token, err2
		}
		//绑定json对象
		err := json.Unmarshal(buf[:n], token)
		if err != nil {
			return *token, err
		}
	}
	return *token, nil
}

func GetUserInfo(accessToken string) (*UserInfo, error) {
	info := new(UserInfo)
	client := &http.Client{}

	request, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return info, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))

	resp, err := client.Do(request)
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	buf := make([]byte, 4096)
	for {
		//读取数据
		n, err2 := resp.Body.Read(buf)
		//读完就退出循环
		if n == 0 {
			break
		}
		if err2 != nil && err2 != io.EOF {

			return info, err2
		}
		//直接绑定对象
		err = json.Unmarshal(buf[:n], info)
		if err != nil {
			return info, err
		}
	}
	return info, nil
}

func CreateRandomString(len int) (string, error) {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, err := rand.Int(rand.Reader, bigInt)
		if err != nil {
			return "", err
		}
		container += string(str[randomInt.Int64()])
	}
	return container, nil
}
