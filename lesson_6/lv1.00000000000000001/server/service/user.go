package service

import (
	"RedRock-Work/lesson_6/lv1.00000000000000001/server/dao"
	"RedRock-Work/lesson_6/lv1.00000000000000001/server/model"
	"RedRock-Work/lesson_6/lv1.00000000000000001/server/tool"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(user model.InputUser) model.RspMsg {
	u := model.User{
		Username: user.UserName,
		Password: user.Password,
	}
	fmt.Println("Register : ", u.Username)
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost) //加密处理
	if err != nil {
		fmt.Printf("error : %v", err)
		return tool.GetRsp(400, "服务器错误")
	}
	err = dao.SearchUserInfoByUserName(&u)
	if err != gorm.ErrRecordNotFound {
		if err == nil {
			return tool.GetRsp(400, "账号已被注册")
		} else {
			fmt.Printf("gorm error : %v", err)
			return tool.GetRsp(400, "服务器错误")
		}
	}
	u.Password = string(password)

	mysqlErr, redisErr := dao.InsertUser(u)
	if mysqlErr != nil || redisErr != nil {
		fmt.Printf("mysql error: %v \n redis error : %v", mysqlErr, redisErr)
		return tool.GetRsp(400, "服务器错误")
	} else {
		return tool.GetRsp(200, "注册成功")
	}
}

func Login(user model.InputUser) model.RspMsg {
	fmt.Println("Login : ", user)
	password := user.Password
	u := model.User{
		Username: user.UserName,
	}
	err := dao.Login(&u)
	if err != nil {
		if err.Error()[8:] == "nil returned" || err == gorm.ErrRecordNotFound {
			return tool.GetRsp(400, "账号不存在")
		}
		fmt.Printf("redis error : %v", err)
		return tool.GetRsp(400, "服务器错误")
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
		if err != nil {
			return tool.GetRsp(400, "密码错误")
		}
		return tool.GetRsp(200, "登录成功")
	}
}

func ChangePassword(user model.InputUser) model.RspMsg {
	fmt.Printf("ChangePassword : %v", user)
	u := model.User{
		Username: user.UserName,
	}
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost) //加密处理
	if err != nil {
		fmt.Printf("server error : %v", err)
		return tool.GetRsp(400, "服务器错误")
	}
	u.Password = string(password)
	err = dao.ChangePassword(u)
	if err != nil {
		fmt.Printf("change password error : %v", err)
		return tool.GetRsp(400, "服务器错误")
	}
	return tool.GetRsp(200, "修改成功")

}
