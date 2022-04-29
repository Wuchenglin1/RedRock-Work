package api

import (
	"RedRock-Work/lesson_6/lv1.00000000000000001/client/model"
	"RedRock-Work/lesson_6/lv1.00000000000000001/client/tool"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	RegisterFunc, err := tool.GetService("Register")
	if err != nil {
		tool.RespErrWithData(c, "服务器错误")
		return
	}
	u := model.User{
		UserName: c.PostForm("userName"),
		Password: c.PostForm("password"),
	}
	if u.UserName == "" {
		tool.RespErrWithData(c, "用户名不能为空")
		return
	}
	if len(u.UserName) < 6 || len(u.UserName) > 20 {
		tool.RespErrWithData(c, "用户名长度不正确")
		return
	}
	if u.Password == "" {
		tool.RespErrWithData(c, "密码不能为空")
		return
	}
	if len(u.Password) < 6 || len(u.Password) > 20 {
		tool.RespErrWithData(c, "密码长度不规范")
		return
	}
	m, err := RegisterFunc(u)
	if err != nil {
		tool.RespErrWithData(c, "服务器错误")
		return
	}
	msg := m.(map[string]interface{})

	tool.RespData(c, msg["Status"].(int), msg["Message"].(string), nil)
}

func Login(c *gin.Context) {
	LoginFunc, err := tool.GetService("Login")
	if err != nil {
		tool.RespErrWithData(c, "服务器错误")
		return
	}
	u := model.User{
		UserName: c.PostForm("userName"),
		Password: c.PostForm("password"),
	}
	if u.UserName == "" {
		tool.RespErrWithData(c, "用户名不能为空")
		return
	}
	if len(u.UserName) < 6 || len(u.UserName) > 20 {
		tool.RespErrWithData(c, "用户名长度不正确")
		return
	}
	if u.Password == "" {
		tool.RespErrWithData(c, "密码不能为空")
		return
	}
	if len(u.Password) < 6 || len(u.Password) > 20 {
		tool.RespErrWithData(c, "密码长度不规范")
		return
	}
	m, err := LoginFunc(u)
	if err != nil {
		fmt.Printf("get loginfunc error : %v", err)
		tool.RespErrWithData(c, "服务器错误")
		return
	}
	msg := m.(map[string]interface{})
	if msg["Status"].(int) == 200 {
		c.SetCookie("userName", u.UserName, 3600, "/", "", false, true)
	}
	tool.RespData(c, msg["Status"].(int), msg["Message"].(string), nil)

}

func ChangePassword(c *gin.Context) {
	changePasswordFunc, err := tool.GetService("ChangePassword")
	u := model.User{Password: c.PostForm("changePassword")}
	//这里本来应该是jwt鉴权，简单地用cookie表示下
	userName, err := c.Cookie("userName")
	if err != nil {
		tool.RespErrWithData(c, "您还没有登录")
		return
	}
	if u.Password == "" {
		tool.RespErrWithData(c, "新密码不能为空")
		return
	}
	if len(u.Password) < 6 || len(u.Password) > 20 {
		tool.RespErrWithData(c, "密码长度不正确")
		return
	}
	u.UserName = userName
	m, err := changePasswordFunc(u)
	if err != nil {
		fmt.Printf("get ChangePasswordFunc error : %v", err)
		tool.RespErrWithData(c, "服务器错误")
		return
	}
	msg := m.(map[string]interface{})
	tool.RespData(c, msg["Status"].(int), msg["Message"].(string), nil)
}
