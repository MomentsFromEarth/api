package router

import (
	account "github.com/MomentsFromEarth/api/internal/controllers/account"
	moment "github.com/MomentsFromEarth/api/internal/controllers/moment"
	tag "github.com/MomentsFromEarth/api/internal/controllers/tag"

	"github.com/gin-gonic/gin"
)

// Init is the entrypoint of router
func Init(e *gin.Engine) {
	e.POST("/tag/:tag_id", tag.Create)
	e.GET("/tag/:tag_id", tag.Read)
	e.DELETE("/tag/:tag_id", tag.Delete)

	e.POST("/moment", moment.Create)
	e.GET("/moment/:moment_id", moment.Read)
	e.PUT("/moment/:moment_id", moment.Update)
	e.DELETE("/moment/:moment_id", moment.Delete)
	e.PUT("/moment/:moment_id/callback", moment.Callback)

	e.POST("/account", account.Create)
	e.GET("/account", account.Read)
	e.PUT("/account", account.Update)
	e.DELETE("/account", account.Delete)
	e.GET("/account/profile/:username", account.Profile)
}
