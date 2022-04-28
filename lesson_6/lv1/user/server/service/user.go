package service

import (
	"RedRock-Work/lesson_6/lv1/user/server/dao"
	"RedRock-Work/lesson_6/lv1/user/server/grpc/user"
	"RedRock-Work/lesson_6/lv1/user/server/model"
	"RedRock-Work/lesson_6/lv1/user/server/tool"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(userInfo *user.ReqUser) (*user.Rsp, error) {
	rsp := new(user.Rsp)
	fmt.Println("Login : ", userInfo)
	password, err := bcrypt.GenerateFromPassword([]byte(userInfo.GetPassword()), bcrypt.DefaultCost) //加密处理
	if err != nil {
		return tool.GetRsp(rsp, "服务器错误", 400), err
	}
	u := model.UserInfo{
		Username: userInfo.GetUserName(),
		Password: string(password),
	}

	err = dao.SearchUserInfoByUserName(&u)
	if err == nil {
		return tool.GetRsp(rsp, "账号已被注册", 400), nil
	}

	mysqlErr, redisErr := dao.InsertUser(u)
	if mysqlErr != nil || redisErr != nil {
		fmt.Printf("mysql error: %v \n redis error : %v", mysqlErr, redisErr)
		return tool.GetRsp(rsp, "服务器错误", 400), err
	} else {

		return tool.GetRsp(rsp, "注册成功", 200), nil
	}
}

func Login(userInfo *user.ReqUser) (*user.Rsp, error) {
	rsp := new(user.Rsp)
	fmt.Println("Register : ", userInfo)
	password := userInfo.GetPassword()
	u := model.UserInfo{
		Username: userInfo.UserName,
	}
	err := dao.Login(&u)
	if err != nil {
		if err.Error()[8:] == "nil returned" {
			return tool.GetRsp(rsp, "账号不存在", 400), err
		}
		fmt.Printf("redis error : %v", err)
		return tool.GetRsp(rsp, "服务器错误", 400), err
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
		if err != nil {
			return tool.GetRsp(rsp, "密码错误", 400), nil
		}
		return tool.GetRsp(rsp, "登录成功", 200), nil
	}
}

func ChangePassword(userInfo *user.ReqUser) (*user.Rsp, error) {
	rsp := new(user.Rsp)
	u := model.UserInfo{
		Model:    gorm.Model{},
		Username: userInfo.UserName,
	}
	//这里正常流程应该是中间件判断token的，默认账号存在(虽然我这里图简便用的cookie)
	//err := dao.SearchUserInfoByUserName(&u)
	//if err != nil {
	//	return tool.GetRsp(rsp, "账号不存在", 400), nil
	//}
	password, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost) //加密处理
	if err != nil {
		return tool.GetRsp(rsp, "服务器错误", 400), err
	}
	u.Password = string(password)
	err = dao.ChangePassword(u)
	if err != nil {
		return tool.GetRsp(rsp, "service error", 400), err
	}
	return tool.GetRsp(rsp, "修改成功", 200), nil
}
