package router

import (
	"github.com/MomentsFromEarth/api/internal/auth"
	account "github.com/MomentsFromEarth/api/internal/controllers/account"
	moment "github.com/MomentsFromEarth/api/internal/controllers/moment"
	tag "github.com/MomentsFromEarth/api/internal/controllers/tag"

	"github.com/gin-gonic/gin"
)

// Init is the entrypoint of router
func Init(e *gin.Engine) {
	a := e.Group("/")
	a.Use(auth.Run())
	{
		a.POST("/tag/:tag_id", tag.Create)
		a.GET("/tag/:tag_id", tag.Read)
		a.DELETE("/tag/:tag_id", tag.Delete)

		a.POST("/moment", moment.Create)
		a.GET("/moment/:moment_id", moment.Read)
		a.PUT("/moment/:moment_id", moment.Update)
		a.DELETE("/moment/:moment_id", moment.Delete)

		a.POST("/account", account.Create)
		a.GET("/account", account.Read)
		a.PUT("/account", account.Update)
		a.DELETE("/account", account.Delete)
		a.GET("/account/profile/:username", account.Profile)
	}
	e.PUT("/moment/:moment_id/callback", moment.Callback)
}
