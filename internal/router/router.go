package router

import (
	account "github.com/MomentsFromEarth/api/internal/controllers"
	"github.com/gin-gonic/gin"
)

// Init is the entrypoint of router
func Init(e *gin.Engine) {
	e.POST("/account", account.Create)
	e.GET("/account", account.Read)
	e.PUT("/account", account.Update)
	e.DELETE("/account", account.Delete)
}
