package moment

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create is an entrypoint of controller
func Create(c *gin.Context) {
	c.Next()
}

// Read is an entrypoint of controller
func Read(c *gin.Context) {
	c.Next()
}

// Update is an entrypoint of controller
func Update(c *gin.Context) {
	c.Next()
}

// Callback is an entrypoint of controller
func Callback(c *gin.Context) {
	momentID, _ := c.Params.Get("moment_id")
	c.JSON(http.StatusOK, gin.H{"echo": momentID})
	c.Next()
}

// Delete is an entrypoint of controller
func Delete(c *gin.Context) {
	c.Next()
}

func errResponse(c *gin.Context, status int, message string) {
	fmt.Printf("[MomentError] %v", message)
	c.AbortWithStatusJSON(status, gin.H{"message": message})
}
