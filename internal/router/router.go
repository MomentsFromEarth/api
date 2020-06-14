package router

import (
	account "github.com/MomentsFromEarth/api/internal/controllers"
	"github.com/gin-gonic/gin"
)

// Init is the entrypoint of router
func Init(e *gin.Engine) {
	e.GET("/echo/:name", func(c *gin.Context) {
		name := c.Param("name")
		message := c.Query("message")
		email, _ := c.Get("email")
		c.JSON(200, gin.H{
			"email":   email,
			"name":    name,
			"message": message,
		})
	})
	e.POST("/account", account.Create)
	e.GET("/account", account.Read)
	e.PUT("/account", account.Update)
	e.DELETE("/account", account.Delete)
}
