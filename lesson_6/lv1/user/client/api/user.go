package api

import (
	"RedRock-Work/lesson_6/lv1/user/client/model"
	"RedRock-Work/lesson_6/lv1/user/client/service"
	"RedRock-Work/lesson_6/lv1/user/client/tool"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	u := model.UserInfo{
		Username: c.PostForm("userName"),
		Password: c.PostForm("password"),
	}
	if u.Username == "" {
		tool.RspErrWithData(c, "用户名不能为空")
		return
	}
	if len(u.Username) < 6 || len(u.Username) > 20 {
		tool.RspErrWithData(c, "用户名长度不合理")
		return
	}
	if u.Password == "" {
		tool.RspErrWithData(c, "密码不能为空")
		return
	}
	if len(u.Password) < 6 || len(u.Password) > 20 {
		tool.RspErrWithData(c, "密码长度不合理")
		return
	}
	status, message, err := service.Register(u)
	if err != nil {
		fmt.Println(err)
		tool.RspErrWithData(c, "服务器错误")
		return
	}
	tool.RspWithData(c, status, message, nil)
}

func Login(c *gin.Context) {
	u := model.UserInfo{
		Username: c.PostForm("userName"),
		Password: c.PostForm("password"),
	}
	if u.Username == "" {
		tool.RspErrWithData(c, "用户名不能为空")
		return
	}
	if len(u.Username) < 6 || len(u.Username) > 20 {
		tool.RspErrWithData(c, "用户名长度不合理")
		return
	}
	if u.Password == "" {
		tool.RspErrWithData(c, "密码不能为空")
		return
	}
	if len(u.Password) < 6 || len(u.Password) > 20 {
		tool.RspErrWithData(c, "密码长度不合理")
		return
	}
	status, message, err := service.Login(u)
	if err != nil {
		if err.Error()[8:] == "nil returned" {
			tool.RspErrWithData(c, "账号不存在")
			return
		}
		fmt.Println(err)
		tool.RspErrWithData(c, "服务器错误")
		return
	}
	if status == 200 {
		c.SetCookie("userName", u.Username, 3600, "/", "", false, true)
	}
	tool.RspWithData(c, status, message, nil)
}

func ChangePassword(c *gin.Context) {
	//这里需要changePassword即需要修改成的密码
	u := model.UserInfo{
		Password: c.PostForm("changePassword"),
	}
	//这里本来应该是jwt鉴权，简单地用cookie表示下
	userName, err := c.Cookie("userName")
	if err != nil {
		tool.RspErrWithData(c, "您还没有登录")
		return
	}
	if u.Password == "" {
		tool.RspErrWithData(c, "新密码不能为空")
		return
	}
	if len(u.Password) < 6 || len(u.Password) > 20 {
		tool.RspErrWithData(c, "密码长度不正确")
		return
	}
	u.Username = userName
	status, messsage, err := service.ChangePassword(u)
	if err != nil {
		tool.RspErrWithData(c, "服务器错误")
	}
	tool.RspWithData(c, status, messsage, nil)
}
