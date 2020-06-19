package app

import (
	"github.com/MomentsFromEarth/api/internal/auth"
	"github.com/MomentsFromEarth/api/internal/router"
	"github.com/gin-gonic/gin"
)

// Init is the entrypoint of app package
func Init() *gin.Engine {
	engine := gin.Default()
	auth.Init()
	router.Init(engine)
	return engine
}
