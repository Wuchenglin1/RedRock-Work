package tool

import "github.com/gin-gonic/gin"

func RespErrWithData(c *gin.Context, message string) {
	c.JSON(200, gin.H{
		"status":  400,
		"message": message,
		"data":    nil,
	})
}

func RespData(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(200, gin.H{
		"status":  status,
		"message": message,
		"data":    data,
	})
}
