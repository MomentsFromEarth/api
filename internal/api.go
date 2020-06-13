package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Init is the entrypoint of api package
func Init() *gin.Engine {
	r := gin.Default()
	r.GET("/pingmetest", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pingmetest",
		})
	})
	return r
}

func init() {
	fmt.Println("api.init")
}
