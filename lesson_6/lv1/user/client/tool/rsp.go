package tool

import "github.com/gin-gonic/gin"

func RspErrWithData(c *gin.Context, message string) {
	c.JSON(200, gin.H{
		"status":  400,
		"message": message,
		"data":    nil,
	})
}

func RspWithData(c *gin.Context, status int32, message string, data interface{}) {
	c.JSON(200, gin.H{
		"status":  status,
		"message": message,
		"data":    data,
	})
}
