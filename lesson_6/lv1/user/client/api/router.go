package api

import "github.com/gin-gonic/gin"

func InitRouter() {
	e := gin.Default()

	user := e.Group("/user")
	{
		user.POST("/register", Register)
		user.POST("/login", Login)
		user.PUT("/changePassword", ChangePassword)
	}

	e.Run()
}
